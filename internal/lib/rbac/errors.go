package rbac

import "errors"

var (
	ErrItemAlreadyExists     = errors.New("item already exists")
	ErrItemNotFound          = errors.New("item not found")
	ErrNoGuestUser           = errors.New("no guest user")
	ErrGuestRoleNameNotExist = errors.New("guest role name does not exist")

	ErrWrongItem                = errors.New("wrong item")
	ErrDefaultRolesNotFound     = errors.New("default roles were not found")
	ErrItemModificationConflict = errors.New("item modification conflict")
	ErrChildAssertFailed        = errors.New("child assert failed")
	ErrAssignForbidden          = errors.New("assign forbidden")
)
