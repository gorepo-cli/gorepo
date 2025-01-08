package terminal

import (
	"errors"
	"fmt"
	"gorepo-cli/pkg/logger"
)

type MockTerminal struct {
	QuestionsAnswersBool   map[string]bool
	QuestionsAnswersString map[string]string
}

var _ Methods = &MockTerminal{}

func NewMockOs(wd string, qABool map[string]bool, qAString map[string]string) *MockTerminal {
	return &MockTerminal{
		QuestionsAnswersBool:   qABool,
		QuestionsAnswersString: qAString,
	}
}

func (m *MockTerminal) AskBool(question, choices, defaultValue string, l logger.Methods) (response bool, err error) {
	if answer, exists := m.QuestionsAnswersBool[question]; exists {
		return answer, nil
	}
	return false, errors.New(fmt.Sprintf("question `%s` not in the test, provide an answer", question))
}

func (m *MockTerminal) AskString(question, choices, defaultValue string, l logger.Methods) (response string, err error) {
	if answer, exists := m.QuestionsAnswersString[question]; exists {
		return answer, nil
	}
	return "", errors.New(fmt.Sprintf("question `%s` not in the test, provide an answer", question))
}
