package service

import (
	"forum/internal/utils"
	"mime/multipart"
)

type UploadService struct {
	uploadDir string
}

func NewUploadService(uploadDir string) *UploadService {
	return &UploadService{uploadDir: uploadDir}
}

func (s *UploadService) SaveImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	return utils.SaveUploadedFile(file, header, s.uploadDir)
}
