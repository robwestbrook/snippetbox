package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check for exact path to "/"
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Initialize a slice containing the path to the
	// template files. The base must be first.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// Parse template files and check for errors
	// Template path is relative to the root of the
	// project directory. If there are errors, use
	// the application error log
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Execute the base template to write to response body
	// The last parameter represents any dynamic data 
	// to pass in
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Get id from URL parameters query string
	// and validate the id as an integer
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// Write response
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Check for only a POST request
	// and return error if not a POST request
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	// Write response
	w.Write([]byte("Create a new snippet..."))
}