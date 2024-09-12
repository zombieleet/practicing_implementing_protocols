package main

import (
	//"database/sql"
	"log/slog"
	"os"

	"github.com/zombieleet/ftp-protocol/internal/server"
)

func main() {

	s := server.FTPServer{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		Port:   21,
	}

	s.Start()

}
