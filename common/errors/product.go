package errors

import "errors"

var (
	ErrSlugAlreadyExists = errors.New("slug đã tồn tại")

	ErrCategoryNotFound = errors.New("không tìm thấy danh mục sản phẩm")

	ErrHasCategoryNotFound = errors.New("có danh mục sản phẩm không tìm thấy")

	ErrProductNotFound = errors.New("không tìm thấy sản phẩm")

	ErrColorAlreadyExists = errors.New("màu sắc đã tồn tại")

	ErrSizeAlreadyExists = errors.New("size đã tồn tại")

	ErrSKUAlreadyExists = errors.New("SKU đã tồn tại")

	ErrColorNotFound = errors.New("không tìm thấy màu sắc")

	ErrSizeNotFound = errors.New("không tìm thấy kích cỡ")

	ErrUnSupportedFileType = errors.New("định dạng file không được hỗ trợ")

	ErrTagAlreadyExists = errors.New("tag sản phẩm đã tồn tại")
)