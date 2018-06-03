package crypto

import (
	"testing"
)

func Test_md5(t *testing.T) {
	t.Log(Md5encode("123456"))
	t.Log("success")
}
