# go-queue
[![LICENSE](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square&label=License)](https://github.com/czasg/go-queue/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/czasg/go-queue/branch/main/graph/badge.svg?token=GMXXOKC4P8)](https://codecov.io/gh/czasg/go-queue)
[![GitHub Stars](https://img.shields.io/github/stars/czasg/go-queue.svg?style=flat-square&label=Stars&logo=github)](https://github.com/czasg/go-queue/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/czasg/go-queue.svg?style=flat-square&label=Forks&logo=github)](https://github.com/czasg/go-queue/network/members)
[![GitHub Issue](https://img.shields.io/github/issues/czasg/go-queue.svg?style=flat-square&label=Issues&logo=github)](https://github.com/czasg/go-queue/issues)

## 背景
在 go 中，存在 `chan` 这种天然的 FIFO 队列。  
但类似**LIFO队列，持久化机制**等，并没有通用的标准库，更多的是借用三方软件来实现。 
 
go-queue 实现了简单的 **FIFO、LIFO、持久化** 等基础能力。

## 目标
1、内存队列和磁盘队列  
- [x] FIFO Memory Queue - 内存队列
- [x] LIFO Memory Queue - 内存队列
- [x] FIFO Disk Queue - 磁盘队列
- [x] LIFO Disk Queue - 磁盘队列

2、内存队列支持阻塞 Get/Put 数据，磁盘队列不支持
- [x] FIFO Block Memory Queue - 内存队列支持阻塞
- [x] LIFO Block Memory Queue - 内存队列支持阻塞

## 使用
1、初始化队列
```go title="初始化队列"
// 依赖
import "github.com/czasg/go-queue"
// FIFO内存队列
_ = queue.NewFifoMemoryQueue() 
// LIFO内存队列
_ = queue.NewLifoMemoryQueue() 

var fifofilename, lifofilename string
_, _ = queue.NewFifoDiskQueue(fifofilename) // FIFO磁盘队列
_, _ = queue.NewLifoDiskQueue(lifofilename) // LIFO磁盘队列
```

2、推送数据
```go
q := queue.NewFifoMemoryQueue() 
// 非阻塞
_ = q.Put(nil, []byte("data"))
// 阻塞
_ = q.Put(context.Background(), []byte("data"))
```
disk queue 不支持阻塞方式

3、获取数据
```go
q := queue.NewFifoMemoryQueue() 
// 非阻塞
_, _ = q.Get(nil)
// 阻塞
_, _ = q.Get(context.Background())
```
disk queue 不支持阻塞方式

## 队列接口
```
type Queue interface {
    Get(ctx context.Context) ([]byte, error)
    Put(ctx context.Context, data []byte) error
    Len() int
    Close() error
}
```
其中上下文`context.Context`用于决定此次`Get / Put`是否阻塞。
