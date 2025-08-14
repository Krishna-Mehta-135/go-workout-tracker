package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Krishna-Mehta-135/go-workout-tracker/internal/store"
	"github.com/go-chi/chi/v5"
)

// WorkoutHandler handles all workout-related HTTP requests
type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

// NewWorkoutHandler creates a new WorkoutHandler with the given WorkoutStore
func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
	}
}

// HandleWorkoutByID handles GET /workouts/{id}
func (wh *WorkoutHandler) HandleWorkoutByID(w http.ResponseWriter, r *http.Request) {
	// Extract "id" from URL
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r) // 404 if no ID provided
		return
	}

	// Convert string ID to int64
	workoutId, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.NotFound(w, r) // 404 if invalid ID
		return
	}

	// Fetch workout from store
	workout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	// Send workout as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// HandleCreateWorkout handles POST /workouts
func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() // ensure body is closed after reading

	var workout store.Workout

	// Decode JSON body into Workout struct
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Save workout using store
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		fmt.Println("Store error: ", err)
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	// Return the created workout as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)
}
