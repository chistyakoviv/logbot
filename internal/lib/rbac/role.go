package rbac

type Role struct{}

func NewRole(name string) ItemInterface {
	return &Item[Role]{
		name: name,
	}
}
