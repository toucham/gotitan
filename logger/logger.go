package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	level int8        // level of logger that filters which logs to output to STDOUT
	debug *log.Logger // level 0
	info  *log.Logger // level 1
	warn  *log.Logger // level 2
	fatal *log.Logger // level 3
}

func New(origin string) *Logger {
	l := Logger{
		level: 0,
		debug: log.New(os.Stdout, fmt.Sprintf("DEBUG [%s]: ", origin), log.Ltime),
		info:  log.New(os.Stdout, fmt.Sprintf("INFO [%s]: ", origin), log.Ltime),
		warn:  log.New(os.Stdout, fmt.Sprintf("WARN [%s]: ", origin), log.Ltime),
		fatal: log.New(os.Stdout, fmt.Sprintf("FATAL [%s]: ", origin), log.Ltime),
	}
	return &l
}

func (l *Logger) Debug(format string, v ...any) {
	if l.level < 1 {
		l.debug.Printf(format, v...)
	}
}

func (l *Logger) Info(format string, v ...any) {
	if l.level < 2 {
		l.info.Printf(format, v...)
	}
}

func (l *Logger) Warn(format string, v ...any) {
	if l.level < 3 {
		l.warn.Printf(format, v...)
	}
}

func (l *Logger) Fatal(format string, v ...any) {
	if l.level < 4 {
		l.fatal.Printf(format, v...)
	}
}
