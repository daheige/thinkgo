/**
采用cas实现乐观锁及TryLock
Compare And Swap 简称CAS，在sync/atomic包种
这类原子操作由名称以‘CompareAndSwap’为前缀的若干个函数代表。
声明如下
    func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
调用函数后，会先判断参数addr指向的被操作值与参数old的值是否相等
仅当此判断得到肯定的结果之后，才会用参数new代表的新值替换掉原先的旧值，否则操作就会被忽略。
so, 需要用for循环不断进行尝试,直到成功为止

使用锁的做法趋于悲观
    我们总假设会有并发的操作要修改被操作的值，并使用锁将相关操作放入临界区中加以保护
使用CAS操作的做法趋于乐观
    总是假设被操作值未曾被改变（即与旧值相等），并一旦确认这个假设的真实性就立即进行值替换。
*/
package common

import (
	"runtime"
	"sync/atomic"
	"time"
)

const (
	unlock = 0
	locked = 1
)

type Lock struct {
	flag       int32
	unlockChan chan bool
}

func NewLock() *Lock {
	l := new(Lock)
	l.unlockChan = make(chan bool, 3)
	return l
}

func (this *Lock) Lock() {
	for {
		if atomic.CompareAndSwapInt32(&this.flag, unlock, locked) {
			return
		}

		<-this.unlockChan
	}
}

// 不断地尝试原子地更新flag的值,直到操作成功为止
func (this *Lock) SpinLock(count int) bool {
	for i := 0; i < count; i++ {
		if atomic.CompareAndSwapInt32(&this.flag, unlock, locked) {
			return true
		}
	}

	return false
}

//尝试枷锁
func (this *Lock) TryLock() bool {
	return this.SpinLock(1)
}

func (this *Lock) TimeoutSpinLock(t time.Duration) bool {
	tk := time.NewTicker(t)
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			return this.SpinLock(1)

		default:
			if this.SpinLock(1024) {
				return true
			}
		}
	}
}

func (this *Lock) SchedLock() {
	for {
		if this.SpinLock(1024) {
			return
		}
		runtime.Gosched()
	}
}

func (this *Lock) TimeoutSchedLock(t time.Duration) bool {
	tk := time.NewTicker(t)
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			return this.SpinLock(1)

		default:
			if this.SpinLock(1024) {
				return true
			}
		}
		runtime.Gosched()
	}
}

//释放锁
func (this *Lock) Unlock() {
	atomic.SwapInt32(&this.flag, unlock)
	select {
	case this.unlockChan <- true:
	default:
	}
}

//分配安全的bool chan
type Semaphore struct {
	flag chan bool
}

func NewSemaphore(capa int) *Semaphore {
	return &Semaphore{make(chan bool, capa)}
}

func (this *Semaphore) Alloc() {
	this.flag <- true
}

func (this *Semaphore) TryAlloc() bool {
	select {
	case this.flag <- true:
		return true
	default:
	}
	return false
}

func (this *Semaphore) Free() bool {
	select {
	case <-this.flag:
		return true
	default:
	}
	return false
}
