package domain

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordTooWeak    = errors.New("password must be at least 8 chars and contain upper, lower, number and special")
)