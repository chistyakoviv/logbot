package rbac

import (
	"slices"
	"testing"
	"time"
)

type allowRule struct{}

func (allowRule) Execute(userId any, item ItemInterface, context RuleContext) bool {
	return true
}

type denyByContextRule struct{}

func (denyByContextRule) Execute(userId any, item ItemInterface, context RuleContext) bool {
	value, _ := context.GetParameterValue("allow").(bool)
	return value
}

type denyRule struct{}

func (denyRule) Execute(userId any, item ItemInterface, context RuleContext) bool {
	return false
}

func TestManagerAssignAndResolveUserItems(t *testing.T) {
	manager := NewManager[string](
		NewRuleFactory(),
		NewItemsStorageInMemory(),
		NewAssignmentsStorageInMemory[string](),
		nil,
	)

	if err := manager.AddRole(NewRole("member")); err != nil {
		t.Fatalf("add role: %v", err)
	}
	if err := manager.AddRole(NewRole("admin")); err != nil {
		t.Fatalf("add role: %v", err)
	}
	if err := manager.AddPermission(NewPermission("read")); err != nil {
		t.Fatalf("add permission: %v", err)
	}
	if err := manager.AddPermission(NewPermission("write")); err != nil {
		t.Fatalf("add permission: %v", err)
	}
	if err := manager.AddChild("admin", "read"); err != nil {
		t.Fatalf("add child: %v", err)
	}
	if err := manager.AddChild("admin", "write"); err != nil {
		t.Fatalf("add child: %v", err)
	}

	manager.SetDefaultRoleNames([]string{"member"})

	createdAt := time.Date(2026, time.May, 27, 10, 0, 0, 0, time.UTC)
	if err := manager.Assign("u1", "admin", createdAt); err != nil {
		t.Fatalf("assign: %v", err)
	}

	items := manager.GetItemsByUserId("u1")
	assertItemNamesMatch(t, items, []string{"member", "admin", "read", "write"})

	roles := manager.GetRolesByUserId("u1")
	assertItemNamesMatch(t, roles, []string{"member", "admin"})

	permissions := manager.GetPermissionsByUserId("u1")
	assertItemNamesMatch(t, permissions, []string{"read", "write"})
}

func TestManagerUserHasPermissionSupportsGuestAndRules(t *testing.T) {
	ruleFactory := NewRuleFactory().
		Add("allow", func() RuleInterface { return allowRule{} }).
		Add("gated", func() RuleInterface { return denyByContextRule{} })

	manager := NewManager[string](
		ruleFactory,
		NewItemsStorageInMemory(),
		NewAssignmentsStorageInMemory[string](),
		&ManagerOpts{includeRolesInAccessChecks: true},
	)

	guest := NewRole("guest")
	member := NewRole("member")
	read := NewPermission("read").WithRuleName("allow")
	publish := NewPermission("publish").WithRuleName("gated")

	if err := manager.AddRole(guest); err != nil {
		t.Fatalf("add guest: %v", err)
	}
	if err := manager.AddRole(member); err != nil {
		t.Fatalf("add member: %v", err)
	}
	if err := manager.AddPermission(read); err != nil {
		t.Fatalf("add read: %v", err)
	}
	if err := manager.AddPermission(publish); err != nil {
		t.Fatalf("add publish: %v", err)
	}
	if err := manager.AddChild("guest", "read"); err != nil {
		t.Fatalf("guest child: %v", err)
	}
	if err := manager.AddChild("member", "publish"); err != nil {
		t.Fatalf("member child: %v", err)
	}
	if err := manager.Assign("u1", "member", time.Time{}); err != nil {
		t.Fatalf("assign member: %v", err)
	}

	manager.SetGuestRoleName("guest")

	if !manager.UserHasPermission("", "read", nil) {
		t.Fatalf("expected guest user to inherit guest permissions")
	}
	if manager.UserHasPermission("u1", "publish", RuleContextParameters{"allow": false}) {
		t.Fatalf("expected rule-gated permission to be denied")
	}
	if !manager.UserHasPermission("u1", "publish", RuleContextParameters{"allow": true}) {
		t.Fatalf("expected rule-gated permission to be allowed")
	}
}

