package logger

import (
	"os"

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
	logMaxAge   = 7       //默认日志保留天数
	logMaxSize  = 512     //默认日志大小，单位为Mb
	logCompress = false   //默认日志不压缩
	logLevel    = "debug" //最低日志级别
	logFileName = ""
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

func initCore() {
	if logFileName == "" {
		logFileName = os.TempDir() + "/zap.log" //默认日志文件名称
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

func Debug(msg string, fields ...zap.Field) {
	fLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	fLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	fLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	fLogger.Error(msg, fields...)
}

//调试模式下的panic，程序不退出，继续运行
func DPanic(msg string, fields ...zap.Field) {
	fLogger.DPanic(msg, fields...)
}

//下面的panic,fatal一般不建议使用，除非不可恢复的panic或致命错误程序必须退出
//抛出panic的时候，先记录日志，然后执行panic,退出当前goroutine
func Panic(msg string, fields ...zap.Field) {
	fLogger.Panic(msg, fields...)
}

//抛出致命错误，然后退出程序
func Fatal(msg string, fields ...zap.Field) {
	fLogger.Fatal(msg, fields...)
}
