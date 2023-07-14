package main

import (
	"broker/rabbitMQ"
	"log"
	"net/http"
)

func main() {
	connection, err := rabbitMQ.Connect()
	if err != nil {
		log.Fatal(err)
	}
	app := NewApp(connection)

	log.Println("Starting Broker service on port 8080")

	//define server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	//start server
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
