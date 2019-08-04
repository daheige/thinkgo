//chan实现trylock乐观锁
package common

import (
	"time"
)

type ChanLock struct {
	ch chan struct{} //空结构体
}

func NewChanLock() *ChanLock {
	return &ChanLock{
		ch: make(chan struct{}, 1), //有缓冲通道
	}
}

func (l *ChanLock) Lock() {
	l.ch <- struct{}{} //这里是一个空结构体
}

func (l *ChanLock) Unlock() {
	<-l.ch
}

//乐观锁实现
func (l *ChanLock) TryLock() bool {
	select {
	case l.ch <- struct{}{}:
		return true
	default:
	}

	return false
}

//指定时间内的乐观锁
func (l *ChanLock) TryLockTimeout(timeout time.Duration) bool {
	t := time.After(timeout)
	select {
	case l.ch <- struct{}{}:
		return true
	case <-t:
		return false //timeout
	}
}
