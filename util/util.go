// fileutil
package util

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	Fileseprater string
)

func init() {
	if os.IsPathSeparator('\\') {
		Fileseprater = "\\"
	} else {
		Fileseprater = "/"
	}
}

func MakeDirAll(filepath string) {
	os.MkdirAll(filepath, 0660)
}

func Lastmodified(filename string) string {
	fi, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	lastmodified := fi.ModTime()
	return lastmodified.Format("20060102")
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 复制文件
// 将src复制到dst
func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)

}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func ExtactFileName(filename string) (s string) {
	rs := []rune(filename)
	index := -1
	for i := len(rs) - 1; i >= 0; i-- {
		if string(rs[i]) == Fileseprater {
			index = i
			break
		}
	}
	return string(rs[index+1 : len(rs)])

}

func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	splitstring := strings.Split(path, "\\")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	ret := strings.Replace(splitstring[0], "\\", "/", size-1)
	return ret
}
