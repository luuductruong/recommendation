package helper

import (
	"math/rand"
	"time"
)

func MapList[S any, D any](source []*S, mapper func(*S) *D) []*D {
	return Mapper(source, func(t *S, _ int) *D {
		return mapper(t)
	})
}

func Mapper[T any, R any](collection []T, iteratee func(T, int) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = iteratee(item, i)
	}

	return result
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
