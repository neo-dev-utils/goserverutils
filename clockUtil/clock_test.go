package clockUtil

import (
	"testing"
	"time"
)

/*
*测试生产者
 */
func TestClock(t *testing.T) {
	timer := CreateClock()
	ch := make(chan bool)
	timer.AddClock(
		time.Second,
		0,
		0,
		func(inter any) {
			// log.Trace("执行重复定时器")
		},
		10,
	)
	timer.AddClock(
		time.Second*5,
		0,
		1,
		func(inter any) {
			// log.Trace("执行单次定时器1")
			timer.AddClock(
				time.Second*5,
				0,
				1,
				func(inter any) {
					// log.Trace("执行单次定时器2")
					ch <- true
				},
				10,
			)
		},
		10,
	)
	<-ch
	timer.Reset()
	timer.AddClock(
		time.Second,
		0,
		0,
		func(inter any) {
			// log.Trace("执行重复定时器2")
		},
		10,
	)
	timer.AddClock(
		time.Second*5,
		0,
		1,
		func(inter any) {
			// log.Trace("执行单次定时器3")
			ch <- true
		},
		10,
	)
	<-ch
	timer.Close()
	time.Sleep(time.Second * 5)
}
