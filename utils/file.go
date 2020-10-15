package utils

import "os"

//IsFileExist 判断文件是否存在
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}