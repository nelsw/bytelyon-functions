package random

import "math/rand"

func Element[T any](a ...T) T {
	return a[rand.Intn(len(a))]
}

func Between(min, max int) int {
	return rand.Intn(max-min) + min
}
