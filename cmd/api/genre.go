package main

import (
	"errors"
	"net/http"
	"strconv"

	genre_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/genre"
	"github.com/go-chi/chi/v5"
)

func (app *application) createGenreHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string  `json:"name"`
		CategoryIds []int64 `json:"categoryIds"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	command := genre_usecase.CreateGenreCommand{
		Name:        input.Name,
		CategoryIds: &input.CategoryIds,
	}

	noti, output := app.useCases.Genre.Create.Execute(command)

	if output != nil {
		app.writeJson(w, http.StatusCreated, envelope{"id": output.ID}, nil)
		return
	}
	err = app.writeError(w, http.StatusBadRequest, "Could not save genre", noti)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *application) listGenresHandler(w http.ResponseWriter, _ *http.Request) {
	output, err := app.useCases.Genre.FindAll.Execute()

	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	app.writeJson(w, http.StatusOK, output, nil)
}

func (app *application) deleteGenreByIdHandler(w http.ResponseWriter, r *http.Request) {
	genreIdStr := chi.URLParam(r, "id")
	genreId, err := strconv.ParseInt(genreIdStr, 10, 64)
	if err != nil {
		app.badRequestResponse(w, errors.New("invalid id"))
		return
	}

	if genreId <= 0 {
		app.notFoundResponse(w)
		return
	}

	command := genre_usecase.DeleteGenreCommand{
		GenreId: genreId,
	}

	err = app.useCases.Genre.DeleteById.Execute(command)

	if err != nil {
		app.notFoundResponse(w)
		return
	}

	err = app.writeJson(w, http.StatusNoContent, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
