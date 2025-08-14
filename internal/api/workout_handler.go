package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Krishna-Mehta-135/go-workout-tracker/internal/store"
	"github.com/Krishna-Mehta-135/go-workout-tracker/internal/utils"
	"github.com/go-chi/chi/v5"
)

// WorkoutHandler handles all workout-related HTTP requests
type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

// NewWorkoutHandler creates a new WorkoutHandler with the given WorkoutStore
func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger: logger,
	}
}

// HandleWorkoutByID handles GET /workouts/{id}
func (wh *WorkoutHandler) HandleWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout id"})
		return
	}

	// Fetch workout from store
	workout, err := wh.workoutStore.GetWorkoutById(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutBYID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Invalid server error"})
	}

	utils.WriteJSON(w,http.StatusOK, utils.Envelope{"workout" : workout })
}

// HandleCreateWorkout handles POST /workouts
func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() // ensure body is closed after reading

	var workout store.Workout

	// Decode JSON body into Workout struct
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: decodingCreateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	// Save workout using store
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: createWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	// Extract "id" from URL
	workoutID, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout update id"})
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutById(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
		wh.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
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
		wh.logger.Printf("ERROR: updatingWorkout: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

func (wh *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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

	err = wh.workoutStore.DeleteWorkout(workoutId)
	if err == sql.ErrNoRows {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Error deleting workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
