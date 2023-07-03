package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Hitesh-Nagothu/vault-service/data"
	"github.com/Hitesh-Nagothu/vault-service/handlers"
	"github.com/Hitesh-Nagothu/vault-service/middlewares"
	"github.com/Hitesh-Nagothu/vault-service/service"
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

	//db setup
	db := data.GetMongoDBInstance()
	fileRepo := data.NewFileRepository(db)
	userRepo := data.NewUserRepository(db)
	chunkRepo := data.NewChunkRepository(db)

	fileService := service.NewFileService(logger, fileRepo)
	userService := service.NewUserService(logger, userRepo)
	chunkService := service.NewChunkService(logger, chunkRepo)

	fileUploadHandler := handlers.NewFileUpload(logger, fileService)

	sm := http.NewServeMux()
	sm.Handle("/file", middlewares.AuthMiddleware(logger, fileUploadHandler))

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
