package rbac

import "maps"

type itemsStorageInMemory struct {
	items map[string]ItemInterface
	// For faster access, but doubles the memory
	permissions map[string]ItemInterface
	roles       map[string]ItemInterface

	// children[parentName][childName] = Item
	children map[string]map[string]ItemInterface
}

func NewItemsStorageInMemory() ItemsStorageInterface {
	return &itemsStorageInMemory{
		items:       make(map[string]ItemInterface),
		permissions: make(map[string]ItemInterface),
		roles:       make(map[string]ItemInterface),
		children:    make(map[string]map[string]ItemInterface),
	}
}

func (i *itemsStorageInMemory) GetAll() []ItemInterface {
	result := make([]ItemInterface, 0, len(i.items))
	for _, item := range i.items {
		result = append(result, item)
	}
	return result
}

func (i *itemsStorageInMemory) GetByNames(names []string) []ItemInterface {
	result := make([]ItemInterface, 0)
	for _, name := range names {
		if item, ok := i.items[name]; ok {
			result = append(result, item)
		}
	}
	return result
}

func (i *itemsStorageInMemory) Get(name string) (ItemInterface, error) {
	if item, ok := i.items[name]; ok {
		return item, nil
	}
	return nil, ErrItemNotFound
}

func (i *itemsStorageInMemory) Exists(name string) bool {
	_, ok := i.items[name]
	return ok
}

func (i *itemsStorageInMemory) RoleExists(name string) bool {
	_, ok := i.roles[name]
	return ok
}

func (i *itemsStorageInMemory) Add(item ItemInterface) {
	i.items[item.GetName()] = item
	t := item.GetType()
	switch t.(type) {
	case Permission:
		i.permissions[item.GetName()] = item
	case Role:
		i.roles[item.GetName()] = item
	}
}

func (i *itemsStorageInMemory) GetRole(name string) (ItemInterface, error) {
	if item, ok := i.roles[name]; ok {
		return item, nil
	}
	return nil, ErrItemNotFound
}

func (i *itemsStorageInMemory) GetRoles() []ItemInterface {
	result := make([]ItemInterface, 0, len(i.roles))
	for _, item := range i.roles {
		result = append(result, item)
	}
	return result
}

func (i *itemsStorageInMemory) GetRolesByNames(names []string) []ItemInterface {
	result := make([]ItemInterface, 0)
	for _, name := range names {
		if item, ok := i.roles[name]; ok {
			result = append(result, item)
		}
	}
	return result
}

func (i *itemsStorageInMemory) GetPermission(name string) (ItemInterface, error) {
	if item, ok := i.permissions[name]; ok {
		return item, nil
	}
	return nil, ErrItemNotFound
}

func (i *itemsStorageInMemory) GetPermissions() []ItemInterface {
	result := make([]ItemInterface, 0, len(i.permissions))
	for _, item := range i.permissions {
		result = append(result, item)
	}
	return result
}

func (i *itemsStorageInMemory) GetPermissionsByNames(names []string) []ItemInterface {
	result := make([]ItemInterface, 0)
	for _, name := range names {
		if item, ok := i.permissions[name]; ok {
			result = append(result, item)
		}
	}
	return result
}

func (i *itemsStorageInMemory) GetParents(name string) []ItemInterface {
	result := make([]ItemInterface, 0)
	i.fillParentsRecursive(name, &result)
	return result
}

func (i *itemsStorageInMemory) fillParentsRecursive(name string, result *[]ItemInterface) {
	for parentName, children := range i.children {
		for _, child := range children {
			if child.GetName() != name {
				continue
			}

			parent, err := i.Get(parentName)
			if err == nil {
				*result = append(*result, parent)
			}

			i.fillParentsRecursive(parentName, result)
		}
	}
}

func (i *itemsStorageInMemory) GetHierarchy(name string) map[string]TreeNode {
	result := make(map[string]TreeNode, 0)

	if _, ok := i.items[name]; !ok {
		return result
	}

	result[name] = TreeNode{
		Item:     i.items[name],
		Children: make(map[string]ItemInterface, 0),
	}

	addedChildItems := make(map[string]ItemInterface, 0)
	i.fillHierarchyRecursive(name, result, addedChildItems)
	return result
}

func (i *itemsStorageInMemory) fillHierarchyRecursive(
	name string,
	result map[string]TreeNode,
	addedChildItems map[string]ItemInterface,
) {
	for parentName, children := range i.children {
		for childName, child := range children {
			if child.GetName() != name {
				continue
			}

			_, err := i.Get(parentName)
			if err == nil {
				addedChildItems[childName] = child
				result[parentName] = TreeNode{
					Item:     i.items[parentName],
					Children: maps.Clone(addedChildItems),
				}
			}

			i.fillHierarchyRecursive(parentName, result, maps.Clone(addedChildItems))
		}
	}
}

func (i *itemsStorageInMemory) GetDirectChildren(name string) []ItemInterface {
	result := make([]ItemInterface, 0)
	if children, ok := i.children[name]; ok {
		result = make([]ItemInterface, 0, len(children))
		for _, child := range children {
			result = append(result, child)
		}
	}
	return result
}

