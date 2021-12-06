# go-queue
[![LICENSE](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square&label=License)](https://github.com/czasg/go-queue/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/czasg/go-queue/branch/main/graph/badge.svg?token=GMXXOKC4P8)](https://codecov.io/gh/czasg/go-queue)
[![GitHub Stars](https://img.shields.io/github/stars/czasg/go-queue.svg?style=flat-square&label=Stars&logo=github)](https://github.com/czasg/go-queue/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/czasg/go-queue.svg?style=flat-square&label=Forks&logo=github)](https://github.com/czasg/go-queue/network/members)
[![GitHub Issue](https://img.shields.io/github/issues/czasg/go-queue.svg?style=flat-square&label=Issues&logo=github)](https://github.com/czasg/go-queue/issues)


go-queue was thread-safe collections for memory/disks queues (FIFO), stacks (LIFO) and priority.


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
type Queue interface {
	Push(data []byte) error
	Pop() ([]byte, error)
	Close() error
	Len() int
}

type PriorityQueue interface {
	Push(data []byte, priority int) error
	Pop() ([]byte, error)
	Close() error
	Len() int
}
```

# Demo
### fifo memory queue
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

### fifo disk queue
disk queue will never full.
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

### lifo memory queue
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

### lifo disk queue
disk queue will never full.
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
`PriorityQueue` based on `Queue` as a queue factory.
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
