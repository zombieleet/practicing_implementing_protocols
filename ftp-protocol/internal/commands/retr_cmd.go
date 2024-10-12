package command

import (
	"context"
	"path/filepath"

	commandErrors "github.com/zombieleet/ftp-protocol/internal/commands/errors"
	"github.com/zombieleet/ftp-protocol/internal/dtp"
	"github.com/zombieleet/ftp-protocol/internal/reply"
)

type RetrCmd struct {
	Params []string
}

func (retrCmd *RetrCmd) Validate(ctx context.Context, e *ExecuteOptions) error {
	if !e.LoggedIn {
		return commandErrors.ErrUnAuthorized
	}

	if len(retrCmd.Params) > 1 {
		return commandErrors.ErrInvalidArgLength
	}

	return nil
}

func (retrCmd *RetrCmd) Execute(ctx context.Context, e *ExecuteOptions) (*reply.ReplyResponse, error) {
	var err error

	if err = retrCmd.Validate(ctx, e); err != nil {
		return nil, err
	}

	e.DTPChannel <- dtp.DTPChannel{
		FileOperationCommand: "RETR",
		FilePath:             filepath.Join(e.CurrentDir, retrCmd.Params[0]),
	}

	return &reply.ReplyResponse{
		Code:    reply.CodeOpeningDataConn,
		Message: "Data connection already opened; transfer starting.",
	}, nil
}

func (retrCmd *RetrCmd) Name() string {
	return "RETR"
}

func (retrCmd *RetrCmd) Args() interface{} {
	return retrCmd.Params
}
