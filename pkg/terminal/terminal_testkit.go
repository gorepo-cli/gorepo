package terminal

import (
	"errors"
	"fmt"
	"gorepo-cli/pkg/logger"
)

type MockMethods interface {
	Methods
}

type MockTerminal struct {
	QuestionsAnswersBool         map[string]bool
	QuestionsAnswersString       map[string]string
	QuestionsAnswersSingleSelect map[string]string
}

var _ MockMethods = &MockTerminal{}

func NewMockTerminal(qABool map[string]bool, qAString map[string]string, qaSingleSelect map[string]string) *MockTerminal {
	return &MockTerminal{
		QuestionsAnswersBool:         qABool,
		QuestionsAnswersString:       qAString,
		QuestionsAnswersSingleSelect: qaSingleSelect,
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
		if answer == "" {
			return defaultValue, nil
		}
		return answer, nil
	}
	return "", errors.New(fmt.Sprintf("question `%s` not in the test, provide an answer", question))
}

func (m *MockTerminal) SingleSelect(question string, choices [][]string, l logger.Methods) (response string, err error) {
	if answer, exists := m.QuestionsAnswersSingleSelect[question]; exists {
		return answer, nil
	}
	return "", errors.New(fmt.Sprintf("question `%s` not in the test, provide an answer", question))
}
