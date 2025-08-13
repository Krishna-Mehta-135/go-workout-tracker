package main

import (
	"FEMProject/internal/app"
	"FEMProject/internal/routes"
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")  //made dynamic port using flags
	flag.Parse()

	app, err := app.NewApplication() //Creates your Application "toolbox" with the logger
	if err != nil {
		panic(err)
	}
	app.Logger.Printf("We are running on port %d\n", port)

	r := routes.SetupRoutes(app)

	//http server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//we assign the server to error because the func serve and listen only returns error, else it keeps running
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
