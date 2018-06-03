package common

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

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
		fmt.Fprintf(os.Stderr, "\n%v %v\n", NumberNow(), err)

		for skip := 1; ; skip++ {
			if pc, file, line, ok := runtime.Caller(skip); ok {
				fn := runtime.FuncForPC(pc).Name()
				fmt.Fprintln(os.Stderr, NumberNow(), fn, fileline(file, line))
			} else {
				break
			}
		}
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

// create a uuid string
func NewUUID() string {
	u := [16]byte{}
	rand.Read(u[:])
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}
