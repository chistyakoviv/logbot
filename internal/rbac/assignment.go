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

func (a *Assignment) GetUserId() any {
	return a.userId
}

func (a *Assignment) GetItemName() string {
	return a.itemName
}

func (a Assignment) WithItemName(itemName string) Assignment {
	a.itemName = itemName
	return a
}

func (a *Assignment) GetCreatedAt() time.Time {
	return a.createdAt
}

func (a *Assignment) GetAttributes() map[string]any {
	return map[string]any{
		"user_id":    a.userId,
		"item_name":  a.itemName,
		"created_at": a.createdAt,
	}
}
