package logging

import (
	"fmt"
	"time"
)

type Logger struct {
	tag string
}

func NewLogger(tag string) *Logger {
	return &Logger{tag}
}

func (l *Logger) Info(message string) {
	l.log("INFO", message)
}

func (l *Logger) Warn(message string) {
	l.log("WARN", message)
}

func (l *Logger) Error(message string) {
	l.log("ERROR", message)
}

func (l *Logger) log(lvl string, message string) {
	fmt.Printf("%s [%s] %s %s \n", l.getTSFormatted(), l.tag, lvl, message)
}

func (l *Logger) getTSFormatted() string {
	return time.Now().Format(time.RFC3339)
}
