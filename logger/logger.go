package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Fatal(string, ...any)
}

type ConsoleLogger struct {
	level int8        // level of logger that filters which logs to output to STDOUT
	debug *log.Logger // level 0
	info  *log.Logger // level 1
	warn  *log.Logger // level 2
	fatal *log.Logger // level 3
}

func New(origin string) *ConsoleLogger {
	l := ConsoleLogger{
		level: 0,
		debug: log.New(os.Stdout, fmt.Sprintf("DEBUG [%s]: ", origin), log.Ltime),
		info:  log.New(os.Stdout, fmt.Sprintf("INFO [%s]: ", origin), log.Ltime),
		warn:  log.New(os.Stdout, fmt.Sprintf("WARN [%s]: ", origin), log.Ltime),
		fatal: log.New(os.Stdout, fmt.Sprintf("FATAL [%s]: ", origin), log.Ltime),
	}
	return &l
}

func (l *ConsoleLogger) Debug(format string, v ...any) {
	if l.level < 1 {
		l.debug.Printf(format, v...)
	}
}

func (l *ConsoleLogger) Info(format string, v ...any) {
	if l.level < 2 {
		l.info.Printf(format, v...)
	}
}

func (l *ConsoleLogger) Warn(format string, v ...any) {
	if l.level < 3 {
		l.warn.Printf(format, v...)
	}
}

func (l *ConsoleLogger) Fatal(format string, v ...any) {
	if l.level < 4 {
		l.fatal.Printf(format, v...)
	}
}
