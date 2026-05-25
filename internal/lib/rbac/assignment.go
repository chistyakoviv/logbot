package rbac

import "time"

// `Assignment` represents an assignment of a role or a permission to a user.
type Assignment[T comparable] struct {
	userId    T
	itemName  string
	createdAt time.Time
}

func NewAssignment[T comparable](userId T, itemName string, createdAt time.Time) *Assignment[T] {
	return &Assignment[T]{
		userId:    userId,
		itemName:  itemName,
		createdAt: createdAt,
	}
}

func (a *Assignment[T]) GetUserId() T {
	return a.userId
}

func (a *Assignment[T]) GetItemName() string {
	return a.itemName
}

func (a Assignment[T]) WithItemName(itemName string) Assignment[T] {
	a.itemName = itemName
	return a
}

func (a *Assignment[T]) GetCreatedAt() time.Time {
	return a.createdAt
}

func (a *Assignment[T]) GetAttributes() map[string]any {
	return map[string]any{
		"user_id":    a.userId,
		"item_name":  a.itemName,
		"created_at": a.createdAt,
	}
}
