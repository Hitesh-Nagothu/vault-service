package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Hitesh-Nagothu/vault-service/service"
	"go.uber.org/zap"
)

type File struct {
	logger      *zap.Logger
	fileService *service.FileService
	ipfsService *service.IPFSService
}

func NewFile(l *zap.Logger, fs *service.FileService, ipfs *service.IPFSService) *File {
	return &File{
		logger:      l,
		fileService: fs,
		ipfsService: ipfs,
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

	fileTypes := handler.fileService.GetFileType(fileHeader)

	if len(fileTypes) == 0 {
		handler.logger.Error("Failed to infer the type of file upload")
		http.Error(w, "Failed to infer the type of file upload", http.StatusBadRequest)
		return
	}

	fileType, isAllowed := handler.fileService.IsAllowedFileType(fileTypes)
	if !isAllowed {
		handler.logger.Error("Invalid file type", zap.Strings("allowed_types", fileTypes))
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	userEmailFromContext, _ := r.Context().Value("email").(string)
	if len(userEmailFromContext) == 0 {
		handler.logger.Error("No user email found. Cannot process the file")
		http.Error(w, "No user email found. Cannot process the file", http.StatusBadRequest)
		return
	}

	filebytes, readErr := io.ReadAll(file)
	if readErr != nil {
		handler.logger.Error("Failed to read file", zap.Error(readErr))
		http.Error(w, "Failed to read file", http.StatusNotAcceptable)
		return
	}

	fileHash, ipfsHashErr := handler.getIPFSHashForFile(filebytes)
	if ipfsHashErr != nil {
		handler.logger.Error("Something went wrong getting ipfs hash for file", zap.Error(ipfsHashErr))
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	fmt.Println(fileHash)

	handler.logger.Info("File uploaded successfully",
		zap.String("user_email", userEmailFromContext),
		zap.String("file_type", fileType))
}

func (handler *File) getFile(w http.ResponseWriter, r *http.Request) {

}

func (handler *File) updateFile(w http.ResponseWriter, r *http.Request) {

}

func (handler *File) deleteFile(w http.ResponseWriter, r *http.Request) {

}

func (handler *File) getIPFSHashForFile(fileData []byte) (string, error) {

	hash, err := handler.ipfsService.AddContent(fileData)
	if err != nil {
		//TODO handle partial upload errors
		fmt.Println("Failed to add a chunk to ipfs %w", err)
		return "", err
	}

	return hash, nil
}
