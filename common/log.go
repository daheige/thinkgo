/**
 * 每天流动式日志实现
 * 操作日志记录到文件，支持info,error,debug,notice,alert等
 * 写日志文件的时候，采用chan实现的乐观锁方式对文件句柄进行加锁
 * 等级参考php Monolog/logger.php
 */
package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

/* 日志级别 从上到下，由高到低 */
const (
	EMERGENCY = "emerg"  // 严重错误: 导致系统崩溃无法使用
	ALERT     = "alter"  // 警戒性错误: 必须被立即修改的错误
	CRIT      = "crit"   // 临界值错误: 超过临界值的错误，例如一天24小时，而输入的是25小时这样
	ERR       = "error"  // 一般错误: 一般性错误
	WARN      = "warn"   // 警告性错误: 需要发出警告的错误
	NOTICE    = "notice" // 通知: 程序可以运行但是还不够完美的错误
	INFO      = "info"   // 信息: 程序输出信息
	DEBUG     = "debug"  // 调试: 调试信息
)

var LogLevelMap = map[string]int{
	EMERGENCY: 600,
	ALERT:     550,
	CRIT:      500,
	ERR:       400,
	WARN:      300,
	NOTICE:    250,
	INFO:      200,
	DEBUG:     100,
}

var (
	logDir  = ""             //日志文件存放目录
	logFile = ""             //日志文件
	logLock = NewMutexLock() //采用sync实现加锁，效率比chan实现的加锁效率高一点
	//logLock         = NewChanLock()               //采用chan实现的乐观锁方式，实现加锁，效率稍微低一点
	logTicker       = time.NewTicker(time.Second) //time一次性触发器
	logDay          = 0                           //当前日期
	logTime         = true                        //默认显示时间和行号，参考 SetLogTime 方法
	logTimeZone     = "Local"                     //time zone default Local/Shanghai
	logtmFmtWithMS  = "2006-01-02 15:04:05.999"
	logtmFmtMissMS  = "2006-01-02 15:04:05"
	logtmFmtTime    = "2006-01-02"
	defaultLogLevel = INFO //默认为INFO级别
	logtmLoc        = &time.Location{}
)

//设置日志记录时区
func SetLogTimeZone(timezone string) {
	logTimeZone = timezone
}

//日志存放目录
func SetLogDir(dir string) {
	if dir == "" {
		logDir = os.TempDir()
	} else {
		logDir = dir
	}

	logtmLoc, _ = time.LoadLocation(logTimeZone)
	now := time.Now().In(logtmLoc)

	//建立日志文件
	newFile(now)
}

//设置是否显示时间和行号
//当logtime=false自定义日志格式
func SetLogTime(logtime bool) {
	logTime = logtime
}

//创建日志文件
func newFile(now time.Time) {
	if len(logDir) == 0 {
		return
	}

	logDay = now.Day()
	filename := filepath.Join(logDir, fmt.Sprintf("%s.%s.log", filepath.Base(os.Args[0]), now.Format(logtmFmtTime)))

	//创建文件
	fp, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(now.Format(logtmFmtMissMS), "open log file", filename, err, "use STDOUT")
	} else {
		fp.Close()
		logFile = filename
	}
}

func checkLogExist() {
	select {
	case <-logTicker.C:
	default:
		return
	}

	if !logLock.TryLock() {
		return
	}

	defer logLock.Unlock()

	loc, _ := time.LoadLocation(logTimeZone)
	now := time.Now().In(loc)
	//判断当天的日志文件是否存在，不存在就创建
	if logDay != now.Day() {
		newFile(now)
	}
}

func writeLog(levelName string, message interface{}) {
	if _, ok := LogLevelMap[levelName]; !ok {
		levelName = defaultLogLevel
	}

	//检测日志文件是否存在
	checkLogExist()

	//对日志内容转换为bytes
	var strBytes []byte
	if logTime {
		_, file, line, _ := runtime.Caller(2)
		now := time.Now().In(logtmLoc)
		strBytes = []byte(fmt.Sprintf("%s %s %s line:[%d]: %v", now.Format(logtmFmtWithMS), levelName, filepath.Base(file), line, message))
	} else {
		if v, ok := message.(string); ok {
			strBytes = []byte(v)
		} else {
			strBytes = []byte(fmt.Sprintf("%v", message))
		}
	}

	//追加换行符
	strBytes = append(strBytes, []byte("\n")...)
	if logFile == "" {
		now := time.Now().In(logtmLoc)
		fmt.Println(now.Format(logtmFmtMissMS), "open log file", "use stdout")
		return
	}

	//开始写日志，这里需要对文件句柄进行加锁
	logLock.Lock()
	defer logLock.Unlock()

	fp, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		now := time.Now().In(logtmLoc)
		fmt.Println(now.Format(logtmFmtMissMS), "open log file", logFile, err, "use stdout")
		return
	}

	defer fp.Close()

	if _, err := fp.Write(strBytes); err != nil {
		fmt.Println("write log to buf error: ", err)
		fmt.Println("current log: ", string(strBytes))
		return
	}

}

func DebugLog(v interface{}) {
	writeLog("debug", v)
}

func InfoLog(v interface{}) {
	writeLog("info", v)
}

func NoticeLog(v interface{}) {
	writeLog("notice", v)
}

func WarnLog(v interface{}) {
	writeLog("warn", v)
}

func ErrorLog(v interface{}) {
	writeLog("error", v)
}

func CritLog(v interface{}) {
	writeLog("crit", v)
}

func AlterLog(v interface{}) {
	writeLog("alter", v)
}

func EmergLog(v interface{}) {
	writeLog("emerg", v)
}
