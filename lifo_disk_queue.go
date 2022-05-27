package queue

import (
    "bytes"
    "context"
    "encoding/binary"
    "fmt"
    "io"
    "os"
    "strconv"
    "sync"
)

func NewLifoDiskQueue(file string) (Queue, error) {
    var err error
    ctx, cancel := context.WithCancel(context.Background())
    queue := LifoDiskQueue{
        ctx:       ctx,
        cancel:    cancel,
    }
    queue.file, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE, os.ModePerm)
    if err != nil {
        return nil, err
    }
    stat, err := queue.file.Stat()
    if err != nil {
        return nil, err
    }
    if stat.Size() == 0 {
        return &queue, nil
    }
    _, err = queue.file.Seek(-4, io.SeekEnd)
    if err != nil {
        return nil, err
    }
    buf := make([]byte, 4)
    _, err = queue.file.Read(buf)
    if err != nil {
        return nil, err
    }
    var length int32
    err = binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &length)
    if err != nil {
        return nil, err
    }
    offset, err := queue.file.Seek(int64(-4-length), io.SeekCurrent)
    if err != nil {
        return nil, err
    }
    buf = make([]byte, length)
    _, err = queue.file.Read(buf)
    if err != nil {
        return nil, err
    }
    queue.index, err = strconv.Atoi(string(buf))
    if err != nil {
        return nil, err
    }
    err = queue.file.Truncate(offset)
    if err != nil {
        return nil, err
    }
    _, err = queue.file.Seek(0, io.SeekEnd)
    if err != nil {
        return nil, err
    }
    return &queue, nil
}

type LifoDiskQueue struct {
    index     int
    file      *os.File
    lock      sync.Mutex
    ctx       context.Context
    cancel    context.CancelFunc
}

func (q *LifoDiskQueue) Get(ctx context.Context) ([]byte, error) {
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
    _, err := q.file.Seek(-4, io.SeekCurrent)
    if err != nil {
        return nil, err
    }
    buf := make([]byte, 4)
    _, err = q.file.Read(buf)
    if err != nil {
        return nil, err
    }
    var length int32
    err = binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &length)
    if err != nil {
        return nil, err
    }
    _, err = q.file.Seek(int64(-4-length), io.SeekCurrent)
    if err != nil {
        return nil, err
    }
    buf = make([]byte, length)
    _, err = q.file.Read(buf)
    if err != nil {
        return nil, err
    }
    _, err = q.file.Seek(int64(-length), io.SeekCurrent)
    if err != nil {
      return nil, err
    }
    q.index--
    return buf, nil
}

func (q *LifoDiskQueue) Put(ctx context.Context, data []byte) error {
    select {
    case <-q.ctx.Done():
        return ErrQueueClosed
    default:
    }
    q.lock.Lock()
    defer q.lock.Unlock()
    buf := new(bytes.Buffer)
    _ = binary.Write(buf, binary.BigEndian, int32(len(data)))
    _, err := q.file.Write(bytes.Join([][]byte{data, buf.Bytes()}, []byte("")))
    if err != nil {
        return err
    }
    q.index++
    return nil
}

func (q *LifoDiskQueue) Close() error {
    select {
    case <-q.ctx.Done():
        return nil
    default:
    }
    q.lock.Lock()
    defer q.lock.Unlock()
    q.cancel()
    defer q.file.Close()
    if q.index < 1 {
        return q.file.Truncate(0)
    }
    offset, err := q.file.Seek(0, io.SeekCurrent)
    if err != nil {
        return err
    }
    err = q.file.Truncate(offset)
    if err != nil {
        return err
    }
    data := fmt.Sprintf("%d", q.index)
    buf := new(bytes.Buffer)
    _ = binary.Write(buf, binary.BigEndian, int32(len(data)))
    _, err = q.file.Write(bytes.Join([][]byte{[]byte(data), buf.Bytes()}, []byte("")))
    return err
}

func (q *LifoDiskQueue) Len() int {
    return q.index
}
