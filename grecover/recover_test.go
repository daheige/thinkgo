package grecover

import (
	"log"
	"testing"
)

func TestCatchStack(t *testing.T) {
	t.Log("test grecover")

	TracePanic = true
	testSay()

	log.Println("ok")
}

func testSay() {
	defer CheckPanic()

	log.Println("11111")
	s := []string{
		"a", "b", "c",
	}

	// mock slice panic
	log.Println("d: ", s[3])
}

/**
=== RUN   TestCatchStack
    TestCatchStack: recover_test.go:9: test grecover
2020/05/23 20:33:15 11111
2020/05/23 20:33:15 panic error:  runtime error: index out of range [3] with length 3
2020/05/23 20:33:15 goroutine 6 [running]:
runtime/debug.Stack(0xc0000680a0, 0xc00000c100, 0x2)
	/usr/local/go/src/runtime/debug/stack.go:24 +0x9d
github.com/daheige/thinkgo/grecover.CatchStack(...)
	/mygo/web/go/thinkgo/grecover/stack_helper.go:28
github.com/daheige/thinkgo/grecover.CheckPanic()
	/mygo/web/go/thinkgo/grecover/stack_helper.go:21 +0xcf
panic(0x1132280, 0xc00001a1c0)
	/usr/local/go/src/runtime/panic.go:969 +0x166
github.com/daheige/thinkgo/grecover.testSay()
	/mygo/web/go/thinkgo/grecover/recover_test.go:24 +0x90
github.com/daheige/thinkgo/grecover.TestCatchStack(0xc00012c120)
	/mygo/web/go/thinkgo/grecover/recover_test.go:11 +0x72
testing.tRunner(0xc00012c120, 0x114c888)
	/usr/local/go/src/testing/testing.go:991 +0xdc
created by testing.(*T).Run
	/usr/local/go/src/testing/testing.go:1042 +0x357

2020/05/23 20:33:15 ok
--- PASS: TestCatchStack (0.00s)
PASS
*/
