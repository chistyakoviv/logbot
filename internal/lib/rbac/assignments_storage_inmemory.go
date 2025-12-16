package rbac

type assignmentsStorageInMemory struct {
	assignments map[any]map[string]*Assignment
}

func NewAssignmentsStorageInMemory() AssignmentsStorageInterface {
	return &assignmentsStorageInMemory{
		assignments: make(map[any]map[string]*Assignment),
	}
}

func (a *assignmentsStorageInMemory) GetAll() map[any]map[string]*Assignment {
	return a.assignments
}

func (a *assignmentsStorageInMemory) GetByUserId(userId any) map[string]*Assignment {
	_, ok := a.assignments[userId]
	if !ok {
		// Return empty map so that we don't have to check for nil
		return make(map[string]*Assignment)
	}
	return a.assignments[userId]
}

func (a *assignmentsStorageInMemory) GetByItemNames(itemNames []string) []*Assignment {
	assignments := make([]*Assignment, 0)
	for _, assignment := range a.assignments {
		for _, itemName := range itemNames {
			if assignment[itemName] != nil {
				assignments = append(assignments, assignment[itemName])
			}
		}
	}
	return assignments
}

func (a *assignmentsStorageInMemory) Get(userId any, itemName string) *Assignment {
	assigments, ok := a.assignments[userId]
	if !ok {
		return nil
	}
	return assigments[itemName]
}

func (a *assignmentsStorageInMemory) Exists(userId any, itemName string) bool {
	assigments, ok := a.assignments[userId]
	if !ok {
		return false
	}
	_, ok = assigments[itemName]
	return ok
}

func (a *assignmentsStorageInMemory) UserHasItem(userId any, itemNames []string) bool {
	assigments, ok := a.assignments[userId]
	if !ok {
		return false
	}
	for _, itemName := range itemNames {
		_, ok = assigments[itemName]
		if ok {
			return true
		}
	}
	return false
}

func (a *assignmentsStorageInMemory) FilterUserItemNames(userId any, itemNames []string) []string {
	result := make([]string, 0)
	assigments, ok := a.assignments[userId]
	if !ok {
		return result
	}
	for _, itemName := range itemNames {
		_, ok = assigments[itemName]
		if ok {
			result = append(result, itemName)
		}
	}
	return result
}

func (a *assignmentsStorageInMemory) Add(assignment *Assignment) {
	_, ok := a.assignments[assignment.GetUserId()]
	if !ok {
		a.assignments[assignment.GetUserId()] = make(map[string]*Assignment)
	}
	a.assignments[assignment.GetUserId()][assignment.GetItemName()] = assignment
}

func (a *assignmentsStorageInMemory) HasItem(itemName string) bool {
	for _, assigments := range a.assignments {
		_, ok := assigments[itemName]
		if ok {
			return true
		}
	}
	return false
}

func (a *assignmentsStorageInMemory) RenameItem(oldName string, newName string) {
	if oldName == newName {
		return
	}
	for _, assigments := range a.assignments {
		if assigments[oldName] != nil {
			clonedAssignment := assigments[oldName].WithItemName(newName)
			assigments[newName] = &clonedAssignment
			delete(assigments, oldName)
		}
	}
}

func (a *assignmentsStorageInMemory) Remove(userId any, itemName string) {
	assigments, ok := a.assignments[userId]
	if !ok {
		return
	}
	delete(assigments, itemName)
}

func (a *assignmentsStorageInMemory) RemoveByUserId(userId any) {
	delete(a.assignments, userId)
}

func (a *assignmentsStorageInMemory) RemoveByItemName(itemName string) {
	for _, assigments := range a.assignments {
		delete(assigments, itemName)
	}
}

func (a *assignmentsStorageInMemory) Clear() {
	a.assignments = make(map[any]map[string]*Assignment)
}
