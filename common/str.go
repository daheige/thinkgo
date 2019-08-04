package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"html"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/daheige/thinkgo/crypto"
	uuid "github.com/satori/go.uuid"
)

//================str,int,int64,float64 conv func=======================
// IntToStr int-->string
func IntToStr(n int) string {
	return strconv.Itoa(n)
}

// StrToInt string-->int
func StrToInt(s string) int {
	if i, err := strconv.Atoi(s); err != nil {
		return 0
	} else {
		return i
	}
}

// Int64ToStr int64-->string
func Int64ToStr(i64 int64) string {
	return strconv.FormatInt(i64, 10)
}

// StrToInt64 string--> int64
func StrToInt64(s string) int64 {
	if i64, err := strconv.ParseInt(s, 10, 64); err != nil {
		return 0
	} else {
		return i64
	}
}

// StrToFloat64 string--->float64
func StrToFloat64(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}

	return f
}

// Float64ToStr float64 to string
// 'e' (-d.dddde±dd，十进制指数)
func Float64ToStr(f64 float64) string {
	return strconv.FormatFloat(f64, 'e', -1, 64)
}

//===================str join func========================
// StrJoin 多个字符串按照指定的分隔符进行拼接
func StrJoin(sep string, str ...string) string {
	return strings.Join(str, sep)
}

// StrJoinByBuf 通过buf缓冲区的方式连接字符串
func StrJoinByBuf(str ...string) string {
	if len(str) == 0 {
		return ""
	}

	var buf bytes.Buffer
	for _, s := range str {
		buf.WriteString(s)
	}

	return buf.String()
}

//==============str md5,md5File,sha1,crc32,bin2hex,hex2bin,hash func=======
// Other advanced functions, please see the thinkgo/crypto package.
// md5 func
func Md5(str string) string {
	return crypto.Md5(str)
}

// Md5File calculates the md5 hash of a given file
func Md5File(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// sha1 string
func Sha1(str string) string {
	return crypto.Sha1(str)
}

// Sha1File calculates the sha1 hash of a file
func Sha1File(path string) (string, error) {
	return crypto.Sha1File(path)
}

// Crc32 calculates the crc32 polynomial of a string
func Crc32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

// Bin2hex converts binary data into hexadecimal representation
func Bin2hex(src []byte) string {
	return hex.EncodeToString(src)
}

// Hex2bin decodes a hexadecimally encoded binary string
func Hex2bin(str string) []byte {
	s, _ := hex.DecodeString(str)
	return s
}

// Hash : []byte to uint64
func Hash(mem []byte) uint64 {
	var hash uint64 = 5381
	for _, b := range mem {
		hash = (hash << 5) + hash + uint64(b)
	}
	return hash
}

//=================str aes/des func==========================
// key = "abcdefghijklmnopqrstuvwxyz123456"
// iv = "0123456789ABCDEF"
func EnAES(in, key, iv []byte) ([]byte, error) {
	cip, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	out := make([]byte, len(in))
	cipher.NewCFBEncrypter(cip, iv).XORKeyStream(out, in)
	return out, nil
}

func DeAES(in, key, iv []byte) ([]byte, error) {
	cip, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	out := make([]byte, len(in))
	cipher.NewCFBDecrypter(cip, iv).XORKeyStream(out, in)
	return out, nil
}

//=================uuid,rnduuid func====================
// NewUUID 通过随机数的方式生成uuid
//如果rand.Read失败,就按照当前时间戳+随机数进行md5方式生成
//该方式生成的uuid有可能存在重复值
//返回格式:7999b726-ca3c-42b6-bda2-259f4ac0879a
func NewUUID() string {
	u := [16]byte{}
	ns := time.Now().UnixNano()
	rand.Seed(ns)
	_, err := rand.Read(u[0:])
	if err != nil {
		rndStr := strconv.FormatInt(ns, 10) + strconv.FormatInt(RandInt64(1000, 9999), 10)
		s := crypto.Md5(rndStr)
		return fmt.Sprintf("%s-%s-%s-%s-%s", s[:8], s[8:12], s[12:16], s[16:20], s[20:])
	}

	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// RndUuid 基于时间ns和随机数实现唯一的uuid
//在单台机器上是不会出现重复的uuid
//如果要在分布式的架构上生成不重复的uuid
// 只需要在rndStr的前面加一些自定义的字符串
//返回格式:eba1e8cd-0460-4910-49c6-44bdf3cf024d
func RndUuid() string {
	s := RndUuidMd5()
	return fmt.Sprintf("%s-%s-%s-%s-%s", s[:8], s[8:12], s[12:16], s[16:20], s[20:])
}

// RndUuidMd5 uuid
func RndUuidMd5() string {
	ns := time.Now().UnixNano()
	rndStr := StrJoin("", strconv.FormatInt(ns, 10), strconv.FormatInt(RandInt64(1000, 9999), 10))
	return crypto.Md5(rndStr)
}

func Uuid() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

//=====================html special characters================
// HTMLSpecialchars converts special characters to HTML entities
func HTMLSpecialchars(str string) string {
	return html.EscapeString(str)
}

// HTMLSpecialcharsDecode converts special HTML entities back to characters
func HTMLSpecialcharsDecode(str string) string {
	return html.UnescapeString(str)
}

//=====================str xss,XssUnescape func=============
// Xss 防止xss攻击
func Xss(str string) string {
	if len(str) == 0 {
		return ""
	}

	return html.EscapeString(str)
}

func XssUnescape(str string) string {
	return html.UnescapeString(str)
}

//================str krand func===========
// Krand 根据kind生成不同风格的指定区间随机字符串
// 纯数字kind=0,小写字母kind=1
// 大写字母kind=2,数字+大小写字母kind=3
func Krand(size int, kind int) string {
	ikind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano()) //随机种子
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}

		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}

	return string(result)
}

