// chan实现trylock乐观锁
package chanlock

import (
	"time"
)

// ChanLock chan lock
type ChanLock struct {
	ch chan struct{} // 空结构体
}

// NewChanLock 实例化一个通道空结构体锁对象
func NewChanLock() *ChanLock {
	return &ChanLock{
		ch: make(chan struct{}, 1), // 有缓冲通道
	}
}

// Lock 通道枷锁
func (l *ChanLock) Lock() {
	l.ch <- struct{}{} // 这里是一个空结构体
}

// Unlock实现通道解锁
func (l *ChanLock) Unlock() {
	<-l.ch
}

// TryLock 乐观锁实现
func (l *ChanLock) TryLock() bool {
	select {
	case l.ch <- struct{}{}:
		return true
	default:
	}

	return false
}

// TryLockTimeout 指定时间内的乐观锁
func (l *ChanLock) TryLockTimeout(timeout time.Duration) bool {
	t := time.After(timeout)
	select {
	case l.ch <- struct{}{}:
		return true
	case <-t:
		return false // timeout
	}
}
