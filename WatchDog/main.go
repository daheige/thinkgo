//看门狗实现，用以监控容易失控的循环或超时
package WatchDog

import (
	"sync/atomic"
	"time"

	"thinkgo/common"
)

type WatchDog struct {
	wait  time.Duration //等待的时间
	hung  func()
	meat  int32
	count int32
}

func NewDog(duration time.Duration, meat int32, hung func()) *WatchDog {
	d := new(WatchDog)
	d.wait = duration
	d.hung = hung
	d.meat = meat

	go d.eat()
	return d
}

func (this *WatchDog) eat() {
	defer common.CheckPanic()

	for this.hung != nil && atomic.LoadInt32(&this.count) < 3 {
		time.Sleep(this.wait)

		m := atomic.LoadInt32(&this.meat)
		if m < 0 {
			return
		}

		if m == 0 {
			this.hung()
			atomic.AddInt32(&this.count, 1)
		} else {
			atomic.StoreInt32(&this.meat, m/2)
			atomic.StoreInt32(&this.count, 0)
		}
	}
}

func (this *WatchDog) Feed(meat uint16) bool {
	defer common.CheckPanic()

	// meat enought
	if atomic.LoadInt32(&this.meat) > 1024*65536 {
		return true
	}

	return atomic.AddInt32(&this.meat, int32(meat)) > 0
}

func (this *WatchDog) Kill() {
	defer common.CheckPanic()
	atomic.StoreInt32(&this.meat, -65536)
}

func (this *WatchDog) Living() bool {
	defer common.CheckPanic()
	return atomic.LoadInt32(&this.count) == 0
}
