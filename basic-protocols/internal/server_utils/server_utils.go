package server_utils

import (
	"log/slog"
	"net"
)

type ServerOptions struct {
	HostPort string
	Logger   *slog.Logger
}

func CreateTCPServer(opt ServerOptions) (*net.TCPListener, error) {
	resolvedTcpAddr, err := net.ResolveTCPAddr("tcp", opt.HostPort)
	logger := opt.Logger

	if err != nil {
		logger.Error("Cannot resolve tcp addr", "error", err.Error())
		return nil, err
	}

	return net.ListenTCP("tcp", resolvedTcpAddr)
}
