package main

import (
	"github.com/zombieelet/basic-protocols/internal/cmd_utils"
	"github.com/zombieelet/basic-protocols/internal/server_utils"
	"github.com/zombieelet/basic-protocols/pkg/ping_pong"
)

func main() {
	hostPort := cmd_utils.GetHostAndPort()
	logger := cmd_utils.GetLogger().With("host:port", hostPort)

	logger = logger.With("host:port", hostPort)

	server, err := server_utils.CreateTCPServer(
		server_utils.ServerOptions{
			HostPort: hostPort,
			Logger:   logger,
		},
	)

	if err != nil {
		logger.Error("cannot start ping-pong-server", "err", err.Error())
		return
	}

	logger.Info("Listening for incomming connections", "address", server.Addr().String())

	for {
		conn, err := server.AcceptTCP()

		if err != nil {
			logger.Error("cannot accept an incomming connection")
			continue
		}

		logger.Info("accepting connection from ", "address", conn.RemoteAddr())

		go ping_pong.ReplyWithPong(logger, conn)
	}

}
