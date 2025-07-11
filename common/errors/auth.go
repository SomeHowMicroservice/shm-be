package errors

import "errors"

var (
	ErrAuthDataNotFound = errors.New("không tìm thấy dữ liệu xác thực")

	ErrTooManyAttempts = errors.New("vượt quá số lần thử OTP")

	ErrInvalidOTP = errors.New("mã OTP không chính xác")

	ErrInvalidToken = errors.New("token không hợp lệ hoặc đã hết hạn")

	ErrUserIdNotFound = errors.New("không tìm thấy user_id")

	ErrRolesNotFound = errors.New("không tìm thấy roles")

	ErrUnAuth = errors.New("bạn chưa đăng nhập")

	ErrForbidden = errors.New("không có quyền truy cập")
)