package common

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// CopyFile 复制文件到另一个文件
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

func StoreGobData(data interface{}, fileName string) error {
	buf := new(bytes.Buffer) //创建写入缓冲区
	//创建gob编码器
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(data) //将data数据编码到缓冲区
	if err != nil {
		log.Println("encode gob data error: ", err.Error())
		return err
	}

	//将缓冲区中已编码的数据写入文件中
	err = ioutil.WriteFile(fileName, buf.Bytes(), 0644)
	if err != nil {
		log.Println("write gob data error: ", err.Error())
		return err
	}

	return nil
}

// LoadGobData 将gob写入的内容，载入到data中
func LoadGobData(data interface{}, fileName string) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println("read gob data error: ", err.Error())
		return
	}

	//根据这些原始数据，创建缓冲区
	buf := bytes.NewBuffer(raw)
	//将数据解码到缓冲区 (为缓冲区创建解码器)
	dec := gob.NewDecoder(buf)
	//解码数据到data中
	err = dec.Decode(data)
	if err != nil {
		log.Println("get gob data error: ", err.Error())
		return
	}
}

// CheckPathExist check file or path exist
func CheckPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

// Filebase 获取文件的名称不带后缀
// Get the name of the file without a suffix
func Filebase(file string) string {
	beg, end := len(file)-1, len(file)
	for ; beg >= 0; beg-- {
		if os.IsPathSeparator(file[beg]) {
			beg++
			break
		} else if file[beg] == '.' {
			end = beg
		}
	}
	return file[beg:end]
}

// Fileline 获取文件名:行数
func Fileline(file string, line int) string {
	beg, end := len(file)-1, len(file)
	for ; beg >= 0; beg-- {
		if os.IsPathSeparator(file[beg]) {
			beg++
			break
		} else if file[beg] == '.' {
			end = beg
		}
	}

	return fmt.Sprint(file[beg:end], ":", line)
}

// Chown 清空文件并保持文件权限不变，并非linux chown操作
// Empty the file and keep the file permissions unchanged
// not the linux chown operation
func Chown(name string, info os.FileInfo) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}

	f.Close()

	return nil
}
