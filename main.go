package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Hitesh-Nagothu/vault-service/handlers"
	"github.com/Hitesh-Nagothu/vault-service/middlewares"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Create a new Zap logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	fileUploadHandler := handlers.NewFileUpload(logger)
	fileRetrieveHandler := handlers.NewFileRetrieve(logger)

	sm := http.NewServeMux()
	sm.Handle("/file", middlewares.AuthMiddleware(logger, fileUploadHandler))
	sm.Handle("/file/all", middlewares.AuthMiddleware(logger, fileRetrieveHandler))

	server := http.Server{
		Handler: sm,
		Addr:    ":8080",
	}

	fmt.Println("Server starting on port 8080")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
