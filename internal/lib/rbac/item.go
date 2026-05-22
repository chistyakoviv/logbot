package rbac

import "time"

type ItemInterface interface {
	GetName() string
	WithName(name string) ItemInterface
	GetDescription() string
	WithDescription(description string) ItemInterface
	GetRuleName() string
	WithRuleName(ruleName string) ItemInterface
	GetUpdatedAt() time.Time
	WithUpdatedAt(updatedAt time.Time) ItemInterface
	GetCreatedAt() time.Time
	WithCreatedAt(createdAt time.Time) ItemInterface
	HasCreatedAt() bool
	HasUpdatedAt() bool
	GetAttributes() map[string]any
	GetType() any
}

type Item[T Permission | Role] struct {
	name        string
	description string
	ruleName    string
	createdAt   time.Time
	updatedAt   time.Time
	kind        T
}

func NewItem[T Permission | Role](name string) ItemInterface {
	return &Item[T]{
		name: name,
	}
}

// Use only value receivers for consistency
func (i Item[T]) GetName() string {
	return i.name
}

func (i Item[T]) WithName(name string) ItemInterface {
	i.name = name
	return &i
}

func (i Item[T]) GetDescription() string {
	return i.description
}

func (i Item[T]) WithDescription(description string) ItemInterface {
	i.description = description
	return &i
}

func (i Item[T]) GetRuleName() string {
	return i.ruleName
}

func (i Item[T]) WithRuleName(ruleName string) ItemInterface {
	i.ruleName = ruleName
	return &i
}

func (i Item[T]) GetUpdatedAt() time.Time {
	return i.updatedAt
}

func (i Item[T]) WithUpdatedAt(updatedAt time.Time) ItemInterface {
	i.updatedAt = updatedAt
	return &i
}

func (i Item[T]) GetCreatedAt() time.Time {
	return i.createdAt
}

func (i Item[T]) WithCreatedAt(createdAt time.Time) ItemInterface {
	i.createdAt = createdAt
	return &i
}

func (i Item[T]) HasCreatedAt() bool {
	return !i.createdAt.IsZero()
}

func (i Item[T]) HasUpdatedAt() bool {
	return !i.updatedAt.IsZero()
}

func (i Item[T]) GetType() any {
	return i.kind
}

func (i Item[T]) GetAttributes() map[string]any {
	return map[string]any{
		"name":        i.name,
		"description": i.description,
		"rule_name":   i.ruleName,
		"created_at":  i.createdAt,
		"updated_at":  i.updatedAt,
	}
}

func IsItem[T Permission | Role](item ItemInterface) bool {
	_, ok := item.GetType().(T)
	return ok
}
