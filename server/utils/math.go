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

// func min(T constraints.Ordered) (x, y T) T {
// 	if x < y {
// 	  return x
// 	}
// 	return y
//   }
