// dtp WIP
// TODO:
// 1. support block mode transfer
package dtp

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/zombieleet/ftp-protocol/internal/reply"
)

type DTPChannel struct {
	FileOperationCommand string
	FilePath             string
}

// TODO:
// 1. parse in some sort of dtp port option that takes all this details + a logger
// 2. compose a new logger out of the logger
func OpenDTPPort(port int, dtpChannel <-chan DTPChannel) (<-chan *reply.ReplyResponse, error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))

	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		return nil, err
	}

	responseChan := make(chan *reply.ReplyResponse)

	go func() {

		defer listener.Close()

		fmt.Println("listening on DTP port ", ":"+strconv.Itoa(port))

		conn, err := listener.AcceptTCP()

		if err != nil {
			return
		}

		go func() {
			defer conn.Close()

			dtpOperation := <-dtpChannel

			if dtpOperation.FileOperationCommand == "RETR" {
				responseChan <- retr(conn, dtpOperation.FilePath)
			}
		}()

	}()

	return responseChan, nil
}

func retr(conn *net.TCPConn, filename string) *reply.ReplyResponse {

	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)

	if err != nil {
		return &reply.ReplyResponse{
			Code:    reply.CodeLocalErrorProcessing,
			Message: err.Error(),
		}
	}

	defer file.Close()

	_, err = conn.ReadFrom(file)

	if err != nil {
		fmt.Println(err)
		return &reply.ReplyResponse{
			Code:    reply.CodeLocalErrorProcessing,
			Message: err.Error(),
		}
	}

	return &reply.ReplyResponse{
		Code:    reply.CodeClosingDataConnection,
		Message: fmt.Sprintf("Transfer complete <%s>", filename),
	}
}
