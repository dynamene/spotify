package spotify

func sameStringSlice(x, y []string) bool {
	// Solution from: https://stackoverflow.com/questions/36000487/check-for-equality-on-slices-without-order
	if len(x) != len(y) {
		return false
	}

	diff := make(map[string]int, len(x))
	for _, _x := range x {
		diff[_x]++
	}

	for _, _y := range y {
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y]--
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}

	if len(diff) != 0 {
		return false
	}

	return true
}
