package helper

/**
import . "xxx/helper"
c:=IfElse(condition, a, b)
*/
func IfElse[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}