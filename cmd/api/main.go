package main

import (
	"log"
	"net/http"
)

func main() {
	app := NewApp()

	log.Println("Starting Broker service on port 80")

	//define server
	srv := &http.Server{
		Addr:    ":80",
		Handler: app.routes(),
	}

	//start server
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
