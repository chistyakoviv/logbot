package rbac

import "errors"

var (
	ErrItemAlreadyExists     = errors.New("item already exists")
	ErrItemNotFound          = errors.New("item not found")
	ErrNoGuestUser           = errors.New("no guest user")
	ErrGuestRoleNameNotExist = errors.New("guest role name does not exist")
)
