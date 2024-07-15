package main

import (
	"encoding/json"
	"net/http"

	"github.com.br/gibranct/admin-do-catalogo/pkg/notification"
)

func (app *application) badRequestResponse(w http.ResponseWriter, err error) {
	app.errorResponse(w, http.StatusBadRequest, err.Error())
}

func (app *application) errorResponse(w http.ResponseWriter, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJson(w, status, env, nil)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) writeError(w http.ResponseWriter, status int, msg string, noti *notification.Notification) error {
	errors := []string{}
	for _, err := range noti.GetErrors() {
		errors = append(errors, err.Error())
	}

	data := envelope{"message": msg, "errors": errors}

	return app.writeJson(w, status, data, nil)
}

func (app *application) writeJson(w http.ResponseWriter, status int, data any, headers http.Header) error {

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) serverErrorResponse(w http.ResponseWriter, err error) {
	message := "the server encountered a problem and could not process your request"
	app.writeError(w, http.StatusInternalServerError, message, nil)
}
