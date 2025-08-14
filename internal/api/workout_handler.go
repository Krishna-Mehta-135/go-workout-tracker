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

func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request){
	// Extract "id" from URL
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r) // 404 if no ID provided
		return
	}

	// Convert string ID to int64
	workoutId, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		fmt.Println("Workout doesn't exist")
		http.NotFound(w, r) // 404 if invalid ID
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		http.Error(w, "Failed to fetch workout", http.StatusInternalServerError)
		return
	}
	
	//if workout is not found
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Since all the items in the struct are pointers and the null value of a pointer is nil. So by using pointers , we can check if the items are empty or not, if not empty , we update the value
	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		fmt.Println("Update workout err: " , err)
		http.Error(w, "Failed to update the workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingWorkout)
}