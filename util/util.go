package util

import (
	"math/rand"
	"time"
)

// used for password generation
var typos = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_!@#")

func GeneratePassword() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	password := make([]rune, 32)
	for i := range password {
		password[i] = typos[rand.Int63n(int64(len(typos)))]
	}

	return string(password)
}
