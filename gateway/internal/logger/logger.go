package logger

import (
	"fmt"
	"os"
)

type Logger struct {
	prefix string
}

func New(prefix string) *Logger {
	var l Logger
	l.prefix = prefix
	return &l
}

func (l *Logger) formatMessage(reqId string, msg string) string {
	return fmt.Sprintf("%s: [%s] %s", l.prefix, reqId, msg)
}

func (l *Logger) Info(reqId string, msg string) {
	fmt.Println(l.formatMessage(reqId, msg))
}

func (l *Logger) Error(reqId string, msg string) {
	fmt.Fprintln(os.Stderr, l.formatMessage(reqId, msg))
}
