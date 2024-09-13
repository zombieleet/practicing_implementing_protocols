package server

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"runtime"
	"strconv"

	"github.com/zombieleet/ftp-protocol/internal/client"

	"github.com/zombieleet/ftp-protocol/internal/storage"
	defaultStorage "github.com/zombieleet/ftp-protocol/internal/storage/default"
)

type FTPServer struct {
	Logger  *slog.Logger
	Storage storage.Storage
	Port    int
}

// panicIfErr purpose is to panic whenever there is a setup error
func panicIfErr(err error) {
	if err != nil {
		panic(err)
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

	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			fServer.Logger.Error(err.Error())
			continue
		}

		cl := client.FTPClient{
			Conn: conn,
			Logger: fServer.Logger.WithGroup("connection_info").With(
				slog.String("client_address", conn.RemoteAddr().String()),
			),
			Storage: fServer.Storage,
		}

		go cl.Handle()
	}
}
