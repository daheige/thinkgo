/**
*基于标准日志log封装而成的slog
*支持指定日志文件大小和日志名称的方式,打印日志
 */
package slog

import (
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	kb int64 = 1 << (10 * iota) //字节大小
	mb
)

var (
	logDir       = os.TempDir()
	logLevel     = log.LstdFlags                           //标准日志格式
	logWriteType = os.O_WRONLY | os.O_APPEND | os.O_CREATE //文件创建标识
	logSize      = 2 * kb                                  //日志文件大小，默认2mb
	logTimeZone  = "PRC"                                   //时区设置
	tmFmtMissMS  = "2006-01-02-15-04-05"
	logLock      = &sync.Mutex{}
)

//设置日志大小，默认2mb
func InitLogSize(size int64) {
	logSize = size * mb
}

func SetLogDir(dir string) {
	logDir = dir
	if _, err := os.Stat(logDir); err != nil { //日志目录不存在
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatalf("创建log目录%s失败,错误信息: %s", logDir, err)
			return
		}
	}
}

//设置日志记录时区
func SetLogTimeZone(zone string) {
	logTimeZone = zone
}

func checkLogDirExist(dir string) {
	if _, err := os.Stat(dir); err != nil { //日志目录不存在
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("创建log目录%s失败,错误信息: %s", dir, err)
			return
		}

		logDir = dir
	}
}

func Save(message interface{}, filename string, prefix string) {
	logLock.Lock() //写日志加锁
	defer logLock.Unlock()

	filename = strings.TrimSuffix(filename, ".log") + ".log"
	checkLogDirExist(logDir)

	//打开日志文件
	filename = logDir + "/" + filename //日志文件全路径名称
	logFile, err := os.OpenFile(filename, logWriteType, 0666)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}
	defer logFile.Close()

	//当文件超过日志指定大小，就重命名
	if file, _ := os.Stat(filename); file.Size() > logSize {
		//重命名
		loc, _ := time.LoadLocation(logTimeZone) //time zone
		newFile := strings.TrimSuffix(filename, ".log") + "_" + time.Now().In(loc).Format(tmFmtMissMS) + ".log"
		os.Rename(filename, newFile)
		OldFile, _ := os.Create(filename)
		defer OldFile.Close()
	}

	//将日志信息写入文件中
	var Info *log.Logger
	Info = log.New(io.MultiWriter(os.Stderr, logFile), prefix, logLevel) //第二个参数是日志内容前缀
	Info.Println(message)
}

//debug模式 不记录日志到文件中
func Debug(message interface{}) {
	go func() {
		var debug *log.Logger
		debug = log.New(os.Stdout, "Debug: ", logLevel)
		debug.Println(message)
	}()
	time.Sleep(1 * time.Millisecond) //防止主线程main退出后，日志操作还未执行
}

func Warn(message interface{}, filename string) {
	go Save(message, filename, "warn: ")
	time.Sleep(2 * time.Millisecond)
}

func Info(message interface{}, filename string) {
	go Save(message, filename, "Info: ")
	time.Sleep(2 * time.Millisecond)
}

func Error(message interface{}, filename string) {
	go Save(message, filename, "Error: ")
	time.Sleep(2 * time.Millisecond)
}
