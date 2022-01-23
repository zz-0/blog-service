package util

import (
	"blog-service/global"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/go-programming-tour-book/blog-service/pkg/util"
)

type FileType int

const TypeImage FileType = iota + 1

//获取文件名称   MD5加密
func GetFileName(name string) string {
	ext := GetFileExt(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

//获取文件后缀
func GetFileExt(name string) string {
	return path.Ext(name)
}

//获取文件保存地址
func GetSavePath() string {
	return global.AppSetting.UploadSavePath
}

//检查保存目录是否存在
func CheckSavePath(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsNotExist(err)
}

//检查文件后缀是否包含在约定的后缀配置项中
func CheckContainExt(t FileType, name string) bool {
	ext := GetFileExt(name)
	switch t {
	case TypeImage:
		for _, allowExt := range global.AppSetting.UploadImageAllowExts {
			if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
				return true
			}
		}
	}
	return false
}

//检查文件大小是否超出限制
func CheckMaxSize(t FileType, f multipart.File) bool {
	content, _ := ioutil.ReadAll(f)
	size := len(content)
	switch t {
	case TypeImage:
		if size >= global.AppSetting.UploadImageMaxSize*1024*1024 {
			return true
		}
	}
	return false
}

//检查文件权限是否足够
func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

//创建保存上传文件的目录
func CreateSavePath(dst string, perm os.FileMode) error {
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}
	return nil
}

//保存上传文件
func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
