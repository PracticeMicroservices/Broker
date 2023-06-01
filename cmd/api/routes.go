package main

import (
	"broker/cmd/api/controllers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type App struct {
	Broker controllers.Broker
}

func NewApp() *App {
	return &App{
		Broker: controllers.NewBrokerController(),
	}
}

func (app *App) routes() http.Handler {
	mux := chi.NewRouter()
	//cors
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Use(middleware.Heartbeat("/healthCheck"))
	//add routes
	mux.Post("/broker", app.Broker.Broker)

	mux.Post("/handle", app.Broker.HandleSubmission)
	return mux
}
