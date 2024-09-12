package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"runtime"
	"strconv"

	command "github.com/zombieleet/ftp-protocol/internal/commands"
	"github.com/zombieleet/ftp-protocol/internal/reply"
	"github.com/zombieleet/ftp-protocol/internal/storage"
	defaultStorage "github.com/zombieleet/ftp-protocol/internal/storage/default"
	commandParser "github.com/zombieleet/ftp-protocol/pkg/command_parser"
)

type FTPServer struct {
	Logger  *slog.Logger
	Storage storage.Storage
	Port    int
}

const WelcomeMessage = "Welcome..."

// panicIfErr purpose is to panic whenever there is a setup error
func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func blockUntilResponse(r reply.Reply, resultChan <-chan reply.ReplyResponse) {
	for {
		select {
		case replyResponse := <-resultChan:
			r.SendResponse(replyResponse.Code, replyResponse.Message)
			return
		}
	}
}

func (fServer *FTPServer) Start() {
	if fServer.Logger == nil {
		fServer.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	fServer.Logger = fServer.Logger.WithGroup("server_info").With(
		slog.String("go_version", runtime.Version()),
		slog.Int("pid", os.Getpid()),
	)

	if fServer.Storage == nil {
		fServer.Storage = defaultStorage.New()
	}

	if fServer.Port < 21 || fServer.Port > 21 {
		panicIfErr(
			fmt.Errorf("Invalid port in FTP_PORT %d only port 21 is allowed", fServer.Port),
		)
	}

	resolvedAddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(fServer.Port))

	panicIfErr(err)

	listener, err := net.ListenTCP("tcp", resolvedAddr)

	panicIfErr(err)

	fServer.Logger.Info("📡 Listening for incomming connections")

	r := reply.Reply{}

	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			fServer.Logger.Error(err.Error())
			continue
		}

		fServer.Logger = fServer.Logger.WithGroup("connection_info").With(
			slog.String("client_address", conn.RemoteAddr().String()),
		)

		fServer.Logger.Info("🤝Accpeted Connection")

		r.Conn = conn
		r.Logger = fServer.Logger

		ctx := context.Background()

		go func() {
			requestBuffer := make([]byte, 1024)

			r.SendResponse(reply.CodeReadyForNewUser, WelcomeMessage)
			for {
				lenOfRead, err := conn.Read(requestBuffer)

				defer func() {
					requestBuffer = requestBuffer[0:]
				}()

				if errors.Is(err, io.EOF) {
					fServer.Logger.Info("Received EOF")
					return
				}

				if err != nil {
					fServer.Logger.Error("Cannot read request from client", "error", err)
					continue
				}

				// resive the buffer to avoid the extra \u0000 unicode character
				// if the buffer is not filled
				requestBuffer = requestBuffer[0:lenOfRead]

				var cmd command.CMD

				rawClientRequestData := string(requestBuffer)

				if cmd, err = commandParser.CommandParser(rawClientRequestData); err != nil {
					errorAsString := err.Error()
					r.Logger = r.Logger.With("raw-command", rawClientRequestData)
					r.SendResponse(r.ToFTPResponseCode(err), errorAsString)
					continue
				}

				fServer.Logger = fServer.Logger.WithGroup("command").With(
					slog.String("cmd", cmd.Name()),
					slog.Any("params", cmd.Args()),
				)

				r.Logger = fServer.Logger

				resultChan, err := cmd.Execute(ctx, &command.ExecuteOptions{
					Storage: fServer.Storage,
					Logger:  fServer.Logger,
					Client:  conn.RemoteAddr().String(),
				})

				if err != nil {
					fServer.Logger.Error("Error executing command", "cmd", cmd.Name())
					r.SendResponse(r.ToFTPResponseCode(err), err.Error())
					continue
				}

				blockUntilResponse(r, resultChan)

			}
		}()
	}

}
