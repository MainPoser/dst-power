package fs

import (
	"errors"
	"io/fs"
	"os"
)

// CreateDir 创建文件夹
func CreateDir(path string) error {
	// 判断文件夹是否存在
	if IsExist(path) {
		return nil
	}
	// 创建多层级目录
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// IsExist 判断文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	// 读取文件信息，判断文件是否存在
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return true
		}
		return false
	}
	return true
}
