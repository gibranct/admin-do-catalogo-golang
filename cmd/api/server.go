package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"
)

func (app *application) server() *http.Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}

	return srv
}
