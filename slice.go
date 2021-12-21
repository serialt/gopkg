package gopkg

// Does slice contain string
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// IndexSlice 查找string在slice的索引
func IndexSlice(slice []string, value string) (int, bool) {
	for index, item := range slice {
		if item == value {
			return index, true
		}
	}
	return 0, false
}

// Is a sub-slice of slice
func subslice(sub, main []string) bool {
	if len(sub) > len(main) {
		return false
	}
	for _, s := range sub {
		if !contains(main, s) {
			return false
		}
	}
	return true
}

// diffSlice sub 是否在main里，返回diff和bool
func diffSlice(sub, main []string) ([]string, bool) {
	var diff []string
	if len(sub) > len(main) {
		return diff, false
	}
	for _, s := range sub {
		if !contains(main, s) {
			diff = append(diff, s)
		}
	}
	return diff, true
}
