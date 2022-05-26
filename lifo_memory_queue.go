package queue

import (
    "context"
    "sync"
)

func NewLifoMemoryQueue(sizes ...int) Queue {
    ctx, cancel := context.WithCancel(context.Background())
    size := 1024
    if len(sizes) > 0 {
        size = sizes[0]
    }
    return &LifoMemoryQueue{
        queue:     make([][]byte, size, size),
        ctx:       ctx,
        cancel:    cancel,
        putNotify: make(chan struct{}, size),
        getNotify: make(chan struct{}, size),
    }
}

var _ Queue = (*LifoMemoryQueue)(nil)

type LifoMemoryQueue struct {
    queue  [][]byte
    ctx    context.Context
    cancel context.CancelFunc
    lock   sync.Mutex
    index  int
    getNotify chan struct{}
    putNotify chan struct{}
}

func (q *LifoMemoryQueue) Get(ctx context.Context) ([]byte, error) {
    select {
    case <-q.ctx.Done():
        return nil, ErrQueueClosed
    default:
    }
    if ctx != nil {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-q.ctx.Done():
            return nil, ErrQueueClosed
        case <-q.getNotify: // 阻塞，则表示当前队列无数据
        }
    }
    q.lock.Lock()
    defer q.lock.Unlock()
    if q.index <= 0 {
        return nil, ErrQueueEmpty
    }
    q.index--
    data := q.queue[q.index]
    q.queue[q.index] = nil
    <-q.putNotify
    if ctx == nil {
        <-q.getNotify
    }
    return data, nil
}

func (q *LifoMemoryQueue) Put(ctx context.Context, data []byte) error {
    select {
    case <-q.ctx.Done():
        return ErrQueueClosed
    default:
    }
    if ctx != nil {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-q.ctx.Done():
            return ErrQueueClosed
        case q.putNotify <- struct{}{}:
        }
    }
    q.lock.Lock()
    defer q.lock.Unlock()
    if q.index >= len(q.queue) {
        return ErrQueueFull
    }
    q.queue[q.index] = data
    q.index++
    q.getNotify <- struct{}{}
    if ctx == nil {
        q.putNotify <- struct{}{}
    }
    return nil
}

func (q *LifoMemoryQueue) Close() error {
    q.cancel()
    return nil
}

func (q *LifoMemoryQueue) Len() int {
    return q.index
}
