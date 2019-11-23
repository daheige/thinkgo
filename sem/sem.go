package sem

// 指定数量的空结构体缓存通道，实现信息号实现互斥锁
// Mutually exclusive by channel semaphore
// A semaphore is a synchronization pattern/primitive
// that imposes mutual exclusion on a limited number of resources.
// 信号量是同步模式/原语，它在有限数量的资源上强加互斥

import (
	"errors"
	"time"
)

var (
	ErrNoTickets      = errors.New("semaphore: could not aquire semaphore")
	ErrIllegalRelease = errors.New("semaphore: can't release the semaphore without acquiring it first")
)

// SemInterface contains the behavior of a semaphore that can be acquired and/or released.
// Other types implement methods in the interface to implement a mutex
type SemInterface interface {
	Acquire() error
	Release() error
}

// sem define
type semaphore struct {
	sem     chan struct{}
	timeout time.Duration //acquire/release timeout
}

//New create semaphonre mutex lock with timeout,tickets: a limited number of resources
func New(tickets int, timeout time.Duration) SemInterface {
	return &semaphore{
		sem:     make(chan struct{}, tickets),
		timeout: timeout,
	}
}

// Acquire get a sem
func (s *semaphore) Acquire() error {
	select {
	case s.sem <- struct{}{}:
		return nil
	case <-time.After(s.timeout):
		return ErrNoTickets
	}
}

// Release release sem
func (s *semaphore) Release() error {
	t := time.After(s.timeout)
	select {
	case <-s.sem:
		return nil
	case <-t: //release error
		return ErrIllegalRelease
	}
}
