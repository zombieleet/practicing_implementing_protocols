package reply

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"regexp"

	commandErrors "github.com/zombieleet/ftp-protocol/internal/commands/errors"
)

type Reply struct {
	Logger *slog.Logger
	Conn   *net.TCPConn
}

type ReplyResponse struct {
	Code    int
	Message string
}

var (
	CodeOk                                = 200 // RFC 959
	CodeUnrecognizedCommand               = 500
	CodeSyntaxErrorInArgs                 = 501
	CodeCommandNotImplemented             = 502
	CodeBadSequenceOfCommands             = 503
	CodeCommandNotImplementedForParameter = 504
	CodeRestartMarkerReply                = 110
	CodeSystemStatus                      = 211
	CodeSystemHelp                        = 211
	CodeDirStatus                         = 212
	CodeFileStatus                        = 213
	// On how to use the server or the meaning of a particular
	// non-standard command.  This reply is useful only to the human user.
	CodeHelpMessage                = 214
	CodeClosingConnection          = 221
	CodeReadyForNewUser            = 220
	CodeServiceUnavailable         = 421
	CodeDataConnAlreadyOpen        = 125
	CodeDataConnOpen               = 225
	CodeUnableToOpenDataConn       = 425
	CodeConnectionClosed           = 426
	CodeTransferAborted            = 426
	CodePassiveMode                = 227
	CodeUserLoggedIn               = 230
	CodeNotLoggedIn                = 530
	CodeUserNameOkay               = 331
	CodeNeedAccountForLogin        = 332
	CodeNeedAccountForStoringFiles = 532
	CodeOpeningDataConn            = 150
	CodeFileActionCompleted        = 250
	CodePathCreated                = 257
	CodeRequestedFileActionPending = 350
	CodeFileBusy                   = 450
	CodeFileNotFound               = 550
	CodeFileNoAccess               = 550
	CodeLocalErrorProcessing       = 451
	CodePageTypeUnknown            = 551
	CodeInsufficentStorage         = 452
	CodeExceedStorageAllocation    = 552
	CodeBadFileName                = 553
)

func (r *Reply) SendResponse(code int, msg string) {
	msg = fmt.Sprintf("%d %s", code, msg)
	responseLogger := r.Logger.With(slog.String("response_message", msg))

	byteMessage := []byte(msg)
	// RFC 959, commands/responses should be marked as completed with CRLF escape sequence
	matched, err := regexp.Match("\r\n$", byteMessage)

	if !matched && err == nil {
		byteMessage = append(byteMessage, '\r', '\n')
	}

	if err != nil {
		responseLogger.Error(err.Error())
		responseLogger.Error("Unable to send response")
		return
	}

	_, err = r.Conn.Write(byteMessage)

	if err != nil {
		responseLogger.Error(err.Error())
		responseLogger.Error("Unable to send response")
		return
	}

	responseLogger.Info("Sent Succesfully to client")
}

// ToFTPResponseCode returns an ftp reply code based on an error message
func (r *Reply) ToFTPResponseCode(err error) int {

	if errors.Is(err, commandErrors.ErrInvalidArgLength) {
		return CodeSyntaxErrorInArgs
	}

	if errors.Is(err, commandErrors.ErrBadCommand) {
		return CodeUnrecognizedCommand
	}

	return 0
}
