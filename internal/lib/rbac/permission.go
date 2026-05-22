package rbac

type Permission struct{}

func NewPermission(name string) ItemInterface {
	return &Item[Permission]{
		name: name,
	}
}
