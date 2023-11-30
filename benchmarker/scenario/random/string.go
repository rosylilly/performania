package random

import "math/rand"

const (
	chars = "abcdefghjkmnpqrstwxyzABCDEFGHJKLMNPRSTWXYZ0123456789"
)

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
