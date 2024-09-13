package command

import (
	"context"
	commandErrors "github.com/zombieleet/ftp-protocol/internal/commands/errors"
	"github.com/zombieleet/ftp-protocol/internal/reply"
)

type PassCmd struct {
	Params      []string
	maskedParam []string
}

func (p *PassCmd) Validate(ctx context.Context, e *ExecuteOptions) error {
	if len(p.Params) > 1 || len(p.Params) < 1 {
		return commandErrors.ErrInvalidArgLength
	}

	password := p.Params[0]

	return e.Storage.Login(e.Username, password)
}

func (p *PassCmd) Execute(ctx context.Context, e *ExecuteOptions) (*reply.ReplyResponse, error) {
	if err := p.Validate(ctx, e); err != nil {
		return nil, err
	}

	return &reply.ReplyResponse{
		Code:    reply.CodeLoggedInOk,
		Message: "Logged In",
	}, nil
}

func (p *PassCmd) Name() string {
	return "PASS"
}

func (p *PassCmd) Args() interface{} {
	if len(p.maskedParam) > 0 {
		return p.maskedParam
	}

	p.maskedParam = make([]string, len(p.Params))

	copy(p.maskedParam, p.Params)

	for k, _ := range p.maskedParam {
		p.maskedParam[k] = "xxxxxx"
	}
	return p.maskedParam
}
