package rbac

type AssignmentStorageInterface interface {
	// Returns all role and permission assignment information.
	GetAll() map[any]map[string]*Assignment

	// Returns all role or permission assignment information for the specified user.
	GetByUserId(userId any) map[string]*Assignment

	// Returns all role or permission assignment information by the specified item names' list.
	GetByItemNames(itemNames []string) []*Assignment

	// Returns role or permission assignment for the specified item name that belongs to user with the specified ID.
	Get(userId any, itemName string) *Assignment

	// Whether assignment with a given item name and user id pair exists.
	Exists(userId any, itemName string) bool

	// Whether at least one item from the given list is assigned to the user.
	UserHasItem(userId any, itemNames []string) bool

	// Filters item names leaving only the ones that are assigned to specific user.
	FilterUserItemNames(userId any, itemNames []string) []string

	// Adds assignment to the storage.
	Add(assignment *Assignment)

	// Returns whether there is assignment for a named role or permission.
	HasItem(itemName string) bool

	// Change the name of an item in assignments.
	RenameItem(oldName string, newName string)

	// Removes assignment of a role or a permission to the user with ID specified.
	Remove(userId any, itemName string)

	// Removes all role or permission assignments for a user with ID specified.
	RemoveByUserId(userId any)

	// Removes all assignments for role or permission.
	RemoveByItemName(itemName string)

	// Removes all role and permission assignments.
	Clear()
}
