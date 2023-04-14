package utils

func List2Set[T string | int](list []T) map[T]struct{} {
	m := make(map[T]struct{})
	for _, item := range list {
		m[item] = struct{}{}
	}
	return m
}

func SetSubList[T string | int](m map[T]struct{}, list []T) map[T]struct{} {
	for _, item := range list {
		delete(m, item)
	}
	return m
}

func Set2List[T string | int](m map[T]struct{}) []T {
	var list []T
	for t := range m {
		list = append(list, t)
	}
	return list
}

func ListIsContains[T string | int](list []T, t T) bool {
	for _, item := range list {
		if item == t {
			return true
		}
	}
	return false
}

func ListRemoveOne[T string | int](list *[]T, t T) {
	for index, item := range *list {
		if item == t {
			*list = append((*list)[:index], (*list)[index+1:]...)
			break
		}
	}
}
