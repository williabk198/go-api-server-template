package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/williabk198/go-api-server-template/controller"
)

// NewRouter maps routes to controller functions and returns the root router
func NewRouter(controls controller.Controller) http.Handler {
	rootRouter := chi.NewRouter()
	rootRouter.Use(middleware.SetHeader("Content-Type", "application/json"))
	rootRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
	}))

	rootRouter.Route("/person", func(r chi.Router) {
		r.Post("/", controls.Person().Add)
		r.Get("/{id}", controls.Person().GetSpecific)
		r.Delete("/{id}", controls.Person().Remove)
		r.Put("/{id}", controls.Person().Update)
	})

	return rootRouter
}
