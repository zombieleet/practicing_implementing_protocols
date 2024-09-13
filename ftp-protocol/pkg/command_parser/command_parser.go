package command_parser

import (
	"errors"
	"regexp"
	"strings"

	command "github.com/zombieleet/ftp-protocol/internal/commands"
)

var (
	ErrEmptyCommand = errors.New("Empty command")
)

var re = regexp.MustCompile(`\r\n|\r|\n`)

func CommandParser(rawString string) (command.CMD, error) {

	if len(rawString) == 0 {
		return nil, ErrEmptyCommand
	}

	parsedString := re.ReplaceAllString(rawString, "")
	splittedRawString := strings.Split(parsedString, " ")
	cmd, err := command.GetCommand(splittedRawString[0], splittedRawString[1:])

	if err != nil {
		return nil, err
	}

	return cmd, nil
}
