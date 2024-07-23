package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com.br/gibranct/admin-do-catalogo/internal/domain"
	categoryUseCase "github.com.br/gibranct/admin-do-catalogo/internal/usecases/category"
	"github.com/go-chi/chi/v5"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	ccc := categoryUseCase.CreateCategoryCommand{
		Name:        input.Name,
		Description: input.Description,
	}

	noti, output := app.useCases.Category.Create.Execute(ccc)

	if output != nil {
		app.writeJson(w, http.StatusCreated, envelope{"id": output.ID}, nil)
		return
	}
	err = app.writeError(w, http.StatusBadRequest, "Could not save category", noti)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *application) getCategoryByIdHandler(w http.ResponseWriter, r *http.Request) {
	categoryIdStr := chi.URLParam(r, "id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		app.badRequestResponse(w, errors.New("invalid id"))
		return
	}

	if categoryId <= 0 {
		app.notFoundResponse(w)
		return
	}

	out, err := app.useCases.Category.FindOne.Execute(categoryId)

	if err != nil {
		app.notFoundResponse(w)
		return
	}

	err = app.writeJson(w, http.StatusOK, out, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *application) deactivateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryIdStr := chi.URLParam(r, "id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		app.badRequestResponse(w, errors.New("invalid id"))
		return
	}

	if categoryId <= 0 {
		app.notFoundResponse(w)
		return
	}

	err = app.useCases.Category.Deactivate.Execute(categoryId)

	if err != nil {
		app.notFoundResponse(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) activateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryIdStr := chi.URLParam(r, "id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		app.badRequestResponse(w, errors.New("invalid id"))
		return
	}

	if categoryId <= 0 {
		app.notFoundResponse(w)
		return
	}

	err = app.useCases.Category.Activate.Execute(categoryId)

	if err != nil {
		app.notFoundResponse(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) listCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	perPage, err := strconv.Atoi(r.URL.Query().Get("perPage"))

	if err != nil {
		app.badRequestResponse(w, err)
	}

	query := domain.SearchQuery{
		Sort:      r.URL.Query().Get("sort"),
		Term:      r.URL.Query().Get("search"),
		Page:      page,
		PerPage:   perPage,
		Direction: r.URL.Query().Get("dir"),
	}

	output, err := app.useCases.Category.FindAll.Execute(query)

	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	app.writeJson(w, http.StatusOK, output, nil)
}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, errors.New("invalid id"))
		return
	}
	command := categoryUseCase.UpdateCategoryCommand{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
	}

	noti := app.useCases.Category.Update.Execute(command)

	if !noti.HasErrors() {
		app.writeJson(w, http.StatusOK, envelope{"id": input.ID}, nil)
		return
	}

	err = app.writeError(w, http.StatusBadRequest, "Could not update category", noti)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
