package service

import (
	"blog-service/global"
	"errors"
	"mime/multipart"
	"os"

	"github.com/go-programming-tour-book/blog-service/pkg/upload"
)

type FileInfo struct {
	Name      string
	AccessUrl string
}

//上传静态文件
func (svc *Service) UploadFile(fileTpye upload.FileType, file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error) {
	fileName := upload.GetFileName(fileHeader.Filename)

	//获取静态文件的储存地址
	uploadSavePath := upload.GetSavePath()
	dst := uploadSavePath + "/" + fileName
	if !upload.CheckContainExt(fileTpye, fileName) {
		return nil, errors.New("file suffix is not supported.")
	}
	if upload.CheckSavePath(uploadSavePath) {
		if err := upload.CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
			return nil, errors.New("failed to create save directory.")
		}
	}
	if upload.CheckMaxSize(fileTpye, file) {
		return nil, errors.New("exceeded maximum file limit.")
	}
	if upload.CheckPermission(uploadSavePath) {
		return nil, errors.New("insufficient file permissions.")
	}
	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}

	//获取路径
	accessUrl := global.AppSetting.UploadServerUrl + "/" + fileName
	return &FileInfo{Name: fileName, AccessUrl: accessUrl}, nil
}
