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

func (l *ChanLock) TryLock(timeout time.Duration) bool {
	t := time.After(timeout)
	select {
	case l.ch <- struct{}{}:
		return true
	case <-t:
		return false //timeout
	}
}
