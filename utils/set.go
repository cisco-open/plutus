package utils

// UniqueConcat concatenates two slices such that there are no repeated items
func UniqueConcat(s1 []string, s2 ...string) []string {
	set := make(map[string]bool, len(s1))

	for _, e1 := range s1 {
		set[e1] = true
	}

	for _, e2 := range s2 {
		set[e2] = true
	}

	slice := make([]string, len(set))

	i := 0
	for e := range set {
		slice[i] = e
		i++
	}
	return slice
}

// Intersection gives the comman elements of both the slices
func Intersection(s1 []string, s2 []string) []string {
	res := make([]string, 0)

	for _, e1 := range s1 {
		for _, e2 := range s2 {
			if e1 == e2 {
				res = append(res, e1)
			}
		}
	}

	return res
}

// ToSet converts the given slice into a set
func ToSet(strs []string) map[string]bool {
	set := make(map[string]bool)

	for _, str := range strs {
		set[str] = true
	}

	return set
}

// ToSlice converts the given set into a slice
func ToSlice(set map[string]bool) []string {
	slice := make([]string, len(set))

	i := 0
	for elem := range set {
		slice[i] = elem
		i++
	}
	return slice
}
