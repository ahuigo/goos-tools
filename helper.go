package helper

/**
import . "xxx/helper"
c:=IfThen(condition, a, b)
*/
func IfThen[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}