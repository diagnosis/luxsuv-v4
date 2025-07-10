package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	info *log.Logger
	err  *log.Logger
	warn *log.Logger
	file *os.File
}

func NewLogger(logfilePath string) (*Logger, error) {
	file, err := os.OpenFile(logfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	multiWriter := io.MultiWriter(os.Stdout, file)
	return &Logger{
		info: log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		err:  log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		warn: log.New(multiWriter, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		file: file,
	}, nil
}

func (l *Logger) Info(msg string) {
	l.info.Printf("%s", msg)
}

func (l *Logger) Err(msg string) {
	l.err.Printf("%s", msg)
}

func (l *Logger) Warn(msg string) {
	l.warn.Printf("%s", msg)
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		log.Printf("Error closing log file: %v", err)
	}
}
