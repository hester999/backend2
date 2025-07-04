package apperr

import "errors"

// client errors
var (
	ErrClientNotFound = errors.New("client not found")
	ErrInsertFailed   = errors.New("failed to insert")
	ErrUpdateFailed   = errors.New("failed to update")
	ErrDeleteFailed   = errors.New("failed to delete")
)

// product errors
var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductInsert   = errors.New("failed to insert product")
	ErrProductUpdate   = errors.New("failed to update product")
	ErrProductDelete   = errors.New("failed to delete product")
)

// supplier errors
var (
	ErrSupplierNotFound = errors.New("supplier not found")
	ErrSupplierInsert   = errors.New("failed to insert supplier")
	ErrSupplierUpdate   = errors.New("failed to update supplier")
	ErrSupplierDelete   = errors.New("failed to delete supplier")
)

// image err
var (
	ErrImageNotFound = errors.New("image not found")
	ErrImageInsert   = errors.New("failed to insert image")
	ErrImageUpdate   = errors.New("failed to update image")
	ErrImageDelete   = errors.New("failed to delete image")
)
