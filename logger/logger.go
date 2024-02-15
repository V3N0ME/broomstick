package logger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Logger interface {
	Info(...any)
	Error(...any)
}

type logger struct {
	wal *bufio.Writer
}

func NewLogger() Logger {

	wal, err := os.OpenFile("broomstick.logs", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	writter := bufio.NewWriter(wal)

	return &logger{
		wal: writter,
	}
}

func (l *logger) log(msg string) {

	_, filename, line, _ := runtime.Caller(2)

	log.Printf("[%s:%d] %v", filename, line, msg)

	l.wal.WriteString(msg + "\n")
	l.wal.Flush()
}

func (l *logger) Info(msgs ...any) {

	combinedMsg := ""

	for _, msg := range msgs {
		combinedMsg += fmt.Sprint(msg)
		combinedMsg += " "
	}

	combinedMsg = strings.TrimSpace(combinedMsg)

	l.log(fmt.Sprint("[INFO]", " ", combinedMsg, " ", time.Now()))
}

func (l *logger) Error(msgs ...any) {

	combinedMsg := ""

	for _, msg := range msgs {
		combinedMsg += fmt.Sprint(msg)
		combinedMsg += " "
	}

	combinedMsg = strings.TrimSpace(combinedMsg)

	l.log(fmt.Sprint("[ERROR]", " ", combinedMsg, " ", time.Now()))
}
