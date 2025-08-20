package api

import (
	"github.com/Krishna-Mehta-135/go-workout-tracker/internal/store"
	"log"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}
