package logger

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger句柄，支持zap logger上的Debug,Info,Error,Panic,Warn,Fatal等方法
var fLogger *zap.Logger

var core zapcore.Core

//日志级别定义，从低到高
var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

var (
	logMaxAge   = 7         //默认日志保留天数
	logMaxSize  = 512       //默认日志大小，单位为Mb
	logCompress = false     //默认日志不压缩
	logLevel    = "debug"   //最低日志级别
	logFileName = "zap.log" //默认日志文件，不包含全路径
	logDir      = ""        //日志文件存放目录
)

func MaxAge(n int) {
	logMaxAge = n
}

func MaxSize(n int) {
	logMaxSize = n
}

func LogCompress(b bool) {
	logCompress = b
}

func LogLevel(lvl string) {
	logLevel = lvl
}

//设置日志文件路径，如果日志文件不存在zap会自动创建文件
func SetLogFile(fileName string) {
	logFileName = fileName
}

//获得日志级别
func getLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}

	return zapcore.InfoLevel
}

//check file or path exist
func checkPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

//日志存放目录
func SetLogDir(dir string) {
	if dir == "" {
		logDir = os.TempDir()
	} else {
		if !checkPathExist(dir) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Println(err)
				return
			}
		}

		logDir = dir
	}
}

func initCore() {
	if logDir == "" {
		logFileName = filepath.Join(os.TempDir(), logFileName) //默认日志文件名称
	} else {
		logFileName = filepath.Join(logDir, logFileName)
	}

	//日志最低级别设置
	level := getLevel(logLevel)
	syncWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:  logFileName, //⽇志⽂件路径
		MaxSize:   logMaxSize,  //单位为MB,默认为512MB
		MaxAge:    logMaxAge,   // 文件最多保存多少天
		LocalTime: true,        //采用本地时间
		Compress:  logCompress, //是否压缩日志
	})

	encoderConf := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConf), syncWriter, zap.NewAtomicLevelAt(level))
}

func InitLogger() {
	if core == nil {
		initCore()
	}

	fLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

//sugar语法糖，支持简单的msg信息打印
//支持Debug,Info,Error,Panic,Warn,Fatal等方法
func LogSugar() *zap.SugaredLogger {
	if core == nil {
		initCore()
	}

	//基于zapcore创建sugar
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
	return logger.Sugar()
}

//debug日志直接输出到终端
func Debug(msg string, options map[string]interface{}) {
	log.Println("msg: ", msg)
	log.Println("log fields: ", options)
}

func Info(msg string, options map[string]interface{}) {
	fields := parseFields(options)
	if len(fields) == 0 {
		fLogger.Info(msg)
		return
	}

	fLogger.Info(msg, fields...)
}

func Warn(msg string, options map[string]interface{}) {
	fields := parseFields(options)
	if len(fields) == 0 {
		fLogger.Warn(msg)
		return
	}

	fLogger.Warn(msg, fields...)
}

func Error(msg string, options map[string]interface{}) {
	fields := parseFields(options)
	if len(fields) == 0 {
		fLogger.Error(msg)
		return
	}

	fLogger.Error(msg, fields...)
}

//调试模式下的panic，程序不退出，继续运行
func DPanic(msg string, options map[string]interface{}) {
	fields := parseFields(options)
	if len(fields) == 0 {
		fLogger.DPanic(msg)
		return
	}

	fLogger.DPanic(msg, fields...)
}

//下面的panic,fatal一般不建议使用，除非不可恢复的panic或致命错误程序必须退出
//抛出panic的时候，先记录日志，然后执行panic,退出当前goroutine
func Panic(msg string, options map[string]interface{}) {
	fields := parseFields(options)
	if len(fields) == 0 {
		fLogger.Panic(msg)
		return
	}

	fLogger.Panic(msg, fields...)
}

//抛出致命错误，然后退出程序
func Fatal(msg string, options map[string]interface{}) {
	fields := parseFields(options)
	if len(fields) == 0 {
		fLogger.Fatal(msg)
		return
	}

	fLogger.Fatal(msg, fields...)
}

// parseFields 解析map[string]interface{}中的字段到zap.Field
func parseFields(fields map[string]interface{}) []zap.Field {
	fLen := len(fields)
	f := make([]zap.Field, 0, fLen+2)

	//当前函数调用的位置和行数
	if _, ok := fields["trace_line"]; !ok {
		_, file, line, _ := runtime.Caller(2)
		f = append(f, zap.String("trace_file", file))
		f = append(f, zap.Int("trace_line", line))
	}

	for k, _ := range fields {
		f = append(f, zap.Any(k, fields[k]))
	}

	return f
}
