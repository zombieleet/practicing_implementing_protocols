package storage

type Storage interface {
	UserExists(string) bool
	// RecordActiveUserConn records connection that has executed the USER command, first param is the remote address of the client connection, second param is the username preceeding the USER command
	RecordActiveUserConn(string, string)
}
