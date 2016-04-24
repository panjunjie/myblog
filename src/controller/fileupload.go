package controller

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"io"
	"mime/multipart"
	"myblog/src/lib"
	"net/http"
	"os"
	"path/filepath"
)

func KindEditorUploadJson(w http.ResponseWriter, req *http.Request) {
	result := make(map[string]interface{})
	u := lib.RandStr()
	//var dir = req.FormValue("dir")

	fileName := u

	//fileTypes := "gif,jpg,jpeg,png,bmp";

	req.ParseMultipartForm(32 << 20)
	file, _, err := req.FormFile("imgFile")
	if file == nil {
		showErr(err, "文件不存在")
		result["error"] = 1
		result["message"] = "没有选择文件"
		r.JSON(w, http.StatusFound, result)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		showErr(err, "图片编码出错")
		result["error"] = 1
		result["message"] = "图片编码出错"
		r.JSON(w, http.StatusFound, result)
		return
	}

	var dst *image.NRGBA
	dst = imaging.Thumbnail(img, 330, 248, imaging.CatmullRom)
	media_path, err := settings.String("admin", "media_path")
	if err != nil {
		showErr(err, "配置里读取文件路径出错")
		result["error"] = 1
		result["message"] = "在配置里读取文件路径出错"
		r.JSON(w, http.StatusFound, result)
		return
	}
	media_url, err := settings.String("admin", "media_url")
	if err != nil {
		showErr(err, "配置里读取文件路径出错")
		result["error"] = 1
		result["message"] = "在配置里读取文件路径出错"
		r.JSON(w, http.StatusFound, result)
		return
	}

	var dir = "article"

	if lib.IsDirExists(filepath.Join(media_path, dir)) == false {
		os.MkdirAll(filepath.Join(media_path, dir), os.ModePerm)
	}

	fileName = filepath.Join(media_path, dir, fileName)

	imaging.Save(dst, fileName+fmt.Sprintf("_%d_%d.jpg", 330, 248))

	dst = imaging.Resize(img, 600, 0, imaging.CatmullRom)
	imaging.Save(dst, fileName+fmt.Sprintf("_%d.jpg", 600))
	result["error"] = 0
	url := fmt.Sprintf("%s%s/%s", media_url, dir, u+"_600.jpg")
	result["url"] = url

	r.JSON(w, http.StatusOK, result)
}

func KindEditorFileManageJson(w http.ResponseWriter, req *http.Request) {
	result := make(map[string]interface{})

	media_path, err := settings.String("admin", "media_path")
	if err != nil {
		showErr(err, "配置里读取文件路径出错")
	}

	media_url, err := settings.String("admin", "media_url")
	if err != nil {
		showErr(err, "配置里读取文件路径出错")
	}

	dir := "article"

	rootPath := filepath.Join(media_path, dir)
	moveupPath := filepath.Join(media_path, "tmp")

	//dirList := Directory.GetDirectories(currentPath);
	//fileList := Directory.GetFiles(currentPath);

	result["moveup_dir_path"] = moveupPath
	result["current_dir_path"] = rootPath
	result["current_url"] = fmt.Sprintf("%s%s/", media_url, dir)

	dirinfo := lib.DirInfo{}
	dirinfo = dirinfo.DirInfoInit(rootPath)

	result["total_count"] = dirinfo.FileNum + dirinfo.DirNum

	file_list := make([]interface{}, 0)

	FileList := dirinfo.FileList //文件列表
	DirList := dirinfo.DirList   //文件夹信息
	/**
	    switch (order){
			case "size":
				Array.Sort(dirList, new NameSorter());
				Array.Sort(fileList, new SizeSorter());
				break;
			case "type":
				Array.Sort(dirList, new NameSorter());
				Array.Sort(fileList, new TypeSorter());
				break;
			case "name":
			default:
				Array.Sort(dirList, new NameSorter());
				Array.Sort(fileList, new NameSorter());
				break;
		}
		**/

	for _, v := range DirList {
		fileInfo := make(map[string]interface{})
		fileInfo["is_dir"] = true
		fileInfo["has_file"] = true
		fileInfo["filesize"] = 0
		fileInfo["is_photo"] = false
		fileInfo["filetype"] = ""
		fileInfo["filename"] = v.(os.FileInfo).Name()
		fileInfo["datetime"] = v.(os.FileInfo).ModTime().Format(TIMEFORMAT)
		file_list = append(file_list, fileInfo)
	}

	for _, v := range FileList {
		fileInfo := make(map[string]interface{})
		fileInfo["is_dir"] = false
		fileInfo["has_file"] = false
		fileInfo["filesize"] = v.(os.FileInfo).Size()
		fileInfo["is_photo"] = true
		fileInfo["filetype"] = "jpg"
		fileInfo["filename"] = v.(os.FileInfo).Name()
		fileInfo["datetime"] = v.(os.FileInfo).ModTime().Format(TIMEFORMAT)
		file_list = append(file_list, fileInfo)
	}

	result["file_list"] = file_list

	r.JSON(w, http.StatusOK, result)
}

