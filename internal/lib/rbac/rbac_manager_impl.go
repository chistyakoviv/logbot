package rbac

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/chistyakoviv/logbot/internal/utils"
)

type manager struct {
	ruleFactory RuleFactoryInterface

	defaultRoleNames []string
	guestRoleName    string

	itemsStorage       ItemsStorageInterface
	assignmentsStorage AssignmentsStorageInterface

	enableDirectPermissions    bool
	includeRolesInAccessChecks bool

	// All getters use RLock and all setters use Lock,
	// so make sure getters do not call setters
	// and setters do not call getters
	// to avoid deadlocks
	mu sync.RWMutex
}

type ManagerOpts struct {
	enableDirectPermissions    bool
	includeRolesInAccessChecks bool
}

func NewManager(
	ruleFactory RuleFactoryInterface,
	itemsStorage ItemsStorageInterface,
	assignmentsStorage AssignmentsStorageInterface,
	opts *ManagerOpts,
) ManagerInterface {
	rbac := &manager{
		ruleFactory:        ruleFactory,
		itemsStorage:       itemsStorage,
		assignmentsStorage: assignmentsStorage,
	}

	if opts != nil {
		rbac.enableDirectPermissions = opts.enableDirectPermissions
		rbac.includeRolesInAccessChecks = opts.includeRolesInAccessChecks
	}

	return rbac
}

func (r *manager) UserHasPermission(
	userId any,
	permissionName string,
	parameters RuleContextParameters,
) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, err := r.itemsStorage.Get(permissionName)
	if err != nil {
		return false
	}

	if !r.includeRolesInAccessChecks && IsRole(item) {
		return false
	}

	var guestRole ItemInterface = nil
	var guestRoleErr error
	if utils.IsEmpty(userId) {
		guestRole, guestRoleErr = r.GetGuestRole()
		if guestRoleErr != nil {
			return false
		}
	}

	hierarchy := r.itemsStorage.GetHierarchy(item.GetName())
	itemNames := make([]string, 0, len(hierarchy))
	for _, treeNode := range hierarchy {
		itemNames = append(itemNames, treeNode.Item.GetName())
	}
	userItemNames := make([]string, 0)
	if guestRole != nil {
		userItemNames = append(userItemNames, guestRole.GetName())
	} else {
		userItemNames = r.filterUserItemNames(userId, itemNames)
	}
	userItemNamesMap := make(map[string]bool, len(userItemNames))
	for _, userItemName := range userItemNames {
		userItemNamesMap[userItemName] = true
	}

	for _, data := range hierarchy {
		if !userItemNamesMap[data.Item.GetName()] || !r.executeRule(userId, data.Item, parameters) {
			continue
		}

		hasPermission := true
		for _, child := range data.Children {
			if !r.executeRule(userId, child, parameters) {
				hasPermission = false
				break
			}
		}

		if hasPermission {
			return true
		}
	}

	return false
}

func (r *manager) CanAddChild(parentName string, childName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.assertFutureChild(parentName, childName) == nil
}

func (r *manager) AddChild(parentName string, childName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	err := r.assertFutureChild(parentName, childName)
	if err != nil {
		return err
	}

	r.itemsStorage.AddChild(parentName, childName)
	return nil
}

func (r *manager) RemoveChild(parentName string, childName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.itemsStorage.RemoveChild(parentName, childName)
}

func (r *manager) RemoveChildren(parentName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.itemsStorage.RemoveChildren(parentName)
}

func (r *manager) HasChild(parentName string, childName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.itemsStorage.HasDirectChild(parentName, childName)
}

func (r *manager) HasChildren(parentName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.itemsStorage.HasChildren(parentName)
}

func (r *manager) Assign(userId any, itemName string, createdAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, err := r.itemsStorage.Get(itemName)
	if err != nil {
		return err
	}

	if !r.enableDirectPermissions && IsPermission(item) {
		return fmt.Errorf("assigning permissions directly is disabled. Prefer assigning roles only")
	}

	if r.assignmentsStorage.Exists(userId, itemName) {
		return nil
	}

	timeCretedAt := time.Now()
	if !createdAt.IsZero() {
		timeCretedAt = createdAt
	}
	assignment := NewAssignment(userId, itemName, timeCretedAt)
	r.assignmentsStorage.Add(assignment)
	return nil
}

func (r *manager) Revoke(userId any, itemName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.assignmentsStorage.Remove(userId, itemName)
}

func (r *manager) RevokeAll(userId any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.assignmentsStorage.RemoveByUserId(userId)
}

