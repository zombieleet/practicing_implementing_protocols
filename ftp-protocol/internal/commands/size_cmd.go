package command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	commandErrors "github.com/zombieleet/ftp-protocol/internal/commands/errors"
	"github.com/zombieleet/ftp-protocol/internal/reply"
)

// SizeCmd is responsible for returning the size of a file
type SizeCmd struct {
	Params []string
}

func (sizeCmd *SizeCmd) Validate(ctx context.Context, e *ExecuteOptions) error {
	if !e.LoggedIn {
		return commandErrors.ErrUnAuthorized
	}

	if len(sizeCmd.Params) > 1 {
		return commandErrors.ErrInvalidArgLength
	}

	return nil
}

func (sizeCmd *SizeCmd) Execute(ctx context.Context, e *ExecuteOptions) (*reply.ReplyResponse, error) {
	var err error

	if err = sizeCmd.Validate(ctx, e); err != nil {
		return nil, err
	}

	absoluteFilePath := filepath.Join(e.CurrentDir, sizeCmd.Params[0])

	fileInfo, err := os.Stat(absoluteFilePath)

	if err != nil {
		return &reply.ReplyResponse{
			Code:    reply.CodeFileNotFound,
			Message: err.Error(),
		}, nil
	}

	return &reply.ReplyResponse{
		Code:    reply.CodeFileStatInfo,
		Message: fmt.Sprintf("%d", fileInfo.Size()),
	}, nil
}

func (sizeCmd *SizeCmd) Name() string {
	return "SIZE"
}

func (sizeCmd *SizeCmd) Args() interface{} {
	return sizeCmd.Params
}
