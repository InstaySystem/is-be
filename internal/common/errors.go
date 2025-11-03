package common

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("username đã tồn tại")
	
	ErrEmailAlreadyExists = errors.New("email đã tồn tại")

	ErrUserNotFound = errors.New("user not found")

	ErrLoginFailed = errors.New("incorrect username or password")

	ErrInvalidToken = errors.New("invalid token")
)