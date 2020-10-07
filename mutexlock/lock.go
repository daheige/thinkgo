/**
* Package mutexlock mutex trylock.
* 在sync.Mutex基础上，实现乐观锁TryLock
 */
package mutexlock

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const mutexLocked = 1 << iota

// NewMutexLock 创建lock实例
func NewMutexLock() *Mutex {
	return &Mutex{}
}

// Mutex mutex
type Mutex struct {
	in sync.Mutex
}

// Lock 加锁
func (m *Mutex) Lock() {
	m.in.Lock()
}

// Unlock 解锁
func (m *Mutex) Unlock() {
	m.in.Unlock()
}

// TryLock 尝试枷锁
func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.in)), 0, mutexLocked)
}
