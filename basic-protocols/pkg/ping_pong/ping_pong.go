package ping_pong

import (
	"errors"
	"io"
	"log/slog"
	"net"
)

// ReplyWithPong replies a client connecting to it with a pong message
// when it recieves a ping
func ReplyWithPong(logger *slog.Logger, conn *net.TCPConn) {
	buffer := make([]byte, 4, 4)
	for {

		numberOfBytesRead, err := conn.Read(buffer)

		if errors.Is(err, io.EOF) {
			logger.Info("Received EOF", "address", conn.RemoteAddr())
			conn.Close()
			return
		}

		if err != nil {
			logger.Error("Recevied error while reading from client", "error", err.Error(), "address", conn.RemoteAddr())
			return
		}

		clientMessageAsString := string(buffer)

		if numberOfBytesRead != 4 || clientMessageAsString != "ping" {
			continue
		}

		numberOfBytesWritten, err := conn.Write([]byte("pong\n"))

		if err != nil {
			logger.Error("Cannot write pong to client", "error", err.Error(), "address", conn.RemoteAddr())
			conn.Close()
			return
		}

		if numberOfBytesWritten != 5 {
			logger.Error("Incomplete/Unexpected content written to client", "address", conn.RemoteAddr())
		}

		buffer = buffer[:0]
	}

}
