package log

import (
	"fmt"
	"os"
)

// TODO: hook this into a *real* logging library (with colour and NOCOLOR)

const (
	Black   = "\033[1;30m%s\033[0m"
	Red     = "\033[1;31m%s\033[0m"
	Green   = "\033[1;32m%s\033[0m"
	Yellow  = "\033[1;33m%s\033[0m"
	Purple  = "\033[1;34m%s\033[0m"
	Magenta = "\033[1;35m%s\033[0m"
	Teal    = "\033[1;36m%s\033[0m"
	White   = "\033[1;37m%s\033[0m"
)

func Debug(format string, a ...any) {
	line := fmt.Sprintf(format, a...)
	fmt.Printf(Teal, "DEBUG: "+line+"\n")
}

func Info(format string, a ...any) {
	line := fmt.Sprintf(format, a...)
	fmt.Printf("INFO: %s\n", line)
}

func Warning(format string, a ...any) {
	line := fmt.Sprintf(format, a...)
	fmt.Printf(Yellow, "WARNING: "+line+"\n")
}

func Error(format string, a ...any) {
	line := fmt.Sprintf(format, a...)
	fmt.Printf(Red, "ERROR: "+line+"\n")
}

func Critical(format string, a ...any) {
	line := fmt.Sprintf(format, a...)
	fmt.Printf(Purple, "CRITICAL: "+line+"\n")
}

func Fatal(exitImmediately bool, format string, a ...any) {
	line := fmt.Sprintf(format, a...)
	fmt.Printf("FATAL: %s\n", line)
	if exitImmediately {
		os.Exit(1)
	}
}