func (i *itemsStorageInMemory) GetAllChildren(names []string) []ItemInterface {
	result := make([]ItemInterface, 0)
	result = i.getAllChildrenInternal(names, result)
	return result
}

func (i *itemsStorageInMemory) GetAllChildRoles(names []string) []ItemInterface {
	result := make([]ItemInterface, 0)
	result = i.getAllChildrenInternal(names, result)
	return i.filterRoles(result)
}

func (i *itemsStorageInMemory) GetAllChildPermissions(names []string) []ItemInterface {
	result := make([]ItemInterface, 0)
	result = i.getAllChildrenInternal(names, result)
	return i.filterPermissions(result)
}

func (i *itemsStorageInMemory) AddChild(parentName string, childName string) {
	if _, ok := i.children[parentName]; !ok {
		i.children[parentName] = make(map[string]ItemInterface, 0)
	}
	child, err := i.Get(childName)
	if err != nil {
		return
	}
	i.children[parentName][childName] = child
}

func (i *itemsStorageInMemory) HasChildren(name string) bool {
	_, ok := i.children[name]
	return ok
}

func (i *itemsStorageInMemory) HasChild(parentName string, childName string) bool {
	if parentName == childName {
		return true
	}

	children := i.GetDirectChildren(parentName)
	if len(children) == 0 {
		return false
	}

	for _, child := range children {
		if i.HasChild(child.GetName(), childName) {
			return true
		}
	}
	return false
}

func (i *itemsStorageInMemory) HasDirectChild(parentName string, childName string) bool {
	_, ok := i.children[parentName][childName]
	return ok
}

func (i *itemsStorageInMemory) RemoveChild(parentName string, childName string) {
	delete(i.children[parentName], childName)
}

func (i *itemsStorageInMemory) RemoveChildren(parentName string) {
	delete(i.children, parentName)
}

func (i *itemsStorageInMemory) Remove(name string) {
	i.clearChildrenFromItem(name)
	i.removeItemByName(name)
}

func (i *itemsStorageInMemory) Update(name string, item ItemInterface) {
	if item.GetName() != name {
		i.updateItemName(name, item)
		i.removeItemByName(name)
	}
	i.Add(item)
}

func (i *itemsStorageInMemory) Clear() {
	i.items = make(map[string]ItemInterface, 0)
	i.roles = make(map[string]ItemInterface, 0)
	i.permissions = make(map[string]ItemInterface, 0)
	i.children = make(map[string]map[string]ItemInterface, 0)
}

func (i *itemsStorageInMemory) ClearPermissions() {
	for permName := range i.permissions {
		delete(i.items, permName)
	}
	i.permissions = make(map[string]ItemInterface, 0)
}

func (i *itemsStorageInMemory) ClearRoles() {
	for roleName := range i.roles {
		delete(i.items, roleName)
	}
	i.roles = make(map[string]ItemInterface, 0)
}

func (i *itemsStorageInMemory) updateItemName(name string, item ItemInterface) {
	i.updateChildrenForItemName(name, item)
}

func (i *itemsStorageInMemory) updateChildrenForItemName(name string, item ItemInterface) {
	// If old item has children, move them to new item
	if i.HasChildren(name) {
		i.children[item.GetName()] = i.children[name]
		delete(i.children, name)
	}

	// if old item has parents, link them to new item
	for _, children := range i.children {
		if _, ok := children[name]; ok {
			children[item.GetName()] = item
			delete(children, name)
		}
	}
}

func (i *itemsStorageInMemory) removeItemByName(name string) {
	delete(i.roles, name)
	delete(i.permissions, name)
	delete(i.items, name)

	for _, children := range i.children {
		delete(children, name)
	}
}

func (i *itemsStorageInMemory) clearChildrenFromItem(name string) {
	delete(i.children, name)
}

func (i *itemsStorageInMemory) getAllChildrenInternal(
	names []string,
	result []ItemInterface,
) []ItemInterface {
	baseNames := make(map[string]bool, len(names))
	for _, name := range names {
		baseNames[name] = true
	}
	for _, name := range names {
		result = i.fillChildrenRecursive(name, result, baseNames)
	}
	return result
}

func (i *itemsStorageInMemory) fillChildrenRecursive(
	name string,
	result []ItemInterface,
	baseNames map[string]bool,
) []ItemInterface {
	for childName := range i.children[name] {
		if baseNames[childName] {
			continue
		}
		baseNames[childName] = true

		child, err := i.Get(childName)
		if err == nil {
			result = append(result, child)
			result = i.fillChildrenRecursive(child.GetName(), result, baseNames)
		}
	}
	return result
}

func (i *itemsStorageInMemory) filterRoles(items []ItemInterface) []ItemInterface {
	result := make([]ItemInterface, 0)
	for _, item := range items {
		if IsItem[Role](item) {
			result = append(result, item)
		}
	}
	return result
}

func (i *itemsStorageInMemory) filterPermissions(items []ItemInterface) []ItemInterface {
	result := make([]ItemInterface, 0)
	for _, item := range items {
		if IsItem[Permission](item) {
			result = append(result, item)
		}
	}
	return result
}
