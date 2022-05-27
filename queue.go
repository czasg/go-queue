package queue

import (
    "context"
    "errors"
)

var (
    ErrQueueClosed = errors.New("queue closed")
    ErrQueueEmpty  = errors.New("queue empty")
    ErrQueueFull   = errors.New("queue full")
)

type Queue interface {
    Get(ctx context.Context) ([]byte, error)
    Put(ctx context.Context, data []byte) error
    Len() int
    Close() error
}
