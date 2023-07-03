package handlers

import (
	"net/http"

	"github.com/Hitesh-Nagothu/vault-service/service"
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
		handler.getUser(w, r)
	case http.MethodPost:
		handler.createUser(w, r)
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

func (handler *User) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handler.logger.Error("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userEmailFromContext, _ := r.Context().Value("email").(string)
	if len(userEmailFromContext) == 0 {
		handler.logger.Error("No user email found. Cannot process the file")
		http.Error(w, "Something went wrong. Cannot process the file", http.StatusBadRequest)
		return
	}

}

func (handler *User) getUser(w http.ResponseWriter, r *http.Request) {

}

func (handler *User) updateUser(w http.ResponseWriter, r *http.Request) {

}

func (handler *User) deletUser(w http.ResponseWriter, r *http.Request) {

}
