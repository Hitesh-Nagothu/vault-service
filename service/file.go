package service

import (
	"mime"
	"mime/multipart"
	"strings"

	"github.com/Hitesh-Nagothu/vault-service/data"
	"go.uber.org/zap"
)

type FileService struct {
	repo   *data.FileRepository
	logger *zap.Logger
}

func NewFileService(logger *zap.Logger, repo *data.FileRepository) *FileService {
	return &FileService{
		logger: logger,
		repo:   repo,
	}
}

func (fs *FileService) GetFileType(fileHeader *multipart.FileHeader) []string {
	contentType := fileHeader.Header.Get("Content-Type")
	extension, _ := mime.ExtensionsByType(contentType)
	extensionsWithoutPrefix := make([]string, len(extension))

	if len(extension) > 0 {
		for _, ext := range extension {
			extensionsWithoutPrefix = append(extensionsWithoutPrefix, strings.TrimPrefix(ext, "."))
		}
		return extensionsWithoutPrefix
	}
	return make([]string, 0)
}

func (fs *FileService) IsAllowedFileType(fileTypes []string) (string, bool) {

	for _, fileType := range fileTypes {
		switch fileType {
		case "jpg", "jpeg", "png", "gif", "pdf", "txt", "doc", "docx":
			return fileType, true
		default:
			return "", false
		}
	}

	return "", false
}
