package utils

import (
	"math/rand"
	"os"
	"runtime"
	"time"
)

func IsRoot() bool {
	if runtime.GOOS != "windows" {
		return os.Getuid() == 0
	}
	return false
}
func RandomStr(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
