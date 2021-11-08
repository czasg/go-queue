package queue

import (
	"testing"
)

func TestNewFifoMemoryQueue(t *testing.T) {
	queue := NewFifoMemoryQueue(3)
	assertErr(t, queue.Push([]byte{1}), nil)
	assertErr(t, queue.Push([]byte{2}), nil)
	assertErr(t, queue.Push([]byte{3}), nil)
	assertErr(t, queue.Push([]byte{4}), ErrFullQueue)
	assertErr(t, queue.Push([]byte{5}), ErrFullQueue)
	data, err := queue.Pop()
	assertErr(t, err, nil)
	assertData(t, data, []byte{1})
	data, err = queue.Pop()
	assertErr(t, err, nil)
	assertData(t, data, []byte{2})
	data, err = queue.Pop()
	assertErr(t, err, nil)
	assertData(t, data, []byte{3})
	data, err = queue.Pop()
	assertErr(t, err, ErrEmptyQueue)
	assertData(t, data, nil)
	data, err = queue.Pop()
	assertErr(t, err, ErrEmptyQueue)
	assertData(t, data, nil)
	err = queue.Close()
	assertErr(t, err, nil)

	queue = NewFifoMemoryQueue(3)
	_ = queue.Close()
	_ = queue.Push([]byte{1})
	_, _ = queue.Pop()
}
