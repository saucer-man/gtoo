package utils

import (
	"os"
	"runtime"
)

func IsRoot() bool {
	if runtime.GOOS != "windows" {
		return os.Getuid() == 0
	}
	return false
}
