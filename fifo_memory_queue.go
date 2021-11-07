package queue

func NewFifoMemoryQueue(size int) Queue {
	return &FifoMemoryQueue{Q: make(chan []byte, size)}
}

var _ Queue = (*FifoMemoryQueue)(nil)

type FifoMemoryQueue struct {
	Q chan []byte
}

func (q *FifoMemoryQueue) Push(data []byte) error {
	select {
	case q.Q <- data:
		return nil
	default:
		return ErrFullQueue
	}
}

func (q *FifoMemoryQueue) Pop() ([]byte, error) {
	select {
	case data := <-q.Q:
		return data, nil
	default:
		return nil, ErrEmptyQueue
	}
}

func (q *FifoMemoryQueue) Close() error {
	return nil
}
