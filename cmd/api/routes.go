package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	//	Initialize a new httprouter router instance
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//	Register the relevant methods and URL patterns with their respective handlers
	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/movies/:id", app.showMovieHandler)

	return router
}
