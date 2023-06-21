package handlers

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Hitesh-Nagothu/vault-service/data"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FileUpload struct {
	logger *zap.Logger
}

func NewFileUpload(l *zap.Logger) *FileUpload {
	return &FileUpload{l}
}

type ChunkedWriter struct {
	file io.Writer
}

func NewChunkedWriter(file io.Writer) *ChunkedWriter {
	return &ChunkedWriter{
		file: file,
	}
}

const (
	MaxFileSize = 5 * 1024 * 1024 // Maximum file size: 5 MB
	ChunkSize   = 1 * 1024 * 1024 // Chunk size: 1 MB
)

func (fh *FileUpload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if the file is an allowed type
	fileType := getFileType(fileHeader)
	if !isAllowedFileType(fileType) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	fileStoreInstace := data.GetFileStore()

	fileUuid := uuid.New()
	fileName := fileHeader.Filename
	fileChunks := []uuid.UUID{}

	fileMetadataObj := data.FileMetadata{
		Name:     fileName,
		Chunks:   fileChunks,
		MimeType: fileType,
	}

	fileStoreInstace = data.GetFileStore()
	fileStoreInstace.Data[fileUuid] = fileMetadataObj

	fmt.Fprintln(w, "File uploaded successfully")
}

func getFileType(fileHeader *multipart.FileHeader) string {
	contentType := fileHeader.Header.Get("Content-Type")
	extension, _ := mime.ExtensionsByType(contentType)
	if len(extension) > 0 {
		return strings.TrimPrefix(extension[0], ".")
	}
	return ""
}

func isAllowedFileType(fileType string) bool {
	switch fileType {
	case "jpg", "jpeg", "png", "gif", "pdf", "txt", "doc", "docx":
		return true
	default:
		return false
	}
}

func SplitIntoChunks(data []byte) [][]byte {
	numChunks := (len(data) + ChunkSize - 1) / ChunkSize
	chunks := make([][]byte, numChunks)

	for i := 0; i < numChunks; i++ {
		start := i * ChunkSize
		end := start + ChunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks[i] = data[start:end]
	}

	return chunks
}
