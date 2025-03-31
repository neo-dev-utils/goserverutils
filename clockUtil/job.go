package clockUtil

import (
	"gocommutils/timeUtil"
	"sort"
	"sync"
	"time"
)

// Job 任务
type Job struct {
	clock         *Clock        // 定时器
	keyType       uint64        // 关键类型(分类)
	count         int           // 执行次数
	executionTime time.Time     // 执行的时间
	intervalTime  time.Duration // 间隔时间
	callBackFunc  func(any)     // 回调
	param         any           // 参数
}

// Cancel 取消
func (j *Job) Cancel() {
	// fmt.Println("===================删除定时器=")
	j.clock.DelJob(j)
}

// GetIntervalTime 执行时间
func (j *Job) GetIntervalTime() time.Duration {
	return j.intervalTime
}

// GetTimeDifference 获取距离执行时间相差多少
func (j *Job) GetTimeDifference() int64 {
	return (j.executionTime.UnixNano() - timeUtil.Time.NowUnixNano()) / 1e6
}

// JobSlice 游戏切片
type JobSlice []*Job

func (js JobSlice) Len() int {
	return len(js)
}

func (js JobSlice) Swap(i int, j int) {
	js[i], js[j] = js[j], js[i]
}

func (js JobSlice) Less(i int, j int) bool {
	return js[i].executionTime.UnixNano() < js[j].executionTime.UnixNano()
}

// JobQueue 任务队列
type JobQueue struct {
	jobs JobSlice
	lock sync.Mutex
}

// NewJobQueue 新建
func NewJobQueue() *JobQueue {
	return &JobQueue{
		jobs: JobSlice{},
	}
}

// Remove 删除
func (jq *JobQueue) Remove(job *Job) {
	jq.lock.Lock()
	defer jq.lock.Unlock()

	for k, v := range jq.jobs {
		if v == job {
			jq.jobs = append(jq.jobs[:k], jq.jobs[k+1:]...)
			break
		}
	}

}

// Insert 添加
func (jq *JobQueue) Insert(job *Job) {
	jq.lock.Lock()
	defer jq.lock.Unlock()

	jq.jobs = append(jq.jobs, job)
}

// Min 小
func (jq *JobQueue) Min() *Job {
	jq.lock.Lock()
	defer jq.lock.Unlock()

	if len(jq.jobs) <= 0 {
		return nil
	}
	sort.Sort(jq.jobs)
	return jq.jobs[0]
}

// Clear 清空
func (jq *JobQueue) Clear() {
	jq.lock.Lock()
	defer jq.lock.Unlock()

	jq.jobs = JobSlice{}
}

// Clear 清空
func (jq *JobQueue) ClearType(keyType uint64) {
	jq.lock.Lock()
	defer jq.lock.Unlock()

	// list := make([]*Job, 0)
	// for _, v := range jq.jobs {
	// 	if keyType == v.keyType {
	// 		list = append(list, v)
	// 	}
	// }

	// for _, v := range list {
	// 	jq.Remove(v)
	// }

	for k, v := range jq.jobs {
		if v.keyType == keyType {
			jq.jobs = append(jq.jobs[:k], jq.jobs[k+1:]...)
			break
		}
	}

}
