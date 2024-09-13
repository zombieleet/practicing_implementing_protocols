package command

import (
	"context"

	commandErrors "github.com/zombieleet/ftp-protocol/internal/commands/errors"

	"github.com/zombieleet/ftp-protocol/internal/reply"
)

type UserCmd struct {
	Params []string
}

// Validates validates USER cmd and it's parameters
func (u *UserCmd) Validate(_ context.Context, e *ExecuteOptions) error {
	if len(u.Params) > 1 || len(u.Params) < 1 {
		return commandErrors.ErrInvalidArgLength
	}

	user := u.Params[0]

	if exists := e.Storage.UserExists(user); !exists {
		return commandErrors.ErrUserDoesNotExists
	}

	return nil
}

// Execute validates and execute operation under USER cmd
func (u *UserCmd) Execute(c context.Context, e *ExecuteOptions) (*reply.ReplyResponse, error) {
	if err := u.Validate(c, e); err != nil {
		return nil, err
	}

	return &reply.ReplyResponse{
		Code:    reply.CodeUserNameOkay,
		Message: "USER ok. Enter PASSWORD",
	}, nil
}

// Name returns the name of USER command
func (u *UserCmd) Name() string {
	return "USER"
}

func (u *UserCmd) Args() interface{} {
	return u.Params
}
