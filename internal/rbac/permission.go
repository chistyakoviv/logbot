package rbac

import "time"

type Permission struct {
	*Item
}

func NewPermission(name string) ItemInterface {
	return &Permission{
		Item: &Item{
			name: name,
		},
	}
}

// Implement all methods that have value receiver,
// to keep the type of the item. Otherwise the type will be
// converted to Item after calling such methods.
func (p Permission) WithName(name string) ItemInterface {
	p.name = name
	return &p
}

func (p Permission) WithDescription(description string) ItemInterface {
	p.description = description
	return &p
}

func (p Permission) WithRuleName(ruleName string) ItemInterface {
	p.ruleName = ruleName
	return &p
}

func (p Permission) WithUpdatedAt(updatedAt time.Time) ItemInterface {
	p.updatedAt = updatedAt
	return &p
}

func (p Permission) WithCreatedAt(createdAt time.Time) ItemInterface {
	p.createdAt = createdAt
	return &p
}

func IsPermission(item ItemInterface) bool {
	_, ok := item.(*Permission)
	return ok
}
