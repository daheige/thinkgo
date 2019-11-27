package gfile

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// Gzip compress data use gzip
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

// Gunzip decompress data user gunzip
func Gunzip(in []byte) ([]byte, error) {
	rd, err := gzip.NewReader(bytes.NewBuffer(in))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(rd)
}
