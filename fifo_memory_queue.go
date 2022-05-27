package queue

import (
	"context"
)

func NewFifoMemoryQueue(sizes ...int) Queue {
	ctx, cancel := context.WithCancel(context.Background())
	size := 1024
	if len(sizes) > 0 {
		size = sizes[0]
	}
	return &FifoMemoryQueue{
		queue:  make(chan []byte, size),
		ctx:    ctx,
		cancel: cancel,
	}
}

var _ Queue = (*FifoMemoryQueue)(nil)

type FifoMemoryQueue struct {
	queue  chan []byte
	ctx    context.Context
	cancel context.CancelFunc
}

func (q *FifoMemoryQueue) Get(ctx context.Context) ([]byte, error) {
	select {
	case <-q.ctx.Done():
		return nil, ErrQueueClosed
	default:
	}
	if ctx == nil {
		select {
		case data := <-q.queue:
			return data, nil
		default:
			return nil, ErrQueueEmpty
		}
	}
	select {
	case <-q.ctx.Done():
		return nil, ErrQueueClosed
	case <-ctx.Done():
		return nil, ctx.Err()
	case data := <-q.queue:
		return data, nil
	}
}

func (q *FifoMemoryQueue) Put(ctx context.Context, data []byte) error {
	select {
	case <-q.ctx.Done():
		return ErrQueueClosed
	default:
	}
	if ctx == nil {
		select {
		case q.queue <- data:
			return nil
		default:
			return ErrQueueFull
		}
	}
	select {
	case <-q.ctx.Done():
		return ErrQueueClosed
	case <-ctx.Done():
		return ctx.Err()
	case q.queue <- data:
		return nil
	}
}

func (q *FifoMemoryQueue) Close() error {
	q.cancel()
	return nil
}

func (q *FifoMemoryQueue) Len() int {
	return len(q.queue)
}
