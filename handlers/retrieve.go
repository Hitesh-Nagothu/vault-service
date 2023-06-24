package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hitesh-Nagothu/vault-service/data"
	"github.com/Hitesh-Nagothu/vault-service/ipfs"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FileRetrieve struct {
	logger *zap.Logger
}

func NewFileRetrieve(l *zap.Logger) *FileRetrieve {
	return &FileRetrieve{l}
}

func (fr *FileRetrieve) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userEmailFromContext, _ := r.Context().Value("email").(string)
	if len(userEmailFromContext) == 0 {
		http.Error(w, "No user email found. Unauthorized request to fetch files", http.StatusBadRequest)
		return
	}

	userStore := data.GetUserStore()
	fileStore := data.GetFileStore()
	chunkStore := data.GetChunkStore()

	userData, userExists := userStore.Data[userEmailFromContext]
	if !userExists {
		fmt.Fprintf(w, "No files found for user %s", userEmailFromContext)
		return
	}

	var response = make(map[string]interface{})
	for _, fileId := range userData.Files {
		fileMetadata := fileStore.Data[fileId]
		var chunkHashes []string
		for _, chunkId := range fileMetadata.Chunks {
			chunkHashes = append(chunkHashes, chunkStore.Data[chunkId])
		}

		response[fileMetadata.Name] = map[string]interface{}{
			"mime": fileMetadata.MimeType,
			"cids": chunkHashes,
		}
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func getContentForFile(chunkIds []uuid.UUID) ([][]byte, error) {
	chunkStore := data.GetChunkStore()
	ipfsInstance := ipfs.GetIPFSInstance()

	var cidsOfFile []string
	for _, id := range chunkIds {
		cidsOfFile = append(cidsOfFile, chunkStore.Data[id])
	}

	fmt.Println(cidsOfFile)

	var response [][]byte
	for _, cid := range cidsOfFile {
		chunkData, err := ipfsInstance.GetContent(cid)
		if err != nil {
			fmt.Println("Failed to get content for a cid")
			return nil, err
		}
		response = append(response, chunkData)
	}

	return response, nil
}
