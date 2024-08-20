package main

import (
	"net/http"

	genre_usecase "github.com.br/gibranct/admin-do-catalogo/internal/usecases/genre"
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
