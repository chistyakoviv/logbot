package rbac

type Role struct {
	item Item
}

func NewRole(name string) *Role {
	return &Role{
		item: *NewItem(name),
	}
}

func (r *Role) GetType() string {
	return TypeRole
}
