package handlers

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Hitesh-Nagothu/vault-service/data"
	"github.com/Hitesh-Nagothu/vault-service/ipfs"
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
	ChunkSize   = 5 * 1024 * 1024 // Chunk size: 1 MB
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
	fileTypes := getFileType(fileHeader)

	if len(fileTypes) == 0 {
		http.Error(w, "Failed to infer the type of file upload.", http.StatusBadRequest)
		return
	}

	fileType, isAllowed := isAllowedFileType(fileTypes)
	if !isAllowed {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	userEmailFromContext, _ := r.Context().Value("email").(string)
	if len(userEmailFromContext) == 0 {
		http.Error(w, "No user email found. Cannot process the file", http.StatusBadRequest)
		return
	}

	filebytes, readErr := io.ReadAll(file)
	if readErr != nil {
		http.Error(w, "Failed to read file", http.StatusNotAcceptable)
		return
	}

	chunksOfFiles := SplitIntoChunks(filebytes)

	chunkIdToChunkData := make(map[uuid.UUID][]byte)
	for _, chunk := range chunksOfFiles {
		newChunkId := uuid.New()
		chunkIdToChunkData[newChunkId] = chunk
	}

	chunkIdToHash, ipfsErr := GetIPFSHashForChunks(chunkIdToChunkData)
	if ipfsErr != nil {
		http.Error(w, "Failed to generate hashes for file data", http.StatusExpectationFailed)
		return
	}

	chunkIdsForFile := make([]uuid.UUID, 0)

	chunkStore := data.GetChunkStore()
	for chunkId, chunkHash := range chunkIdToHash {
		chunkStore.Data[chunkId] = chunkHash
		chunkIdsForFile = append(chunkIdsForFile, chunkId)
	}

	fileUuid := uuid.New()
	fileName := fileHeader.Filename
	fileChunks := chunkIdsForFile

	fileMetadataObj := data.FileMetadata{
		Name:     fileName,
		Chunks:   fileChunks,
		MimeType: fileType,
	}

	fileStoreInstance := data.GetFileStore()
	userStoreInstance := data.GetUserStore()

	fileStoreInstance.Data[fileUuid] = fileMetadataObj
	userMetadata := userStoreInstance.Data[userEmailFromContext]

	userMetadata.Files = append(userStoreInstance.Data[userEmailFromContext].Files, fileUuid)
	userStoreInstance.Data[userEmailFromContext] = userMetadata

	fmt.Fprintln(w, "File uploaded successfully")
}

func getFileType(fileHeader *multipart.FileHeader) []string {
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

func isAllowedFileType(fileTypes []string) (string, bool) {

	for _, fileType := range fileTypes {
		switch fileType {
		case "jpg", "jpeg", "png", "gif", "pdf", "txt", "doc", "docx":
			return fileType, true
		default:
			//DO nothing
		}
	}

	return "", false
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

func GetIPFSHashForChunks(chunkData map[uuid.UUID][]byte) (map[uuid.UUID]string, error) {

	ipfsInstance := ipfs.GetIPFSInstance()

	chunkIdToHash := make(map[uuid.UUID]string)
	for chunkId, chunkBytes := range chunkData {

		_, err := rand.Read(chunkBytes)
		if err != nil {
			fmt.Println("Failed to generate random bytes:", err)
			return nil, errors.New("failed to generate random bytes")
		}

		hash, err := ipfsInstance.AddContent(chunkBytes)
		if err != nil {
			//TODO handle partial upload errors
			fmt.Println("Failed to add a chunk to ipfs %w", err)
			return nil, err
		}
		chunkIdToHash[chunkId] = hash
		fmt.Println(hash)
	}

	return chunkIdToHash, nil
}
