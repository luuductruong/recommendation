package helper

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func MapList[S any, D any](source []S, mapper func(S) D) []D {
	rs := make([]D, len(source))
	for i := range source {
		rs[i] = mapper(source[i])
	}
	return rs
}

func RandString(n int) string {
	var source = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	max := len(source)
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = source[rand.Intn(max)]
	}

	return string(b)
}

func NewStringUUID() string {
	return uuid.New().String()
}

// Unique ...
func Unique[T comparable](arr []T) []T {
	u := make([]T, 0)
	m := make(map[T]bool)
	for _, val := range arr {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}
