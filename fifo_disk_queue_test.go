package queue

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewFifoDiskQueue(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	queue, err := NewFifoDiskQueueWithChunk(dir, 3)
	if err != nil {
		panic(err)
	}
	assertErr(t, queue.Push([]byte{1}), nil)
	assertErr(t, queue.Push([]byte{2}), nil)
	assertErr(t, queue.Push([]byte{3}), nil)
	assertErr(t, queue.Push([]byte{4}), nil)
	assertErr(t, queue.Push([]byte{5}), nil)
	err = queue.Close()
	if err != nil {
		panic(err)
	}
	_, err = NewFifoDiskQueueWithChunk(dir, 4)
	assertErr(t, err, ErrChunkSizeInconsistency)
	queue, err = NewFifoDiskQueueWithChunk(dir, 3)
	if err != nil {
		panic(err)
	}
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
	assertErr(t, err, nil)
	assertData(t, data, []byte{4})
	data, err = queue.Pop()
	assertErr(t, err, nil)
	assertData(t, data, []byte{5})
	data, err = queue.Pop()
	assertErr(t, err, ErrEmptyQueue)
	assertData(t, data, nil)
	err = queue.Close()
	assertErr(t, err, nil)
	err = queue.Close()
	assertErr(t, err, nil)
	err = queue.Push([]byte{1})
	assertErr(t, err, ErrClosed)
	_, err = queue.Pop()
	assertErr(t, err, ErrClosed)
}

func TestFifoDiskQueue_openData(t *testing.T) {
	q := FifoDiskQueue{Dir: ">0<"}
	_, err := q.openData(0, 0)
	if err == nil {
		t.Error("failure")
	}
}

func TestFifoDiskQueue_saveStat(t *testing.T) {
	q := FifoDiskQueue{Dir: ">0<"}
	err := q.saveStat()
	if err == nil {
		t.Error("failure")
	}

	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	_, _ = NewFifoDiskQueueWithChunk(dir, 10)

	_, _ = NewFifoDiskQueue(">0<")
}

func TestFifoDiskQueue_openStat(t *testing.T) {
	q := FifoDiskQueue{Dir: ">0<"}
	_, err := q.openStat()
	if err == nil {
		t.Error("failure")
	}
}
