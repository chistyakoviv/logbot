package rbac

type StorableItemInterface interface {
	GetName() string
	GetAttributes() map[string]any
}

type ItemStorageInterface interface {
	// Removes all roles and permissions.
	Clear()

	// Returns all roles and permissions in the system.
	GetAll() map[string]StorableItemInterface

	// Returns roles and permission by the given names' list.
	GetByNames(names []string) map[string]StorableItemInterface

	// Returns the named role or permission.
	Get(name string) (StorableItemInterface, error)

	// Whether a named role or permission exists.
	Exists(name string) bool

	// Whether a named role exists.
	RoleExists(name string) bool

	// Adds the role or the permission to RBAC system.
	Add(item StorableItemInterface)

	// Updates the specified role or permission in the system.
	Update(name string, item StorableItemInterface)

	// Removes a role or permission from the RBAC system.
	Remove(name string)

	// Returns all roles in the system.
	GetRoles() map[string]*Role

	// Returns roles by the given names' list.
	GetRolesByNames(names []string) map[string]*Role

	// Returns the named role.
	GetRole(name string) (*Role, error)

	// Removes all roles.
	ClearRoles()

	// Returns all permissions in the system.
	GetPermissions() map[string]*Permission

	// Returns permissions by the given names' list.
	GetPermissionsByNames(names []string) map[string]*Permission

	// Returns the named permission.
	GetPermission(name string) (*Permission, error)

	// Removes all permissions.
	ClearPermissions()

	// Returns the parent permissions and/or roles.
	GetParents(name string) map[string]StorableItemInterface

	// Returns the parents tree for a single item which additionally contains children for each parent (only among the
	// found items). The base item is included too, its children list is always empty.
	GetHierarchy(name string) map[string]map[string]map[string]StorableItemInterface

	// Returns direct child permissions and/or roles.
	GetDirectChildren(name string) map[string]StorableItemInterface

	// Returns all child permissions and/or roles.
	GetAllChildren(names []string) map[string]StorableItemInterface

	// Returns all child roles.
	GetAllChildRoles(names []string) map[string]*Role

	// Returns all child permissions.
	GetAllChildPermissions(names []string) map[string]*Permission

	// Returns whether named parent has children.
	HasChildren(name string) bool

	// Returns whether selected parent has a child with a given name.
	HasChild(parentName string, childName string) bool

	// Returns whether selected parent has a direct child with a given name.
	HasDirectChild(parentName string, childName string) bool

	// Adds a role or a permission as a child of another role or permission.
	AddChild(parentName string, childName string)

	// Removes a child from its parent.
	// Note, the child role or permission is not deleted. Only the parent-child relationship is removed.
	RemoveChild(parentName string, childName string)

	// Removes all children form their parent.
	// Note, the children roles or permissions are not deleted. Only the parent-child relationships are removed.
	RemoveChildren(parentName string)
}
