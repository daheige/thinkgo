package gutils

import (
	"log"
	"testing"
)

/**
$ go test -v -test.run TestUuid
--- PASS: TestUuid (24.51s)
PASS
ok      github.com/daheige/thinkgo/common  24.517s
*/
func TestUuid(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		// log.Println("current newuuid", NewUUID())
		log.Println("current uuid of version4: ", Uuid())
	}
}

func TestRndUuid(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		log.Println("current rnd uuid", RndUuid())
	}
}

/**
go test -v
2019/11/27 22:28:53 current rnd uuid 9ddf1430-deb8-ef32-dfea-4a3f71f45404
2019/11/27 22:28:53 current rnd uuid aa32b47a-02f6-714b-a593-fdb407838330
2019/11/27 22:28:53 current rnd uuid 9076963a-cbb5-d6c4-789d-1c7eb6751966
--- PASS: TestRndUuid (36.11s)
PASS
ok      github.com/daheige/thinkgo/gutils       75.358s
*/
