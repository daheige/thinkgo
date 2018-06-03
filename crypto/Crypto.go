package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

//md5 crypto
func Md5encode(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	return strings.ToLower(hex.EncodeToString(cipherStr)) // 输出加密结果
}
