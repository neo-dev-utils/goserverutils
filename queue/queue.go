package queue

import (
	"time"
)

type Queue struct {
	data chan interface{}
	len  int
	exit bool
}

// 初始化队列的长度
func NewQueue(max_queue_len int) (dc *Queue) {
	dc = &Queue{exit: false}
	dc.len = max_queue_len
	dc.data = make(chan interface{}, max_queue_len)
	return dc
}

func (q *Queue) Push(data interface{}, waittime time.Duration) bool {
	if q.exit {
		return false
	}
	click := time.After(waittime)
	select {
	case q.data <- data:
		return true
	case <-click:
		return false
	}
}

func (q *Queue) Pop(waittime time.Duration) (data interface{}) {
	if q.exit {
		return nil
	}
	click := time.After(waittime)
	select {
	case data = <-q.data:
		return data
	case <-click:
		return nil
	}
}

func (q *Queue) Close() {
	q.exit = true
	close(q.data)
}
