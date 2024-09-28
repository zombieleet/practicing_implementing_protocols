package command

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"

	commandErrors "github.com/zombieleet/ftp-protocol/internal/commands/errors"
	"github.com/zombieleet/ftp-protocol/internal/reply"
)

// This command requests the server-DTP to "listen" on a data
// port (which is not its default data port) and to wait for a
// transfer command.  The response to this command includes the
// connection rather than initiate one upon receipt of a
// host and port address this server is listening on.
//
// 1. generate random port > 1024 and < 65535
// 2. check if the generated port is available to be listened on

type PasvCmd struct{}

func (pasvCmd *PasvCmd) Validate(ctx context.Context, e *ExecuteOptions) error {
	return nil
}

func (pasvCmd *PasvCmd) Execute(context.Context, *ExecuteOptions) (*reply.ReplyResponse, error) {
	var port []byte
	var err error

	if port, err = getAvailablePort(); err != nil {
		return nil, err
	}

	return &reply.ReplyResponse{
		Code:    reply.CodePassiveMode,
		Message: fmt.Sprintf("Entered passive mode. %d,%d,%d,%d,%d,%d", 127, 0, 0, 1, port[0], port[1]),
	}, nil
}

func (pasvCmd *PasvCmd) Name() string {
	return "PASV"
}

func (pasvCmd *PasvCmd) Args() interface{} {
	return ""
}

// getAvailablePort returns an available port as described in RFC 959 (https://datatracker.ietf.org/doc/html/rfc959)
func getAvailablePort() ([]byte, error) {
	// 1024 to 65535 -> registered port user applications should listen to
	// the range also includes ephemeral ports, but sine we are not using localhost:0
	// there is no need to be bothered.
	for port := 1024; port <= 65535; port++ {
		addr := fmt.Sprintf("localhost:%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			// close the connection
			defer listener.Close()

			buf := new(bytes.Buffer)

			// we use `binary.BigEndian` we want the most significant byte to be stored first (at the lowest location in the mem buf)
			// also this is what TCP/IP rfc recommends.
			// for example if we have a port of 5000, will be equivalent to [19 136]
			// but if we deicde to use binary.LittleEndian or binary.NativeEndian, it will be [136 19
			err = binary.Write(buf, binary.BigEndian, uint16(port))

			if err != nil {
				continue
			}
			fmt.Println(port)
			return buf.Bytes(), nil
		}
	}
	return []byte{}, commandErrors.ErrNoAvailablePortForPassiveDTP
}
