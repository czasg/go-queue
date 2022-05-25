package queue

import (
    "context"
    "errors"
    "reflect"
    "testing"
    "time"
)

func test_queue(name string, queue Queue, t *testing.T) {
    {
        if length := queue.Len(); length != 0 {
            t.Error(name, "空队列-获取长度为0", length)
        }
        if data, err := queue.Get(nil); data != nil || !errors.Is(err, ErrQueueEmpty) {
            t.Error(name, "空队列-Get数据返回ErrQueueEmpty", data, err)
        }
    }
    {
        if err := queue.Put(nil, []byte{}); err != nil {
            t.Error(name, "空队列-Put数据返回nil", err)
        }
        if data, err := queue.Get(nil); !reflect.DeepEqual(data, []byte{}) || err != nil {
            t.Error(name, "非空队列-Get数据返回nil", data, err)
        }
    }
    {
        for i := 0; i < 1024; i++ {
            if err := queue.Put(nil, []byte{}); err != nil {
                t.Error(name, "空队列-推送1024条数据返回nil", err)
            }
        }
        if length := queue.Len(); length != 1024 {
            t.Error(name, "满队列获取长度为1024", length)
        }
        if err := queue.Put(nil, []byte{}); !errors.Is(err, ErrQueueFull) {
            t.Error(name, "满队列推送数据返回ErrQueueFull", err)
        }
        for i := 0; i < 1024; i++ {
            if data, err := queue.Get(nil); !reflect.DeepEqual(data, []byte{}) || err != nil {
                t.Error(name, "满队列Get数据返回nil", data, err)
            }
        }
        if data, err := queue.Get(nil); data != nil || !errors.Is(err, ErrQueueEmpty) {
            t.Error(name, "1024空队列-Get数据返回ErrQueueEmpty", data, err)
        }
    }
    {
        go func() {
            time.Sleep(time.Millisecond)
            if err := queue.Put(context.Background(), []byte{}); err != nil {
                t.Error(name, "阻塞队列Put数据返回nil", err)
            }
        }()
        if data, err := queue.Get(context.Background()); !reflect.DeepEqual(data, []byte{}) || err != nil {
            t.Error(name, "阻塞队列Get数据返回nil", data, err)
        }
        if data, err := queue.Get(nil); data != nil || !errors.Is(err, ErrQueueEmpty) {
            t.Error(name, "阻塞队列-Get数据返回ErrQueueEmpty", data, err)
        }
    }
    {
        ctx, cancel := context.WithCancel(context.Background())
        go func() {
            if data, err := queue.Get(ctx); data != nil || !errors.Is(err, ctx.Err()) {
                t.Error(name, "ctx失效-阻塞队列Get数据返回cancel", data, err)
            }
        }()
        time.Sleep(time.Millisecond)
        cancel()
    }
    {
        for i := 0; i < 1024; i++ {
            if err := queue.Put(nil, []byte{}); err != nil {
                t.Error(name, "阻塞队列-推送1024条数据返回nil", err)
            }
        }
        ctx, cancel := context.WithCancel(context.Background())
        go func() {
            if err := queue.Put(ctx, []byte{}); !errors.Is(err, ctx.Err()) {
                t.Error(name, "ctx失效-阻塞队列-推送1024条数据返回cancel", err)
            }
        }()
        time.Sleep(time.Millisecond)
        cancel()
        time.Sleep(time.Millisecond * 2)
        for i := 0; i < 1024; i++ {
            if data, err := queue.Get(nil); !reflect.DeepEqual(data, []byte{}) || err != nil {
                t.Error(name, "满队列Get数据返回nil", data, err)
            }
        }
    }
    // 队列关闭测试
    {
        if err := queue.Close(); err != nil {
            t.Error(name, "队列关闭返回nil", err)
        }
    }
    {
        if data, err := queue.Get(nil); data != nil || !errors.Is(err, ErrQueueClosed) {
            t.Error(name, "关闭队列-Get数据返回ErrQueueClosed", data, err)
        }
        if data, err := queue.Get(context.Background()); data != nil || !errors.Is(err, ErrQueueClosed) {
            t.Error(name, "关闭队列-Get数据返回ErrQueueClosed", data, err)
        }
        if err := queue.Put(nil, []byte{}); !errors.Is(err, ErrQueueClosed) {
            t.Error(name, "关闭队列-推送数据返回ErrQueueClosed", err)
        }
        if err := queue.Put(context.Background(), []byte{}); !errors.Is(err, ErrQueueClosed) {
            t.Error(name, "关闭队列-推送数据返回ErrQueueClosed", err)
        }
    }
}
