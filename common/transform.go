package common

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"io/ioutil"
)

// hash : []byte to uint64
func Hash(mem []byte) uint64 {
	var hash uint64 = 5381
	for _, b := range mem {
		hash = (hash << 5) + hash + uint64(b)
	}
	return hash
}

// compress data use gzip
func Gzip(in []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	if wt, err := gzip.NewWriterLevel(buf, gzip.BestCompression); err != nil {
		return nil, err
	} else {
		if _, err := wt.Write(in); err != nil {
			return nil, err
		}
		wt.Close()
	}
	return buf.Bytes(), nil
}

// decompress data user gunzip
func Gunzip(in []byte) ([]byte, error) {
	rd, err := gzip.NewReader(bytes.NewBuffer(in))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(rd)
}

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
