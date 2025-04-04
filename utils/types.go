package utils

func ToPtr[T any](i T) *T {
	return &i
}
