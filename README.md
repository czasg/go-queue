# go-queue
[![LICENSE](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square&label=License)](https://github.com/czasg/go-queue/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/czasg/go-queue/branch/main/graph/badge.svg?token=GMXXOKC4P8)](https://codecov.io/gh/czasg/go-queue)
[![GitHub Stars](https://img.shields.io/github/stars/czasg/go-queue.svg?style=flat-square&label=Stars&logo=github)](https://github.com/czasg/go-queue/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/czasg/go-queue.svg?style=flat-square&label=Forks&logo=github)](https://github.com/czasg/go-queue/network/members)
[![GitHub Issue](https://img.shields.io/github/issues/czasg/go-queue.svg?style=flat-square&label=Issues&logo=github)](https://github.com/czasg/go-queue/issues)

## 1.背景
在 go 中，存在 `chan` 这种天然的 FIFO 队列。  
但类似 **LIFO、持久化** 等特殊结构/能力，并没有通用的标准库，更多的是借用三方软件来实现。 
 
go-queue 定义了简单的队列标准，提供了 **FIFO、LIFO、持久化** 等能力。

## 2.目标
1、内存队列、磁盘队列  
- [x] FIFO Memory Queue - 内存队列
- [x] LIFO Memory Queue - 内存队列
- [x] FIFO Disk Queue - 磁盘队列
- [x] LIFO Disk Queue - 磁盘队列

2、Get/Put 支持阻塞，磁盘队列可不支持
- [x] FIFO Block Memory Queue - 内存队列支持阻塞
- [x] LIFO Block Memory Queue - 内存队列支持阻塞

## 3.使用
1、初始化队列
```go
// 依赖
import "github.com/czasg/go-queue"

// 初始化内存队列
_ = queue.NewFifoMemoryQueue() 
_ = queue.NewLifoMemoryQueue(2048) 

// 初始化磁盘队列，需要指定目标文件
var fifofilename, lifofilename string
_, _ = queue.NewFifoDiskQueue(fifofilename)
_, _ = queue.NewLifoDiskQueue(lifofilename)
```

2、推送数据
```go
q := queue.NewFifoMemoryQueue() 
// 非阻塞
_ = q.Put(nil, []byte("data"))
// 阻塞
_ = q.Put(context.Background(), []byte("data"))
```

3、获取数据
```go
q := queue.NewFifoMemoryQueue() 
// 非阻塞
_, _ = q.Get(nil)
// 阻塞
_, _ = q.Get(context.Background())
```

4、关闭队列
```go
q := queue.NewFifoMemoryQueue() 
q.Close()
```
特别是磁盘队列，使用完后务必确保关闭。否则会出现**文件损坏**问题。

## 4.队列接口
```
type Queue interface {
    Get(ctx context.Context) ([]byte, error)
    Put(ctx context.Context, data []byte) error
    Len() int
    Close() error
}
```
其中上下文`context.Context`用于决定此次`Get / Put`是否阻塞。

## 5.Demo
### FIFO Memory Queue
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

### LIFO Memory Queue
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

### FIFO Disk Queue
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

### LIFO Disk Queue
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
