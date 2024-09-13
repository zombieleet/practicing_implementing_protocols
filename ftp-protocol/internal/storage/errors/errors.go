package storage_errors

import "errors"

var (
	ErrActiveConnUserNotFound = errors.New("No previous USER command has been iniated")
	ErrUserDoesNotExists      = errors.New("User does not exists")
	ErrBadUserNameAndPassword = errors.New("Invalid login details, cannot auth")
)
