package queue

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewLifoDiskQueue(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	queue, err := NewLifoDiskQueue(dir)
	if err != nil {
		panic(err)
	}
	queue.Len()
	assertErr(t, queue.Push([]byte{1}), nil)
	assertErr(t, queue.Push([]byte{2}), nil)
	assertErr(t, queue.Push([]byte{3}), nil)
	assertErr(t, queue.Push([]byte{4}), nil)
	assertErr(t, queue.Push([]byte{5}), nil)
	err = queue.Close()
	if err != nil {
		panic(err)
	}
	_ = queue.Push([]byte{1})
	_, _ = queue.Pop()
	queue, err = NewLifoDiskQueue(dir)
	if err != nil {
		panic(err)
	}
	data, err := queue.Pop()
	assertErr(t, err, nil)
	assertData(t, data, []byte{5})
	data, err = queue.Pop()
	assertErr(t, err, nil)
	assertData(t, data, []byte{4})
	data, err = queue.Pop()
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
	queue.Close()
}
