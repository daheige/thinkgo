package crypto

import (
	"testing"
)

func Test_md5(t *testing.T) {
	t.Log(Md5("123456"))
	t.Log("success")
}

func TestHmac256(t *testing.T) {
	t.Log(Hmac256("123456", ""))
	key := GetIteratorStr(16)
	t.Log("key: ", key)
	t.Log(Hmac256("123456", key))
}

func TestSha256(t *testing.T) {
	t.Log("test Sha256")
	t.Log(Sha256("123456"))
}

var k = GetIteratorStr(16)
var iv = GetIteratorStr(16)

func TestCbc256(t *testing.T) {
	t.Log(AesEncrypt("123456", k, iv))
}

func TestDecodeCbc256(t *testing.T) {
	b, _ := AesEncrypt("123456", k, iv)
	bytes, _ := AesDecrypt(b, k, iv)
	t.Log(string(bytes))
}

//test ecb
func TestAesEbc(t *testing.T) {
	k := GetIteratorStr(8)
	b, _ := EncryptEcb("123456", k)
	t.Log("ebc加密后:", b)

	s, _ := DecryptEcb(b, k)
	t.Log("ebc解密:", s)
}

/*
$ go test -v
=== RUN   Test_md5
--- PASS: Test_md5 (0.00s)
    md5_test.go:8: e10adc3949ba59abbe56e057f20f883e
    md5_test.go:9: success
=== RUN   TestHmac256
--- PASS: TestHmac256 (0.00s)
    md5_test.go:13: b8f19d151b14d384f924c369db08c04e
    md5_test.go:15: key:  f1574e46f3bf8ee1
    md5_test.go:16: b17018fab2fe58bc5913e347ae295cc9
=== RUN   TestSha256
--- PASS: TestSha256 (0.00s)
    md5_test.go:20: test Sha256
    md5_test.go:21: 8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92
=== RUN   TestCbc256
--- PASS: TestCbc256 (0.00s)
    md5_test.go:28: rJIiGlFEikKMkA2hv86Ubg== <nil>
=== RUN   TestDecodeCbc256
--- PASS: TestDecodeCbc256 (0.00s)
    md5_test.go:34: 123456
=== RUN   TestAesEbc
--- PASS: TestAesEbc (0.00s)
    md5_test.go:41: ebc加密后: 036B626594702270
    md5_test.go:42: ebc解密: 123456
PASS
ok      thinkgo/crypto  0.002s
*/
