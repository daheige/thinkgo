package common

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/daheige/thinkgo/crypto"
)

// os_Chown is a var so we can mock it out during tests.
var os_Chown = os.Chown

// EmptyStruct zero size, empty struct
type EmptyStruct struct{}

// EmptyArray 兼容其他语言的[]空数组,一般在tojson的时候转换为[]
type EmptyArray []struct{}

// H 对map[string]interface{}别名类型，简化书写
type H map[string]interface{}

// dotask返回的结果
type TaskRes struct {
	Err      error
	Result   chan interface{}
	CostTime float64
}

// CheckPanic check panic when exit
func CheckPanic() {
	if err := recover(); err != nil {
		loc, _ := time.LoadLocation("Local")
		fmt.Fprintf(os.Stderr, "\n%s %+v\n", FormatTime19(time.Now().In(loc)), err)
		fmt.Fprintf(os.Stderr, "full stack info: \n%s", CatchStack())
	}
}

// CatchStack 捕获指定stack信息,一般在处理panic/recover中处理
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

// Stack 获取完整的堆栈信息
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

func Md5(str string) string {
	return crypto.Md5(str)
}

func Sha1(str string) string {
	return crypto.Sha1(str)
}

// NewUUID 通过随机数的方式生成uuid
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

// RndUuid 基于时间ns和随机数实现唯一的uuid
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

// DoTask 在独立携程中运行fn
// 这里返回结果设计为interface{},因为有时候返回结果可以是error
func DoTask(fn func() interface{}) *TaskRes {
	t := time.Now()
	done := make(chan struct{}, 1)
	res := &TaskRes{
		Result: make(chan interface{}, 1),
	}

	go func() {
		defer func() {
			close(res.Result)
			close(done)
			if err := recover(); err != nil {
				log.Println("task exec panic error: ", err)
				res.Err = errors.New(fmt.Sprintf("%v", err))
			}
		}()

		r := fn()
		res.Result <- r
	}()

	<-done

	res.CostTime = time.Now().Sub(t).Seconds()
	return res
}

// DoTaskWithArgs 在独立携程中执行有参数的fn
func DoTaskWithArgs(fn func(args ...interface{}) interface{}, args ...interface{}) *TaskRes {
	t := time.Now()
	done := make(chan struct{}, 1)
	res := &TaskRes{
		Result: make(chan interface{}, 1),
	}

	go func() {
		defer func() {
			close(res.Result)
			close(done)
			if err := recover(); err != nil {
				log.Println("task exec panic error: ", err)
				res.Err = errors.New(fmt.Sprintf("%v", err))
			}
		}()

		r := fn(args...)
		res.Result <- r
	}()

	<-done

	res.CostTime = time.Now().Sub(t).Seconds()
	return res
}

// DoTaskWithTimeout 采用done+select+time.After实现goroutine超时调用
func DoTaskWithTimeout(fn func() interface{}, timeout time.Duration) *TaskRes {
	t := time.Now()
	done := make(chan struct{}, 1)
	res := &TaskRes{
		Result: make(chan interface{}, 1),
	}

	go func() {
		defer func() {
			close(res.Result)
			close(done)
			if err := recover(); err != nil {
				res.Err = errors.New(fmt.Sprintf("%v", err))
			}
		}()

		r := fn()
		res.Result <- r
	}()

	select {
	case <-done:
		log.Println("task has done")
	case <-time.After(timeout):
		if res.Err == nil { //当执行过程中没有发生了panic的话，这里设置为任务超时错误
			res.Err = errors.New("task timeout")
		}
	}

	res.CostTime = time.Now().Sub(t).Seconds()
	return res
}

// DoTaskWithContext 通过上下文context+done+select实现goroutine超时调用
func DoTaskWithContext(ctx context.Context, fn func() interface{}, timeout time.Duration) *TaskRes {
	t := time.Now()
	done := make(chan struct{}, 1)
	res := &TaskRes{
		Result: make(chan interface{}, 1),
	}

	ctx2, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	go func() {
		defer func() {
			close(res.Result)
			close(done)
			if err := recover(); err != nil {
				res.Err = errors.New(fmt.Sprintf("%v", err))
			}
		}()

		r := fn()
		res.Result <- r
	}()

	select {
	case <-done:
		log.Println("task has done")
	case <-ctx2.Done(): //超时了
		if res.Err == nil {
			res.Err = errors.New("task timeout")
		}
	}

	res.CostTime = time.Now().Sub(t).Seconds()
	return res
}

// DoTaskWithTimeoutArgs 采用done+select+time.After实现goroutine超时调用
func DoTaskWithTimeoutArgs(fn func(args ...interface{}) interface{}, timeout time.Duration, args ...interface{}) *TaskRes {
	t := time.Now()
	done := make(chan struct{}, 1)
	res := &TaskRes{
		Result: make(chan interface{}, 1),
	}

	go func() {
		defer func() {
			close(res.Result)
			close(done)
			if err := recover(); err != nil {
				res.Err = errors.New(fmt.Sprintf("%v", err))
			}
		}()

		r := fn(args...)
		res.Result <- r
	}()

	select {
	case <-done:
		log.Println("task has done")
	case <-time.After(timeout):
		if res.Err == nil {
			res.Err = errors.New("task timeout")
		}
	}

	res.CostTime = time.Now().Sub(t).Seconds()
	return res
}

// DoTaskWithContextArgs 通过上下文context+done+select实现goroutine超时调用
func DoTaskWithContextArgs(ctx context.Context, fn func(args ...interface{}) interface{}, timeout time.Duration, args ...interface{}) *TaskRes {
	t := time.Now()
	done := make(chan struct{}, 1)
	res := &TaskRes{
		Result: make(chan interface{}, 1),
	}

	ctx2, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	go func() {
		defer func() {
			close(res.Result)
			close(done)
			if err := recover(); err != nil {
				res.Err = errors.New(fmt.Sprintf("%v", err))
			}
		}()

		r := fn(args...)
		res.Result <- r
	}()

	select {
	case <-done:
		log.Println("task has done")
	case <-ctx2.Done(): //超时了
		if res.Err == nil {
			res.Err = ctx2.Err()
		}
	}

	res.CostTime = time.Now().Sub(t).Seconds()
	return res
}
