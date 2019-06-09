/**
 * 每天流动式日志实现
 * 操作日志记录到文件，支持info,error,debug,notice,alert等
 * 写日志文件的时候，采用乐观锁方式对文件句柄进行加锁
 * 等级参考php Monolog/logger.php
 * 日志切割机制参考lumberjack包实现
 * json encode采用jsoniter库快速json encode处理
 */
package common

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/daheige/thinkgo/internal/json"
)

/* 日志级别 从上到下，由高到低 */
const (
	EMERGENCY = "emerg"  // 严重错误: 导致系统崩溃无法使用
	ALTER     = "alter"  // 警戒性错误: 必须被立即修改的错误
	CRIT      = "crit"   // 临界值错误: 超过临界值的错误，例如一天24小时，而输入的是25小时这样
	ERR       = "error"  // 一般错误:比如类型错误，数据库连接不上等等
	WARN      = "warn"   // 警告性错误: 需要发出警告的错误
	NOTICE    = "notice" // 通知: 程序可以运行但是还不够完美的错误
	INFO      = "info"   // 信息: 程序输出信息
	DEBUG     = "debug"  // 调试: 调试信息
)

var LogLevelMap = map[string]int{
	EMERGENCY: 600,
	ALTER:     550,
	CRIT:      500,
	ERR:       400,
	WARN:      300,
	NOTICE:    250,
	INFO:      200,
	DEBUG:     100,
}

var (
	logDir                = ""             //日志文件存放目录
	logFile               = ""             //日志文件
	logLock               = NewMutexLock() //采用sync实现加锁，效率比chan实现的加锁效率高一点
	logDay                = 0              //当前日期
	logTimeZone           = "Local"        //time zone default Local "Asia/Shanghai"
	logtmFmtWithMS        = "2006-01-02 15:04:05.999"
	logtmFmtMissMS        = "2006-01-02 15:04:05"
	logtmFmtTime          = "2006-01-02"
	defaultLogLevel       = INFO //默认为INFO级别
	logtmLoc              = &time.Location{}
	megabyte        int64 = 1024 * 1024               //1MB = 1024 * 1024byte
	defaultMaxSize  int64 = 512                       //默认单个日志文件大小,单位为mb
	currentTime           = time.Now                  //当前时间函数
	logSplit              = false                     //默认日志不分割,当设置为true可以加快写入的速度
	logtmSplit            = "2006-01-02-15-04-05.999" //日志备份文件名时间格式
)

//日志内容结构体
type logContent struct {
	Level     int                    `json:"level"`
	LevelName string                 `json:"level_name"`
	LocalTime string                 `json:"local_time"` //当前时间
	Msg       interface{}            `json:"msg"`
	LineNo    int                    `json:"line_no"`   //当前行号
	FilePath  string                 `json:"file_path"` //当前文件
	Context   map[string]interface{} `json:"context"`
}

//设置日志记录时区
func SetLogTimeZone(timezone string) {
	logTimeZone = timezone
}

//日志分割设置
func LogSplit(b bool) {
	logSplit = b
}

//日志存放目录
func SetLogDir(dir string) {
	if dir == "" {
		logDir = os.TempDir()
	} else {
		if !CheckPathExist(dir) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Println(err)
				return
			}
		}

		logDir = dir
	}

	logtmLoc, _ = time.LoadLocation(logTimeZone)
	now := currentTime().In(logtmLoc)
	newFile(now) //建立日志文件
}

func LogSize(n int64) {
	defaultMaxSize = n
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
		log.Println("open log file", filename, err, "use stdout")
		return
	}

	fp.Close()
	logFile = filename
}

//判断当天的日志文件是否存在，不存在就创建
func checkLogExist() {
	if now := currentTime().In(logtmLoc); logDay != now.Day() {
		newFile(now)
	}
}

// backupName creates a new filename from the given name, inserting a timestamp
// between the filename and the extension
// 创建备份文件名称
func backupName(name string) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]

	timestamp := currentTime().Format(logtmSplit)
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, timestamp, ext))
}

//当日志文件超过了指定大小，对其进行分割处理
func splitLog() {
	if logFile == "" {
		return
	}

	//检测文件大小是否超过指定大小
	fileInfo, err := os.Stat(logFile)
	if err != nil {
		log.Println("get file stat error: ", err)
		return
	}

	if fileInfo.Size() >= defaultMaxSize*megabyte {
		newName := backupName(logFile)
		if err := os.Rename(logFile, newName); err != nil {
			log.Printf("can't rename log file: %s\n", err)
			return
		}

		// this is a no-op anywhere but linux
		if err := Chown(logFile, fileInfo); err != nil {
			log.Printf("can't chown log file: %s\n", err)
			return
		}
	}

}

func writeLog(levelName string, msg interface{}, options map[string]interface{}) {
	if _, ok := LogLevelMap[levelName]; !ok {
		levelName = defaultLogLevel
	}

	//对日志内容转换为bytes
	var strBytes []byte
	_, file, line, _ := runtime.Caller(2)

	c := &logContent{
		LevelName: levelName,
		Level:     LogLevelMap[levelName],
		LocalTime: currentTime().In(logtmLoc).Format(logtmFmtWithMS),
		Msg:       msg,
		LineNo:    line,
		FilePath:  file,
	}

	if len(options) > 0 {
		c.Context = options
	}

	//序列化为json格式
	strBytes, err := json.Marshal(c)
	if err != nil {
		log.Println("write log file,use stdout")
		log.Println("log content: ", string(strBytes))
		return
	}

	//追加换行符
	strBytes = append(strBytes, []byte("\n")...)

	//开始写日志，这里需要对文件句柄进行加锁
	logLock.Lock()
	defer logLock.Unlock()

	//检测日志文件是否存在
	checkLogExist()
	if logFile == "" {
		log.Println("write log file,use stdout")
		log.Println("log content:", string(strBytes))
		return
	}

	//检测日志是否需要分割
	if logSplit {
		splitLog()
	}

	fp, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("open log file: %s error: %s\n", logFile, err)
		log.Println("log content:", string(strBytes))
		return
	}

	defer fp.Close()

	if _, err := fp.Write(strBytes); err != nil {
		log.Printf("write log file: %s error: %s\n", logFile, err)
		return
	}
}

func DebugLog(v interface{}, options map[string]interface{}) {
	writeLog("debug", v, options)
}

func InfoLog(v interface{}, options map[string]interface{}) {
	writeLog("info", v, options)
}

func NoticeLog(v interface{}, options map[string]interface{}) {
	writeLog("notice", v, options)
}

func WarnLog(v interface{}, options map[string]interface{}) {
	writeLog("warn", v, options)
}

func ErrorLog(v interface{}, options map[string]interface{}) {
	writeLog("error", v, options)
}

func CritLog(v interface{}, options map[string]interface{}) {
	writeLog("crit", v, options)
}

func AlterLog(v interface{}, options map[string]interface{}) {
	writeLog("alter", v, options)
}

func EmergLog(v interface{}, options map[string]interface{}) {
	writeLog("emerg", v, options)
}

// RecoverLog 异常捕获处理，对于异常或者panic进行捕获处理
// 记录到日志中，方便定位问题
func RecoverLog() {
	defer func() {
		if err := recover(); err != nil {
			EmergLog("exec panic", map[string]interface{}{
				"error":       err,
				"error_trace": string(CatchStack()),
			})
		}
	}()
}
