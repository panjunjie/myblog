package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

func GetFileName(fullPath string) string {
	if fullPath != "" {
		return path.Base(fullPath)
	}
	return fullPath
}

/*
**获取文件的后缀(前带点)
 */
func GetFileExt(fullPath string) string {
	if fullPath != "" {
		fileName := path.Base(fullPath)
		return path.Ext(fileName)
	}
	return fullPath
}

func GetFileNameNotExt(fullPath string) string {
	if fullPath != "" {
		fileName := path.Base(fullPath)
		fileExt := path.Ext(fileName)
		return strings.TrimSuffix(fileName, fileExt)
	}
	return fullPath
}

func CreateFileName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), r.Intn(100))
	//return fmt.Sprintf("%d%s", time.Now().UnixNano(), RandomNumeric(4))
}

func IsDirExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func IsFileExists(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func WriteStr2File(filename, content string) {
	var file *os.File
	defer file.Close()
	if IsFileExists(filename) { //如果文件存在
		f, err := os.OpenFile(filename, os.O_APPEND, 0666)
		defer f.Close()
		if err != nil {
			fmt.Printf("打开文件发生错误：", err)
		}
		file = f
	} else {
		f, err := os.Create(filename) //创建文件
		defer f.Close()
		if err != nil {
			fmt.Printf("创建文件发生错误：", err)
		}
		file = f
	}

	if file != nil {
		_, err := io.WriteString(file, content) //写入文件(字符串)
		if err != nil {
			fmt.Printf("写入文件发生错误：", err)
		}
	}
}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func MoveFile(fromPath, targetPath string) {
	err := os.Rename(fromPath, targetPath)
	fmt.Println(err)
}

func RemoveFile(path string) bool {
	err := os.Remove(path)
	if err != nil {
		return false
	}
	return true
}
