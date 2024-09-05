package cmd_utils

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

var port = flag.Int("port", 1337, "The port the ping pong server should run on")
var logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

func GetHostAndPort() string {
	flag.Parse()
	if port == nil {
		pingPongPort := os.Getenv("PING_PONG_PORT")
		parsedPortAsInt, err := strconv.ParseInt(pingPongPort, 10, 32)
		if err != nil {
			logger.Error("ping ping address not set")
		}
		*port = int(parsedPortAsInt)
	}
	return fmt.Sprintf(":%d", *port)
}

func GetLogger() *slog.Logger {
	return logger
}
