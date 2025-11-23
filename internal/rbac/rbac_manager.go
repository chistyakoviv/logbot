package rbac

import (
	"time"
)

type ManagerInterface interface {
	// Checks whether the user has the specified permission.
	UserHasPermission(userId any, permissionName string, parameters RuleContextParameters) bool

	// Checks the possibility of adding a child to a parent.
	CanAddChild(parentName string, childName string) bool

	// Adds an item as a child of another item.
	AddChild(parentName string, childName string) error

	// Removes a child from its parent.
	RemoveChild(parentName string, childName string)

	// Removes all children form their parent.
	RemoveChildren(parentName string)

	// Returns a value indicating whether the child already exists for the parent.
	HasChild(parentName string, childName string) bool

	// Returns whether named parent has children.
	HasChildren(name string) bool

	// Assigns a role or permission to a user.
	Assign(userId any, itemName string, createdAt time.Time) error

	// Revokes a role or a permission from a user.
	Revoke(userId any, itemName string)

	// Revokes all roles and permissions from a user.
	RevokeAll(userId any)

	// Returns the items that are assigned to the user via assign().
	GetItemsByUserId(userId any) (map[string]ItemInterface, error)

	// Returns the roles that are assigned to the user via assign().
	GetRolesByUserId(userId any) (map[string]ItemInterface, error)

	// Returns child roles of the role specified. Depth isn't limited.
	GetChildRoles(name string) (map[string]ItemInterface, error)

	// Returns all permissions that the specified role represents.
	GetPermissionsByRoleName(name string) map[string]ItemInterface

	// Returns all permissions that the user has.
	GetPermissionsByUserId(userId any) map[string]ItemInterface

	// Returns all user IDs assigned to the role specified.
	GetUserIdsByRoleName(name string) []any

	// Adds role or permission to RBAC system.
	// Panics if the permission already exists.
	AddRole(role ItemInterface) error

	// Gets role by name.
	GetRole(name string) (ItemInterface, error)

	// Updates role in RBAC system.
	UpdateRole(name string, role ItemInterface) error

	// Removes role from RBAC system.
	RemoveRole(name string)

	// Adds permission to RBAC system.
	// Panics if the permission already exists.
	AddPermission(permission ItemInterface) error

	// Gets permission by name.
	GetPermission(name string) (ItemInterface, error)

	// Updates permission in RBAC system.
	UpdatePermission(name string, permission ItemInterface) error

	// Removes permission from RBAC system.
	RemovePermission(name string)

	// Sets default role names.
	SetDefaultRoleNames(roleNames []string)

	// Returns default role names.
	GetDefaultRoleNames() []string

	// Returns default roles.
	GetDefaultRoles() (map[string]ItemInterface, error)

	// Set guest role name.
	SetGuestRoleName(roleName string)

	// Get guest role name.
	GetGuestRoleName() string

	// Get a guest role.
	GetGuestRole() (ItemInterface, error)
}
