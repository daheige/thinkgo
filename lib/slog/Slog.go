/**
*file 日志操作库
 */
package slog

import (
    "io"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    kb int64 = 1 << (10 * iota) //字节大小
    mb
)

var (
    logLevel     = log.LstdFlags                           //标准日志格式
    logWriteType = os.O_CREATE | os.O_WRONLY | os.O_APPEND //文件创建标识
    logSize      = 2 * kb                                  //日志文件大小，默认2mb
)

//设置日志大小，默认2mb
func InitLogSize(size int64) {
    logSize = size * mb
}

func Save(message interface{}, filename string, prefix string) {
    filename = strings.TrimSuffix(filename, ".log") + ".log"
    path := filepath.Dir(filename)           //得到日志目录
    if _, err := os.Stat(path); err != nil { //日志目录不存在
        if err := os.Mkdir(path, 0755); err != nil {
            log.Fatalf("创建目录%s失败,错误信息: %s", path, err)
        }
    }

    //打开日志文件
    errFile, err := os.OpenFile(filename, logWriteType, 0666)
    if err != nil {
        log.Fatalln("打开日志文件失败：", err)
    }
    defer errFile.Close()

    //当文件超过日志指定大小，就重命名
    if file, _ := os.Stat(filename); file.Size() > logSize {
        //重命名
        newFile := strings.TrimSuffix(filename, ".log") + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".log"
        os.Rename(filename, newFile)
        OldFile, _ := os.Create(filename)
        defer OldFile.Close()
    }

    //将日志信息写入文件中
    var Info *log.Logger
    Info = log.New(io.MultiWriter(os.Stderr, errFile), prefix, logLevel) //第二个参数是日志内容前缀
    Info.Println(message)
}

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
