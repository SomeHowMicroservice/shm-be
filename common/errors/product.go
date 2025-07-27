package errors

import "errors"

var (
	ErrSlugAlreadyExists = errors.New("slug đã tồn tại")

	ErrCategoryNotFound = errors.New("không tìm thấy danh mục sản phẩm")

	ErrHasCategoryNotFound = errors.New("có danh mục sản phẩm không tìm thấy")

	ErrProductNotFound = errors.New("không tìm thấy sản phẩm")

	ErrColorAlreadyExists = errors.New("màu sắc đã tồn tại")

	ErrSizeAlreadyExists = errors.New("size đã tồn tại")
)