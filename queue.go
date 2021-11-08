package queue

import "errors"

type Queue interface {
	Push(data []byte) error
	Pop() ([]byte, error)
	Close() error
	Len() int
}

type PriorityQueue interface {
	Push(data []byte, priority int) error
	Pop() ([]byte, error)
	Close() error
	Len() int
}

var (
	ErrClosed                 = errors.New("queue closed")
	ErrEmptyQueue             = errors.New("queue empty")
	ErrFullQueue              = errors.New("queue full")
	ErrChunkSizeInconsistency = errors.New("queue chunk size inconsistency")
)
