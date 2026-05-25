package rbac

type assignmentsStorageInMemory[T comparable] struct {
	assignments map[T]map[string]*Assignment[T]
}

func NewAssignmentsStorageInMemory[T comparable]() AssignmentsStorageInterface[T] {
	return &assignmentsStorageInMemory[T]{
		assignments: make(map[T]map[string]*Assignment[T]),
	}
}

func (a *assignmentsStorageInMemory[T]) GetAll() []*Assignment[T] {
	res := make([]*Assignment[T], 0)
	for _, userAssignments := range a.assignments {
		for _, assignment := range userAssignments {
			res = append(res, assignment)
		}
	}
	return res
}

func (a *assignmentsStorageInMemory[T]) GetByUserId(userId T) []*Assignment[T] {
	userAssignments, ok := a.assignments[userId]
	if !ok {
		// Return empty slice so that we don't have to check for nil
		return make([]*Assignment[T], 0)
	}
	res := make([]*Assignment[T], 0, len(userAssignments))
	for _, assignment := range userAssignments {
		res = append(res, assignment)
	}
	return res
}

func (a *assignmentsStorageInMemory[T]) GetByItemNames(itemNames []string) []*Assignment[T] {
	assignments := make([]*Assignment[T], 0)
	for _, assignment := range a.assignments {
		for _, itemName := range itemNames {
			if assignment[itemName] != nil {
				assignments = append(assignments, assignment[itemName])
			}
		}
	}
	return assignments
}

func (a *assignmentsStorageInMemory[T]) Get(userId T, itemName string) *Assignment[T] {
	assigments, ok := a.assignments[userId]
	if !ok {
		return nil
	}
	return assigments[itemName]
}

func (a *assignmentsStorageInMemory[T]) Exists(userId T, itemName string) bool {
	assigments, ok := a.assignments[userId]
	if !ok {
		return false
	}
	_, ok = assigments[itemName]
	return ok
}

func (a *assignmentsStorageInMemory[T]) UserHasItem(userId T, itemNames []string) bool {
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

func (a *assignmentsStorageInMemory[T]) FilterUserItemNames(userId T, itemNames []string) []string {
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

func (a *assignmentsStorageInMemory[T]) Add(assignment *Assignment[T]) {
	_, ok := a.assignments[assignment.GetUserId()]
	if !ok {
		a.assignments[assignment.GetUserId()] = make(map[string]*Assignment[T])
	}
	a.assignments[assignment.GetUserId()][assignment.GetItemName()] = assignment
}

func (a *assignmentsStorageInMemory[T]) HasItem(itemName string) bool {
	for _, assigments := range a.assignments {
		_, ok := assigments[itemName]
		if ok {
			return true
		}
	}
	return false
}

func (a *assignmentsStorageInMemory[T]) RenameItem(oldName string, newName string) {
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

func (a *assignmentsStorageInMemory[T]) Remove(userId T, itemName string) {
	assigments, ok := a.assignments[userId]
	if !ok {
		return
	}
	delete(assigments, itemName)
}

func (a *assignmentsStorageInMemory[T]) RemoveByUserId(userId T) {
	delete(a.assignments, userId)
}

func (a *assignmentsStorageInMemory[T]) RemoveByItemName(itemName string) {
	for _, assigments := range a.assignments {
		delete(assigments, itemName)
	}
}

func (a *assignmentsStorageInMemory[T]) Clear() {
	a.assignments = make(map[T]map[string]*Assignment[T])
}
