package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Exists 判断档案是否存在
func Exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

// FileExists 判断文件是否存在
func FileExists(filename string) bool {
	fi, err := os.Stat(filename)
	return (err == nil || os.IsExist(err)) && !fi.IsDir()
}

// DirExists 判断目录是否存在
func DirExists(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}

// CreatedDir 创建文件夹
func CreatedDir(dir string) error {
	exist := DirExists(dir)
	if !exist {
		// 创建文件夹
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveDir 移除文件夹
func RemoveDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}

// CopyDir 复制文件夹
func CopyDir(fpathname, tpathname string, copyTo func(tfilename, filename string) error) error {
	rd, err := ioutil.ReadDir(fpathname)
	fmt.Println(err)
	fmt.Println(fpathname)
	fmt.Println(tpathname)
	for _, fi := range rd {
		if fi.IsDir() {
			if err := CreatedDir(filepath.Join(tpathname, fi.Name())); err != nil {
				return nil
			}
			if err := CopyDir(filepath.Join(fpathname, fi.Name()), filepath.Join(tpathname, fi.Name()), copyTo); err != nil {
				return err
			}
		} else {
			if err := copyTo(filepath.Join(fpathname, fi.Name()), filepath.Join(tpathname, fi.Name())); err != nil {
				return err
			}
		}
	}
	return err
}
