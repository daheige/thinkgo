package common

import (
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"os"
)

//复制文件到另一个文件
func CopyFile(distName, srcName string) (w int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dist, err := os.OpenFile(distName, os.O_CREATE|os.O_WRONLY, 0644) //打开文件，如果不存在就创建一个
	if err != nil {
		return
	}
	defer dist.Close() //关闭文件句柄

	return io.Copy(dist, src)
}

func StoreGobData(data interface{}, fileName string) {
	buf := new(bytes.Buffer) //创建写入缓冲区
	//创建gob编码器
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(data) //将data数据编码到缓冲区
	if err != nil {
		panic(err)
	}

	//将缓冲区中已编码的数据写入文件中
	err = ioutil.WriteFile(fileName, buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

//将gob写入的内容，载入到data中
func LoadGobData(data interface{}, fileName string) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	//根据这些原始数据，创建缓冲区
	buf := bytes.NewBuffer(raw)
	//将数据解码到缓冲区 (为缓冲区创建解码器)
	dec := gob.NewDecoder(buf)
	//解码数据到data中
	err = dec.Decode(data)
	if err != nil {
		panic(err)
	}
}
