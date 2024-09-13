package logger

import (
	"fmt"
	"os"
)

func checkNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}
func checkPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}
func isNotExistMkDir(src string) error {
	if checkNotExist(src) {
		err := os.MkdirAll(src, os.ModePerm)
		return err
	}
	return nil
}

func mustOpen(fileName, dir string) (*os.File, error) {
	if checkPermission(dir) {
		return nil, fmt.Errorf("permission denied dir: %s", dir)
	}
	if err := isNotExistMkDir(dir); err != nil {
		return nil, fmt.Errorf("error during making dir %s, err:%s", dir, err)
	}
	f, err := os.OpenFile(dir+string(os.PathSeparator)+fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s, err:%s", fileName, err)
	}
	return f, nil
}
