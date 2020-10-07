package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"math/rand"
	"os"
	"time"
)

// randStr 用于生成随机字符串
var randStr = "0123456789abcdef"

// Md5 md5 string
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// Sha1 sha1 string
func Sha1(s string) string {
	r := sha1.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}

// Sha1File sha1 file
func Sha1File(fName string) (string, error) {
	f, err := os.Open(fName)
	if err != nil {
		return "", err
	}

	defer f.Close()

	h := sha1.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// Hmac256 hmac256算法
func Hmac256(data, key string) string {
	if len(key) != 16 {
		key = GetIteratorStr(16)
	}

	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// HmacSha1 实现php hmac_sha1
func HmacSha1(str string, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(str))
	return hex.EncodeToString(mac.Sum(nil))
}

// GetIteratorStr 得到指定16进制的数字
func GetIteratorStr(length int) string {
	b := []byte(randStr)
	bLen := len(b)
	res := make([]byte, 0, length+1)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		res = append(res, b[r.Intn(bLen)])
	}

	return string(res)
}

// Sha256 sha256得到的值是一个固定值
func Sha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// EncryptEcb ECB加密
func EncryptEcb(src, key string) (string, error) {
	data := []byte(src)
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return "", err
	}

	bs := block.BlockSize()
	// 对明文数据进行补码
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		return "", errors.New("need a multiple of the block size")
	}

	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		// 对明文按照blocksize进行分块加密
		// 必要时可以使用go关键字进行并行加密
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}

	return hex.EncodeToString(out), nil
}

/**
ECB（电子密本方式）就是将数据按照8个字节一段进行DES加密或解密得到一段8个字节的密文或者明文，
最后一段不足8个字节，按照需求补足8个字节进行计算，之后按照顺序将计算所得的数据连在一起即可，
各段数据之间互不影响。
特点：
简单，有利于并行计算，误差不会被传送；
不能隐藏明文的模式；在密文中出现明文消息的重复
可能对明文进行主动攻击；加密消息块相互独立成为被攻击的弱点
*/

// DecryptEcb ECB解密 key必须是8位
func DecryptEcb(src, key string) (string, error) {
	data, err := hex.DecodeString(src)
	if err != nil {
		return "", err
	}

	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return "", err
	}

	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return "", errors.New("crypto/cipher: input not full blocks")
	}

	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}

	out = PKCS5UnPadding(out)
	return string(out), nil
}

/**
概念
CBC（密文分组链接方式）有向量的概念, 它的实现机制使加密的各段数据之间有了联系。
加密步骤：
    首先将数据按照8个字节一组进行分组得到D1D2......Dn（若数据不是8的整数倍，用指定的PADDING数据补位）
    第一组数据D1与初始化向量I异或后的结果进行DES加密得到第一组密文C1（初始化向量I为全零）
    第二组数据D2与第一组的加密结果C1异或以后的结果进行DES加密，得到第二组密文C2
    之后的数据以此类推，得到Cn
    按顺序连为C1C2C3......Cn即为加密结果。
解密是加密的逆过程：
    首先将数据按照8个字节一组进行分组得到C1C2C3......Cn
    将第一组数据进行解密后与初始化向量I进行异或得到第一组明文D1（注意：一定是先解密再异或）
    将第二组数据C2进行解密后与第一组密文数据进行异或得到第二组数据D2
    之后依此类推，得到Dn
    按顺序连为D1D2D3......Dn即为解密结果。
特点
    不容易主动攻击,安全性好于ECB,适合传输长度长的报文,是SSL、IPSec的标准。
    每个密文块依赖于所有的信息明文消息中一个改变会影响所有密文块
    发送方和接收方都需要知道初始化向量
    加密过程是串行的，无法被并行化(在解密时，从两个邻接的密文块中即可得到一个平文块。因此，解密过程可以被并行化
*/
// AesEncrypt CBC加密 key
// iv必须是16位
// 当key 16位的时候 相当于php base64_encode(openssl_encrypt($str, 'aes-128-cbc', $key, true, $iv));
// 当key 24位的时候 相当于php base64_encode(openssl_encrypt($str, 'aes-192-cbc', $key, true, $iv));
// 当key 32位的时候 相当于php base64_encode(openssl_encrypt($str, 'aes-256-cbc', $key, true, $iv));
func AesEncrypt(encodeStr, key string, iv string) (string, error) {
	// 根据key 生成密文
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	encodeBytes := []byte(encodeStr)
	blockSize := block.BlockSize()
	encodeBytes = PKCS5Padding(encodeBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	c := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(c, encodeBytes)

	return base64.StdEncoding.EncodeToString(c), nil
}

// AesDecrypt CBC解密key
// iv必须是16位
// 对应php 解密方式
// 当key 16位的时候 相当于php openssl_decrypt(base64_decode($strEncode), 'aes-128-cbc', $key, true, $iv)
// 当key 24位的时候 相当于php openssl_decrypt(base64_decode($strEncode), 'aes-192-cbc', $key, true, $iv)
// 当key 32位的时候 相当于php openssl_decrypt(base64_decode($strEncode), 'aes-256-cbc', $key, true, $iv)
func AesDecrypt(decodeStr, key, iv string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr) // 先解密base64
	if err != nil {
		return "", err
	}

	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	origData := make([]byte, len(decodeBytes))

	blockMode.CryptBlocks(origData, decodeBytes)
	origData = PKCS5UnPadding(origData)

	return string(origData), nil
}

// PKCS5Padding 明文补码算法
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// PKCS5UnPadding 明文减码算法
func PKCS5UnPadding(origData []byte) []byte {
	oLen := len(origData)
	unPadLen := int(origData[oLen-1])
	return origData[:(oLen - unPadLen)]
}
