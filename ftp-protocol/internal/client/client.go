package client

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"

	command "github.com/zombieleet/ftp-protocol/internal/commands"
	"github.com/zombieleet/ftp-protocol/internal/reply"
	"github.com/zombieleet/ftp-protocol/internal/storage"
	commandParser "github.com/zombieleet/ftp-protocol/pkg/command_parser"
)

type FTPClient struct {
	Conn   *net.TCPConn
	Logger *slog.Logger

	Storage  storage.Storage
	loggedIn bool

	currentDir string
	username   string
}

func (ftpClient *FTPClient) Handle() {
	r := reply.Reply{
		Conn:   ftpClient.Conn,
		Logger: ftpClient.Logger,
	}

	ftpClient.Logger.Info("🤝Accpeted Connection")

	ctx := context.Background()
	r.SendResponse(reply.CodeReadyForNewUser, "Welcome...")

	reader := bufio.NewReader(ftpClient.Conn)

	homedir, _ := os.UserHomeDir()

	ftpClient.currentDir = homedir

	for {

		rawClientRequestData, err := reader.ReadString('\n')

		if errors.Is(err, io.EOF) {
			ftpClient.Logger.Info("Received EOF")
			return
		}

		if err != nil {
			ftpClient.Logger.Error("Cannot read request from client", "error", err)
			continue
		}

		rawClientRequestData = strings.TrimSpace(rawClientRequestData)
		if rawClientRequestData == "" {
			continue
		}

		var cmd command.CMD

		if cmd, err = commandParser.CommandParser(rawClientRequestData); err != nil {
			errorAsString := err.Error()
			r.Logger = r.Logger.With("raw_command", rawClientRequestData)
			r.SendResponse(r.ToFTPResponseCode(err), errorAsString)
			continue
		}

		ftpClient.Logger = ftpClient.Logger.WithGroup("command").With(
			slog.String("cmd", cmd.Name()),
			slog.Any("params", cmd.Args()),
		)

		r.Logger = ftpClient.Logger

		replyResponse, err := cmd.Execute(ctx, &command.ExecuteOptions{
			LoggedIn:   ftpClient.loggedIn,
			Username:   ftpClient.username,
			CurrentDir: ftpClient.currentDir,
			Storage:    ftpClient.Storage,
			Logger:     ftpClient.Logger,
			Client:     ftpClient.Conn.RemoteAddr().String(),
			RootDir:    homedir,
		})

		if err != nil {
			ftpClient.Logger.Error("Error executing command", "cmd", cmd.Name())
			r.SendResponse(r.ToFTPResponseCode(err), err.Error())
			continue
		}

		r.SendResponse(replyResponse.Code, replyResponse.Message)

		switch cmd.Name() {
		case "PASS":
			ftpClient.loggedIn = true
		case "USER":
			params, ok := cmd.Args().([]string)
			if !ok {
				r.SendResponse(reply.CodeLocalErrorProcessing, "Cannot set username specified")
				continue
			}
			ftpClient.username = params[0]
		}
	}

}
