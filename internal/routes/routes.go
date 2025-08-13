package routes

import (
	"FEMProject/internal/app"

	"github.com/go-chi/chi/v5"
)

//we use app from application struct. we write routes from app.funcName
func SetupRoutes(app *app.Application) *chi.Mux{
	r := chi.NewRouter()

	//since Health check func was a method of application struct, we can use it here without importing
	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleWorkoutByID)

	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	return r
}