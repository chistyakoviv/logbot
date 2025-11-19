package rbac

import "time"

// `Assignment` represents an assignment of a role or a permission to a user.
type Assignment struct {
	userId    any
	itemName  string
	createdAt time.Time
}

func NewAssignment(userId any, itemName string) *Assignment {
	return &Assignment{
		userId:    userId,
		itemName:  itemName,
		createdAt: time.Now(),
	}
}

func (a *Assignment) UserId() any {
	return a.userId
}

func (a *Assignment) ItemName() string {
	return a.itemName
}

func (a *Assignment) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Assignment) Attributes() map[string]any {
	return map[string]any{
		"user_id":    a.userId,
		"item_name":  a.itemName,
		"created_at": a.createdAt,
	}
}
