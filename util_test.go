package main

func toPtr[T any](val T) *T {
	return &val
}
