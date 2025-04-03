package storage

import "errors"

var (
	ErrAdminExists = errors.New("admin already exists")
	ErrAdminNotFound = errors.New("admin not found")
)