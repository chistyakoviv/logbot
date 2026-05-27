package rbac

import "testing"

func TestItemsStorageInMemoryChildTraversalAndHierarchy(t *testing.T) {
	storage := NewItemsStorageInMemory()
	admin := NewRole("admin")
	editor := NewRole("editor")
	read := NewPermission("read")
	write := NewPermission("write")

	for _, item := range []ItemInterface{admin, editor, read, write} {
		storage.Add(item)
	}

	storage.AddChild("admin", "editor")
	storage.AddChild("editor", "read")
	storage.AddChild("admin", "write")

	assertItemNamesMatch(t, storage.GetDirectChildren("admin"), []string{"editor", "write"})
	assertItemNamesMatch(t, storage.GetAllChildren([]string{"admin"}), []string{"editor", "read", "write"})
	assertItemNamesMatch(t, storage.GetAllChildRoles([]string{"admin"}), []string{"editor"})
	assertItemNamesMatch(t, storage.GetAllChildPermissions([]string{"admin"}), []string{"read", "write"})
	assertItemNamesMatch(t, storage.GetParents("read"), []string{"admin", "editor"})

	hierarchy := storage.GetHierarchy("read")
	if len(hierarchy) != 3 {
		t.Fatalf("expected 3 hierarchy nodes, got %d", len(hierarchy))
	}

	baseNode, ok := hierarchy["read"]
	if !ok {
		t.Fatalf("expected hierarchy to contain base item")
	}
	if baseNode.Item == nil || baseNode.Item.GetName() != "read" {
		t.Fatalf("expected base node to contain read item")
	}
	if len(baseNode.Children) != 0 {
		t.Fatalf("expected base node to have no children, got %d", len(baseNode.Children))
	}

	editorNode, ok := hierarchy["editor"]
	if !ok {
		t.Fatalf("expected hierarchy to contain editor node")
	}
	assertMapItemNamesMatch(t, editorNode.Children, []string{"read"})

	adminNode, ok := hierarchy["admin"]
	if !ok {
		t.Fatalf("expected hierarchy to contain admin node")
	}
	assertMapItemNamesMatch(t, adminNode.Children, []string{"editor", "read"})
}

func TestItemsStorageInMemoryGetHierarchyFindsParentsAndChildrenForEachParentNode(t *testing.T) {
	storage := NewItemsStorageInMemory()
	root := NewRole("root")
	admin := NewRole("admin")
	auditor := NewRole("auditor")
	editor := NewRole("editor")
	read := NewPermission("read")

	for _, item := range []ItemInterface{root, admin, auditor, editor, read} {
		storage.Add(item)
	}

	storage.AddChild("root", "admin")
	storage.AddChild("root", "auditor")
	storage.AddChild("admin", "editor")
	storage.AddChild("editor", "read")
	storage.AddChild("auditor", "read")

	assertItemNamesMatch(t, storage.GetParents("read"), []string{"admin", "auditor", "editor", "root"})

	hierarchy := storage.GetHierarchy("read")
	if len(hierarchy) != 5 {
		t.Fatalf("expected 5 hierarchy nodes, got %d", len(hierarchy))
	}

	assertHierarchyNodeChildren(t, hierarchy, "read", []string{})
	assertHierarchyNodeChildren(t, hierarchy, "editor", []string{"read"})
	assertHierarchyNodeChildren(t, hierarchy, "admin", []string{"editor", "read"})
	assertHierarchyNodeChildren(t, hierarchy, "auditor", []string{"read"})
	assertHierarchyNodeChildren(t, hierarchy, "root", []string{"admin", "auditor", "editor", "read"})
}

func TestItemsStorageInMemoryUpdateRenamesItemAndPreservesRelations(t *testing.T) {
	storage := NewItemsStorageInMemory()
	admin := NewRole("admin")
	editor := NewRole("editor")
	read := NewPermission("read")

	for _, item := range []ItemInterface{admin, editor, read} {
		storage.Add(item)
	}

	storage.AddChild("admin", "editor")
	storage.AddChild("editor", "read")
	storage.Update("editor", editor.WithName("author"))

	if storage.Exists("editor") {
		t.Fatalf("expected old item name to be removed")
	}
	if !storage.Exists("author") {
		t.Fatalf("expected renamed item to exist")
	}
	if !storage.HasDirectChild("admin", "author") {
		t.Fatalf("expected parent relation to be updated")
	}
	if !storage.HasDirectChild("author", "read") {
		t.Fatalf("expected child relation to be preserved")
	}
}

func assertHierarchyNodeChildren(t *testing.T, hierarchy map[string]TreeNode, nodeName string, expectedChildren []string) {
	t.Helper()

	node, ok := hierarchy[nodeName]
	if !ok {
		t.Fatalf("expected hierarchy to contain %q node", nodeName)
	}
	if node.Item == nil || node.Item.GetName() != nodeName {
		t.Fatalf("expected hierarchy node %q to contain matching item", nodeName)
	}

	assertMapItemNamesMatch(t, node.Children, expectedChildren)
}
