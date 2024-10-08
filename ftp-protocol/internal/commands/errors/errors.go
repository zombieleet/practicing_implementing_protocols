package command_errors

import "errors"

var (
	ErrInvalidArgLength             = errors.New("Invalid number of arguments")
	ErrBadCommand                   = errors.New("Command is not allowed")
	ErrUserDoesNotExists            = errors.New("User does not exists")
	ErrUnAuthorized                 = errors.New("Unauthorized operation, need login")
	ErrNoAvailablePortForPassiveDTP = errors.New("No availaable prots in the range 1024 - 65535")
)
