package app

import (
	"FEMProject/internal/api"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Why we need this: Instead of using global variables, we bundle our app's components (like logger) into a single struct. This is a common Go pattern for organizing application state.
type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

// A logger is a tool that helps you track what's happening in your program by printing messages. Think of it like a diary for your application.
// Instead of using fmt.Println() everywhere, we use a logger
// example->
// We use a logger:
// app.Logger.Println("Server started")


// This is like a constructor because it:
func NewApplication() (*Application, error) {
	// 1. Creates a new logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//stores->

	//handler->
	workoutHandler := api.NewWorkoutHandler()

	// 2. Creates a new Application with that logger
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
	}

	// 3. Returns the ready-to-use Application
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available")
}
