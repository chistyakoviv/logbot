package rbac

import "time"

const (
	TypeRole       = "role"
	TypePermission = "permission"
)

type ItemInterface[T any] interface {
	GetType() string
	GetName() string
	WithName(name string) T
	GetDescription() string
	WithDescription(description string) T
	GetRuleName() string
	WithRuleName(ruleName string) T
	GetUpdatedAt() time.Time
	WithUpdatedAt(updatedAt time.Time) T
	GetCreatedAt() time.Time
	WithCreatedAt(createdAt time.Time) T
	HasCreatedAt() bool
	HasUpdatedAt() bool
	GetAttributes() map[string]any
}

type StorableItemInterface interface {
	GetName() string
	GetAttributes() map[string]any
}

// Check if Item implements the ItemInterface interface
var _ ItemInterface[Item] = (*Item)(nil)

type Item struct {
	name        string
	description string
	ruleName    string
	createdAt   time.Time
	updatedAt   time.Time
}

func NewItem(name string) *Item {
	return &Item{
		name: name,
	}
}

func (i *Item) GetType() string {
	return ""
}

func (i *Item) GetName() string {
	return i.name
}

func (i Item) WithName(name string) Item {
	i.name = name
	return i
}

func (i *Item) GetDescription() string {
	return i.description
}

func (i Item) WithDescription(description string) Item {
	i.description = description
	return i
}

func (i *Item) GetRuleName() string {
	return i.ruleName
}

func (i Item) WithRuleName(ruleName string) Item {
	i.ruleName = ruleName
	return i
}

func (i *Item) GetUpdatedAt() time.Time {
	return i.updatedAt
}

func (i Item) WithUpdatedAt(updatedAt time.Time) Item {
	i.updatedAt = updatedAt
	return i
}

func (i *Item) GetCreatedAt() time.Time {
	return i.createdAt
}

func (i Item) WithCreatedAt(createdAt time.Time) Item {
	i.createdAt = createdAt
	return i
}

func (i *Item) HasCreatedAt() bool {
	return !i.createdAt.IsZero()
}

func (i *Item) HasUpdatedAt() bool {
	return !i.updatedAt.IsZero()
}

func (i *Item) GetAttributes() map[string]any {
	return map[string]any{
		"name":        i.name,
		"description": i.description,
		"rule_name":   i.ruleName,
		"created_at":  i.createdAt,
		"updated_at":  i.updatedAt,
	}
}
