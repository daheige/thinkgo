package common

import (
	"log"
	"testing"
)

/**
$ go test -v -test.run TestUuid
--- PASS: TestUuid (24.51s)
PASS
ok      thinkgo/common  24.517s
*/
func TestUuid(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		log.Println("current newuuid", NewUUID())
	}
}

func TestRndUuid(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		log.Println("current rnd uuid", RndUuid())
	}
}

/**
$ go test -v -test.run TestRndUuid
--- PASS: TestRndUuid (25.18s)
PASS
ok      thinkgo/common  25.184s

*/
