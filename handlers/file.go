package handlers

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"
)

type FileUpload struct {
	logger *zap.Logger
}

const (
	MaxFileSize = 10 * 1024 * 1024 // Maximum file size: 10 MB
	ChunkSize   = 4 * 1024 * 1024  // Chunk size: 4 MB
)

func NewFileUpload(l *zap.Logger) *FileUpload {
	return &FileUpload{l}
}

func (fh *FileUpload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	/*
		Pull the user email
		Find if the user email exists in store
		if user does not exist create a user

		get the user uuid
		create a file uuid for the file and store the key  in map
		the values is an object of file metadata key and array of chunk uuids
		create a chunk map where property value is ipfs hash

	*/

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

	// Create a new file on the server to store the uploaded file
	uploadedFile, err := os.Create(fileHeader.Filename)
	if err != nil {
		http.Error(w, "Failed to create file on the server", http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	// Chunk and write the file data to the server
	chunkedWriter := NewChunkedWriter(uploadedFile, ChunkSize)
	_, err = io.Copy(chunkedWriter, file)
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

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

type ChunkedWriter struct {
	file    io.Writer
	remains int64
}

func NewChunkedWriter(file io.Writer, chunkSize int64) *ChunkedWriter {
	return &ChunkedWriter{
		file:    file,
		remains: chunkSize,
	}
}

func (w *ChunkedWriter) Write(data []byte) (int, error) {
	length := int64(len(data))
	if length > w.remains {
		data = data[:w.remains]
		length = w.remains
	}

	written, err := w.file.Write(data)
	if err != nil {
		return written, err
	}

	w.remains -= length
	if w.remains <= 0 {
		return written, io.EOF
	}

	return written, nil
}
