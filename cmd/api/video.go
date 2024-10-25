package main

import (
	"net/http"

	video_usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/video"
)

func (app *application) createVideoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		LaunchedAt  int     `json:"yearLaunched"`
		Duration    float64 `json:"duration"`
		Opened      bool    `json:"opened"`
		Published   bool    `json:"published"`
		Rating      string  `json:"rating"`
		CategoryIds []int64 `json:"categoryIds"`
		GenreIds    []int64 `json:"genreIds"`
		MemberIds   []int64 `json:"memberIds"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	command := video_usecase.CreateVideoCommand{
		Title:       input.Title,
		Description: input.Description,
		LaunchedAt:  input.LaunchedAt,
		Duration:    input.Duration,
		Opened:      input.Opened,
		Published:   input.Published,
		Rating:      input.Rating,
		CategoryIds: input.CategoryIds,
		GenreIds:    input.GenreIds,
		MemberIds:   input.MemberIds,
	}

	noti, output := app.useCases.Video.Create.Execute(command)

	if output != nil {
		app.writeJson(w, http.StatusCreated, envelope{"id": output.ID}, nil)
		return
	}
	err = app.writeError(w, http.StatusBadRequest, "Could not save video", noti)
	if err != nil {
		app.serverErrorResponse(w, err)
	}

}
