package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Hitesh-Nagothu/vault-service/data"
	"github.com/Hitesh-Nagothu/vault-service/handlers"
	"github.com/Hitesh-Nagothu/vault-service/middlewares"
	"github.com/Hitesh-Nagothu/vault-service/service"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Database struct {
		URL  string `mapstructure:"url"`
		Name string `mapstructure:"name"`
	} `mapstructure:"database"`
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

func main() {

	env := flag.String("env", "default", "The environment to run the server in")
	flag.Parse()

	configFile := fmt.Sprintf("config/%s/%s.yaml", *env, *env)

	viper.Set("env", *env)
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)

	fmt.Println("Reading config from", viper.ConfigFileUsed())
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found")
		} else {
			fmt.Println("Error reading config file:", err)
		}
		log.Fatal("Failed to read config")
	}

	var config Config
	unmarshallErr := viper.Unmarshal(&config)
	if unmarshallErr != nil {
		log.Fatal("Failed to unmarshal configuration")
	}

	fmt.Println("Config loaded successfully")

	logger, err := GetLogger(viper.GetString("env"))
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	//db setup
	db := data.GetMongoDBInstance(config.Database.URL)

	//ipfs
	ipfsService := service.NewIPFSService(logger)

	//chunk
	chunkRepo := data.NewChunkRepository(db, logger)
	chunkService := service.NewChunkService(logger, chunkRepo)

	//user
	userRepo := data.NewUserRepository(db, logger)
	userService := service.NewUserService(logger, userRepo)
	userHandler := handlers.NewUser(logger, userService)

	//file
	fileRepo := data.NewFileRepository(db, logger)
	fileService := service.NewFileService(logger, fileRepo, ipfsService, chunkService, userService)
	fileHandler := handlers.NewFile(logger, fileService)

	handler := middlewares.NewMiddlewareHandler()
	handler.Use(middlewares.AuthMiddleware)
	handler.Handle("/file", fileHandler)
	handler.Handle("/user", userHandler)

	serverAddr := fmt.Sprintf(":%s", strconv.Itoa(config.Server.Port))
	serverErr := http.ListenAndServe(serverAddr, handler)
	if serverErr != nil {
		log.Fatal("Server error: ", err)
	}
}

func GetLogger(env string) (*zap.Logger, error) {
	switch env {
	case "default", "alp":
		return zap.NewDevelopment()
	case "bet", "prd":
		return zap.NewProduction()
	default:
		log.Fatal("unknown env found ")
	}
	return nil, errors.New("Failed to identify a logger")
}
