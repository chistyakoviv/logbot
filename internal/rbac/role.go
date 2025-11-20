package rbac

import "time"

type Role struct {
	*Item
}

func NewRole(name string) ItemInterface {
	return &Role{
		Item: &Item{
			name: name,
		},
	}
}

// Implement all methods that have value receiver,
// to keep the type of the item. Otherwise the type will be
// converted to Item after calling such methods.
func (r Role) WithName(name string) ItemInterface {
	r.name = name
	return &r
}

func (r Role) WithDescription(description string) ItemInterface {
	r.description = description
	return &r
}

func (r Role) WithRuleName(ruleName string) ItemInterface {
	r.ruleName = ruleName
	return &r
}

func (r Role) WithUpdatedAt(updatedAt time.Time) ItemInterface {
	r.updatedAt = updatedAt
	return &r
}

func (r Role) WithCreatedAt(createdAt time.Time) ItemInterface {
	r.createdAt = createdAt
	return &r
}

func (r *Role) GetType() string {
	return TypeRole
}
