package utils

type OrderedType interface {
	int | float32 | float64
}

func Max[T OrderedType](x, y T) T {
	if x > y {
		return x
	}

	return y
}

func Min[T OrderedType](x, y T) T {
	if x > y {
		return y
	}

	return x
}

// 取差集
func DifferenceSet(a []uint, b []uint) []uint {
	var c []uint
	temp := map[uint]struct{}{}

	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}

	for _, val := range a {
		if _, ok := temp[val]; !ok {
			c = append(c, val)
		}
	}

	return c
}
