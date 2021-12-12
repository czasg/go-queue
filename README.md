# go-queue
[![LICENSE](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square&label=License)](https://github.com/czasg/go-queue/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/czasg/go-queue/branch/main/graph/badge.svg?token=GMXXOKC4P8)](https://codecov.io/gh/czasg/go-queue)
[![GitHub Stars](https://img.shields.io/github/stars/czasg/go-queue.svg?style=flat-square&label=Stars&logo=github)](https://github.com/czasg/go-queue/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/czasg/go-queue.svg?style=flat-square&label=Forks&logo=github)](https://github.com/czasg/go-queue/network/members)
[![GitHub Issue](https://img.shields.io/github/issues/czasg/go-queue.svg?style=flat-square&label=Issues&logo=github)](https://github.com/czasg/go-queue/issues)

**go-queue** 实现了常用的队列结构，包括 **FIFO**、**LIFO**、**PriorityQueue** 等，并且支持将数据**持久化**到磁盘。


```text
                  |—————————|                   |——————————————————|               
                  |  queue  | -----factory----- |  priority queue  |
                  |—————————|                   |——————————————————|
           ____________|___________
          |                        |      
      |————————|              |————————|
      |  fifo  |              |  lifo  |
      |————————|              |————————|
     _____|_____              _____|_____
    |           |            |           |
|————————||——————————|   |————————||——————————|
|  disk  ||  memory  |   |  disk  ||  memory  |
|————————||——————————|   |————————||——————————|
```

### plan
- [x] fifo memory queue
- [x] fifo disk queue
- [x] lifo memory queue
- [x] lifo disk queue
- [x] priority queue

### interface
```golang
// 普通队列
type Queue interface {
	Push(data []byte) error
	Pop() ([]byte, error)
	Close() error
	Len() int
}
// 优先级队列
type PriorityQueue interface {
	Push(data []byte, priority int) error
	Pop() ([]byte, error)
	Close() error
	Len() int
}
```

# Queue
## fifo memory queue
基于内存的 `FIFO` 队列，初始化时需要指定队列大小，程序退出后数据丢失。
```golang
package main

import (
	"github.com/czasg/go-queue"
)

func main() {
	maxQueueSize := 2 // out of max size, push data will return (nil, ErrFullQueue)
	q := queue.NewFifoMemoryQueue(maxQueueSize)
	defer q.Close()

	q.Push([]byte("test1")) // nil
	q.Pop()                 // test, nil

	q.Push([]byte("test1")) // nil
   	q.Push([]byte("test1")) // nil
   	q.Push([]byte("test1")) // ErrFullQueue

   	q.Pop() // test, nil
    	q.Pop() // test, nil
    	q.Pop() // ErrEmptyQueue
}
```

## fifo disk queue
基于磁盘文件的 `FIFO` 队列，初始化时需要指定存储目录，程序退出后需要执行 `Close` 方法，以确保数据不会丢失。   
设计上 `fifo disk queue` 没有存储上限，故不需要显示指定队列大小。
```golang
package main

import (
	"github.com/czasg/go-queue"
	"io/ioutil"
	"os"
)

func main() {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	
	q, _ := queue.NewFifoDiskQueue(dir)
	_ = q.Push([]byte("test1"))
	_ = q.Push([]byte("test2"))
	_ = q.Push([]byte("test3"))
	_ = q.Close()

	q, _ = queue.NewFifoDiskQueue(dir) // open again
	_, _ = q.Pop() // test1
	_, _ = q.Pop() // test2
	_, _ = q.Pop() // test3
	_ = q.Close()
}
```

## lifo memory queue
基于内存的 `LIFO` 队列，初始化时需要指定队列大小，程序退出后数据丢失。
```golang
package main

import (
	"github.com/czasg/go-queue"
)

func main() {
	maxQueueSize := 2 // out of max size, push data will return (nil, ErrFullQueue)
	q := queue.NewLifoMemoryQueue(maxQueueSize)
	defer q.Close()

	q.Push([]byte("test1")) // nil
	q.Pop()                 // test, nil

	q.Push([]byte("test1")) // nil
	q.Push([]byte("test1")) // nil
	q.Push([]byte("test1")) // ErrFullQueue

	q.Pop() // test, nil
	q.Pop() // test, nil
	q.Pop() // ErrEmptyQueue
}
```

## lifo disk queue
基于磁盘文件的 `LIFO` 队列，初始化时需要指定存储目录，程序退出后需要执行 `Close` 方法，以确保数据不会丢失。   
设计上 `lifo disk queue` 没有存储上限，故不需要显示指定队列大小。
```golang
package main

import (
	"github.com/czasg/go-queue"
	"io/ioutil"
	"os"
)

func main() {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	q, _ := queue.NewLifoDiskQueue(dir)
	_ = q.Push([]byte("test1"))
	_ = q.Push([]byte("test2"))
	_ = q.Push([]byte("test3"))
	_ = q.Close()

	q, _ = queue.NewLifoDiskQueue(dir) // open again
	_, _ = q.Pop() // test3
	_, _ = q.Pop() // test2
	_, _ = q.Pop() // test1
	_ = q.Close()
}
```

### priority queue
优先级队列，是基于常规队列实现的，故初始化时，需要指定**工厂函数**，用于内部创建新的优先级队列。  
工厂函数返回一个新的 `Queue` 类型，可以指定返回内存队列，也可以返回磁盘队列。
```golang
package main

import (
	"fmt"
	"github.com/czasg/go-queue"
	"reflect"
)

func main() {
	factory := func(priority int) queue.Queue {
		return queue.NewFifoMemoryQueue(1024)
	}
	q := queue.NewPriorityQueueFactory(nil, factory)
	_ = q.Push([]byte("v1"), 1)
	_ = q.Push([]byte("v2"), 10)
	_ = q.Push([]byte("v3"), 5)

	data, _ := q.Pop()
	fmt.Println(reflect.DeepEqual(data, []byte("v2")))
	data, _ = q.Pop()
	fmt.Println(reflect.DeepEqual(data, []byte("v3")))
	data, _ = q.Pop()
	fmt.Println(reflect.DeepEqual(data, []byte("v1")))
}
```
