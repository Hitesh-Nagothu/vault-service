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

	ipfsService := service.NewIPFSService(logger)

	//file
	fileRepo := data.NewFileRepository(db, logger)
	fileService := service.NewFileService(logger, fileRepo)
	fileHandler := handlers.NewFile(logger, fileService, ipfsService)

	//user
	userRepo := data.NewUserRepository(db, logger)
	userService := service.NewUserService(logger, userRepo)
	userHandler := handlers.NewUser(logger, userService)

	handler := middlewares.NewMiddlewareHandler()
	handler.Use(middlewares.AuthMiddleware)
	handler.Handle("/file", fileHandler)
	handler.Handle("/user", userHandler)

	serverErr := http.ListenAndServe(":8080", handler)
	if serverErr != nil {
		log.Fatal("Server error: ", err)
	}
}
