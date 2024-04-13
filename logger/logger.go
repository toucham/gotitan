package logger

import (
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	debug *log.Logger
	warn  *log.Logger
	fatal *log.Logger
}

func New() *Logger {
	l := Logger{
		info:  log.New(os.Stdout, "INFO: ", log.Ltime),
		debug: log.New(os.Stdout, "DEBUG: ", log.Ltime),
		warn:  log.New(os.Stdout, "WARN: ", log.Ltime),
		fatal: log.New(os.Stdout, "FATAL: ", log.Ltime),
	}
	return &l
}

func (l *Logger) Info(format string, v ...any) {
	l.info.Printf(format, v...)
}

func (l *Logger) Fatal(format string, v ...any) {
	l.fatal.Printf(format, v...)
}

func (l *Logger) Warn(format string, v ...any) {
	l.warn.Printf(format, v...)
}

func (l *Logger) Debug(format string, v ...any) {
	l.debug.Printf(format, v...)
}
