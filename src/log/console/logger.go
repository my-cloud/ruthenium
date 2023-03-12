package console

import (
	"log"
	"strings"
)

type Level uint32

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
)

func ParseLevel(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return Debug
	case "info":
		return Info
	case "warn":
		return Warn
	case "error":
		return Error
	case "fatal":
		return Fatal
	}
	return Info
}

type Logger struct {
	level Level
}

func NewLogger(level Level) *Logger {
	return &Logger{level}
}

func (logger *Logger) Debug(msg string) {
	if logger.level <= Debug {
		log.Println("DEBUG:", msg)
	}
}

func (logger *Logger) Info(msg string) {
	if logger.level <= Info {
		log.Println("INFO:", msg)
	}
}

func (logger *Logger) Warn(msg string) {
	if logger.level <= Warn {
		log.Println("WARN:", msg)
	}
}

func (logger *Logger) Error(msg string) {
	if logger.level <= Error {
		log.Println("ERROR:", msg)
	}
}

func (logger *Logger) Fatal(msg string) {
	if logger.level <= Fatal {
		log.Panicln("FATAL:", msg)
	}
}
