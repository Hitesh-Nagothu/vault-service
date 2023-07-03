package main

import (
	"log"
	"net/http"

	"github.com/Hitesh-Nagothu/vault-service/data"
	"github.com/Hitesh-Nagothu/vault-service/handlers"
	"github.com/Hitesh-Nagothu/vault-service/middlewares"
	"github.com/Hitesh-Nagothu/vault-service/service"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	//db setup
	db := data.GetMongoDBInstance()
	fileRepo := data.NewFileRepository(db, logger)
	userRepo := data.NewUserRepository(db, logger)
	chunkRepo := data.NewChunkRepository(db, logger)

	fileService := service.NewFileService(logger, fileRepo)
	userService := service.NewUserService(logger, userRepo)
	chunkService := service.NewChunkService(logger, chunkRepo)

	fileUploadHandler := handlers.NewFileUpload(logger, fileService)

	handler := middlewares.NewMiddlewareHandler()
	handler.Use(middlewares.AuthMiddleware)

	handler.Handle("/file", fileUploadHandler)
	serverErr := http.ListenAndServe(":8080", handler)
	if serverErr != nil {
		log.Fatal("Server error: ", err)
	}
}
