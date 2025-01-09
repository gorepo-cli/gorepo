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
	WarningColor = color.New(color.FgHiMagenta).SprintFunc()
	VerboseColor = color.New(color.FgHiBlack).SprintFunc()
	InfoColor    = color.New(color.FgCyan).SprintFunc()
	SuccessColor = color.New(color.FgGreen).SprintFunc()
)

type Logger struct {
	*log.Logger
}

var _ Methods = &Logger{}

func NewLogger() *Logger {
	return &Logger{Logger: log.New(os.Stdout, "", 0)}
}

func (l *Logger) FatalLn(msg string) {
	l.Println(FatalColor("the blue llama ran into an error: " + msg))
}

func (l *Logger) WarningLn(msg string) {
	l.Logger.Println(WarningColor(msg))
}

func (l *Logger) VerboseLn(msg string) {
	l.Logger.Println(VerboseColor("the blue llama gossips: " + msg))
}

func (l *Logger) InfoLn(msg string) {
	l.Logger.Println(InfoColor(msg))
}

func (l *Logger) SuccessLn(msg string) {
	l.Logger.Println(SuccessColor(msg))
}

func (l *Logger) DefaultLn(msg string) {
	l.Logger.Println(msg)
}

func (l *Logger) DefaultInline(msg string) {
	_, _ = l.Writer().Write([]byte(msg))
}
