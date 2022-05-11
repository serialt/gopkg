/*
 * @Description   : IMAU of Serialt
 * @Author        : serialt
 * @Email         : serialt@qq.com
 * @Github        : https://github.com/serialt
 * @Created Time  : 2022-04-17 23:22:12
 * @Last modified : 2022-05-11 23:37:43
 * @FilePath      : /gopkg/slice.go
 * @Other         :
 * @              :
 *
 *                 人和代码，有一个能跑就行
 *
 */

package gopkg

/**
 * @description: 判断字符串是否在切片里
 * @author: Serialt
 * @param slice {[]string}
 * @param value {string}
 * @return {bool}
 */
// Contains Does slice contain string
func Contains(slice []string, value string) bool {
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

// IsSubSlice Is a sub-slice of slice
func IsSubSlice(sub, main []string) bool {
	if len(sub) > len(main) {
		return false
	}
	for _, s := range sub {
		if !Contains(main, s) {
			return false
		}
	}
	return true
}

// DiffSlice sub是否在main里,返回sub里不在main里的元素，返回diff和bool
func DiffSlice(sub, main []string) ([]string, bool) {
	var diff []string
	if len(sub) > len(main) {
		return diff, false
	}
	for _, s := range sub {
		if !Contains(main, s) {
			diff = append(diff, s)
		}
	}
	if len(diff) > 0 {
		return diff, false
	}
	return diff, true
}

func IsEqualStringSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}

	if (x == nil) != (y == nil) {
		return false
	}

	for i, v := range x {
		if v != x[i] {
			return false
		}
	}

	return true
}

func IsEqualIntSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}

	if (x == nil) != (y == nil) {
		return false
	}

	for i, v := range x {
		if v != x[i] {
			return false
		}
	}

	return true
}
