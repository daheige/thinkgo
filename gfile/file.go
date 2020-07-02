// file 文件相关的一些辅助函数
package gfile

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

//========================= Directory/Filesystem Functions=======

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

// StoreGobData store gob data
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

// Stat stat()
func Stat(filename string) (os.FileInfo, error) {
	return os.Stat(filename)
}

// Pathinfo pathinfo()
// -1: all; 1: dirname; 2: basename; 4: extension; 8: filename
// Usage:
// Pathinfo("/home/go/path/src/php2go/php2go.go", 1|2|4|8)
func Pathinfo(path string, options int) map[string]string {
	if options == -1 {
		options = 1 | 2 | 4 | 8
	}

	info := make(map[string]string)
	if (options & 1) == 1 {
		info["dirname"] = filepath.Dir(path)
	}

	if (options & 2) == 2 {
		info["basename"] = filepath.Base(path)
	}

	if ((options & 4) == 4) || ((options & 8) == 8) {
		var basename string
		if (options & 2) == 2 {
			basename = info["basename"]
		} else {
			basename = filepath.Base(path)
		}

		p := strings.LastIndex(basename, ".")

		var filename string
		var extension string
		if p > 0 {
			filename, extension = basename[:p], basename[p+1:]
		} else if p == -1 {
			filename = basename
		} else if p == 0 {
			extension = basename[p+1:]
		}

		if (options & 4) == 4 {
			info["extension"] = extension
		}

		if (options & 8) == 8 {
			info["filename"] = filename
		}
	}

	return info
}

// FileExists file_exists()
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// IsFile is_file()
func IsFile(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// IsDir is_dir()
func IsDir(filename string) (bool, error) {
	fd, err := os.Stat(filename)
	if err != nil {
		return false, err
	}

	return fd.Mode().IsDir(), nil
}

// FileSize filesize()
func FileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return 0, err
	}

	return info.Size(), nil
}

// FilePutContents file_put_contents()
func FilePutContents(filename string, data string, mode os.FileMode) error {
	return ioutil.WriteFile(filename, []byte(data), mode)
}

// FileGetContents file_get_contents()
func FileGetContents(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)

	return string(data), err
}

// Unlink unlink()
func Unlink(filename string) error {
	return os.Remove(filename)
}

// IsReadable is_readable()
func IsReadable(filename string) bool {
	_, err := syscall.Open(filename, syscall.O_RDONLY, 0)

	return err == nil
}

// IsWriteable is_writeable()
func IsWriteable(filename string) bool {
	_, err := syscall.Open(filename, syscall.O_WRONLY, 0)

	return err == nil
}

// Rename rename()
func Rename(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

// Touch touch()
func Touch(filename string) (bool, error) {
	fd, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return false, err
	}

	fd.Close()

	return true, nil
}

// Mkdir mkdir()
func Mkdir(filename string, mode os.FileMode) error {
	return os.Mkdir(filename, mode)
}

// Getcwd getcwd()
func Getcwd() (string, error) {
	return os.Getwd()
}

// Realpath realpath()
func Realpath(path string) (string, error) {
	return filepath.Abs(path)
}

// Basename basename()
func Basename(path string) string {
	return filepath.Base(path)
}

// Chmod chmod()
func Chmod(filename string, mode os.FileMode) bool {
	return os.Chmod(filename, mode) == nil
}

// Chown chown()
func FileChown(filename string, uid, gid int) bool {
	return os.Chown(filename, uid, gid) == nil
}

// Fclose fclose()
func Fclose(handle *os.File) error {
	return handle.Close()
}

// Filemtime filemtime()
func Filemtime(filename string) (int64, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return 0, err
	}

	defer fd.Close()

	fileinfo, err := fd.Stat()
	if err != nil {
		return 0, err
	}

	return fileinfo.ModTime().Unix(), nil
}

// Fgetcsv fgetcsv()
func Fgetcsv(handle *os.File, length int, delimiter rune) ([][]string, error) {
	reader := csv.NewReader(handle)
	reader.Comma = delimiter

	// TODO length limit
	return reader.ReadAll()
}

// Glob glob()
func Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}
