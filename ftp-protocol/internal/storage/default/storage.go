package default_storage

import (
	"github.com/zombieleet/ftp-protocol/internal/storage"
	storageErrors "github.com/zombieleet/ftp-protocol/internal/storage/errors"
)

type DefaultStorage struct {
	store map[string]string
}

func New() storage.Storage {
	store := make(map[string]string)
	store["victory"] = "password"
	return &DefaultStorage{
		store: store,
	}
}

func (d *DefaultStorage) UserExists(user string) bool {
	if _, ok := d.store[user]; !ok {
		return false
	}
	return true
}

func (d *DefaultStorage) Login(user string, password string) error {
	dbPassword, ok := d.store[user]

	if !ok {
		return storageErrors.ErrUserDoesNotExists
	}

	if dbPassword != password {
		return storageErrors.ErrBadUserNameAndPassword
	}

	return nil
}
