package storage

// Storage to be implemented by storage mechanisms
type Storage interface {
	// UserExists checks if a user exists
	UserExists(string) bool
	// Login the user giving username and password
	Login(string, string) error
}
