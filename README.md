# go-queue
[![LICENSE](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square&label=License)](https://github.com/czasg/go-queue/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/czasg/go-queue/branch/main/graph/badge.svg?token=GMXXOKC4P8)](https://codecov.io/gh/czasg/go-queue)
[![GitHub Stars](https://img.shields.io/github/stars/czasg/go-queue.svg?style=flat-square&label=Stars&logo=github)](https://github.com/czasg/go-queue/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/czasg/go-queue.svg?style=flat-square&label=Forks&logo=github)](https://github.com/czasg/go-queue/network/members)
[![GitHub Issue](https://img.shields.io/github/issues/czasg/go-queue.svg?style=flat-square&label=Issues&logo=github)](https://github.com/czasg/go-queue/issues)


go-queue was collections for queues (FIFO), stacks (LIFO) and priority.

```text
      |—————————|          |——————————————————|               
      |  queue  | -------- |  priority queue  |
      |—————————|          |——————————————————|
     ______|______
    |             |
|————————|   |————————|
|  fifo  |   |  lifo  |
|————————|   |————————|
```

### plan
- [x] fifo memory queue
- [x] fifo disk queue
- [x] lifo memory queue
- [ ] lifo disk queue
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
	"fmt"
	"github.com/czasg/go-queue"
	"reflect"
)

func assert(v1, v2 interface{}) {
	if !reflect.DeepEqual(v1, v2) {
		panic(fmt.Sprintf("%v != %v", v1, v2))
	}
}

func main() {
	maxQueueSize := 2 // out of max size, push data will return (nil, ErrFullQueue)
	q := queue.NewFifoMemoryQueue(maxQueueSize)
	_ = q.Push([]byte("test1"))
	_, _ = q.Pop()

	assert(q.Push([]byte("test1")), nil)
	assert(q.Push([]byte("test2")), nil)
	assert(q.Push([]byte("test3")), queue.ErrFullQueue)
	data, err := q.Pop()
	assert(data, []byte("test1"))
	assert(err, nil)
	data, err = q.Pop()
	assert(data, []byte("test2"))
	assert(err, nil)
	_, err = q.Pop()
	assert(err, queue.ErrEmptyQueue)
	_ = q.Close()
}
```

### fifo disk queue
disk queue will never full.
```golang
package main

import (
	"fmt"
	"github.com/czasg/go-queue"
	"io/ioutil"
	"os"
	"reflect"
)

func assert(v1, v2 interface{}) {
	if !reflect.DeepEqual(v1, v2) {
		panic(fmt.Sprintf("%v != %v", v1, v2))
	}
}

func main() {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	q, _ := queue.NewFifoDiskQueue(dir)
	_ = q.Push([]byte("test1"))
	_, _ = q.Pop()

	assert(q.Push([]byte("test1")), nil)
	assert(q.Push([]byte("test2")), nil)
	data, err := q.Pop()
	assert(data, []byte("test1"))
	assert(err, nil)
	data, err = q.Pop()
	assert(data, []byte("test2"))
	assert(err, nil)
	_ = q.Close()
}
```

### lifo memory queue
```golang
package main

import (
	"fmt"
	"github.com/czasg/go-queue"
	"reflect"
)

func assert(v1, v2 interface{}) {
	if !reflect.DeepEqual(v1, v2) {
		panic(fmt.Sprintf("%v != %v", v1, v2))
	}
}

func main() {
	maxQueueSize := 2 // out of max size, push data will return (nil, ErrFullQueue)
	q := queue.NewLifoMemoryQueue(maxQueueSize)
	_ = q.Push([]byte("test1"))
	_, _ = q.Pop()

	assert(q.Push([]byte("test1")), nil)
	assert(q.Push([]byte("test2")), nil)
	assert(q.Push([]byte("test3")), queue.ErrFullQueue)
	data, err := q.Pop()
	assert(data, []byte("test2"))
	assert(err, nil)
	data, err = q.Pop()
	assert(data, []byte("test1"))
	assert(err, nil)
	_, err = q.Pop()
	assert(err, queue.ErrEmptyQueue)
	_ = q.Close()
}
```

### priority queue
```golang
package main

import (
	"fmt"
	"github.com/czasg/go-queue"
	"reflect"
)

func main() {
	factory := func() queue.Queue {
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
