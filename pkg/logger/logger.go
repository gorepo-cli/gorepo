package logger

import (
	"github.com/fatih/color"
	"log"
	"os"
)

type Methods interface {
	FatalLn(msg string)
	WarningLn(msg string)
	VerboseLn(msg string)
	SuccessLn(msg string)
	InfoLn(msg string)
	DefaultLn(msg string)
	DefaultInline(msg string)
}

var (
	FatalColor   = color.New(color.FgRed).SprintFunc()
	WarningColor = color.New(color.FgYellow).SprintFunc()
	VerboseColor = color.New(color.FgHiBlack).SprintFunc()
	SuccessColor = color.New(color.FgGreen).SprintFunc()
	InfoColor    = color.New(color.FgCyan).SprintFunc()
)

type Logger struct {
	*log.Logger
}

var _ Methods = &Logger{}

func NewLogger() *Logger {
	return &Logger{Logger: log.New(os.Stdout, "", 0)}
}

func (l *Logger) FatalLn(msg string) {
	l.Println(FatalColor(msg))
}

func (l *Logger) WarningLn(msg string) {
	l.Logger.Println(WarningColor(msg))
}

func (l *Logger) VerboseLn(msg string) {
	l.Logger.Println(VerboseColor(msg))
}

func (l *Logger) SuccessLn(msg string) {
	l.Logger.Println(SuccessColor(msg))
}

func (l *Logger) InfoLn(msg string) {
	l.Logger.Println(InfoColor(msg))
}

func (l *Logger) DefaultLn(msg string) {
	l.Logger.Println(msg)
}

func (l *Logger) DefaultInline(msg string) {
	_, _ = l.Writer().Write([]byte(msg))
}