//内部方法 上传文摘的小图片
func uploadArticlePic(file multipart.File) string {
	if file != nil {
		img, _, _ := image.Decode(file)
		var dst = imaging.Thumbnail(img, 330, 248, imaging.CatmullRom)

		media_path, err := settings.String("admin", "media_path")
		if err != nil {
			showErr(err, "配置里读取文件路径出错")
		}
		media_url, err := settings.String("admin", "media_url")
		if err != nil {
			showErr(err, "配置里读取文件路径出错")
		}

		var dir = "article"

		if lib.IsDirExists(filepath.Join(media_path, dir)) == false {
			os.MkdirAll(filepath.Join(media_path, dir), os.ModePerm)
		}

		u := lib.RandStr()
		fileName := u

		fileName = filepath.Join(media_path, dir, fileName)

		err = imaging.Save(dst, fileName+fmt.Sprintf("_%d_%d.jpg", 330, 248))

		if err != nil {
			fmt.Printf("保存文摘小图片出错：%s", err)
		} else {
			return fmt.Sprintf("%s%s/%s", media_url, dir, fmt.Sprintf("%s_%d_%d.jpg", u, 330, 248))
		}
	}
	return ""
}

//内部方法 上传头像的小图片
func uploadHatHead(file multipart.File) string {
	if file != nil {
		img, _, _ := image.Decode(file)
		var dst = imaging.Resize(img, 250, 0, imaging.Lanczos) //imaging.Thumbnail(img, 250, 250, imaging.CatmullRom)

		media_path, err := settings.String("admin", "media_path")
		if err != nil {
			showErr(err, "配置里读取文件路径出错")
		}
		media_url, err := settings.String("admin", "media_url")
		if err != nil {
			showErr(err, "配置里读取文件路径出错")
		}

		var dir = "account"

		if lib.IsDirExists(filepath.Join(media_path, dir)) == false {
			os.MkdirAll(filepath.Join(media_path, dir), os.ModePerm)
		}

		u := lib.RandStr()
		fileName := u

		fileName = filepath.Join(media_path, dir, fileName)

		err = imaging.Save(dst, fileName+fmt.Sprintf("_%d.jpg", 250))

		if err != nil {
			fmt.Printf("保存头像图片出错：%s", err)
		} else {
			return fmt.Sprintf("%s%s/%s", media_url, dir, fmt.Sprintf("%s_%d.jpg", u, 250))
		}
	}
	return ""
}

func uploadDonatePng(file multipart.File) string {
	if file != nil {
		img, _, _ := image.Decode(file)
		var dst = imaging.Resize(img, 250, 250, imaging.Lanczos) //imaging.Thumbnail(img, 250, 250, imaging.CatmullRom)

		media_path, err := settings.String("admin", "media_path")
		if err != nil {
			showErr(err, "配置里读取文件路径出错")
		}
		media_url, err := settings.String("admin", "media_url")
		if err != nil {
			showErr(err, "配置里读取文件路径出错")
		}

		dir := "account"

		if lib.IsDirExists(filepath.Join(media_path, dir)) == false {
			os.MkdirAll(filepath.Join(media_path, dir), os.ModePerm)
		}

		u := lib.RandStr()
		fileName := u

		fileName = filepath.Join(media_path, dir, fileName)

		err = imaging.Save(dst, fileName+fmt.Sprintf("_%d.png", 250))

		if err != nil {
			fmt.Printf("保存头像图片出错：%s", err)
		} else {
			return fmt.Sprintf("%s%s/%s", media_url, dir, fmt.Sprintf("%s_%d.png", u, 250))
		}
	}
	return ""
}

//上传文件
func uploadFile(file *multipart.FileHeader, folder string) (string, error) {
	path := ""
	src, err := file.Open()

	if err != nil {
		showErr(err, "上传文件出错")
		return path, err
	}
	defer src.Close()

	media_path, err := settings.String("admin", "media_path")
	if err != nil {
		showErr(err, "在配置里读取文件路径出错")
		return path, err
	}

	media_url, err := settings.String("admin", "media_url")
	if err != nil {
		showErr(err, "在配置里读取文件路径出错")
		return path, err
	}

	if lib.IsDirExists(filepath.Join(media_path, folder)) == false {
		os.MkdirAll(filepath.Join(media_path, folder), os.ModePerm)
	}

	fileExt := lib.GetFileExt(file.Filename)
	fileName := lib.CreateFileName() + fileExt

	filePath := filepath.Join(media_path, folder, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		showErr(err, "文件创建出错")
		return path, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		showErr(err, "上传拷贝出错")
		return path, err
	}

	path = media_url + folder + "/" + fileName

	return path, nil
}
