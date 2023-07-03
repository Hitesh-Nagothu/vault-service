package handlers

import (
	"io"
	"net/http"

	"github.com/Hitesh-Nagothu/vault-service/service"
	"go.uber.org/zap"
)

type FileUpload struct {
	logger  *zap.Logger
	service *service.FileService
}

func NewFileUpload(l *zap.Logger, s *service.FileService) *FileUpload {
	return &FileUpload{
		logger:  l,
		service: s,
	}
}

type ChunkedWriter struct {
	file io.Writer
}

func NewChunkedWriter(file io.Writer) *ChunkedWriter {
	return &ChunkedWriter{
		file: file,
	}
}

func (fh *FileUpload) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
