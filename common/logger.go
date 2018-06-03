package common

// 每日滚动的LOG实现
// 日志生成时间默认采用PRC时间,如需要更改,请调用SetLogTimeZone
import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var (
	logDir      = os.TempDir()
	logDay      = 0
	logLock     = &sync.Mutex{}
	logFile     *os.File
	logTimeZone = "PRC" //time zone default prc
)

func init() {
	if proc, err := filepath.Abs(os.Args[0]); err == nil {
		SetLogDir(filepath.Dir(proc))
	}
}

//设置日志记录时区
func SetLogTimeZone(zone string) {
	logTimeZone = zone
}

func SetLogDir(dir string) {
	logDir = dir
	if _, err := os.Stat(logDir); err != nil { //日志目录不存在
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatalf("创建目录%s失败,错误信息: %s", logDir, err)
			return
		}
	}
}

func check() {
	logLock.Lock()
	defer logLock.Unlock()

	loc, _ := time.LoadLocation(logTimeZone)
	now := time.Now().In(loc)
	if logDay == now.Day() {
		return
	}

	logDay = now.Day()
	logProc := filepath.Base(os.Args[0])
	filename := filepath.Join(logDir, fmt.Sprintf("%s.%s.log", logProc, now.Format(tmFmtTime)))

	newlog, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logFile = os.Stderr
		fmt.Fprintln(os.Stderr, now.Format(tmFmtMissMS), "open log file", err, "use STDOUT")
	} else {
		logFile.Sync()
		logFile.Close()
		logFile = newlog
	}
}

func fileline(file string, line int) string {
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

func offset() string {
	_, file, line, _ := runtime.Caller(2)
	return fileline(file, line)
}

func DropLog(v ...interface{}) {}

func DebugLog(v ...interface{}) {
	check()
	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, GetTimeByTimeZone(logTimeZone), offset(), "debug", v)
}

func InfoLog(v ...interface{}) {
	check()
	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, GetTimeByTimeZone(logTimeZone), offset(), "info", v)
}

func WarningLog(v ...interface{}) {
	check()
	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, GetTimeByTimeZone(logTimeZone), offset(), "warning", v)
}

func ErrorLog(v ...interface{}) {
	check()
	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, GetTimeByTimeZone(logTimeZone), offset(), "error", v)
}

func CustomLog(v ...interface{}) {
	check()
	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, GetTimeByTimeZone(logTimeZone), v)
}
