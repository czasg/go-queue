package queue

import (
    "context"
    "errors"
    "io/ioutil"
    "os"
    "reflect"
    "testing"
)

func TestNewFifoDiskQueue1(t *testing.T) {
    file, err := ioutil.TempFile("", "")
    if err != nil {
        panic(err)
    }
    defer os.RemoveAll(file.Name())
    file.Close()
    queue, err := NewFifoDiskQueue(file.Name())
    if err != nil {
        panic(err)
    }
    name := "TestNewFifoDiskQueue"
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
    {
        if err := queue.Close(); err != nil {
            t.Error(name, "队列关闭返回nil", err)
        }
    }
}

func TestNewFifoDiskQueue2(t *testing.T) {
    file, err := ioutil.TempFile("", "")
    if err != nil {
        panic(err)
    }
    defer os.RemoveAll(file.Name())
    file.Close()
    queue, err := NewFifoDiskQueue(file.Name())
    if err != nil {
        panic(err)
    }
    name := "TestNewFifoDiskQueue2"
    if err := queue.Put(nil, []byte{}); err != nil {
        t.Error(name, "空队列-Put数据返回nil", err)
    }
    if err := queue.Close(); err != nil {
        t.Error(name, "队列关闭返回nil", err)
    }
    queue, err = NewFifoDiskQueue(file.Name())
    if err != nil {
        panic(err)
    }
    if data, err := queue.Get(nil); !reflect.DeepEqual(data, []byte{}) || err != nil {
        t.Error(name, "非空队列-Get数据返回nil", data, err)
    }
}
