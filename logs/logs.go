package logs

import (
	"fmt"
	"log"
	"os"
)

var stdErr = log.New(os.Stderr, "", log.LstdFlags)
var stdOut = log.New(os.Stdout, "", log.LstdFlags)

func Error(args ...interface{}) {
	_ = stdErr.Output(2, "\033[31mERROR\033[0m "+fmt.Sprint(args...))
}

func Errorf(format string, args ...interface{}) {
	_ = stdErr.Output(2, "\033[31mERROR\033[0m "+fmt.Sprintf(format, args...))
}

func Info(args ...interface{}) {
	_ = stdOut.Output(2, "\033[32mINFO\033[0m "+fmt.Sprint(args...))
}

func Infof(format string, args ...interface{}) {
	_ = stdOut.Output(2, "\033[32mINFO\033[0m "+fmt.Sprintf(format, args...))
}

func Warn(args ...interface{}) {
	_ = stdOut.Output(2, "\033[33mWARN\033[0m "+fmt.Sprint(args...))
}

func Fatal(args ...interface{}) {
	_ = stdErr.Output(2, fmt.Sprint(args...))
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	_ = stdErr.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func Debug(args ...interface{}) {
	_ = stdOut.Output(2, "DEBUG "+fmt.Sprint(args...))
}
