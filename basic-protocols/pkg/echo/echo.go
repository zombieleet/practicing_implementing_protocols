package echo

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"net"
)

func EchoClientMessage(logger *slog.Logger, conn *net.TCPConn) {
	connIOReader := bufio.NewReader(conn)

	for {
		lineByte, _, err := connIOReader.ReadLine()
		if errors.Is(err, io.EOF) {
			logger.Info("got EOF from client", "client", conn.RemoteAddr().String())
			continue
		}

		if err != nil {
			logger.Error("got an error while reading client message", "error", err.Error(), "client", conn.RemoteAddr().String())
			continue
		}
		conn.Write(append(lineByte, '\n'))
	}
}
