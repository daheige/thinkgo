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

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1(s string) string {
	r := sha1.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}

func Sha1File(fName string) string {
	f, e := os.Open(fName)
	if e != nil {
		return ""
	}

	h := sha1.New()
	_, e = io.Copy(h, f)
	if e != nil {
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}

//hmac256算法
func Hmac256(data, key string) string {
	if len(key) != 16 {
		key = GetIteratorStr(16)
	}

	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

//得到指定16进制的数字
func GetIteratorStr(length int) string {
	str := "0123456789abcdef"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}

//sha256得到的值是一个固定值
func Sha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//ECB加密
func EncryptEcb(src, key string) (string, error) {
	data := []byte(src)
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return "", err
	}

	bs := block.BlockSize()
	//对明文数据进行补码
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		panic("Need a multiple of the blocksize")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//对明文按照blocksize进行分块加密
		//必要时可以使用go关键字进行并行加密
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}

	return hex.EncodeToString(out), nil
}

/**
ECB（电子密本方式）就是将数据按照8个字节一段进行DES加密或解密得到一段8个字节的密文或者明文，最后一段不足8个字节，按照需求补足8个字节进行计算，之后按照顺序将计算所得的数据连在一起即可，各段数据之间互不影响。
特点
简单，有利于并行计算，误差不会被传送；
不能隐藏明文的模式；在密文中出现明文消息的重复
可能对明文进行主动攻击；加密消息块相互独立成为被攻击的弱点
*/
//ECB解密 key必须是8位
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
//CBC加密 key,iv必须是16位
func AesEncrypt(encodeStr, key string, iv string) (string, error) {
	encodeBytes := []byte(encodeStr)
	//根据key 生成密文
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	encodeBytes = PKCS5Padding(encodeBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	crypted := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(crypted, encodeBytes)

	return base64.StdEncoding.EncodeToString(crypted), nil
}

//CBC解密key,iv必须是16位
func AesDecrypt(decodeStr, key, iv string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr) //先解密base64
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	origData := make([]byte, len(decodeBytes))

	blockMode.CryptBlocks(origData, decodeBytes)
	origData = PKCS5UnPadding(origData)

	return string(origData), nil
}

//明文补码算法
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//明文减码算法
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
