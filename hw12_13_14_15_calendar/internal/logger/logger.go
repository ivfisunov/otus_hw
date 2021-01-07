package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logg struct {
	logger *log.Logger
	level  string
}

type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
}

var mapping = map[string]int{
	"DEBUG": 1,
	"INFO":  2,
	"WARN":  3,
	"ERROR": 4,
}

func New(level string, filePath string) (*Logg, error) {
	if _, ok := mapping[level]; !ok {
		return nil, fmt.Errorf("invalid logger level type: %v", level)
	}

	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening logfile: %w", err)
	}

	mw := io.MultiWriter(logFile, os.Stdout)
	logger := log.New(mw, "", log.Ldate|log.Ltime)
	return &Logg{logger: logger, level: level}, nil
}

func (l Logg) Debug(msg string) {
	if mapping["DEBUG"] >= mapping[l.level] {
		l.logger.SetPrefix("[DEBUG] ")
		l.logger.Println(msg)
	}
}

func (l Logg) Info(msg string) {
	if mapping["INFO"] >= mapping[l.level] {
		l.logger.SetPrefix("[INFO] ")
		l.logger.Println(msg)
	}
}

func (l Logg) Warn(msg string) {
	if mapping["WARN"] >= mapping[l.level] {
		l.logger.SetPrefix("[WARN] ")
		l.logger.Println(msg)
	}
}

func (l Logg) Error(msg string) {
	if mapping["ERROR"] >= mapping[l.level] {
		l.logger.SetPrefix("[ERROR] ")
		l.logger.Println(msg)
	}
}
