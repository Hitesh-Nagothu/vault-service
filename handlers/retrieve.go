package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hitesh-Nagothu/vault-service/data"
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

	userData, userExists := userStore.Data[userEmailFromContext]
	if !userExists {
		fmt.Fprintf(w, "No files found for user %s", userEmailFromContext)
		return
	}

	// Encode the UUIDs into JSON
	jsonData, err := json.Marshal(userData.Files)
	if err != nil {
		http.Error(w, "Failed to encode UUIDs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
