package main

import (
	"log"
	"net/http"

	"github.com/Hitesh-Nagothu/vault-service/handlers"
	"github.com/Hitesh-Nagothu/vault-service/middlewares"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	handler := middlewares.NewMiddlewareHandler()
	handler.Use(middlewares.AuthMiddleware)

	fileUploadHandler := handlers.NewFileUpload(logger)
	fileRetrieveHandler := handlers.NewFileRetrieve(logger)

	handler.Handle("/file", fileUploadHandler)
	handler.Handle("/file/all", fileRetrieveHandler)

	serverErr := http.ListenAndServe(":8080", handler)
	if serverErr != nil {
		log.Fatal("Server error: ", err)
	}

}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
