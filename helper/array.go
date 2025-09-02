package helper

func SliceContains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
func SliceRemove(slice *[]int, val int) {
	s := *slice
	j := 0
	for _, v := range s {
		if v != val {
			s[j] = v
			j++
		}
	}
	*slice = s[:j] // reslice in place
}
func AddUnique(slice []int, v int) ([]int, bool) {
	for _, elem := range slice {
		if elem == v {
			// already exists, return unchanged
			return slice, false
		}
	}
	// not found, append
	return append(slice, v), true
}

func SliceShift(slice *[]int) (int, bool) {
	if slice == nil || len(*slice) == 0 {
		return 0, false // slice is empty or nil
	}

	val := (*slice)[0]
	*slice = (*slice)[1:] // remove the first element

	return val, true
}
