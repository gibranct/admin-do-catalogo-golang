package main

import (
	"net/http"

	category_usecase "github.com.br/gibranct/admin-do-catalogo/internal/usecases/category"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsActive    bool   `json:"active"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	ccc := category_usecase.CreateCategoryCommand{
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
