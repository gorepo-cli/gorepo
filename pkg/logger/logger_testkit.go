package logger

type MockLogger struct {
	Messages []string
}

var _ Methods = &MockLogger{}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		Messages: []string{},
	}
}

func (l *MockLogger) FatalLn(msg string) {
	l.Messages = append(l.Messages, "FATAL: "+msg)
}

func (l *MockLogger) WarningLn(msg string) {
	l.Messages = append(l.Messages, "WARNING: "+msg)
}

func (l *MockLogger) VerboseLn(msg string) {
	l.Messages = append(l.Messages, "VERBOSE: "+msg)
}

func (l *MockLogger) SuccessLn(msg string) {
	l.Messages = append(l.Messages, "SUCCESS: "+msg)
}

func (l *MockLogger) InfoLn(msg string) {
	l.Messages = append(l.Messages, "INFO: "+msg)
}

func (l *MockLogger) DefaultLn(msg string) {
	l.Messages = append(l.Messages, "DEFAULT: "+msg)
}

func (l *MockLogger) DefaultInline(msg string) {
	l.Messages = append(l.Messages, "DEFAULT: "+msg)
}

func (l *MockLogger) Output() []string {
	return l.Messages
}
