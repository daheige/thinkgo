package gutils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"html"
	"io"
	"math/rand"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	uuid "github.com/satori/go.uuid"

	"github.com/daheige/thinkgo/crypto"
	"github.com/daheige/thinkgo/gnum"
)

// Addslashes addslashes()
func Addslashes(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '\'', '"', '\\':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// Stripslashes stripslashes()
func Stripslashes(str string) string {
	var buf bytes.Buffer
	l, skip := len(str), false
	for i, char := range str {
		if skip {
			skip = false
		} else if char == '\\' {
			if i+1 < l && str[i+1] == '\\' {
				skip = true
			}
			continue
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// ===================str join func========================
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

// ==============str md5,md5File,sha1,crc32,bin2hex,hex2bin,hash func=======
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

// Decbin decbin()
func Decbin(number int64) string {
	return strconv.FormatInt(number, 2)
}

// Bindec bindec()
func Bindec(str string) (string, error) {
	i, err := strconv.ParseInt(str, 2, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, 10), nil
}

// Hash : []byte to uint64
func Hash(mem []byte) uint64 {
	var hash uint64 = 5381
	for _, b := range mem {
		hash = (hash << 5) + hash + uint64(b)
	}
	return hash
}

// =================str aes/des func==========================
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

// =================uuid,rnduuid,uniqid func====================
// NewUUID 通过随机数的方式生成uuid
// 如果rand.Read失败,就按照当前时间戳+随机数进行md5方式生成
// 该方式生成的uuid有可能存在重复值
// 返回格式:7999b726-ca3c-42b6-bda2-259f4ac0879a
func NewUUID() string {
	u := [16]byte{}
	ns := time.Now().UnixNano()
	rand.Seed(ns)
	_, err := rand.Read(u[0:])
	if err != nil {
		rndStr := strconv.FormatInt(ns, 10) + strconv.FormatInt(gnum.RandInt64(1000, 9999), 10)
		s := crypto.Md5(rndStr)
		return fmt.Sprintf("%s-%s-%s-%s-%s", s[:8], s[8:12], s[12:16], s[16:20], s[20:])
	}

	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// RndUuid 基于时间ns和随机数实现唯一的uuid
// 在单台机器上是不会出现重复的uuid
// 如果要在分布式的架构上生成不重复的uuid
// 只需要在rndStr的前面加一些自定义的字符串
// 返回格式:eba1e8cd-0460-4910-49c6-44bdf3cf024d
func RndUuid() string {
	s := RndUuidMd5()
	return fmt.Sprintf("%s-%s-%s-%s-%s", s[:8], s[8:12], s[12:16], s[16:20], s[20:])
}

// RndUuidMd5 uuid
func RndUuidMd5() string {
	ns := time.Now().UnixNano()
	rndStr := StrJoin("", strconv.FormatInt(ns, 10), strconv.FormatInt(gnum.RandInt64(1000, 9999), 10))
	return crypto.Md5(rndStr)
}

func Uuid() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

// Uniqid uniqid()
func Uniqid(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%08x%05x", prefix, now.Unix(), now.UnixNano()%0x100000)
}

// =====================html special characters================
// HTMLSpecialchars converts special characters to HTML entities
func HTMLSpecialchars(str string) string {
	return html.EscapeString(str)
}

// HTMLSpecialcharsDecode converts special HTML entities back to characters
func HTMLSpecialcharsDecode(str string) string {
	return html.UnescapeString(str)
}

// HTMLEntities htmlentities()
func HTMLEntities(str string) string {
	return html.EscapeString(str)
}

// HTMLEntityDecode html_entity_decode()
func HTMLEntityDecode(str string) string {
	return html.UnescapeString(str)
}

// =====================str xss,XssUnescape func=============
// Xss 防止xss攻击
func Xss(str string) string {
	if len(str) == 0 {
		return ""
	}

	return html.EscapeString(str)
}

// XssUnescape 还原xss字符串
func XssUnescape(str string) string {
	return html.UnescapeString(str)
}

// ================str krand func===========
// Krand 根据kind生成不同风格的指定区间随机字符串
// 纯数字kind=0,小写字母kind=1
// 大写字母kind=2,数字+大小写字母kind=3
func Krand(size int, kind int) string {
	oldKind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano()) // 随机种子
	for i := 0; i < size; i++ {
		if is_all { // random oldKind
			oldKind = rand.Intn(3)
		}

		scope, base := kinds[oldKind][0], kinds[oldKind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}

	return string(result)
}

// ================str chr,ord func======================
// Chr returns a one-character string containing the character specified by ascii
// go1.15 return string(rune(ascii))
func Chr(ascii int) string {
	for ascii < 0 {
		ascii += 256
	}

	ascii %= 256
	return string(rune(ascii))
}

// Ord return ASCII value of character
func Ord(character string) rune {
	return []rune(character)[0]
}

// =================str explode,implode,strlen================
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

// MbStrlen mb_strlen()
func MbStrlen(str string) int {
	return utf8.RuneCountInString(str)
}

// =================str strpos,Strrpos,stripos,Strripos func====================
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

// StrReplace str_replace()
func StrReplace(search, replace, subject string, count int) string {
	return strings.Replace(subject, search, replace, count)
}

// StrRepeat str_repeat()
func StrRepeat(input string, multiplier int) string {
	return strings.Repeat(input, multiplier)
}

// Strstr strstr()
func Strstr(haystack string, needle string) string {
	if needle == "" {
		return ""
	}
	idx := strings.Index(haystack, needle)
	if idx == -1 {
		return ""
	}
	return haystack[idx+len([]byte(needle))-1:]
}

// Substr substr()
func Substr(str string, start uint, length int) string {
	if length < -1 {
		return str
	}

	switch {
	case length == -1:
		return str[start:]
	case length == 0:
		return ""
	}

	end := int(start) + length
	if end > len(str) {
		end = len(str)
	}

	return str[start:end]
}

// ==================str upper/lower==============================

// Strtoupper strtoupper(str) makes a string uppercase
func Strtoupper(str string) string {
	return strings.ToUpper(str)
}

// Strtolower strtolower(str) makes a string lowercase
func Strtolower(str string) string {
	return strings.ToLower(str)
}

// StrShuffle str_shuffle(str)
func StrShuffle(str string) string {
	runes := []rune(str)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := make([]rune, len(runes))
	for i, v := range r.Perm(len(runes)) {
		s[i] = runes[v]
	}
	return string(s)
}

// Trim trim()
func Trim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimSpace(str)
	}
	return strings.Trim(str, characterMask[0])
}

// Ltrim ltrim()
func Ltrim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimLeftFunc(str, unicode.IsSpace)
	}

	return strings.TrimLeft(str, characterMask[0])
}

// Rtrim rtrim()
func Rtrim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimRightFunc(str, unicode.IsSpace)
	}

	return strings.TrimRight(str, characterMask[0])
}

