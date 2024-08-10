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
		r.Get("/categories", app.listCategoriesHandler)
		r.Get("/categories/{id}", app.getCategoryByIdHandler)
		r.Put("/categories/{id}", app.updateCategoryHandler)
		r.Post("/categories/{id}/activate", app.activateCategoryHandler)
		r.Post("/categories/{id}/deactivate", app.deactivateCategoryHandler)

		r.Post("/cast-members", app.createCastMemberHandler)
		r.Get("/cast-members", app.listCastMemberHandler)
		r.Put("/cast-members/{id}", app.updateCastMemberHandler)
	})

	return router
}
