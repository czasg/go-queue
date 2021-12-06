package queue

import (
	"testing"
)

func TestNewPriorityQueueFactory(t *testing.T) {
	factory := func(priority int) Queue {
		return NewFifoMemoryQueue(1024)
	}
	queue := NewPriorityQueueFactory(nil, factory)
	_ = queue.Push([]byte{11}, 1)
	_ = queue.Push([]byte{12}, 1)
	_ = queue.Push([]byte{13}, 1)
	_ = queue.Push([]byte{21}, 2)
	_ = queue.Push([]byte{22}, 2)
	_ = queue.Push([]byte{23}, 2)
	_ = queue.Push([]byte{31}, 3)
	_ = queue.Push([]byte{32}, 3)
	_ = queue.Push([]byte{33}, 3)

	data, err := queue.Pop()
	assertData(t, data, []byte{31})
	assertErr(t, err, nil)
	data, err = queue.Pop()
	assertData(t, data, []byte{32})
	assertErr(t, err, nil)
	data, err = queue.Pop()
	assertData(t, data, []byte{33})
	assertErr(t, err, nil)
	data, err = queue.Pop()
	assertData(t, data, []byte{21})
	assertErr(t, err, nil)
	data, err = queue.Pop()
	assertData(t, data, []byte{22})
	assertErr(t, err, nil)
	data, err = queue.Pop()
	assertData(t, data, []byte{23})
	assertErr(t, err, nil)
	data, err = queue.Pop()
	assertData(t, data, []byte{11})
	assertErr(t, err, nil)
	data, err = queue.Pop()
	assertData(t, data, []byte{12})
	assertErr(t, err, nil)
	data, err = queue.Pop()
	assertData(t, data, []byte{13})
	assertErr(t, err, nil)

	if queue.Len() != 0 {
		t.Error("failure")
	}

	_, err = queue.Pop()
	assertErr(t, err, ErrEmptyQueue)

	_ = queue.Close()
	_ = queue.Push([]byte{1}, 0)
	_, _ = queue.Pop()


	queue = NewPriorityQueueFactory(map[int]Queue{
		1:factory(1),
		2:factory(2),
	}, factory)
	_ = queue.Push([]byte{11}, 1)
	_ = queue.Push([]byte{12}, 1)
	_ = queue.Push([]byte{13}, 1)
	_ = queue.Push([]byte{21}, 2)
	_ = queue.Push([]byte{22}, 2)
	_ = queue.Push([]byte{23}, 2)
	_ = queue.Push([]byte{31}, 3)
	_ = queue.Push([]byte{32}, 3)
	_ = queue.Push([]byte{33}, 3)
	queue.Len()
	_ = queue.Close()
}