// =======================str Ucfirst/Lcfirst/Ucwords==================

// UcFirst ucfirst(str) make a string's first character uppercase
func UcFirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return u + str[len(u):]
	}
	return ""
}

// LcFirst lcfirst(str) make a string's first character lowercase
func LcFirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		return u + str[len(u):]
	}
	return ""
}

// Ucwords ucwords(str)
// uppercases the first character of each word in a string
func Ucwords(str string) string {
	return strings.Title(str)
}

// ======================url encode/decode=======================

// URLEncode urlencode(str) url encode
func URLEncode(str string) string {
	return url.QueryEscape(str)
}

// URLDecode urldecode(str) url decode
func URLDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// Rawurlencode rawurlencode(str)
func Rawurlencode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// Rawurldecode rawurldecode(str)
func Rawurldecode(str string) (string, error) {
	return url.QueryUnescape(strings.Replace(str, "%20", "+", -1))
}

// HTTPBuildQuery http_build_query() url a=1&b=2
func HTTPBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}

// =================base64 encode/decode=========================

// Base64Encode base64_encode(str)
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode base64_decode(str)
func Base64Decode(str string) (string, error) {
	switch len(str) % 4 {
	case 2:
		str += "=="
	case 3:
		str += "="
	}

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Empty empty()
func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// ========================ip convert===========
// IP2long ip2long()
func IP2long(ipAddress string) uint32 {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 0
	}
	return binary.BigEndian.Uint32(ip.To4())
}

// Long2ip long2ip()
func Long2ip(properAddress uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, properAddress)
	ip := net.IP(ipByte)
	return ip.String()
}
