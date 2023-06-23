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

	userData, userExists := userStore.Data[userEmailFromContext]
	if !userExists {
		fmt.Fprintf(w, "No files found for user %s", userEmailFromContext)
		return
	}

	sampelFileId := userData.Files[0]
	sampleChunkIds := fileStore.Data[sampelFileId].Chunks

	dataChunks, ipfsReadErr := getContentForFile(sampleChunkIds)
	if ipfsReadErr != nil {
		return
	}

	// Encode the UUIDs into JSON
	jsonData, err := json.Marshal(dataChunks)
	if err != nil {
		http.Error(w, "Failed to encode UUIDs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(jsonData)
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
