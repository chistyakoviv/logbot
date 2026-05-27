package rbac

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/chistyakoviv/logbot/internal/utils"
)

type manager[T comparable] struct {
	ruleFactory RuleFactoryInterface

	defaultRoleNames []string
	guestRoleName    string

	itemsStorage       ItemsStorageInterface
	assignmentsStorage AssignmentsStorageInterface[T]

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

func NewManager[T comparable](
	ruleFactory RuleFactoryInterface,
	itemsStorage ItemsStorageInterface,
	assignmentsStorage AssignmentsStorageInterface[T],
	opts *ManagerOpts,
) ManagerInterface[T] {
	rbac := &manager[T]{
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

func (r *manager[T]) UserHasPermission(
	userId T,
	permissionName string,
	parameters RuleContextParameters,
) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, err := r.itemsStorage.Get(permissionName)
	if err != nil {
		return false
	}

	if !r.includeRolesInAccessChecks && IsItem[Role](item) {
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

	return r.userHasPermissionViaPath(userId, item, guestRole, parameters, make(map[string]bool))
}

func (r *manager[T]) CanAddChild(parentName string, childName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.assertFutureChild(parentName, childName) == nil
}

func (r *manager[T]) AddChild(parentName string, childName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	err := r.assertFutureChild(parentName, childName)
	if err != nil {
		return err
	}

	r.itemsStorage.AddChild(parentName, childName)
	return nil
}

func (r *manager[T]) RemoveChild(parentName string, childName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.itemsStorage.RemoveChild(parentName, childName)
}

func (r *manager[T]) RemoveChildren(parentName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.itemsStorage.RemoveChildren(parentName)
}

func (r *manager[T]) HasChild(parentName string, childName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.itemsStorage.HasDirectChild(parentName, childName)
}

func (r *manager[T]) HasChildren(parentName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.itemsStorage.HasChildren(parentName)
}

func (r *manager[T]) Assign(userId T, itemName string, createdAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, err := r.itemsStorage.Get(itemName)
	if err != nil {
		return err
	}

	if !r.enableDirectPermissions && IsItem[Permission](item) {
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

func (r *manager[T]) Revoke(userId T, itemName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.assignmentsStorage.Remove(userId, itemName)
}

func (r *manager[T]) RevokeAll(userId T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.assignmentsStorage.RemoveByUserId(userId)
}

func (r *manager[T]) GetItemsByUserId(userId T) ([]ItemInterface, error) {
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
	result := make([]ItemInterface, 0, len(defaultRoles)+len(itemsByNames)+len(children))
	for _, item := range defaultRoles {
		result = append(result, item)
	}
	for _, item := range itemsByNames {
		result = append(result, item)
	}
	for _, item := range children {
		result = append(result, item)
	}

	return result, nil
}

func (r *manager[T]) GetRolesByUserId(userId T) ([]ItemInterface, error) {
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
	result := make([]ItemInterface, 0, len(defaultRoles)+len(rolesByNames)+len(childRoles))
	for _, role := range defaultRoles {
		result = append(result, role)
	}
	for _, role := range rolesByNames {
		result = append(result, role)
	}
	for _, role := range childRoles {
		result = append(result, role)
	}

	return result, nil
}

func (r *manager[T]) GetChildRoles(roleName string) ([]ItemInterface, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if !r.itemsStorage.RoleExists(roleName) {
		return nil, fmt.Errorf("Role %s not found", roleName)
	}

	return r.itemsStorage.GetAllChildRoles([]string{roleName}), nil
}

func (r *manager[T]) GetPermissionsByRoleName(name string) []ItemInterface {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.itemsStorage.GetAllChildPermissions([]string{name})
}

func (r *manager[T]) GetPermissionsByUserId(userId T) []ItemInterface {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]ItemInterface, 0)
	assignments := r.assignmentsStorage.GetByUserId(userId)
	if len(assignments) == 0 {
		return result
	}

	assignmentNames := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		assignmentNames = append(assignmentNames, assignment.GetItemName())
	}

	permissionsByNames := r.itemsStorage.GetPermissionsByNames(assignmentNames)
	childPermissions := r.itemsStorage.GetAllChildPermissions(assignmentNames)

	result = make([]ItemInterface, 0, len(permissionsByNames)+len(childPermissions))
	for _, permission := range permissionsByNames {
		result = append(result, permission)
	}
	for _, permission := range childPermissions {
		result = append(result, permission)
	}

	return result
}

func (r *manager[T]) GetUserIdsByRoleName(roleName string) []T {
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
	userIds := make([]T, 0, len(assignments))
	for _, assignment := range assignments {
		userIds = append(userIds, assignment.GetUserId())
	}

	return userIds
}

func (r *manager[T]) AddRole(role ItemInterface) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.addItem(role)
}

func (r *manager[T]) GetRole(name string) (ItemInterface, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.itemsStorage.GetRole(name)
}

func (r *manager[T]) UpdateRole(name string, role ItemInterface) error {
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

func (r *manager[T]) RemoveRole(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removeItem(name)
}

func (r *manager[T]) AddPermission(permission ItemInterface) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.addItem(permission)
}

func (r *manager[T]) GetPermission(name string) (ItemInterface, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.itemsStorage.GetPermission(name)
}

func (r *manager[T]) UpdatePermission(name string, permission ItemInterface) error {
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

func (r *manager[T]) RemovePermission(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removeItem(name)
}

func (r *manager[T]) SetDefaultRoleNames(roleNames []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Copy the original slice to avoid modifying it outside
	roleNamesCopy := make([]string, len(roleNames))
	copy(roleNamesCopy, roleNames)
	r.defaultRoleNames = roleNamesCopy
}

func (r *manager[T]) GetDefaultRoleNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Copy the original slice to avoid modifying it outside
	roleNames := make([]string, len(r.defaultRoleNames))
	copy(roleNames, r.defaultRoleNames)
	return roleNames
}

func (r *manager[T]) SetGuestRoleName(guestRoleName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.guestRoleName = guestRoleName
}

func (r *manager[T]) GetGuestRoleName() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.guestRoleName
}

