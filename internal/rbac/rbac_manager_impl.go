package rbac

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/chistyakoviv/logbot/internal/utils"
)

type rbacManager struct {
	ruleFactory RuleFactoryInterface

	defaultRoleNames []string
	guestRoleName    string

	itemsStorage       ItemsStorageInterface
	assignmentsStorage AssignmentsStorageInterface

	enableDirectPermissions    bool
	includeRolesInAccessChecks bool
}

type RBACManagerOpts struct {
	enableDirectPermissions    bool
	includeRolesInAccessChecks bool
}

func NewRbacManager(
	ruleFactory RuleFactoryInterface,
	itemsStorage ItemsStorageInterface,
	assignmentsStorage AssignmentsStorageInterface,
	opts RBACManagerOpts,
) RBACManagerInterface {
	return &rbacManager{
		ruleFactory:                ruleFactory,
		itemsStorage:               itemsStorage,
		assignmentsStorage:         assignmentsStorage,
		enableDirectPermissions:    opts.enableDirectPermissions,
		includeRolesInAccessChecks: opts.includeRolesInAccessChecks,
	}
}

func (r *rbacManager) UserHasPermission(
	userId any,
	permissionName string,
	parameters RuleContextParameters,
) bool {
	item, err := r.itemsStorage.Get(permissionName)
	if err != nil {
		return false
	}

	if !r.includeRolesInAccessChecks && IsRole(item) {
		return false
	}

	var guestRole *Role = nil
	var guestRoleErr error
	if utils.IsEmpty(userId) {
		guestRole, guestRoleErr = r.GetGuestRole()
		if guestRoleErr != nil {
			return false
		}
	}

	hierarchy := r.itemsStorage.GetHierarchy(item.GetName())
	itemNames := make([]string, 0)
	for _, treeItem := range hierarchy {
		itemNames = append(itemNames, treeItem.Item.GetName())
	}
	userItemNames := make([]string, 0)
	if guestRole != nil {
		userItemNames = append(userItemNames, guestRole.GetName())
	} else {
		userItemNames = r.filterUserItemNames(userId, itemNames)
	}
	userItemNamesMap := make(map[string]ItemInterface, 0)
	for _, userItemName := range userItemNames {
		userItemNamesMap[userItemName] = nil
	}

	for _, data := range hierarchy {
		_, itemExists := userItemNamesMap[data.Item.GetName()]
		if !itemExists || r.executeRule(userId, data.Item, parameters) {
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

func (r *rbacManager) CanAddChild(parentName string, childName string) bool {
	return r.assertFutureChild(parentName, childName) == nil
}

func (r *rbacManager) AddChild(parentName string, childName string) error {
	err := r.assertFutureChild(parentName, childName)
	if err != nil {
		return err
	}

	r.itemsStorage.AddChild(parentName, childName)
	return nil
}

func (r *rbacManager) RemoveChild(parentName string, childName string) {
	r.itemsStorage.RemoveChild(parentName, childName)
}

func (r *rbacManager) RemoveChildren(parentName string) {
	r.itemsStorage.RemoveChildren(parentName)
}

func (r *rbacManager) HasChild(parentName string, childName string) bool {
	return r.itemsStorage.HasDirectChild(parentName, childName)
}

func (r *rbacManager) HasChildren(parentName string) bool {
	return r.itemsStorage.HasChildren(parentName)
}

func (r *rbacManager) Assign(userId any, itemName string, createdAt time.Time) error {
	item, err := r.itemsStorage.Get(itemName)
	if err != nil {
		return err
	}

	if !r.enableDirectPermissions && IsPermission(item) {
		return fmt.Errorf("Assigning permissions directly is disabled. Prefer assigning roles only.")
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

func (r *rbacManager) Revoke(userId any, itemName string) {
	r.assignmentsStorage.Remove(userId, itemName)
}

func (r *rbacManager) RevokeAll(userId any) {
	r.assignmentsStorage.RemoveByUserId(userId)
}

func (r *rbacManager) GetItemsByUserId(userId any) (map[string]ItemInterface, error) {
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
	result := make(map[string]ItemInterface, len(rolesByNames)+len(defaultRoles)+len(childRoles))
	for _, role := range r.itemsStorage.GetRolesByNames(assignmentNames) {
		result[role.GetName()] = role
	}
	for _, role := range defaultRoles {
		result[role.GetName()] = role
	}
	for _, role := range childRoles {
		result[role.GetName()] = role
	}

	return result, nil
}

func (r *rbacManager) GetDefaultRoles() (map[string]ItemInterface, error) {
	return r.filterStoredRoles(r.defaultRoleNames)
}

func (r *rbacManager) filterStoredRoles(roleNames []string) (map[string]ItemInterface, error) {
	storedRoles := r.itemsStorage.GetRolesByNames(roleNames)
	missingRoles := make([]string, 0)
	for _, roleName := range roleNames {
		if _, ok := storedRoles[roleName]; !ok {
			missingRoles = append(missingRoles, roleName)
		}
	}

	if len(missingRoles) > 0 {
		return nil, fmt.Errorf("The following default roles were not found: %s", strings.Join(missingRoles, ", "))
	}

	return storedRoles, nil
}

func (r *rbacManager) assertFutureChild(parentName string, childName string) error {
	if parentName == childName {
		return fmt.Errorf("Cannot add %s as a child of itself.", parentName)
	}

	parent, err := r.itemsStorage.Get(parentName)
	if err != nil {
		return fmt.Errorf("Parent %s does not exist.", parentName)
	}

	child, err := r.itemsStorage.Get(childName)
	if err != nil {
		return fmt.Errorf("Child %s does not exist.", childName)
	}

	if IsPermission(parent) && IsRole(child) {
		return fmt.Errorf("Can not add %s role as a child of %s permission.", childName, parentName)
	}

	if r.itemsStorage.HasDirectChild(parentName, childName) {
		return fmt.Errorf("The item %s already has a child %s.", parentName, childName)
	}

	if r.itemsStorage.HasChild(parentName, childName) {
		return fmt.Errorf("Cannot add %s as a child of %s. A loop has been detected.", childName, parentName)
	}

	return nil
}

func (r *rbacManager) executeRule(userId any, item ItemInterface, parameters RuleContextParameters) bool {
	if item.GetRuleName() == "" {
		return true
	}
	return r.ruleFactory.
		Create(item.GetRuleName()).
		Execute(userId, item, NewRuleContext(r.ruleFactory, parameters))
}

func (r *rbacManager) filterUserItemNames(userId any, itemNames []string) []string {
	names := r.assignmentsStorage.FilterUserItemNames(userId, itemNames)
	for _, roleName := range r.defaultRoleNames {
		if slices.Contains(names, roleName) {
			continue
		}
		names = append(names, roleName)
	}

	return names
}

func (r *rbacManager) GetGuestRole() (ItemInterface, error) {
	if r.guestRoleName == "" {
		return nil, ErrNoGuestUser
	}

	role, err := r.GetRole(r.guestRoleName)
	if err != nil {
		return nil, ErrGuestRoleNameNotExist
	}

	return role, nil
}

func (r *rbacManager) GetRole(name string) (ItemInterface, error) {
	return r.itemsStorage.GetRole(name)
}
