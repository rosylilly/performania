package random

import "math/rand"

func RandomItem[T any](items []T) T {
	return items[rand.Intn(len(items))]
}
