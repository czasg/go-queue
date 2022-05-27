package queue

import (
    "bytes"
    "context"
    "encoding/binary"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
    "sync"
)

func NewFifoDiskQueue(file string) (Queue, error) {
    var err error
    ctx, cancel := context.WithCancel(context.Background())
    queue := FifoDiskQueue{
        ctx:    ctx,
        cancel: cancel,
    }
    queue.writeFile, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE, os.ModePerm)
    if err != nil {
        return nil, err
    }
    queue.readFile, err = os.OpenFile(file, os.O_RDONLY, os.ModePerm)
    if err != nil {
        return nil, err
    }
    stat, err := queue.readFile.Stat()
    if err != nil {
        return nil, err
    }
    if stat.Size() == 0 {
        return &queue, nil
    }
    _, err = queue.writeFile.Seek(-4, io.SeekEnd)
    if err != nil {
        return nil, err
    }
    buf := make([]byte, 4)
    _, err = queue.writeFile.Read(buf)
    if err != nil {
        return nil, err
    }
    var length int32
    err = binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &length)
    if err != nil {
        return nil, err
    }
    offset, err := queue.writeFile.Seek(int64(-4-length), io.SeekCurrent)
    if err != nil {
        return nil, err
    }
    buf = make([]byte, length)
    _, err = queue.writeFile.Read(buf)
    if err != nil {
        return nil, err
    }
    bufs := strings.Split(string(buf), ",")
    if len(bufs) != 2 {
        return nil, fmt.Errorf("%s 格式异常", string(buf))
    }
    indexString := bufs[0]
    offsetString := bufs[1]
    queue.index, err = strconv.Atoi(indexString)
    if err != nil {
        return nil, err
    }
    queue.offset, err = strconv.Atoi(offsetString)
    if err != nil {
        return nil, err
    }
    err = queue.writeFile.Truncate(offset)
    if err != nil {
        return nil, err
    }
    _, err = queue.writeFile.Seek(0, io.SeekEnd)
    if err != nil {
        return nil, err
    }
    _, err = queue.readFile.Seek(int64(queue.offset), io.SeekStart)
    if err != nil {
        return nil, err
    }
    return &queue, nil
}

var _ Queue = (*FifoDiskQueue)(nil)

type FifoDiskQueue struct {
    index     int
    offset    int
    readFile  *os.File
    writeFile *os.File
    lock      sync.Mutex
    ctx       context.Context
    cancel    context.CancelFunc
}

func (q *FifoDiskQueue) Get(ctx context.Context) ([]byte, error) {
    select {
    case <-q.ctx.Done():
        return nil, ErrQueueClosed
    default:
    }
    q.lock.Lock()
    defer q.lock.Unlock()
    if q.index <= 0 {
        return nil, ErrQueueEmpty
    }
    buf := make([]byte, 4)
    _, err := q.readFile.Read(buf)
    if err != nil {
        return nil, err
    }
    var length int32
    err = binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &length)
    if err != nil {
        return nil, err
    }
    buf = make([]byte, length)
    _, err = q.readFile.Read(buf)
    if err != nil {
        return nil, err
    }
    q.index--
    q.offset += int(length) + 4
    return buf, nil
}

func (q *FifoDiskQueue) Put(ctx context.Context, data []byte) error {
    select {
    case <-q.ctx.Done():
        return ErrQueueClosed
    default:
    }
    q.lock.Lock()
    defer q.lock.Unlock()
    buf := new(bytes.Buffer)
    _ = binary.Write(buf, binary.BigEndian, int32(len(data)))
    _, err := q.writeFile.Write(bytes.Join([][]byte{buf.Bytes(), data}, []byte("")))
    if err != nil {
        return err
    }
    q.index++
    return nil
}

func (q *FifoDiskQueue) Close() error {
    select {
    case <-q.ctx.Done():
        return nil
    default:
    }
    q.lock.Lock()
    defer q.lock.Unlock()
    q.cancel()
    defer func() {
        q.readFile.Close()
        q.writeFile.Close()
    }()
    if q.index < 1 {
        return q.writeFile.Truncate(0)
    }
    offset, err := q.writeFile.Seek(0, io.SeekCurrent)
    if err != nil {
        return err
    }
    err = q.writeFile.Truncate(offset)
    if err != nil {
        return err
    }
    data := fmt.Sprintf("%d,%d", q.index, q.offset)
    buf := new(bytes.Buffer)
    _ = binary.Write(buf, binary.BigEndian, int32(len(data)))
    _, err = q.writeFile.Write(bytes.Join([][]byte{[]byte(data), buf.Bytes()}, []byte("")))
    return err
}

func (q *FifoDiskQueue) Len() int {
    return q.index
}
