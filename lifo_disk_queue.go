package queue

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const (
	LIFO_STAT = "lifo.stat.json"
	LIFO_DATA = "lifo.data"
)

func NewLifoDiskQueue(dir string) (Queue, error) {
	dir = filepath.Join(dir, "_lifo")
	queue := LifoDiskQueue{Dir: dir}
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	stat, err := queue.openStat()
	if err != nil {
		return nil, err
	}
	if stat == nil {
		stat = &LifoStat{}
	}
	queue.Stat = stat
	head, err := queue.openHead()
	if err != nil {
		return nil, err
	}
	_, err = head.Seek(0, os.SEEK_END)
	if err != nil {
		return nil, err
	}
	queue.HeadFile = head
	return &queue, nil
}

var _ Queue = (*LifoDiskQueue)(nil)

type LifoDiskQueue struct {
	Dir      string
	Stat     *LifoStat
	HeadFile *os.File
	Lock     sync.Mutex
	Closed   bool
}

type LifoStat struct {
	Size   int
	Offset int
}

func (q *LifoDiskQueue) Push(data []byte) error {
	if q.Closed {
		return ErrClosed
	}
	q.Lock.Lock()
	defer q.Lock.Unlock()
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(len(data)))
	if err != nil {
		return err
	}
	_, err = q.HeadFile.Write(bytes.Join([][]byte{data, buf.Bytes()}, []byte("")))
	if err != nil {
		return err
	}
	q.Stat.Size++
	q.Stat.Offset += 4 + len(data)
	return nil
}

func (q *LifoDiskQueue) Pop() ([]byte, error) {
	if q.Closed {
		return nil, ErrClosed
	}
	if q.Stat.Size < 1 {
		return nil, ErrEmptyQueue
	}
	q.Lock.Lock()
	defer q.Lock.Unlock()
	_, err := q.HeadFile.Seek(-4, os.SEEK_CUR)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 4)
	_, err = q.HeadFile.Read(buf)
	if err != nil {
		return nil, err
	}
	var length int32
	err = binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}
	_, err = q.HeadFile.Seek(int64(-4-length), os.SEEK_CUR)
	if err != nil {
		return nil, err
	}
	buf = make([]byte, length)
	_, err = q.HeadFile.Read(buf)
	if err != nil {
		return nil, err
	}
	_, err = q.HeadFile.Seek(int64(-length), os.SEEK_CUR)
	if err != nil {
		return nil, err
	}
	q.Stat.Size--
	q.Stat.Offset -= 4 + int(length)
	return buf, nil
}

func (q *LifoDiskQueue) Close() error {
	q.Lock.Lock()
	defer q.Lock.Unlock()
	if q.Closed {
		return nil
	}
	q.Closed = true
	// Truncate file
	err := q.HeadFile.Truncate(int64(q.Stat.Offset))
	if err != nil {
		return err
	}
	// Close file
	err = q.HeadFile.Close()
	if err != nil {
		return err
	}
	// Save stat file
	err = q.saveStat()
	if err != nil {
		return err
	}
	if q.Stat.Size > 0 {
		return nil
	}
	return os.RemoveAll(q.Dir)
}

func (q *LifoDiskQueue) Len() int {
	return q.Stat.Size
}

func (q *LifoDiskQueue) openStat() (*LifoStat, error) {
	statJsonPath := filepath.Join(q.Dir, LIFO_STAT)
	_, err := os.Stat(statJsonPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	f, err := os.Open(statJsonPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var stat LifoStat
	err = json.Unmarshal(body, &stat)
	if err != nil {
		return nil, err
	}
	return &stat, nil
}

func (q *LifoDiskQueue) saveStat() error {
	statJsonPath := filepath.Join(q.Dir, LIFO_STAT)
	f, err := os.Create(statJsonPath)
	if err != nil {
		return err
	}
	defer f.Close()
	body, _ := json.MarshalIndent(q.Stat, "", "  ")
	_, err = f.Write(body)
	return err
}

func (q *LifoDiskQueue) openHead() (*os.File, error) {
	dataPath := filepath.Join(q.Dir, LIFO_DATA)
	f, err := os.OpenFile(dataPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return f, nil
}
