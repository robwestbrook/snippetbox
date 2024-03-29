package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
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

// new template data function
// Returns a pointer to a templateData struct.
// The struct is initialized with the following
// fields:
//	1. The current year - adds the current year to a template
//	2. Flash message - adds a flash message to a template
//	3. IsAuthenticated - checks for authenticated user
//	4. CSRFToken - Adds a CSRF token
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: 			time.Now().Year(),
		Flash: 						app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: 	app.isAuthenticated(r),
		CSRFToken: 				nosurf.Token(r),
	}
}

// render function
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Retrieve the template set from the cache based on
	// page name. If no entry exists in the cache, create
	// a new error and call serverError()
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Initialize a new buffer for the templates
	buf := new(bytes.Buffer)

	// Execute the template and write to buffer. Any
	// error calls the serverError() helper function.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template writes to the buffer without
	// errors, it is safe. Write HTTP status code 
	// to response writer.
	w.WriteHeader(status)

	// Write contents of buffer to the response writer,
	// by passing the http.ResponseWriter to a function
	// that takes an io.Writer.
	buf.WriteTo(w)
}

// decodePostForm function. The second parameter, dst,
// is the target destination to decode the form data into
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call Decode() on the decoder instance, passing the
	// target destination as the first parameter
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// With an invalid target destination, Decode()
		// will return an error with the type
		// *form.InvalidDecodeError. Use errors.As()
		// to check for this error and raise a panic rather
		// than returning the error
		var invalidDecodeError *form.InvalidDecoderError

		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}

		// For all errors, they are returned as normal
		return err
	}
	return nil
}

// isAuthenticated function returns true if the
// current request is from an authenticated user,
// otherwise returns false.
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)

	if !ok {
		return false
	}

	return isAuthenticated
}