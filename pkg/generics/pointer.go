package generics

func Pointer[T any](val T) *T {
	return &val
}
