package console

import (
	"log"
	"strings"
)

type Level uint32

const (
	debug Level = iota
	info
	warn
	err
	fatal
)

func parseLevel(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return debug
	case "info":
		return info
	case "warn":
		return warn
	case "error":
		return err
	case "fatal":
		return fatal
	}
	return info
}

type Logger struct {
	level Level
}

func NewLogger(level string) *Logger {
	return &Logger{parseLevel(level)}
}

func NewFatalLogger() *Logger {
	return &Logger{fatal}
}

func (logger *Logger) Debug(msg string) {
	if logger.level <= debug {
		log.Println("DEBUG:", msg)
	}
}

func (logger *Logger) Info(msg string) {
	if logger.level <= info {
		log.Println("INFO:", msg)
	}
}

func (logger *Logger) Warn(msg string) {
	if logger.level <= warn {
		log.Println("WARN:", msg)
	}
}

func (logger *Logger) Error(msg string) {
	if logger.level <= err {
		log.Println("ERROR:", msg)
	}
}

func (logger *Logger) Fatal(msg string) {
	if logger.level <= fatal {
		log.Panicln("FATAL:", msg)
	}
}
