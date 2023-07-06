package handlers

import (
	"fmt"
	"net/http"

	"github.com/Hitesh-Nagothu/vault-service/service"
	"go.uber.org/zap"
)

type File struct {
	logger      *zap.Logger
	fileService *service.FileService
	ipfsService *service.IPFSService
}

func NewFile(l *zap.Logger, fs *service.FileService) *File {
	return &File{
		logger:      l,
		fileService: fs,
	}
}

func (handler *File) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		handler.getFile(w, r)
	case http.MethodPost:
		handler.uploadFile(w, r)
	case http.MethodPut:
		handler.updateFile(w, r)
	case http.MethodDelete:
		handler.deleteFile(w, r)
	default:
		handler.logger.Error("Received bad POST request", zap.String("HTTP Method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func (handler *File) uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handler.logger.Error("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		handler.logger.Error("Failed to retrieve file from request", zap.Error(err))
		http.Error(w, "Failed to retrieve file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	userEmailFromContext, _ := r.Context().Value("email").(string)
	if len(userEmailFromContext) == 0 {
		handler.logger.Error("No user email found. Cannot process the file")
		http.Error(w, "No user email found. Cannot process the file", http.StatusBadRequest)
		return
	}

	uploadFileErr := handler.fileService.CreateFile(file, fileHeader, userEmailFromContext)
	if uploadFileErr != nil {
		http.Error(w, "Failed to upload file "+uploadFileErr.Error(), http.StatusBadRequest)
	}

	fmt.Fprint(w, "File upload complete")
}

func (handler *File) getFile(w http.ResponseWriter, r *http.Request) {

}

func (handler *File) updateFile(w http.ResponseWriter, r *http.Request) {

}

func (handler *File) deleteFile(w http.ResponseWriter, r *http.Request) {

}
