package command

import (
	"context"
	"log/slog"

	commandErrors "github.com/zombieleet/ftp-protocol/internal/commands/errors"
	"github.com/zombieleet/ftp-protocol/internal/dtp"
	"github.com/zombieleet/ftp-protocol/internal/reply"
	"github.com/zombieleet/ftp-protocol/internal/storage"
)

type ExecuteOptions struct {
	Storage                   storage.Storage
	Logger                    *slog.Logger
	Client                    string
	Username                  string
	RootDir                   string
	CurrentDir                string
	LoggedIn                  bool
	DTPChannel                chan dtp.DTPChannel
	DTPControlResponseChannel <-chan *reply.ReplyResponse
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
	case "SYST":
		cmd = &SystCmd{}
	case "PASV", "EPSV", "LPSV":
		cmd = &PasvCmd{
			command: command,
		}
	case "SIZE":
		cmd = &SizeCmd{
			Params: params,
		}
	case "RETR":
		cmd = &RetrCmd{
			Params: params,
		}
	default:
		return nil, commandErrors.ErrBadCommand
	}

	return cmd, nil
}
