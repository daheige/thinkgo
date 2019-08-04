package common

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const mutexLocked = 1 << iota

//创建lock实例
func NewMutexLock() *Mutex {
	return &Mutex{}
}

type Mutex struct {
	in sync.Mutex
}

//枷锁
func (m *Mutex) Lock() {
	m.in.Lock()
}

//解锁
func (m *Mutex) Unlock() {
	m.in.Unlock()
}

// 尝试枷锁
func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.in)), 0, mutexLocked)
}
