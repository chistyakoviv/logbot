package rbac

import (
	"fmt"
	"testing"
	"time"
)

// Before optimization (1)
// BenchmarkManagerUserHasPermissionDeepHierarchy-16                    762           1511110 ns/op         1465244 B/op       1527 allocs/op
// BenchmarkManagerUserHasPermissionMultiParentHierarchy-16          254254              4542 ns/op            5152 B/op         36 allocs/op
// BenchmarkManagerUserHasPermissionGuestRole-16                       2740            396286 ns/op          384418 B/op        752 allocs/op
// PASS
// ok      github.com/chistyakoviv/logbot/internal/lib/rbac        3.654s

// With building direct parents (2)
// BenchmarkManagerUserHasPermissionDeepHierarchy-16                   3620            322262 ns/op           15304 B/op        139 allocs/op
// BenchmarkManagerUserHasPermissionMultiParentHierarchy-16         1680080               723.2 ns/op            78 B/op          3 allocs/op
// BenchmarkManagerUserHasPermissionGuestRole-16                      13634             92673 ns/op            7720 B/op         73 allocs/op
// PASS
// ok      github.com/chistyakoviv/logbot/internal/lib/rbac        5.303s

// Witch cached direct parents (3)
// BenchmarkManagerUserHasPermissionDeepHierarchy-16                  35102             34431 ns/op           15304 B/op        139 allocs/op
// BenchmarkManagerUserHasPermissionMultiParentHierarchy-16         2295321               530.0 ns/op            62 B/op          2 allocs/op
// BenchmarkManagerUserHasPermissionGuestRole-16                      73736             16209 ns/op            7720 B/op         73 allocs/op
// PASS
// ok      github.com/chistyakoviv/logbot/internal/lib/rbac        4.673s

func BenchmarkManagerUserHasPermissionDeepHierarchy(b *testing.B) {
	manager, permissionName := buildBenchmarkManagerWithRoleChain(b, 128)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if !manager.UserHasPermission("u1", permissionName, nil) {
			b.Fatalf("expected permission %q to be granted", permissionName)
		}
	}
}

func BenchmarkManagerUserHasPermissionMultiParentHierarchy(b *testing.B) {
	manager := NewManager[string](
		NewRuleFactory(),
		NewItemsStorageInMemory(),
		NewAssignmentsStorageInMemory[string](),
		nil,
	)

	root := NewRole("root")
	admin := NewRole("admin")
	auditor := NewRole("auditor")
	editor := NewRole("editor")
	read := NewPermission("read")

	for _, item := range []ItemInterface{root, admin, auditor, editor, read} {
		switch item.GetType().(type) {
		case Role:
			if err := manager.AddRole(item); err != nil {
				b.Fatalf("add role %q: %v", item.GetName(), err)
			}
		case Permission:
			if err := manager.AddPermission(item); err != nil {
				b.Fatalf("add permission %q: %v", item.GetName(), err)
			}
		}
	}

	for _, edge := range [][2]string{
		{"root", "admin"},
		{"root", "auditor"},
		{"admin", "editor"},
		{"editor", "read"},
		{"auditor", "read"},
	} {
		if err := manager.AddChild(edge[0], edge[1]); err != nil {
			b.Fatalf("add child %q -> %q: %v", edge[0], edge[1], err)
		}
	}

	if err := manager.Assign("u1", "root", time.Time{}); err != nil {
		b.Fatalf("assign root: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if !manager.UserHasPermission("u1", "read", nil) {
			b.Fatalf("expected permission %q to be granted", "read")
		}
	}
}

func BenchmarkManagerUserHasPermissionGuestRole(b *testing.B) {
	manager, permissionName := buildBenchmarkManagerWithRoleChain(b, 64)
	manager.SetGuestRoleName("role_0")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if !manager.UserHasPermission("", permissionName, nil) {
			b.Fatalf("expected guest permission %q to be granted", permissionName)
		}
	}
}

func buildBenchmarkManagerWithRoleChain(b *testing.B, depth int) (ManagerInterface[string], string) {
	b.Helper()

	manager := NewManager[string](
		NewRuleFactory(),
		NewItemsStorageInMemory(),
		NewAssignmentsStorageInMemory[string](),
		nil,
	)

	for i := 0; i < depth; i++ {
		roleName := fmt.Sprintf("role_%d", i)
		if err := manager.AddRole(NewRole(roleName)); err != nil {
			b.Fatalf("add role %q: %v", roleName, err)
		}
		if i > 0 {
			parentName := fmt.Sprintf("role_%d", i-1)
			if err := manager.AddChild(parentName, roleName); err != nil {
				b.Fatalf("add child %q -> %q: %v", parentName, roleName, err)
			}
		}
	}

	permissionName := fmt.Sprintf("permission_%d", depth)
	if err := manager.AddPermission(NewPermission(permissionName)); err != nil {
		b.Fatalf("add permission %q: %v", permissionName, err)
	}
	if err := manager.AddChild(fmt.Sprintf("role_%d", depth-1), permissionName); err != nil {
		b.Fatalf("add child role -> permission: %v", err)
	}
	if err := manager.Assign("u1", "role_0", time.Time{}); err != nil {
		b.Fatalf("assign role_0: %v", err)
	}

	return manager, permissionName
}