func (r *manager[T]) GetDefaultRoles() ([]ItemInterface, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.filterStoredRoles(r.defaultRoleNames)
}

func (r *manager[T]) GetGuestRole() (ItemInterface, error) {
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

func (r *manager[T]) removeItem(name string) {
	if r.itemsStorage.Exists(name) {
		r.itemsStorage.Remove(name)
		r.assignmentsStorage.RemoveByItemName(name)
	}
}

func (r *manager[T]) assertItemNameForUpdate(name string, item ItemInterface) error {
	if name == item.GetName() || !r.itemsStorage.Exists(item.GetName()) {
		return nil
	}

	return fmt.Errorf(
		"unable to change the role or the permission name, the name %s is already used by another role or permission",
		item.GetName(),
	)
}

func (r *manager[T]) addItem(item ItemInterface) error {
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

func (r *manager[T]) filterStoredRoles(roleNames []string) ([]ItemInterface, error) {
	roles := r.itemsStorage.GetRolesByNames(roleNames)

	storedRoles := make(map[string]ItemInterface, len(roles))
	for _, role := range roles {
		storedRoles[role.GetName()] = role
	}

	missingRoles := make([]string, 0)
	for _, roleName := range roleNames {
		if _, ok := storedRoles[roleName]; !ok {
			missingRoles = append(missingRoles, roleName)
		}
	}

	if len(missingRoles) > 0 {
		return nil, fmt.Errorf("the following default roles were not found: %s", strings.Join(missingRoles, ", "))
	}

	return roles, nil
}

func (r *manager[T]) assertFutureChild(parentName string, childName string) error {
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

	if IsItem[Permission](parent) && IsItem[Role](child) {
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

func (r *manager[T]) executeRule(userId T, item ItemInterface, parameters RuleContextParameters) bool {
	if item.GetRuleName() == "" {
		return true
	}
	return r.ruleFactory.
		Create(item.GetRuleName()).
		Execute(userId, item, NewRuleContext(r.ruleFactory, parameters))
}

func (r *manager[T]) userHasPermissionViaPath(
	userId T,
	item ItemInterface,
	guestRole ItemInterface,
	parameters RuleContextParameters,
	visited map[string]bool,
) bool {
	if visited[item.GetName()] {
		return false
	}
	visited[item.GetName()] = true

	if !r.executeRule(userId, item, parameters) {
		return false
	}

	if guestRole != nil && guestRole.GetName() == item.GetName() {
		return true
	}

	if guestRole == nil && r.userHasItem(userId, item.GetName()) {
		return true
	}

	for _, parent := range r.itemsStorage.GetDirectParents(item.GetName()) {
		if r.userHasPermissionViaPath(userId, parent, guestRole, parameters, visited) {
			return true
		}
	}

	return false
}

func (r *manager[T]) userHasItem(userId T, itemName string) bool {
	if r.assignmentsStorage.Exists(userId, itemName) {
		return true
	}

	return slices.Contains(r.defaultRoleNames, itemName)
}

func (r *manager[T]) filterUserItemNames(userId T, itemNames []string) []string {
	names := r.assignmentsStorage.FilterUserItemNames(userId, itemNames)
	for _, roleName := range r.defaultRoleNames {
		if slices.Contains(names, roleName) {
			continue
		}
		names = append(names, roleName)
	}

	return names
}
