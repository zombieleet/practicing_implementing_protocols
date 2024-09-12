package default_storage

import (
	"github.com/zombieleet/ftp-protocol/internal/storage"
)

type DefaultStorage struct {
	store      map[string]string
	activeConn map[string]string
}

func New() storage.Storage {
	store := make(map[string]string)
	store["victory"] = "password"
	return &DefaultStorage{
		store:      store,
		activeConn: make(map[string]string),
	}
}

func (d *DefaultStorage) UserExists(user string) bool {
	if _, ok := d.store[user]; !ok {
		return false
	}
	return true
}

func (d *DefaultStorage) RecordActiveUserConn(remoteClientAddr string, user string) {
	d.activeConn[remoteClientAddr] = user
}
