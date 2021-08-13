package utils

import "os"

// Exists 判断文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// IsFile 判断所给路径是否是文件
func IsFile(path string) bool {
	if Exists(path) {
		return !IsDir(path)
	}
	return false
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// FileExists checks if a file exists and is not a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || err != nil || info == nil {
		return false
	}
	return !info.IsDir()
}

// FolderExists checks if a folder exists
func FolderExists(folderpath string) bool {
	_, err := os.Stat(folderpath)
	return !os.IsNotExist(err)
}
