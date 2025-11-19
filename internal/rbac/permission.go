package rbac

type Permission struct {
	item Item
}

func NewPermission(name string) *Permission {
	return &Permission{
		item: *NewItem(name),
	}
}

func (p *Permission) GetType() string {
	return TypePermission
}
