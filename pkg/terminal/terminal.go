package terminal

import (
	"bufio"
	"fmt"
	"gorepo-cli/pkg/logger"
	"os"
	"strconv"
	"strings"
)

type Methods interface {
	AskBool(question, choices, defaultValue string, l logger.Methods) (response bool, err error)
	AskString(question, choices, defaultValue string, l logger.Methods) (response string, err error)
	SingleSelect(question string, choices [][]string, l logger.Methods) (response string, err error)
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

func (o *Terminal) SingleSelect(question string, choices [][]string, l logger.Methods) (response string, err error) {
	l.DefaultLn(question)
	choice := 0
	for i, choice := range choices {
		l.DefaultLn(fmt.Sprintf("%s. %s - %s", logger.InfoColor(strconv.Itoa(i+1)), logger.InfoColor(choice[0]), logger.VerboseColor(choice[1])))
	}
	for choice < 1 || choice > len(choices) {
		choiceStr, err := o.AskString("choice", "1-"+strconv.Itoa(len(choices)), "", l)
		if err != nil {
			return "", err
		}
		choice, err = strconv.Atoi(choiceStr)
		if err != nil {
			return "", err
		}
	}
	return strings.TrimSpace(choices[choice-1][0]), nil
}
