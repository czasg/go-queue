package queue

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewFifoMemoryQueue(t *testing.T) {
	test_queue("FifoMemoryQueue", NewFifoMemoryQueue(1024), t)
}

func TestNewFifoMemoryQueueGetClose(t *testing.T) {
	name := "FifoMemoryQueue"
	queue := NewFifoMemoryQueue(1024)
	go func() {
		if data, err := queue.Get(context.Background()); data != nil || !errors.Is(err, ErrQueueClosed) {
			t.Error(name, "TestNewFifoMemoryQueuePutClose-ctx失效-阻塞队列-推送1024条数据返回cancel", err)
		}
	}()
	time.Sleep(time.Millisecond)
	queue.Close()
	time.Sleep(time.Millisecond * 2)
}

func TestNewFifoMemoryQueuePutClose(t *testing.T) {
	name := "FifoMemoryQueue"
	queue := NewFifoMemoryQueue(1024)
	for i := 0; i < 1024; i++ {
		if err := queue.Put(nil, []byte{}); err != nil {
			t.Error(name, "阻塞队列-推送1024条数据返回nil", err)
		}
	}
	go func() {
		if err := queue.Put(context.Background(), []byte{}); !errors.Is(err, ErrQueueClosed) {
			t.Error(name, "TestNewFifoMemoryQueuePutClose-ctx失效-阻塞队列-推送1024条数据返回cancel", err)
		}
	}()
	time.Sleep(time.Millisecond)
	queue.Close()
	time.Sleep(time.Millisecond * 2)
}
