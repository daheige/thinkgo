package common

import (
	"testing"
)

func TestTryLock(t *testing.T) {
	t.Log("test trylock")
	var m = NewMutexLock()
	if m.TryLock() {
		t.Log("加锁成功!")
		m.Unlock()
	}

	m.Lock()
	t.Log("haha")
	m.Unlock()

}

func TestRace(t *testing.T) {
	var mu Mutex
	var x int
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			go func() {
				if mu.TryLock() {
					x++
					mu.Unlock()

				}
			}()
			continue
		}
		go func() {
			mu.Lock()
			x++
			mu.Unlock()

		}()
	}

}
