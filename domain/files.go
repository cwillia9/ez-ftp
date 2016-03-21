package domain

import (
	"math/rand"
	"time"
)

type FileRepository interface {
	FindByIdentifier(string) (File, error)
	Store(File) error
}

// type File interface {
// 	Open()
// 	Close()
// }

// GenerateObsureId will return a random string of length = 32
func GenerateObsureId() string {
	return randomString(32)
}

func randomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
