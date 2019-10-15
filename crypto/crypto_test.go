package crypto

import (
	"log"
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

/**
测试aes-256-cbc加密
$ go test -v -test.run=TestAesCbc
=== RUN   TestAesCbc
2019/10/15 23:37:41 /fxQRPGIHJ9AFsG67MSVDvLFSDp+/ZFGkHT+Y46h4jln9IzORfsEhR6L2qh5mDDQ
2019/10/15 23:37:41 HRHtimkjsJktwu6AzH2ji9MP9OLpRBRf35Xcm7zFNmr5Lj8X1rxxJiCIQJqnLC8r
2019/10/15 23:37:41 Sj1ENtUBam7C6PglPZgLZGy9lC8bppce7NS8RExuVa+xWow04Trnlc+kJh+Wz9LL
2019/10/15 23:37:41 中文数字123字母ABC符号!@#$%^&*() <nil>
--- PASS: TestAesCbc (0.00s)
PASS
ok      github.com/daheige/thinkgo/crypto       0.002s
*/
func TestAesCbc(t *testing.T) {
	str := `中文数字123字母ABC符号!@#$%^&*()`
	k2 := `abcdefghijklmnop`
	iv2 := `1234567890123456`
	b, _ := AesEncrypt(str, k2, iv2)
	log.Println(string(b))

	// log.Println(AesDecrypt(`/fxQRPGIHJ9AFsG67MSVDvLFSDp+/ZFGkHT+Y46h4jln9IzORfsEhR6L2qh5mDDQ`, k2, iv2))

	k2 = `abcdefghijklmnop1234567890123456`
	iv2 = `1234567890123456`
	b, _ = AesEncrypt(str, k2, iv2)
	log.Println(b)

	k2 = `abcdefghijklmnop12345678`
	iv2 = `1234567890123456`
	b, _ = AesEncrypt(str, k2, iv2)
	log.Println(b)

	log.Println(AesDecrypt(b, k2, iv2))
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
ok      github.com/daheige/thinkgo/crypto  0.002s
*/
