package command

import (
	"context"
	"github.com/zombieleet/ftp-protocol/internal/reply"
	"runtime"
)

type SystCmd struct{}

func (syst *SystCmd) Execute(context.Context, *ExecuteOptions) (*reply.ReplyResponse, error) {
	return &reply.ReplyResponse{
		Code:    reply.CodeOk,
		Message: runtime.GOOS,
	}, nil
}

func (syst *SystCmd) Validate(context.Context, *ExecuteOptions) error {
	return nil
}

func (syst *SystCmd) Name() string {
	return "SYST"
}

func (syst *SystCmd) Args() interface{} {
	return ""
}