func (r *manager) GetItemsByUserId(userId any) (map[string]ItemInterface, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	assignments := r.assignmentsStorage.GetByUserId(userId)

	assignmentNames := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		assignmentNames = append(assignmentNames, assignment.GetItemName())
	}

	defaultRoles, err := r.GetDefaultRoles()
	if err != nil {
		return nil, err
	}

	itemsByNames := r.itemsStorage.GetByNames(assignmentNames)
	children := r.itemsStorage.GetAllChildren(assignmentNames)
	result := make(map[string]ItemInterface, len(defaultRoles)+len(itemsByNames)+len(children))
	for _, item := range defaultRoles {
		result[item.GetName()] = item
	}
	for _, item := range itemsByNames {
		result[item.GetName()] = item
	}
	for _, item := range children {
		result[item.GetName()] = item
	}

	return result, nil
}

func (r *manager) GetRolesByUserId(userId any) (map[string]ItemInterface, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	assignments := r.assignmentsStorage.GetByUserId(userId)

	assignmentNames := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		assignmentNames = append(assignmentNames, assignment.GetItemName())
	}

	defaultRoles, err := r.GetDefaultRoles()
	if err != nil {
		return nil, err
	}

	rolesByNames := r.itemsStorage.GetRolesByNames(assignmentNames)
	childRoles := r.itemsStorage.GetAllChildRoles(assignmentNames)
	result := make(map[string]ItemInterface, len(defaultRoles)+len(rolesByNames)+len(childRoles))
	for _, role := range defaultRoles {
		result[role.GetName()] = role
	}
	for _, role := range rolesByNames {
		result[role.GetName()] = role
	}
	for _, role := range childRoles {
		result[role.GetName()] = role
	}

	return result, nil
}

func (r *manager) GetChildRoles(roleName string) (map[string]ItemInterface, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if !r.itemsStorage.RoleExists(roleName) {
		return nil, fmt.Errorf("Role %s not found", roleName)
	}

	return r.itemsStorage.GetAllChildRoles([]string{roleName}), nil
}

func (r *manager) GetPermissionsByRoleName(name string) map[string]ItemInterface {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.itemsStorage.GetAllChildPermissions([]string{name})
}

func (r *manager) GetPermissionsByUserId(userId any) map[string]ItemInterface {
	r.mu.RLock()
	defer r.mu.RUnlock()
	assignments := r.assignmentsStorage.GetByUserId(userId)
	if len(assignments) == 0 {
		return make(map[string]ItemInterface)
	}

	assignmentNames := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		assignmentNames = append(assignmentNames, assignment.GetItemName())
	}

	permissionsByNames := r.itemsStorage.GetPermissionsByNames(assignmentNames)
	childPermissions := r.itemsStorage.GetAllChildPermissions(assignmentNames)

	result := make(map[string]ItemInterface, len(permissionsByNames)+len(childPermissions))
	for _, permission := range permissionsByNames {
		result[permission.GetName()] = permission
	}
	for _, permission := range childPermissions {
		result[permission.GetName()] = permission
	}

	return result
}

func (r *manager) GetUserIdsByRoleName(roleName string) []any {
	r.mu.RLock()
	defer r.mu.RUnlock()
	parents := r.itemsStorage.GetParents(roleName)
	parentNames := make([]string, 0, len(parents))
	for _, parent := range parents {
		parentNames = append(parentNames, parent.GetName())
	}

	roleNames := make([]string, 0, len(parentNames)+1)
	roleNames = append(roleNames, roleName)
	roleNames = append(roleNames, parentNames...)

	assignments := r.assignmentsStorage.GetByItemNames(roleNames)
	userIds := make([]any, 0, len(assignments))
	for _, assignment := range assignments {
		userIds = append(userIds, assignment.GetUserId())
	}

	return userIds
}

func (r *manager) AddRole(role ItemInterface) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.addItem(role)
}

func (r *manager) GetRole(name string) (ItemInterface, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.itemsStorage.GetRole(name)
}

func (r *manager) UpdateRole(name string, role ItemInterface) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	err := r.assertItemNameForUpdate(name, role)
	if err != nil {
		return err
	}

	r.itemsStorage.Update(name, role)
	r.assignmentsStorage.RenameItem(name, role.GetName())

	return nil
}

func (r *manager) RemoveRole(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removeItem(name)
}

func (r *manager) AddPermission(permission ItemInterface) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.addItem(permission)
}

func (r *manager) GetPermission(name string) (ItemInterface, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.itemsStorage.GetPermission(name)
}

