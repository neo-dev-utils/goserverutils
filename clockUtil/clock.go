package clockUtil

import (
	"gocommutils/timeUtil"
	"sync"
	"time"
)

const (
	exit = 0
	open = 1
)

// Clock 定时器
type Clock struct {
	update chan bool      // 任务
	reset  chan chan bool // 重置
	exit   chan chan bool // 退出
	isExit uint32         // 是否退出
	jobs   *JobQueue      // 任务
	// rbTree *rbtree.Rbtree // 红黑树
	goWait sync.WaitGroup // 等待退出
	lock   sync.Mutex     // 锁
}

// CreateClock 创建一个定时器
func CreateClock() *Clock {
	clock := &Clock{
		update: make(chan bool),
		reset:  make(chan chan bool),
		exit:   make(chan chan bool),
		isExit: open,
		jobs:   NewJobQueue(),
		// rbTree: rbtree.New(),
	}

	go clock.Start()
	return clock
}

// Start 开启一个定时器
func (c *Clock) Start() {
	for {
		// 获取新的时间
		job := c.jobs.Min()
		if job == nil {
			job := &Job{
				count:        0,
				intervalTime: 9999 * time.Second,
			}
			job.executionTime = timeUtil.Time.NowTime().Add(job.intervalTime)

			c.jobs.Insert(job)
			continue
		}
		// 添加定时
		executionTime := job.executionTime.Sub(timeUtil.Time.NowTime())
		select {
		case <-time.After(executionTime):
			if job == nil {
				break
			}
			if job.callBackFunc != nil {
				c.goWait.Add(1)
				go func() {
					if job.keyType != 0 {
						// log.Debug("-->>>>> Clock-Start 处理")
					}
					defer c.goWait.Done()
					job.callBackFunc(job.param)
				}()
			}
			if job.count == 0 {
				job.executionTime = timeUtil.Time.NowTime().Add(job.intervalTime)
				break
			}
			if job.count == 1 {
				c.jobs.Remove(job)
				break
			}
			job.count = job.count - 1
			job.executionTime = timeUtil.Time.NowTime().Add(job.intervalTime)
			break
		case <-c.update:
			break
		case ch := <-c.reset:
			c.jobs.Clear()
			ch <- true
			break
		case ch := <-c.exit:
			c.jobs.Clear()
			c.goWait.Wait()
			ch <- true
			return
		}
	}
}

// AddClock 添加定时任务 纳秒级别
func (c *Clock) AddClock(intervalTime time.Duration, keyType uint64, count int, function func(any), param any) (*Job, bool) {
	// log.Debug("-->>>>> Clock-AddClock 开始")
	c.lock.Lock()
	// log.Debug("-->>>>> Clock-AddClock 开始  111")
	if c.isExit > exit {

		job := &Job{
			clock:        c,
			keyType:      keyType,
			count:        count,
			intervalTime: intervalTime,
			callBackFunc: function,
			param:        param,
		}

		job.executionTime = timeUtil.Time.NowTime().Add(job.intervalTime)
		// log.Debug("-->>>>> Clock-AddClock 开始  222")
		c.jobs.Insert(job)
		c.update <- true
		c.lock.Unlock()

		// log.Debug("-->>>>> Clock-AddClock 开始  3333")
		return job, true
	}
	// log.Debug("-->>>>> Clock-AddClock 开始  4444")
	c.lock.Unlock()
	// log.Debug("-->>>>> Clock-AddClock 开始  5555")
	return nil, false
}

// DelJob 删除定时器
func (c *Clock) DelJob(job *Job) {
	c.jobs.Remove(job)
	c.update <- true
	// for range c.update {
	// 	// fmt.Println(err)
	// 	c.update <- true
	// }
}

func (c *Clock) ClearType(keyType uint64) {
	c.jobs.ClearType(keyType)
}

// Close 关闭定时器
func (c *Clock) Close() {
	c.lock.Lock()
	c.isExit = exit
	ch := make(chan bool)
	c.exit <- ch
	<-ch
	close(c.update)
	close(c.reset)
	close(c.exit)
	c.lock.Unlock()
}

// Reset 清空定时任务
func (c *Clock) Reset() {
	c.lock.Lock()
	if c.isExit == exit {
		c.lock.Unlock()
		return
	}
	c.isExit = exit
	ch := make(chan bool)
	c.reset <- ch
	<-ch
	c.isExit = open
	c.lock.Unlock()
}
