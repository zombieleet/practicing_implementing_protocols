package command

import (
	"context"

	commandErrors "github.com/zombieleet/ftp-protocol/internal/commands/errors"
	"github.com/zombieleet/ftp-protocol/internal/reply"
)

type PwdCmd struct{}

func (pwd *PwdCmd) Validate(ctx context.Context, e *ExecuteOptions) error {
	if !e.LoggedIn {
		return commandErrors.ErrUnAuthorized
	}
	return nil
}

func (pwd *PwdCmd) Execute(ctx context.Context, e *ExecuteOptions) (*reply.ReplyResponse, error) {
	if err := pwd.Validate(ctx, e); err != nil {
		return nil, err
	}
	return &reply.ReplyResponse{
		Code:    reply.CodeOk,
		Message: e.CurrentDir,
	}, nil
}

func (pwd *PwdCmd) Name() string {
	return "PWD"
}

func (pwd *PwdCmd) Args() interface{} {
	return ""
}
