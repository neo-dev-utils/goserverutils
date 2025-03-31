package lockUtil

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var now time.Time

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

type TRWlocker struct {
	w, r        time.Time
	wpath       string
	rpath       string
	rlpath      []string
	rulpath     []string
	l           *sync.RWMutex
	dev         bool
	readerCount int32 // number of pending readers
	dur         time.Duration
	cb          func(str string)
}

var TimeNull = time.Time{}

//
//func NewTRWLocker(dev bool, dur time.Duration, cb func(str string)) *TRWlocker {
//	trw := new(TRWlocker)
//	trw.l = new(sync.RWMutex)
//	trw.r = TimeNull
//	trw.dur = dur
//	trw.cb = cb
//	trw.readerCount = 0
//	done := make(chan int, 0)
//	go trw.timetick(done)
//	<-done
//	fmt.Println("done")
//	return trw
//}
//
//func (tl *TRWlocker) timetick(done chan int) {
//	now = timeUtil.Time.NowTime()
//	done<-1
//	for true {
//		now = timeUtil.Time.NowTime()
//		time.Sleep(2*time.Millisecond)
//	}
//}

//func (tl *TRWlocker) Lock() {
//	_, file, line, _ := runtime.Caller(1)
//	tl.wpath = fmt.Sprintf("PID[%d]:Lock Called by %s:%d", getGID(), file, line)
//	tl.w = now
//	tl.l.Lock()
//}
//func (tl *TRWlocker) RCount() int32 {
//	return 	atomic.LoadInt32(&tl.readerCount)
//
//}
//func (tl *TRWlocker) Unlock() {
//	tl.l.Unlock()
//	_, file, line, _ := runtime.Caller(1)
//	//if tl.w.Add(tl.dur).Before(now) {
//		if tl.cb != nil {
//			tl.cb(fmt.Sprintf("time cost %d from %s to %s ", (now.UnixNano() - tl.w.UnixNano())/1000, tl.wpath, fmt.Sprintf("Lock Called by %s:%d", file, line)))
//		}
//	//}
//}
//
//func (tl *TRWlocker) RLock() {
//	//_, file, line, _ := runtime.Caller(1)
//	runtime.Caller(1)
//	if tl.r == TimeNull {
//		tl.r = now
//	}
//	//tl.rlpath = append(tl.rlpath, fmt.Sprintf("RLock Called by %s:%d at %s", file, line, now.Format(time.StampNano)))
//	atomic.AddInt32(&tl.readerCount, 1)
//	tl.l.RLock()
//}
//
//func (tl *TRWlocker) RUnlock() {
//	tl.l.RUnlock()
//	//_, file, line, _ := runtime.Caller(1)
//	runtime.Caller(1)
//
//	//tl.rulpath = append(tl.rulpath, fmt.Sprintf("RULock Called by %s:%d at %s", file, line, now.Format(time.StampNano)))
//	atomic.AddInt32(&tl.readerCount, -1)
//	if tl.readerCount == 0 {
//		//all reader lock release
//		if tl.r != TimeNull && tl.r.Add(tl.dur).Before(now) {
//			if tl.cb != nil {
//				str := "\n"
//				for k := 0; k < len(tl.rlpath);k++ {
//					str +=tl.rlpath[k]+ "\n"
//				}
//				for k := 0; k < len(tl.rulpath);k++ {
//					str += tl.rulpath[k] + "\n"
//				}
//				if tl.cb != nil {
//					tl.cb(str)
//				}
//			}
//		}
//		tl.rulpath = []string{}
//		tl.rlpath = []string{}
//		tl.r = TimeNull
//	}
//}

func NewTRWLocker(dev bool, dur time.Duration, cb func(str string)) *TRWlocker {
	trw := new(TRWlocker)
	trw.l = new(sync.RWMutex)
	trw.r = TimeNull
	trw.dur = dur
	trw.cb = cb
	trw.readerCount = 0
	return trw
}

func (tl *TRWlocker) Lock() {
	tl.l.Lock()
}

func (tl *TRWlocker) Unlock() {
	tl.l.Unlock()

}

func (tl *TRWlocker) RLock() {
	tl.l.RLock()
}

func (tl *TRWlocker) RUnlock() {
	tl.l.RUnlock()

}
