package main

import (
	"errors"
	"net/http"
	"strconv"

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
		app.serverErrorResponse(w)
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
		app.serverErrorResponse(w)
	}
}
