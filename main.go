package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Krishna-Mehta-135/go-workout-tracker/internal/app"
	"github.com/Krishna-Mehta-135/go-workout-tracker/internal/routes"
)

func main() {
	// Default port
	port := 8080

	// CLI flag to override port
	flag.IntVar(&port, "port", port, "Go backend server port")
	flag.Parse()

	// Override with environment variable if present
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		if p, err := strconv.Atoi(portEnv); err == nil {
			port = p
		}
	}

	// Initialize application (DB, logger, handlers)
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	app.Logger.Printf("Server starting on port %d\n", port)

	// Configure HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes.SetupRoutes(app),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start server
	if err := server.ListenAndServe(); err != nil {
		app.Logger.Fatal(err)
	}
}
