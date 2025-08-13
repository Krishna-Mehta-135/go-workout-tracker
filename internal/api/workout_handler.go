package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// This struct defines the type details for the WorkoutHandler. It's empty now, but will later include properties for the store and logger.
type WorkoutHandler struct {
}

//constructor function
func NewWorkoutHandler() *WorkoutHandler { 
	return &WorkoutHandler{}
}

//controllers of type WorkoutHandler struct
func (wh *WorkoutHandler) HandleWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")  //take param from url
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}

	workoutId, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
	}

	fmt.Fprintf(w, "this is the workout id %d\n", workoutId)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Created a workout\n")
}
