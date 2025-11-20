package rbac

import "errors"

var (
	ErrItemAlreadyExists = errors.New("item already exists")
	ErrItemNotFound      = errors.New("item not found")
)
