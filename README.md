# go-queue
[![LICENSE](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square&label=License)](https://github.com/czasg/go-queue/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/czasg/go-queue/branch/main/graph/badge.svg?token=GMXXOKC4P8)](https://codecov.io/gh/czasg/go-queue)
[![GitHub Stars](https://img.shields.io/github/stars/czasg/go-queue.svg?style=flat-square&label=Stars&logo=github)](https://github.com/czasg/go-queue/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/czasg/go-queue.svg?style=flat-square&label=Forks&logo=github)](https://github.com/czasg/go-queue/network/members)
[![GitHub Issue](https://img.shields.io/github/issues/czasg/go-queue.svg?style=flat-square&label=Issues&logo=github)](https://github.com/czasg/go-queue/issues)

## 背景
在 go 中，存在 `chan` 这种天然的 FIFO 队列。  
但类似 LIFO 队列，持久化机制等，并没有通用的标准库，更多的是借用三方软件来实现。 
 
go-queue 实现了简单的 FIFO、LIFO、持久化 等基础能力。

## 目标
1、内存队列和磁盘队列  
- [x] FIFO Memory Queue - 内存队列
- [x] LIFO Memory Queue - 内存队列
- [x] FIFO Disk Queue - 磁盘队列
- [x] LIFO Disk Queue - 磁盘队列

2、内存队列支持阻塞 Get/Put 数据，磁盘队列不支持
- [x] FIFO Block Memory Queue - 内存队列支持阻塞
- [x] LIFO Block Memory Queue - 内存队列支持阻塞

## 模块
设计通用队列接口 `Queue`    
```
type Queue interface {
    Get(ctx context.Context) ([]byte, error)
    Put(ctx context.Context, data []byte) error
    Len() int
    Close() error
}
```
其中上下文`context.Context`用于决定此次`Get / Put`是否阻塞。

## FIFO Memory Queue
```go
package main

import (
    "context"
    "fmt"
    "github.com/czasg/go-queue"
    "time"
)

func main() {
    q := queue.NewFifoMemoryQueue()
    _ = q.Put(nil, []byte("go-queue"))
    data, _ := q.Get(nil)
    fmt.Println("非阻塞获取数据", string(data))

    go func() {
        data, _ := q.Get(context.Background())
        fmt.Println("阻塞获取数据", string(data))
    }()
    time.Sleep(time.Second * 2)
    _ = q.Put(nil, []byte("go-queue"))
    time.Sleep(time.Millisecond)
    q.Close()
}
```


## LIFO Memory Queue
```go
package main

import (
    "context"
    "fmt"
    "github.com/czasg/go-queue"
    "time"
)

func main() {
    q := queue.NewLifoMemoryQueue()
    _ = q.Put(nil, []byte("go-queue"))
    data, _ := q.Get(nil)
    fmt.Println("非阻塞获取数据", string(data))

    go func() {
        data, _ := q.Get(context.Background())
        fmt.Println("阻塞获取数据", string(data))
    }()
    time.Sleep(time.Second * 2)
    _ = q.Put(nil, []byte("go-queue"))
    time.Sleep(time.Millisecond)
    q.Close()
}
```

## FIFO Disk Queue
```go
package main

import (
    "fmt"
    "github.com/czasg/go-queue"
    "io/ioutil"
    "os"
)

func main() {
    file, err := ioutil.TempFile("", "")
    if err != nil {
        panic(err)
    }
    defer os.RemoveAll(file.Name())
    file.Close()

    q, _ := queue.NewFifoDiskQueue(file.Name())
    _ = q.Put(nil, []byte("go-queue"))
    q.Close()

    q, _ = queue.NewFifoDiskQueue(file.Name())
    data, _ := q.Get(nil)
    fmt.Println("获取数据", string(data))
    q.Close()
}
```

## LIFO Disk Queue
```go
package main

import (
    "fmt"
    "github.com/czasg/go-queue"
    "io/ioutil"
    "os"
)

func main() {
    file, err := ioutil.TempFile("", "")
    if err != nil {
        panic(err)
    }
    defer os.RemoveAll(file.Name())
    file.Close()

    q, _ := queue.NewLifoDiskQueue(file.Name())
    _ = q.Put(nil, []byte("go-queue"))
    q.Close()

    q, _ = queue.NewLifoDiskQueue(file.Name())
    data, _ := q.Get(nil)
    fmt.Println("获取数据", string(data))
    q.Close()
}
```