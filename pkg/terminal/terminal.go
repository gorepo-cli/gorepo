package terminal

import (
	"bufio"
	"gorepo-cli/pkg/logger"
	"os"
	"strings"
)

type Methods interface {
	AskBool(question, choices, defaultValue string, l logger.Methods) (response bool, err error)
	AskString(question, choices, defaultValue string, l logger.Methods) (response string, err error)
}

type Terminal struct{}

var _ Methods = &Terminal{}

func (o *Terminal) AskBool(question, choices, defaultValue string, l logger.Methods) (response bool, err error) {
	questionFormated := question
	choicesFormated := ""
	if choices != "" {
		choicesFormated = logger.InfoColor("(" + choices + ")")
	}
	defaultValueFormated := ""
	if defaultValue != "" {
		defaultValueFormated = logger.VerboseColor("default: " + defaultValue)
	}
	l.DefaultInline(questionFormated + " " + choicesFormated + " " + defaultValueFormated + ": ")
	reader := bufio.NewReader(os.Stdin)
	responseStr, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	responseStr = strings.TrimSpace(strings.ToLower(responseStr))
	if responseStr == "" {
		responseStr = defaultValue
	}
	return responseStr == "y" || responseStr == "yes", nil
}

func (o *Terminal) AskString(question, choices, defaultValue string, l logger.Methods) (response string, err error) {
	questionFormated := question
	choicesFormated := ""
	if choices != "" {
		choicesFormated = logger.InfoColor("(" + choices + ")")
	}
	defaultValueFormated := ""
	if defaultValue != "" {
		defaultValueFormated = logger.VerboseColor("default: " + defaultValue)
	}
	l.DefaultInline(questionFormated + " " + choicesFormated + " " + defaultValueFormated + ": ")
	reader := bufio.NewReader(os.Stdin)
	responseStr, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	responseStr = strings.TrimSpace(responseStr)
	if responseStr == "" {
		responseStr = defaultValue
	}
	return responseStr, nil
}