func TestManagerUserHasPermissionAllowsAlternativeParentBranchWhenRuleFails(t *testing.T) {
	ruleFactory := NewRuleFactory().
		Add("allow", func() RuleInterface { return allowRule{} }).
		Add("deny", func() RuleInterface { return denyRule{} })

	manager := NewManager[string](
		ruleFactory,
		NewItemsStorageInMemory(),
		NewAssignmentsStorageInMemory[string](),
		nil,
	)

	for _, role := range []ItemInterface{
		NewRole("root"),
		NewRole("admin").WithRuleName("deny"),
		NewRole("auditor").WithRuleName("allow"),
	} {
		if err := manager.AddRole(role); err != nil {
			t.Fatalf("add role %q: %v", role.GetName(), err)
		}
	}

	if err := manager.AddPermission(NewPermission("read")); err != nil {
		t.Fatalf("add permission: %v", err)
	}

	for _, edge := range [][2]string{
		{"root", "admin"},
		{"root", "auditor"},
		{"admin", "read"},
		{"auditor", "read"},
	} {
		if err := manager.AddChild(edge[0], edge[1]); err != nil {
			t.Fatalf("add child %q -> %q: %v", edge[0], edge[1], err)
		}
	}

	if err := manager.Assign("u1", "root", time.Time{}); err != nil {
		t.Fatalf("assign root: %v", err)
	}

	if !manager.UserHasPermission("u1", "read", nil) {
		t.Fatalf("expected access to be granted via the alternative valid parent branch")
	}
}

func TestManagerRejectsInvalidChildAndDirectPermissionAssignment(t *testing.T) {
	manager := NewManager[string](
		NewRuleFactory(),
		NewItemsStorageInMemory(),
		NewAssignmentsStorageInMemory[string](),
		nil,
	)

	if err := manager.AddRole(NewRole("admin")); err != nil {
		t.Fatalf("add role: %v", err)
	}
	if err := manager.AddPermission(NewPermission("read")); err != nil {
		t.Fatalf("add permission: %v", err)
	}

	if err := manager.AddChild("read", "admin"); err == nil {
		t.Fatalf("expected permission -> role child relation to be rejected")
	}
	if err := manager.Assign("u1", "read", time.Time{}); err == nil {
		t.Fatalf("expected direct permission assignment to be rejected")
	}
}

func TestManagerUpdateRoleRenamesAssignmentsAndHierarchy(t *testing.T) {
	manager := NewManager[string](
		NewRuleFactory(),
		NewItemsStorageInMemory(),
		NewAssignmentsStorageInMemory[string](),
		&ManagerOpts{includeRolesInAccessChecks: true},
	)

	if err := manager.AddRole(NewRole("editor")); err != nil {
		t.Fatalf("add role: %v", err)
	}
	if err := manager.AddPermission(NewPermission("publish")); err != nil {
		t.Fatalf("add permission: %v", err)
	}
	if err := manager.AddChild("editor", "publish"); err != nil {
		t.Fatalf("add child: %v", err)
	}
	if err := manager.Assign("u1", "editor", time.Time{}); err != nil {
		t.Fatalf("assign role: %v", err)
	}

	updatedRole := NewRole("author").WithDescription("renamed role")
	if err := manager.UpdateRole("editor", updatedRole); err != nil {
		t.Fatalf("update role: %v", err)
	}

	items := manager.GetItemsByUserId("u1")
	assertItemNamesMatch(t, items, []string{"author", "publish"})

	if !manager.UserHasPermission("u1", "publish", nil) {
		t.Fatalf("expected renamed role assignment to keep permission access")
	}

	userIDs := manager.GetUserIdsByRoleName("author")
	if len(userIDs) != 1 || userIDs[0] != "u1" {
		t.Fatalf("unexpected user ids after rename: got %v want [u1]", userIDs)
	}
}

func TestManagerRemoveChildAndRevokeAccess(t *testing.T) {
	manager := NewManager[string](
		NewRuleFactory(),
		NewItemsStorageInMemory(),
		NewAssignmentsStorageInMemory[string](),
		nil,
	)

	if err := manager.AddRole(NewRole("member")); err != nil {
		t.Fatalf("add member: %v", err)
	}
	if err := manager.AddRole(NewRole("admin")); err != nil {
		t.Fatalf("add admin: %v", err)
	}
	if err := manager.AddPermission(NewPermission("read")); err != nil {
		t.Fatalf("add read: %v", err)
	}
	if err := manager.AddPermission(NewPermission("write")); err != nil {
		t.Fatalf("add write: %v", err)
	}
	if err := manager.AddChild("member", "read"); err != nil {
		t.Fatalf("member child: %v", err)
	}
	if err := manager.AddChild("admin", "member"); err != nil {
		t.Fatalf("admin member child: %v", err)
	}
	if err := manager.AddChild("admin", "write"); err != nil {
		t.Fatalf("admin write child: %v", err)
	}
	if err := manager.Assign("u1", "admin", time.Time{}); err != nil {
		t.Fatalf("assign admin: %v", err)
	}

	if !manager.HasChild("admin", "member") || !manager.HasChildren("admin") {
		t.Fatalf("expected admin to have children before removal")
	}

	manager.RemoveChild("admin", "write")
	if manager.UserHasPermission("u1", "write", nil) {
		t.Fatalf("expected write permission to be removed with direct child")
	}

	manager.RemoveChildren("admin")
	if manager.HasChildren("admin") {
		t.Fatalf("expected admin children to be removed")
	}
	if manager.UserHasPermission("u1", "read", nil) {
		t.Fatalf("expected inherited read permission to be removed")
	}

	manager.Revoke("u1", "admin")
	items := manager.GetItemsByUserId("u1")
	assertItemNamesMatch(t, items, nil)
}

