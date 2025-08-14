package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Krishna-Mehta-135/go-workout-tracker/internal/api"
	"github.com/Krishna-Mehta-135/go-workout-tracker/internal/migrations"
	"github.com/Krishna-Mehta-135/go-workout-tracker/internal/store"
)

// Application bundles together all core dependencies of the app.
// This avoids using global variables and makes it easier to pass
// dependencies (like logger, DB, handlers) around the codebase.
type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
}

// NewApplication sets up and returns a fully initialized Application instance.
// Responsibilities:
//  1. Connect to the PostgreSQL database.
//  2. Create a logger for consistent logging.
//  3. Initialize request handlers.
//  4. Return the ready-to-use Application.
func NewApplication() (*Application, error) {
	// Connect to PostgreSQL
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Create a logger that writes to stdout with date + time format
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//Initialize stores
	workoutStore := store.NewPostgresWorkoutStore(pgDB)

	// Initialize handlers
	workoutHandler := api.NewWorkoutHandler(workoutStore)

	// Bundle dependencies into Application
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             pgDB,
	}

	return app, nil
}

// HealthCheck is a simple route to verify that the server is running.
func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available")
}
