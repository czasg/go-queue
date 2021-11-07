package queue

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

const (
	FIFO_STAT = "fifo.stat.json"
	FIFO_DATA = "fifo.data.%04d"
)

func NewFifoDiskQueue(dir string) (Queue, error) {
	return NewFifoDiskQueueWithChunk(dir, 50000)
}

func NewFifoDiskQueueWithChunk(dir string, chunk int) (Queue, error) {
	dir = filepath.Join(dir, "_fifo")
	queue := FifoDiskQueue{Dir: dir}
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	stat, err := queue.openStat()
	if err != nil {
		return nil, err
	}
	if stat == nil {
		stat = &FifoStat{ChunkSize: chunk}
	}
	if stat.ChunkSize != chunk {
		return nil, ErrChunkSizeInconsistency
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
	tail, err := queue.openTail()
	if err != nil {
		return nil, err
	}
	_, err = tail.Seek(int64(queue.Stat.Tail.Offset), os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	queue.TailFile = tail
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		_ = queue.Close()
	}()
	return &queue, nil
}

var _ Queue = (*FifoDiskQueue)(nil)

type FifoDiskQueue struct {
	Dir      string
	Stat     *FifoStat
	HeadFile *os.File
	TailFile *os.File
	Lock     sync.Mutex
	Closed   bool
}

type FifoStat struct {
	Size      int
	ChunkSize int
	Head      struct {
		Index int
		Count int
	}
	Tail struct {
		Index  int
		Count  int
		Offset int
	}
}

func (q *FifoDiskQueue) Push(data []byte) error {
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
	_, err = q.HeadFile.Write(bytes.Join([][]byte{buf.Bytes(), data}, []byte("")))
	if err != nil {
		return err
	}
	q.Stat.Head.Count++
	if q.Stat.Head.Count >= q.Stat.ChunkSize {
		q.Stat.Head.Index++
		q.Stat.Head.Count = 0
		_ = q.HeadFile.Close()
		head, err := q.openHead()
		if err != nil {
			return err
		}
		q.HeadFile = head
	}
	q.Stat.Size++
	return err
}

func (q *FifoDiskQueue) Pop() ([]byte, error) {
	if q.Closed {
		return nil, ErrClosed
	}
	if q.Stat.Tail.Index >= q.Stat.Head.Index && q.Stat.Tail.Count >= q.Stat.Head.Count {
		return nil, ErrEmptyQueue
	}
	q.Lock.Lock()
	defer q.Lock.Unlock()
	buf := make([]byte, 4)
	_, err := q.TailFile.Read(buf)
	if err != nil {
		return nil, err
	}
	var length int32
	err = binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}
	buf = make([]byte, length)
	_, err = q.TailFile.Read(buf)
	if err != nil {
		return nil, err
	}
	q.Stat.Tail.Count++
	q.Stat.Tail.Offset += 4 + int(length)
	if q.Stat.Tail.Count == q.Stat.ChunkSize && q.Stat.Tail.Index <= q.Stat.Head.Index {
		_ = q.TailFile.Close()
		_ = os.Remove(filepath.Join(q.Dir, fmt.Sprintf(FIFO_DATA, q.Stat.Tail.Index)))
		q.Stat.Tail.Count = 0
		q.Stat.Tail.Offset = 0
		q.Stat.Tail.Index++
		tail, err := q.openTail()
		if err != nil {
			return nil, err
		}
		q.TailFile = tail
	}
	q.Stat.Size--
	return buf, nil
}

func (q *FifoDiskQueue) Close() error {
	q.Lock.Lock()
	defer q.Lock.Unlock()
	if q.Closed {
		return nil
	}
	q.Closed = true
	_ = q.HeadFile.Close()
	_ = q.TailFile.Close()
	err := q.saveStat()
	if err != nil {
		return err
	}
	if q.Stat.Size > 0 {
		return nil
	}
	return os.RemoveAll(q.Dir)
}

func (q *FifoDiskQueue) openStat() (*FifoStat, error) {
	statJsonPath := filepath.Join(q.Dir, FIFO_STAT)
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
	var stat FifoStat
	err = json.Unmarshal(body, &stat)
	if err != nil {
		return nil, err
	}
	return &stat, nil
}

func (q *FifoDiskQueue) saveStat() error {
	statJsonPath := filepath.Join(q.Dir, FIFO_STAT)
	f, err := os.Create(statJsonPath)
	if err != nil {
		return err
	}
	defer f.Close()
	body, err := json.MarshalIndent(q.Stat, "", "  ")
	if err != nil {
		return err
	}
	_, err = f.Write(body)
	return err
}

func (q *FifoDiskQueue) openData(index, flag int) (*os.File, error) {
	dataPath := filepath.Join(q.Dir, fmt.Sprintf(FIFO_DATA, index))
	f, err := os.OpenFile(dataPath, flag, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (q *FifoDiskQueue) openHead() (*os.File, error) {
	return q.openData(q.Stat.Head.Index, os.O_RDWR|os.O_CREATE|os.O_APPEND)
}

func (q *FifoDiskQueue) openTail() (*os.File, error) {
	return q.openData(q.Stat.Tail.Index, os.O_RDONLY)
}
