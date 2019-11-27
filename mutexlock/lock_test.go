package mutexlock

import (
	"testing"
)

func TestTryLock(t *testing.T) {
	t.Log("test trylock")
	var mutex = NewMutexLock()
	if mutex.TryLock() {
		t.Log("加锁成功!")
		mutex.Unlock()
	}

	mutex.Lock()
	t.Log("haha")
	mutex.Unlock()

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

/**
$ go test -v
=== RUN   TestTryLock
--- PASS: TestTryLock (0.00s)
    trylock_test.go:8: test trylock
    trylock_test.go:11: 加锁成功!
    trylock_test.go:16: haha
=== RUN   TestRace
--- PASS: TestRace (0.00s)
PASS
ok      github.com/daheige/thinkgo/mutexlock    0.003s
*/
