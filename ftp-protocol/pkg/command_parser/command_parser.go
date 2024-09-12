package command_parser

import (
	"errors"
	"strings"

	command "github.com/zombieleet/ftp-protocol/internal/commands"
)

var (
	ErrEmptyCommand = errors.New("Empty command")
)

func CommandParser(rawString string) (command.CMD, error) {

	if len(rawString) == 0 {
		return nil, ErrEmptyCommand
	}

	rr := strings.ReplaceAll(rawString, "\r\n", "")

	splittedRawString := strings.Split(rr, " ")

	cmd, err := command.GetCommand(splittedRawString[0], splittedRawString[1:])

	if err != nil {
		return nil, err
	}

	return cmd, nil
}
