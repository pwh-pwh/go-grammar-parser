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

func SetIntersects[T string | int](set1 map[T]struct{}, set2 map[T]struct{}) bool {
	flag := false
	for t := range set2 {
		if _, ok := set1[t]; ok {
			flag = true
			break
		}
	}
	return flag
}

func SetIntersect[T string | int](set1 map[T]struct{}, set2 map[T]struct{}) map[T]struct{} {
	result := make(map[T]struct{})
	for t := range set1 {
		if _, ok := set2[t]; ok {
			result[t] = struct{}{}
		}
	}
	return result
}

func ListIndexOf[T string | int](list []T, data T) int {
	res := -1
	for index, item := range list {
		if data == item {
			res = index
			break
		}
	}
	return res
}
