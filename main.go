package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lnquy/eventsourcing/config"
	"github.com/lnquy/eventsourcing/router"
	"github.com/rs/cors"
	"net/http"
)

func main() {
	// Configuration
	cfg := config.LoadEnvConfig()

	// New HTTP router
	r := chi.NewRouter()
	// Middlewares
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Recoverer)
	// HTTP routing
	r.Route("/api/v1/es/people", func(r chi.Router) {
		r.Post("/", router.CreatePerson)
		r.Get("/{pid}", router.GetPerson)
		r.Patch("/{pid}", router.UpdatePerson)
	})

	// Start server
	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: cors.Default().Handler(r), // CORS
	}
	logrus.Infof("server: Serving REST HTTP server at %s", cfg.ServerAddr)
	if err := server.ListenAndServe(); err != nil {
		logrus.Errorf("server: Failed to serve HTTP server: %s", err.Error())
	}
}
