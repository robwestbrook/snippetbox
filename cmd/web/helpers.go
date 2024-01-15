package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// ServerError helper.
// Writes an error message and stack trace to the
// errorLog, then sends a generic 500 Internal
// Server Error to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// ClientError helper.
// Sends a specific status code and corresponding
// description to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// NotFound helper.
// A convenience wrapper around clientError which
// sends a 404 Not Found response to user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}