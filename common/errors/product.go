package errors

import "errors"

var (
	ErrSlugAlreadyExists = errors.New("slug đã tồn tại")

	ErrCategoryNotFound = errors.New("không tìm thấy danh mục sản phẩm")

	ErrHasCategoryNotFound = errors.New("có danh mục sản phẩm không tìm thấy")
)