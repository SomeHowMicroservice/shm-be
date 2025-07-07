package errors

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email đã tồn tại")

	ErrUsernameAlreadyExists = errors.New("username đã tồn tại")

	ErrUserNotFound = errors.New("không tìm thấy người dùng")
)