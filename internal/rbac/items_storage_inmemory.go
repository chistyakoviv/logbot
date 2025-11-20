package rbac

import "slices"

type itemsStorageInMemory struct {
	items map[string]StorableItemInterface
	// For faster access, but doubles the memory
	permissions map[string]*Permission
	roles       map[string]*Role

	// children[parentName][childName] = Item
	children map[string]map[string]StorableItemInterface
}

func NewItemsStorageInMemory() ItemStorageInterface {
	return &itemsStorageInMemory{
		items:       make(map[string]StorableItemInterface),
		permissions: make(map[string]*Permission),
		roles:       make(map[string]*Role),
		children:    make(map[string]map[string]StorableItemInterface),
	}
}

func (i *itemsStorageInMemory) GetAll() map[string]StorableItemInterface {
	result := make(map[string]StorableItemInterface, len(i.items))
	for name, item := range i.items {
		result[name] = item
	}
	return result
}

func (i *itemsStorageInMemory) GetByNames(names []string) map[string]StorableItemInterface {
	result := make(map[string]StorableItemInterface, 0)
	for _, name := range names {
		if item, ok := i.items[name]; ok {
			result[name] = item
		}
	}
	return result
}

func (i *itemsStorageInMemory) Get(name string) (StorableItemInterface, error) {
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

func (i *itemsStorageInMemory) Add(item StorableItemInterface) {
	i.items[item.GetName()] = item
	switch item := item.(type) {
	case *Permission:
		i.permissions[item.GetName()] = item
	case *Role:
		i.roles[item.GetName()] = item
	}
}

func (i *itemsStorageInMemory) GetRole(name string) (*Role, error) {
	if item, ok := i.roles[name]; ok {
		return item, nil
	}
	return nil, ErrItemNotFound
}

func (i *itemsStorageInMemory) GetRoles() map[string]*Role {
	result := make(map[string]*Role, len(i.roles))
	for name, item := range i.roles {
		result[name] = item
	}
	return result
}

func (i *itemsStorageInMemory) GetRolesByNames(names []string) map[string]*Role {
	result := make(map[string]*Role, 0)
	for _, name := range names {
		if item, ok := i.roles[name]; ok {
			result[name] = item
		}
	}
	return result
}

func (i *itemsStorageInMemory) GetPermission(name string) (*Permission, error) {
	if item, ok := i.permissions[name]; ok {
		return item, nil
	}
	return nil, ErrItemNotFound
}

func (i *itemsStorageInMemory) GetPermissions() map[string]*Permission {
	result := make(map[string]*Permission, len(i.permissions))
	for name, item := range i.permissions {
		result[name] = item
	}
	return result
}

func (i *itemsStorageInMemory) GetPermissionsByNames(names []string) map[string]*Permission {
	result := make(map[string]*Permission, 0)
	for _, name := range names {
		if item, ok := i.permissions[name]; ok {
			result[name] = item
		}
	}
	return result
}

func (i *itemsStorageInMemory) GetParents(name string) map[string]StorableItemInterface {
	result := make(map[string]StorableItemInterface, 0)
	i.fillParentsRecursive(name, result)
	return result
}

func (i *itemsStorageInMemory) fillParentsRecursive(name string, result map[string]StorableItemInterface) {
	for parentName, children := range i.children {
		for _, child := range children {
			if child.GetName() != name {
				continue
			}

			parent, err := i.Get(parentName)
			if err != nil {
				result[parentName] = parent
			}

			i.fillParentsRecursive(parentName, result)
		}
	}
}

func (i *itemsStorageInMemory) GetHierarchy(name string) map[string]map[string]map[string]StorableItemInterface {
	result := make(map[string]map[string]map[string]StorableItemInterface, 0)
	addedChildItems := make(map[string]StorableItemInterface, 0)
	i.fillHierarchyRecursive(name, result, addedChildItems)
	return result
}

func (i *itemsStorageInMemory) fillHierarchyRecursive(
	name string,
	result map[string]map[string]map[string]StorableItemInterface,
	addedChildItems map[string]StorableItemInterface,
) {
	for parentName, children := range i.children {
		for childName, child := range children {
			if child.GetName() != name {
				continue
			}

			_, err := i.Get(parentName)
			if err != nil {
				result[parentName]["item"] = children

				addedChildItems[childName] = child
				result[parentName]["children"] = addedChildItems
			}

			i.fillHierarchyRecursive(parentName, result, addedChildItems)
		}
	}
}

func (i *itemsStorageInMemory) GetDirectChildren(name string) map[string]StorableItemInterface {
	if children, ok := i.children[name]; ok {
		return children
	}
	return make(map[string]StorableItemInterface, 0)
}

func (i *itemsStorageInMemory) GetAllChildren(names []string) map[string]StorableItemInterface {
	result := make(map[string]StorableItemInterface, 0)
	i.getAllChildrenInternal(names, result)
	return result
}

func (i *itemsStorageInMemory) GetAllChildRoles(names []string) map[string]*Role {
	result := make(map[string]StorableItemInterface, 0)
	i.getAllChildrenInternal(names, result)
	return i.filterRoles(result)
}

func (i *itemsStorageInMemory) GetAllChildPermissions(names []string) map[string]*Permission {
	result := make(map[string]StorableItemInterface, 0)
	i.getAllChildrenInternal(names, result)
	return i.filterPermissions(result)
}

func (i *itemsStorageInMemory) AddChild(parentName string, childName string) {
	if _, ok := i.children[parentName]; !ok {
		i.children[parentName] = make(map[string]StorableItemInterface, 0)
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

func (i *itemsStorageInMemory) Update(name string, item StorableItemInterface) {
	if item.GetName() != name {
		i.updateItemName(name, item)
		i.removeItemByName(item.GetName())
	}
	i.Add(item)
}

func (i *itemsStorageInMemory) Clear() {
	i.items = make(map[string]StorableItemInterface, 0)
	i.roles = make(map[string]*Role, 0)
	i.permissions = make(map[string]*Permission, 0)
	i.children = make(map[string]map[string]StorableItemInterface, 0)
}

func (i *itemsStorageInMemory) ClearPermissions() {
	for permName, _ := range i.permissions {
		delete(i.items, permName)
	}
	i.permissions = make(map[string]*Permission, 0)
}

func (i *itemsStorageInMemory) ClearRoles() {
	for roleName, _ := range i.roles {
		delete(i.items, roleName)
	}
	i.roles = make(map[string]*Role, 0)
}

func (i *itemsStorageInMemory) updateItemName(name string, item StorableItemInterface) {
	i.updateChildrenForItemName(name, item)
}

func (i *itemsStorageInMemory) updateChildrenForItemName(name string, item StorableItemInterface) {
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
	result map[string]StorableItemInterface,
) map[string]StorableItemInterface {
	for _, name := range names {
		i.fillChildrenRecursive(name, result, names)
	}
	return result
}

func (i *itemsStorageInMemory) fillChildrenRecursive(
	name string,
	result map[string]StorableItemInterface,
	baseNames []string,
) {
	for childName, _ := range i.children[name] {
		if slices.Contains(baseNames, childName) {
			continue
		}
		child, err := i.Get(childName)
		if err != nil {
			result[childName] = child
		}

		i.fillChildrenRecursive(child.GetName(), result, baseNames)
	}
}

func (i *itemsStorageInMemory) filterRoles(items map[string]StorableItemInterface) map[string]*Role {
	result := make(map[string]*Role, 0)
	for name, item := range items {
		if role, ok := item.(*Role); ok {
			result[name] = role
		}
	}
	return result
}

func (i *itemsStorageInMemory) filterPermissions(items map[string]StorableItemInterface) map[string]*Permission {
	result := make(map[string]*Permission, 0)
	for name, item := range items {
		if permission, ok := item.(*Permission); ok {
			result[name] = permission
		}
	}
	return result
}
