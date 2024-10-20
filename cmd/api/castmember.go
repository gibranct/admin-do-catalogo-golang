package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com.br/gibranct/admin_do_catalogo/internal/domain"
	"github.com.br/gibranct/admin_do_catalogo/internal/domain/castmember"
	castmemberUsecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases/castmember"
	"github.com/go-chi/chi/v5"
)

func (app *application) createCastMemberHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	castType, err := castmember.TypeFromString(input.Type)
	if err != nil {
		app.badRequestResponse(w, err)
	}

	ccc := castmemberUsecase.CreateCastMemberCommand{
		Name: input.Name,
		Type: castType,
	}

	noti, output := app.useCases.CastMember.Create.Execute(ccc)

	if output != nil {
		app.writeJson(w, http.StatusCreated, envelope{"id": output.ID}, nil)
		return
	}
	err = app.writeError(w, http.StatusBadRequest, "Could not save cast member", noti)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *application) listCastMemberHandler(w http.ResponseWriter, r *http.Request) {
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

	output, err := app.useCases.CastMember.FindAll.Execute(query)

	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	app.writeJson(w, http.StatusOK, output, nil)
}

func (app *application) updateCastMemberHandler(w http.ResponseWriter, r *http.Request) {
	castMemberIdStr := chi.URLParam(r, "id")
	castMemberId, err := strconv.ParseInt(castMemberIdStr, 10, 64)
	if err != nil {
		app.badRequestResponse(w, errors.New("invalid id"))
		return
	}
	if castMemberId <= 0 {
		app.notFoundResponse(w)
		return
	}
	var input struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	err = app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, errors.New("invalid id"))
		return
	}

	command := castmemberUsecase.UpdateCastMemberCommand{
		ID:   castMemberId,
		Name: input.Name,
		Type: input.Type,
	}

	noti := app.useCases.CastMember.Update.Execute(command)

	if noti == nil || !noti.HasErrors() {
		app.writeJson(w, http.StatusOK, envelope{"id": castMemberId}, nil)
		return
	}

	err = app.writeError(w, http.StatusBadRequest, "Could not update cast member", noti)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
