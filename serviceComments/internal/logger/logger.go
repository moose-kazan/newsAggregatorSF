package logger

import (
	"fmt"
	"net/http"
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

func (l *Logger) formatMessage(req *http.Request, msg string) string {
	return fmt.Sprintf("%s: [%s] [%s] %s", l.prefix, req.RemoteAddr, req.Header.Get("X-Request-Id"), msg)
}

func (l *Logger) Info(req *http.Request, msg string) {
	fmt.Println(l.formatMessage(req, msg))
}

func (l *Logger) Error(req *http.Request, msg string) {
	fmt.Fprintln(os.Stderr, l.formatMessage(req, msg))
}