//================str chr,ord func======================
// Chr returns a one-character string containing the character specified by ascii
func Chr(ascii int) string {
	for ascii < 0 {
		ascii += 256
	}
	ascii %= 256
	return string(ascii)
}

// Ord return ASCII value of character
func Ord(character string) rune {
	return []rune(character)[0]
}

//=================str explode,implode,strlen================
// Explode returns an slice of strings, each of which is a substring of str
// formed by splitting it on boundaries formed by the string delimiter.
func Explode(delimiter, str string) []string {
	return strings.Split(str, delimiter)
}

// Implode returns a string containing a string representation of all the slice
// elements in the same order, with the glue string between each element.
func Implode(glue string, pieces []string) string {
	return strings.Join(pieces, glue)
}

// Strlen get string length
// A multi-byte character is counted as 1
func Strlen(str string) int {
	return len([]rune(str))
}

//=================str strpos,Strrpos,stripos,Strripos func====================
// Strpos find position of first occurrence of string in a string
// It's multi-byte safe. return -1 if can not find the substring
func Strpos(haystack, needle string) int {

	pos := strings.Index(haystack, needle)

	if pos < 0 {
		return pos
	}

	rs := []rune(haystack[0:pos])

	return len(rs)
}

// Strrpos find the position of the last occurrence of a substring in a string
func Strrpos(haystack, needle string) int {

	pos := strings.LastIndex(haystack, needle)

	if pos < 0 {
		return pos
	}

	rs := []rune(haystack[0:pos])

	return len(rs)
}

// Stripos find position of the first occurrence of a case-insensitive substring in a string
func Stripos(haystack, needle string) int {
	return Strpos(strings.ToLower(haystack), strings.ToLower(needle))
}

// Strripos find the position of the last occurrence of a case-insensitive substring in a string
func Strripos(haystack, needle string) int {
	return Strrpos(strings.ToLower(haystack), strings.ToLower(needle))
}
