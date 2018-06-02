package lib

// 每日滚动的LOG实现
import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"
    "sync"
    "time"
)

var (
    logDir  = os.TempDir()
    logDay  = 0
    logLock = sync.Mutex{}
    logFile *os.File
)

func init() {
    if proc, err := filepath.Abs(os.Args[0]); err == nil {
        SetLogDir(filepath.Dir(proc))
    }
}

func SetLogDir(dir string) {
    logDir = dir
}

func check() {
    logLock.Lock()
    defer logLock.Unlock()

    now := time.Now().UTC()
    if logDay == now.Day() {
        return
    }

    logDay = now.Day()
    logProc := filepath.Base(os.Args[0])
    filename := filepath.Join(logDir, fmt.Sprintf("%s.%s.log", logProc, now.Format("2006-01-02")))

    newlog, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
    if err != nil {
        logFile = os.Stderr
        fmt.Fprintln(os.Stderr, NumberUTC(), "open log file", err, "use STDOUT")
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
    fmt.Fprintln(logFile, NumberUTC(), offset(), "debug", v)
}

func InfoLog(v ...interface{}) {
    check()
    logLock.Lock()
    defer logLock.Unlock()
    fmt.Fprintln(logFile, NumberUTC(), offset(), "info", v)
}

func WarningLog(v ...interface{}) {
    check()
    logLock.Lock()
    defer logLock.Unlock()
    fmt.Fprintln(logFile, NumberUTC(), offset(), "warning", v)
}

func ErrorLog(v ...interface{}) {
    check()
    logLock.Lock()
    defer logLock.Unlock()
    fmt.Fprintln(logFile, NumberUTC(), offset(), "error", v)
}

func CustomLog(v ...interface{}) {
    check()
    logLock.Lock()
    defer logLock.Unlock()
    fmt.Fprintln(logFile, NumberUTC(), v)
}