func TestManagerDefaultAndGuestRoleValidation(t *testing.T) {
	manager := NewManager[string](
		NewRuleFactory(),
		NewItemsStorageInMemory(),
		NewAssignmentsStorageInMemory[string](),
		nil,
	)

	if _, err := manager.GetGuestRole(); err != ErrNoGuestUser {
		t.Fatalf("expected ErrNoGuestUser, got %v", err)
	}

	manager.SetGuestRoleName("guest")
	if _, err := manager.GetGuestRole(); err != ErrGuestRoleNameNotExist {
		t.Fatalf("expected ErrGuestRoleNameNotExist, got %v", err)
	}

	if err := manager.AddRole(NewRole("guest")); err != nil {
		t.Fatalf("add guest: %v", err)
	}
	if err := manager.AddRole(NewRole("member")); err != nil {
		t.Fatalf("add member: %v", err)
	}

	manager.SetDefaultRoleNames([]string{"member"})
	roleNames := manager.GetDefaultRoleNames()
	roleNames[0] = "changed"

	roleNamesAfterMutation := manager.GetDefaultRoleNames()
	if len(roleNamesAfterMutation) != 1 || roleNamesAfterMutation[0] != "member" {
		t.Fatalf("expected default role names copy to be isolated, got %v", roleNamesAfterMutation)
	}

	manager.SetDefaultRoleNames([]string{"member", "missing"})
	if _, err := manager.GetDefaultRoles(); err == nil {
		t.Fatalf("expected missing default role to return an error")
	}
}

func TestManagerUpdatePermissionRenamesAssignmentAndPreservesDirectAccess(t *testing.T) {
	manager := NewManager[string](
		NewRuleFactory(),
		NewItemsStorageInMemory(),
		NewAssignmentsStorageInMemory[string](),
		&ManagerOpts{enableDirectPermissions: true},
	)

	if err := manager.AddPermission(NewPermission("publish")); err != nil {
		t.Fatalf("add permission: %v", err)
	}
	if err := manager.Assign("u1", "publish", time.Time{}); err != nil {
		t.Fatalf("assign permission: %v", err)
	}

	if err := manager.UpdatePermission("publish", NewPermission("release")); err != nil {
		t.Fatalf("update permission: %v", err)
	}

	permissions := manager.GetPermissionsByUserId("u1")
	assertItemNamesMatch(t, permissions, []string{"release"})

	if !manager.UserHasPermission("u1", "release", nil) {
		t.Fatalf("expected renamed permission assignment to preserve access")
	}
	if manager.UserHasPermission("u1", "publish", nil) {
		t.Fatalf("expected old permission name to no longer grant access")
	}
}

func assertItemNamesMatch(t *testing.T, items []ItemInterface, expected []string) {
	t.Helper()
	got := make([]string, 0, len(items))
	for _, item := range items {
		if item == nil {
			t.Fatalf("expected item to be non-nil")
		}
		got = append(got, item.GetName())
	}
	slices.Sort(got)
	slices.Sort(expected)
	if !slices.Equal(got, expected) {
		t.Fatalf("unexpected item names: got %v want %v", got, expected)
	}
}

func assertMapItemNamesMatch(t *testing.T, items map[string]ItemInterface, expected []string) {
	t.Helper()
	got := make([]string, 0, len(items))
	for key, item := range items {
		if item == nil {
			t.Fatalf("expected map item %q to be non-nil", key)
		}
		got = append(got, item.GetName())
	}
	slices.Sort(got)
	slices.Sort(expected)
	if !slices.Equal(got, expected) {
		t.Fatalf("unexpected item names: got %v want %v", got, expected)
	}
}
