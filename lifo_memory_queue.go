package queue

import "sync"

func NewLifoMemoryQueue(size int) Queue {
	return &LifoMemoryQueue{Q: make([][]byte, size, size)}
}

var _ Queue = (*LifoMemoryQueue)(nil)

type LifoMemoryQueue struct {
	Q      [][]byte
	Index  int
	Lock   sync.Mutex
	Closed bool
}

func (q *LifoMemoryQueue) Push(data []byte) error {
	if q.Closed {
		return ErrClosed
	}
	q.Lock.Lock()
	defer q.Lock.Unlock()
	if q.Index >= len(q.Q) {
		return ErrFullQueue
	}
	q.Q[q.Index] = data
	q.Index++
	return nil
}

func (q *LifoMemoryQueue) Pop() ([]byte, error) {
	if q.Closed {
		return nil, ErrClosed
	}
	q.Lock.Lock()
	defer q.Lock.Unlock()
	if q.Index <= 0 {
		return nil, ErrEmptyQueue
	}
	q.Index--
	data := q.Q[q.Index]
	q.Q[q.Index] = nil
	return data, nil
}

func (q *LifoMemoryQueue) Close() error {
	q.Closed = true
	return nil
}
