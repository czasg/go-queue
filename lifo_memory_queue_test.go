package queue

import (
	"reflect"
	"testing"
)

func assertErr(t *testing.T, err1, err2 error) {
	if err1 != err2 {
		t.Error("failure")
	}
}

func assertData(t *testing.T, data1, data2 []byte) {
	if !reflect.DeepEqual(data1, data2) {
		t.Error("failure")
	}
}

func TestNewLifoMemoryQueue(t *testing.T) {
	queue := NewLifoMemoryQueue(3)
	assertErr(t, queue.Push([]byte{1}), nil)
	assertErr(t, queue.Push([]byte{2}), nil)
	assertErr(t, queue.Push([]byte{3}), nil)
	assertErr(t, queue.Push([]byte{4}), ErrFullQueue)
	assertErr(t, queue.Push([]byte{5}), ErrFullQueue)
	data, err := queue.Pop()
	assertErr(t, err, nil)
	assertData(t, data, []byte{3})
	data, err = queue.Pop()
	assertErr(t, err, nil)
	assertData(t, data, []byte{2})
	data, err = queue.Pop()
	assertErr(t, err, nil)
	assertData(t, data, []byte{1})
	data, err = queue.Pop()
	assertErr(t, err, ErrEmptyQueue)
	assertData(t, data, nil)
	data, err = queue.Pop()
	assertErr(t, err, ErrEmptyQueue)
	assertData(t, data, nil)
	err = queue.Close()
	assertErr(t, err, nil)
}
