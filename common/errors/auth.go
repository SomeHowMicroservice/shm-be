package errors

import "errors"

var (
	ErrAuthDataNotFound = errors.New("không tìm thấy dữ liệu xác thực")

	ErrTooManyAttempts = errors.New("vượt quá số lần thử OTP")

	ErrInvalidOTP = errors.New("mã OTP không chính xác")
)