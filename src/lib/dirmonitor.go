package lib

import (
	"os"
)

/*
Dir Info struct
*/
type DirInfo struct {
	Name string //文件夹名称
	//Path     string    //路径
	FileNum  int           //文件数量
	FileList []os.FileInfo //文件列表
	DirNum   int           //文件夹数量
	DirList  []os.FileInfo //文件夹信息
}

/*
Dir Info init
*/
func (dir DirInfo) DirInfoInit(dirName string) (dirInfo DirInfo) {
	f, _ := os.Open(dirName)
	defer f.Close()

	dirInfo.Name = dirName

	infoD, _ := f.Readdirnames(0)
	dirNum := 0
	fileNum := 0
	fileList := make([]os.FileInfo, 0)
	dirList := make([]os.FileInfo, 0)

	for _, v := range infoD {
		cf, err := os.Stat(dirName + "/" + v)
		if err == nil {
			if cf.IsDir() {
				dirList = append(dirList, cf)
				dirNum++
			} else {
				fileList = append(fileList, cf)
				fileNum++
			}
		}
	}

	dirInfo.FileNum = fileNum
	dirInfo.DirNum = dirNum
	dirInfo.FileList = fileList
	dirInfo.DirList = dirList

	return
}
