package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/Krishna-Mehta-135/go-workout-tracker/internal/app"
	"github.com/Krishna-Mehta-135/go-workout-tracker/internal/routes"
)

func main() {
	// CLI flag for port: go run main.go -port=9090
	var port int
	flag.IntVar(&port, "port", 8080, "Go backend server port")
	flag.Parse()

	// Initialize application
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	
	//Closes Db at the end
	defer app.DB.Close()

	app.Logger.Printf("Server starting on port %d\n", port)

	// Configure and start HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes.SetupRoutes(app),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		app.Logger.Fatal(err)
	}
}
