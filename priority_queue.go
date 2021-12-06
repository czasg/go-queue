package queue

import "sync"

type QueueFactory func(priority int) Queue

func NewPriorityQueueFactory(queues map[int]Queue, factory QueueFactory) PriorityQueue {
	cursor := 0
	if queues == nil {
		queues = map[int]Queue{}
	}
	for priority := range queues {
		if cursor < priority {
			cursor = priority
		}
	}
	return &PriorityQueueFactory{
		QueueFactory: factory,
		Q:            queues,
		Cursor:       cursor,
	}
}

var _ PriorityQueue = (*PriorityQueueFactory)(nil)

type PriorityQueueFactory struct {
	QueueFactory
	Cursor int
	Q      map[int]Queue
	Lock   sync.Mutex
	Closed bool
}

func (q *PriorityQueueFactory) Push(data []byte, priority int) error {
	q.Lock.Lock()
	defer q.Lock.Unlock()
	if q.Closed {
		return ErrClosed
	}
	queue, ok := q.Q[priority]
	if !ok {
		queue = q.QueueFactory(priority)
	}
	err := queue.Push(data)
	if err != nil {
		return err
	}
	q.Q[priority] = queue
	if q.Cursor < priority {
		q.Cursor = priority
	}
	return nil
}

func (q *PriorityQueueFactory) Pop() ([]byte, error) {
	q.Lock.Lock()
	defer q.Lock.Unlock()
	if q.Closed {
		return nil, ErrClosed
	}
	queue, ok := q.Q[q.Cursor]
	if !ok {
		return nil, ErrEmptyQueue
	}
	data, err := queue.Pop()
	if err != nil {
		return nil, err
	}
	if queue.Len() < 1 {
		_ = queue.Close()
		delete(q.Q, q.Cursor)
		q.Cursor = 0
		for priority := range q.Q {
			if q.Cursor < priority {
				q.Cursor = priority
			}
		}
	}
	return data, nil
}

func (q *PriorityQueueFactory) Close() error {
	q.Lock.Lock()
	defer q.Lock.Unlock()
	q.Closed = true
	for _, queue := range q.Q {
		err := queue.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (q *PriorityQueueFactory) Len() (count int) {
	q.Lock.Lock()
	defer q.Lock.Unlock()
	for _, queue := range q.Q {
		count += queue.Len()
	}
	return
}
