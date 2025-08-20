package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/Krishna-Mehta-135/go-workout-tracker/internal/store"
	"github.com/Krishna-Mehta-135/go-workout-tracker/internal/utils"
)

// registerUserRequest represents the incoming request payload
// when a user registers through the API.
// The struct tags (`json:"..."`) tell the JSON decoder how to map
// JSON keys to struct fields.
type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

// UserHandler is an HTTP handler that deals with user-related endpoints.
// It depends on:
// - userStore: interface for persistence (DB operations for users)
// - logger: logging errors & info
type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

// NewUserHandler is a constructor for UserHandler.
// It takes in a userStore and logger and returns a handler instance.
func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

// validateRegisterRequest checks if the incoming register request is valid.
// It ensures required fields are provided and meet basic format rules.
func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if len(req.Username) > 50 {
		return errors.New("username cannot be greater than 50 chars")
	}

	if req.Email == "" {
		return errors.New("Email is required")
	}

	// Use regex to validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("Invalid email format")
	}

	if req.Password == "" {
		return errors.New("Password is empty")
	}

	return nil
}

// HandleRegisterUser handles HTTP POST requests to register a new user.
// Steps:
// 1. Decode incoming JSON into a struct.
// 2. Validate request payload.
// 3. Create a User model and hash the password.
// 4. Store the user in the database.
// 5. Return a JSON response with status.
func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	// Decode JSON body into request struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// Panicf logs and crashes the app (be careful with this in prod!)
		h.logger.Printf("Error: decoding register request :%v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	// Validate request fields
	err = h.validateRegisterRequest(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	// Create a new User model (domain object)
	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	// Bio is optional
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	// Hash the password securely before storing
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("Error: hashing password error %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	// Save user in the database via the store layer
	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("Error: User creation err %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	// Respond with created user (JSON-encoded)
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}
