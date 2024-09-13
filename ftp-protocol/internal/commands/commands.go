package command

import (
	"context"
	"log/slog"

	commandErrors "github.com/zombieleet/ftp-protocol/internal/commands/errors"
	"github.com/zombieleet/ftp-protocol/internal/reply"
	"github.com/zombieleet/ftp-protocol/internal/storage"
)

type ExecuteOptions struct {
	Storage    storage.Storage
	Logger     *slog.Logger
	Client     string
	Username   string
	RootDir    string
	CurrentDir string
	LoggedIn   bool
}

type CMD interface {
	Execute(context.Context, *ExecuteOptions) (*reply.ReplyResponse, error)
	Validate(context.Context, *ExecuteOptions) error
	Name() string
	Args() interface{}
}

func GetCommand(command string, params []string) (CMD, error) {
	var cmd CMD

	switch command {
	case "USER":
		cmd = &UserCmd{
			Params: params,
		}
	case "PASS":
		cmd = &PassCmd{
			Params: params,
		}
	case "PWD":
		cmd = &PwdCmd{}
	default:
		return nil, commandErrors.ErrBadCommand
	}

	return cmd, nil
}
