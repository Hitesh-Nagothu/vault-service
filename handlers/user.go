package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Hitesh-Nagothu/vault-service/service"
	"github.com/Hitesh-Nagothu/vault-service/utility"
	"go.uber.org/zap"
)

type User struct {
	logger      *zap.Logger
	userService *service.UserService
}

func NewUser(logger *zap.Logger, userService *service.UserService) *User {
	return &User{
		logger:      logger,
		userService: userService,
	}
}

func (handler *User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handler.GetUser(w, r)
	case http.MethodPost:
		handler.CreateUser(w, r)
	case http.MethodPut:
		handler.updateUser(w, r)
	case http.MethodDelete:
		handler.deletUser(w, r)
	default:
		handler.logger.Error("Received bad POST request", zap.String("HTTP Method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (handler *User) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handler.logger.Error("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userEmailFromContext, _ := r.Context().Value("email").(string)
	if len(userEmailFromContext) == 0 {
		handler.logger.Error("No user email found. Failed authentication")
		http.Error(w, "Something went wrong. Failed to identify user", http.StatusBadRequest)
		return
	}

	//create the user
	user, err := handler.userService.CreateUser(userEmailFromContext)
	if err != nil {
		handler.logger.Error("Failed to create a new user", zap.String("email", userEmailFromContext), zap.Error(err))
		http.Error(w, "Something went wrong. Try again", http.StatusInternalServerError)
		return
	}

	// Convert User object to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		handler.logger.Error("Failed to encode user data", zap.String("email", user.Email))
		http.Error(w, "Failed to encode user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(userJSON)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (handler *User) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handler.logger.Error("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userEmailFromContext, _ := r.Context().Value("email").(string)
	if len(userEmailFromContext) == 0 {
		handler.logger.Error("No user email found. Failed authentication")
		http.Error(w, "Something went wrong. Failed to identify user", http.StatusBadRequest)
		return
	}

	user, err := handler.userService.GetUser(userEmailFromContext)
	if err != nil {
		handler.logger.Error("User not found", zap.String("email", userEmailFromContext))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if utility.IsStructEmpty(user) {
		//empty user, create one
		createdUser, err := handler.userService.CreateUser(userEmailFromContext)
		user = createdUser
		if err != nil {
			http.Error(w, "Somethign went wrong", http.StatusInternalServerError)
			return
		}
	}

	// Convert User object to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		handler.logger.Error("Failed to encode user data", zap.String("email", user.Email))
		http.Error(w, "Failed to encode user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(userJSON)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (handler *User) updateUser(w http.ResponseWriter, r *http.Request) {

}

func (handler *User) deletUser(w http.ResponseWriter, r *http.Request) {

}
