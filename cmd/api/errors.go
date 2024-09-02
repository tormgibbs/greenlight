package main

import (
	"fmt"
	"net/http"
)

// logError() logs an error message using the application's logger.
// This method is typically called when an internal error occurs that
// needs to be recorded in the server logs.
func (app *application) logError(_ *http.Request, err error) {
	app.logger.Println(err)
}

// errorResponse() sends a JSON-formatted error message and a specified HTTP status code
// to the client. If there's an issue writing the JSON response, it logs the error
// and sends a 500 Internal Server Error status.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	// Write the JSON response. Logs error and return empty response with status code 500 if any
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// serverErrorResponse() logs the error and sends a 500 Internal Server Error
// response to the client. This method is used for unexpected server-side issues.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse() sends a 404 Not Found response to the client when a resource
// cannot be located. This method is used when a requested URL does not match
// any route in the application.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse() sends a 405 Method Not Allowed response to the client
// when the HTTP method used (e.g., POST, GET) is not supported for the requested resource.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// badRequestResponse() sends a 400 Bad Request Response to the client with provided error message in
// response body
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
