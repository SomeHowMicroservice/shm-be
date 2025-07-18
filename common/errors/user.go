package errors

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email đã tồn tại")

	ErrUsernameAlreadyExists = errors.New("username đã tồn tại")

	ErrUserNotFound = errors.New("không tìm thấy người dùng")

	ErrInvalidPassword = errors.New("mật khẩu không chính xác")

	ErrRoleNotFound = errors.New("không tìm thấy quyền")

	ErrProfileNotFound = errors.New("không tìm thấy hồ sơ người dùng")

	ErrMeasurementNotFound = errors.New("không tìm thấy độ đo người dùng")
)