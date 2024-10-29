package service

import (
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateCode(n int) string {
    rand.Seed(time.Now().UnixNano())
    code := make([]byte, n)
    for i := range code {
        code[i] = letters[rand.Intn(len(letters))]
    }
    return string(code)
}
