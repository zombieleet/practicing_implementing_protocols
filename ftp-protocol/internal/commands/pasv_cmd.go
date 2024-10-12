package command

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"

	commandErrors "github.com/zombieleet/ftp-protocol/internal/commands/errors"
	"github.com/zombieleet/ftp-protocol/internal/dtp"
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

type PasvCmd struct {
	command string
}

func (pasvCmd *PasvCmd) Validate(ctx context.Context, e *ExecuteOptions) error {
	if !e.LoggedIn {
		return commandErrors.ErrUnAuthorized
	}
	return nil
}

// Execute get's avaialable port on the server to listen to
func (pasvCmd *PasvCmd) Execute(ctx context.Context, e *ExecuteOptions) (*reply.ReplyResponse, error) {
	var port []byte
	var normalizedPort int
	var err error

	if err = pasvCmd.Validate(ctx, e); err != nil {
		return nil, err
	}

	if normalizedPort, port, err = getAvailablePort(); err != nil {
		return nil, err
	}

	if e.DTPControlResponseChannel, err = dtp.OpenDTPPort(normalizedPort, e.DTPChannel); err != nil {
		return &reply.ReplyResponse{
			Code:    reply.CodeUnableToOpenDataConn,
			Message: err.Error(),
		}, nil
	}


	var message string
	var responseCode int

	switch pasvCmd.command {
	case "EPSV":
		// According to RFC-2428, we can also specify the network protocol, ip address and port
		// but it is expect that the server connects with the same protocol and ip address the control connection is using
		// the format is (|<Protocol>|<Address>|<Port>|)
		//
		// The client can also issue an EPSV command to tell the server what protocol it needs to open connection on
		// EPSV<space><net-protocol> -> EPSV udp6
		//
		// TODO:
		// 1. implement support for allowing client specify networkprotocol
		//
		// There is no need to handle EPSV ALL for now
		message = fmt.Sprintf("Entering Extended Passive Mode (|||%d|)", normalizedPort)
		responseCode = reply.CodeExtendedPassiveMode
	case "PASV":
		message = fmt.Sprintf("Entering Passive Mode %d,%d,%d,%d,%d,%d", 127, 0, 0, 1, port[0], port[1])
		responseCode = reply.CodePassiveMode
	case "LPSV":
		// According to RFC-1639, we also need to specify the internet protocol verson to use, the ip address and the port
		message = fmt.Sprintf("Entering Long Passive Mode %d,%d,%d,%d,%d,%d", 4, 127, 0, 0, 1, normalizedPort)
		responseCode = reply.CodeLongPassiveMode
	}
	return &reply.ReplyResponse{
		Code:    responseCode,
		Message: message,
	}, nil
}

func (pasvCmd *PasvCmd) Name() string {
	return pasvCmd.command
}

func (pasvCmd *PasvCmd) Args() interface{} {
	return ""
}

// getAvailablePort returns an available port as described in RFC 959 (https://datatracker.ietf.org/doc/html/rfc959)
func getAvailablePort() (int, []byte, error) {
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
			return port, buf.Bytes(), nil
		}
	}
	return 0, []byte{}, commandErrors.ErrNoAvailablePortForPassiveDTP
}
