package helper

func IfThen[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}
