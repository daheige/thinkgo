package common

import (
	"bytes"
	"flag"
	"fmt"
	"html"
	"io"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/daheige/thinkgo/crypto"
)

var localTimeZone = "Local" //time zone default prc

// os_Chown is a var so we can mock it out during tests.
var os_Chown = os.Chown

// zero size, empty struct
type EmptyStruct struct{}

// parse flag and print usage/value to writer
func Init(writer io.Writer) {
	flag.Parse()

	if writer != nil {
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(writer, "-%s=%v \n", f.Name, f.Value)
		})
	}
}

// check panic when exit
func CheckPanic() {
	if err := recover(); err != nil {
		loc, _ := time.LoadLocation(localTimeZone)
		fmt.Fprintf(os.Stderr, "\n%s %+v\n", FormatTime19(time.Now().In(loc)), err)
		fmt.Fprintf(os.Stderr, "full stack info: \n%s", CatchStack())
	}

}

//捕获指定stack信息,一般在处理panic/recover中处理
//返回完整的堆栈信息和函数调用信息
func CatchStack() []byte {
	buf := &bytes.Buffer{}

	//完整的堆栈信息
	stack := Stack()
	buf.WriteString("full stack:\n")
	buf.WriteString(string(stack))

	//完整的函数调用信息
	buf.WriteString("full fn call info:\n")

	// skip为0时，打印当前调用文件及行数。
	// 为1时，打印上级调用的文件及行数，依次类推
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc).Name()
		buf.WriteString(fmt.Sprintf("error Stack file: %s:%d call func:%s\n", file, line, fn))
	}

	return buf.Bytes()
}

//获取完整的堆栈信息
// Stack returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
func Stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false) //当第二个参数为true，一次获取所有的堆栈信息
		if n < len(buf) {
			return buf[:n]
		}

		buf = make([]byte, 2*len(buf))
	}
}

// reload signal
func HupSignal() <-chan os.Signal {
	signals := make(chan os.Signal, 3)
	signal.Notify(signals, syscall.SIGHUP)
	return signals
}

// recive quit signal
func QuitSignal() <-chan os.Signal {
	signals := make(chan os.Signal, 3)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	return signals
}

func Md5(str string) string {
	return crypto.Md5(str)
}

func Sha1(str string) string {
	return crypto.Sha1(str)
}

//通过随机数的方式生成uuid
//如果rand.Read失败,就按照当前时间戳+随机数进行md5方式生成
//该方式生成的uuid有可能存在重复值
//返回格式:7999b726-ca3c-42b6-bda2-259f4ac0879a
func NewUUID() string {
	u := [16]byte{}
	ns := time.Now().UnixNano()
	rand.Seed(ns)
	_, err := rand.Read(u[0:])
	if err != nil {
		rndStr := strconv.FormatInt(ns, 10) + strconv.FormatInt(RandInt64(1000, 9999), 10)
		s := crypto.Md5(rndStr)
		return fmt.Sprintf("%s-%s-%s-%s-%s", s[:8], s[8:12], s[12:16], s[16:20], s[20:])
	}

	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

//基于时间ns和随机数实现唯一的uuid
//在单台机器上是不会出现重复的uuid
//如果要在分布式的架构上生成不重复的uuid
// 只需要在rndStr的前面加一些自定义的字符串
//返回格式:eba1e8cd-0460-4910-49c6-44bdf3cf024d
func RndUuid() string {
	s := RndUuidMd5()
	return fmt.Sprintf("%s-%s-%s-%s-%s", s[:8], s[8:12], s[12:16], s[16:20], s[20:])
}

func RndUuidMd5() string {
	ns := time.Now().UnixNano()
	rndStr := StrJoin("", strconv.FormatInt(ns, 10), strconv.FormatInt(RandInt64(1000, 9999), 10))
	return crypto.Md5(rndStr)
}

//获取文件的名称不带后缀
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

//获取文件名:行数
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

//运行shell脚本
func RunShell(exeStr string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", exeStr)
	bytes, err := cmd.CombinedOutput()
	return string(bytes), err
}

//防止xss攻击
func Xss(str string) string {
	if len(str) == 0 {
		return ""
	}

	return html.EscapeString(str)
}

func XssUnescape(str string) string {
	return html.UnescapeString(str)
}

//对浮点数进行四舍五入操作比如 12.125保留2位小数==>12.13
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

//生成m-n之间的随机数
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}

	//随机种子
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

// 根据kind生成不同风格的指定区间随机字符串
// 纯数字kind=0,小写字母kind=1
// 大写字母kind=2,数字+大小写字母kind=3
func Krand(size int, kind int) string {
	ikind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano()) //随机种子
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}

		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}

	return string(result)
}

// int-->string
func IntToStr(n int) string {
	return strconv.Itoa(n)
}

// string-->int
func StrToInt(s string) int {
	if i, err := strconv.Atoi(s); err != nil {
		return 0
	} else {
		return i
	}
}

// int64-->string
func Int64ToStr(i64 int64) string {
	return strconv.FormatInt(i64, 10)
}

// string--> int64
func StrToInt64(s string) int64 {
	if i64, err := strconv.ParseInt(s, 10, 64); err != nil {
		return 0
	} else {
		return i64
	}
}

//多个字符串按照指定的分隔符进行拼接
func StrJoin(sep string, str ...string) string {
	return strings.Join(str, sep)
}

//通过buf缓冲区的方式连接字符串
func StrJoinByBuf(str ...string) string {
	if len(str) == 0 {
		return ""
	}

	var buf bytes.Buffer
	for _, s := range str {
		buf.WriteString(s)
	}

	return buf.String()
}

func Chown(name string, info os.FileInfo) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}

	f.Close()
	stat := info.Sys().(*syscall.Stat_t)

	return os_Chown(name, int(stat.Uid), int(stat.Gid))
}