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
	ErrClosed                 = errors.New("queue closed")                   // 队列已关闭
	ErrEmptyQueue             = errors.New("queue empty")                    // 空队列
	ErrFullQueue              = errors.New("queue full")                     // 满队列
	ErrChunkSizeInconsistency = errors.New("queue chunk size inconsistency") // 仅 fifo-disk-queue 校验时返回
)
