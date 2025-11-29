// Package logger
package logger

import "log"

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
}

type StdLogger struct{}

func New(cfg any) Logger {
	return &StdLogger{}
}

func (l *StdLogger) Info(msg string, args ...any)  { log.Println("[INFO]", msg, args) }
func (l *StdLogger) Error(msg string, args ...any) { log.Println("[ERROR]", msg, args) }
func (l *StdLogger) Fatal(msg string, args ...any) { log.Fatal("[FATAL]", msg, args) }
