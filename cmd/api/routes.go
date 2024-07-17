package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Route("/v1", func(r chi.Router) {
		r.Post("/categories", app.createCategoryHandler)
		r.Get("/categories/{id}", app.getCategoryByIdHandler)
		r.Post("/categories/{id}/activate", app.activateCategoryHandler)
		r.Post("/categories/{id}/deactivate", app.deactivateCategoryHandler)
	})

	return router
}
