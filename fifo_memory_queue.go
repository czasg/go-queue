package queue

func NewFifoMemoryQueue(size int) Queue {
	return &FifoMemoryQueue{Q: make(chan []byte, size)}
}

var _ Queue = (*FifoMemoryQueue)(nil)

type FifoMemoryQueue struct {
	Q      chan []byte
	Closed bool
}

func (q *FifoMemoryQueue) Len() int {
	return len(q.Q)
}

func (q *FifoMemoryQueue) Push(data []byte) error {
	if q.Closed {
		return ErrClosed
	}
	select {
	case q.Q <- data:
		return nil
	default:
		return ErrFullQueue
	}
}

func (q *FifoMemoryQueue) Pop() ([]byte, error) {
	if q.Closed {
		return nil, ErrClosed
	}
	select {
	case data := <-q.Q:
		return data, nil
	default:
		return nil, ErrEmptyQueue
	}
}

func (q *FifoMemoryQueue) Close() error {
	q.Closed = true
	return nil
}