func (r *manager) UpdatePermission(name string, permission ItemInterface) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	err := r.assertItemNameForUpdate(name, permission)
	if err != nil {
		return err
	}

	r.itemsStorage.Update(name, permission)
	r.assignmentsStorage.RenameItem(name, permission.GetName())

	return nil
}

func (r *manager) RemovePermission(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removeItem(name)
}

func (r *manager) SetDefaultRoleNames(roleNames []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Copy the original slice to avoid modifying it outside
	roleNamesCopy := make([]string, len(roleNames))
	copy(roleNamesCopy, roleNames)
	r.defaultRoleNames = roleNamesCopy
}

func (r *manager) GetDefaultRoleNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Copy the original slice to avoid modifying it outside
	roleNames := make([]string, len(r.defaultRoleNames))
	copy(roleNames, r.defaultRoleNames)
	return roleNames
}

func (r *manager) SetGuestRoleName(guestRoleName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.guestRoleName = guestRoleName
}

func (r *manager) GetGuestRoleName() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.guestRoleName
}

func (r *manager) GetDefaultRoles() (map[string]ItemInterface, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.filterStoredRoles(r.defaultRoleNames)
}

func (r *manager) GetGuestRole() (ItemInterface, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.guestRoleName == "" {
		return nil, ErrNoGuestUser
	}

	// Do not use r.GetRole(r.guestRoleName) to avoid second RLock
	role, err := r.itemsStorage.GetRole(r.guestRoleName)
	if err != nil {
		return nil, ErrGuestRoleNameNotExist
	}

	return role, nil
}

func (r *manager) removeItem(name string) {
	if r.itemsStorage.Exists(name) {
		r.itemsStorage.Remove(name)
		r.assignmentsStorage.RemoveByItemName(name)
	}
}

func (r *manager) assertItemNameForUpdate(name string, item ItemInterface) error {
	if name == item.GetName() || !r.itemsStorage.Exists(item.GetName()) {
		return nil
	}

	return fmt.Errorf(
		"unable to change the role or the permission name, the name %s is already used by another role or permission",
		item.GetName(),
	)
}

func (r *manager) addItem(item ItemInterface) error {
	if r.itemsStorage.Exists(item.GetName()) {
		return fmt.Errorf("Item %s already exists", item.GetName())
	}

	time := time.Now()
	if !item.HasCreatedAt() {
		item = item.WithCreatedAt(time)
	}
	if !item.HasUpdatedAt() {
		item = item.WithUpdatedAt(time)
	}

	r.itemsStorage.Add(item)

	return nil
}

func (r *manager) filterStoredRoles(roleNames []string) (map[string]ItemInterface, error) {
	storedRoles := r.itemsStorage.GetRolesByNames(roleNames)
	missingRoles := make([]string, 0)
	for _, roleName := range roleNames {
		if _, ok := storedRoles[roleName]; !ok {
			missingRoles = append(missingRoles, roleName)
		}
	}

	if len(missingRoles) > 0 {
		return nil, fmt.Errorf("the following default roles were not found: %s", strings.Join(missingRoles, ", "))
	}

	return storedRoles, nil
}

func (r *manager) assertFutureChild(parentName string, childName string) error {
	if parentName == childName {
		return fmt.Errorf("cannot add %s as a child of itself", parentName)
	}

	parent, err := r.itemsStorage.Get(parentName)
	if err != nil {
		return fmt.Errorf("parent %s does not exist", parentName)
	}

	child, err := r.itemsStorage.Get(childName)
	if err != nil {
		return fmt.Errorf("child %s does not exist", childName)
	}

	if IsPermission(parent) && IsRole(child) {
		return fmt.Errorf("can not add %s role as a child of %s permission", childName, parentName)
	}

	if r.itemsStorage.HasDirectChild(parentName, childName) {
		return fmt.Errorf("the item %s already has a child %s", parentName, childName)
	}

	if r.itemsStorage.HasChild(parentName, childName) {
		return fmt.Errorf("cannot add %s as a child of %s. A loop has been detected", childName, parentName)
	}

	return nil
}

func (r *manager) executeRule(userId any, item ItemInterface, parameters RuleContextParameters) bool {
	if item.GetRuleName() == "" {
		return true
	}
	return r.ruleFactory.
		Create(item.GetRuleName()).
		Execute(userId, item, NewRuleContext(r.ruleFactory, parameters))
}

func (r *manager) filterUserItemNames(userId any, itemNames []string) []string {
	names := r.assignmentsStorage.FilterUserItemNames(userId, itemNames)
	for _, roleName := range r.defaultRoleNames {
		if slices.Contains(names, roleName) {
			continue
		}
		names = append(names, roleName)
	}

	return names
}
