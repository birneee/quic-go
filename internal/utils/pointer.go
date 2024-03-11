package utils

// allocate new memory and initiate it with the given value
func New[T any](value T) *T {
	n := new(T)
	*n = value
	return n
}
