package common

import (
	"testing"
)

func TestUuid(t *testing.T) {
	t.Log(NewUUID())
}

/**
=== RUN   TestUuid
--- PASS: TestUuid (0.00s)
    uuid_test.go:8: f1e022cd73f9e9e5de480275edcaa133
PASS
ok      thinkgo/common  0.005s
*/
