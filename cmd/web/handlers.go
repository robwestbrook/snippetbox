package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/robwestbrook/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check for exact path to "/"
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
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
	
	// Use SnippetModel's GET method to retrieve data
	// for a record by ID. Return 404 if not found.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	
	// Initialize a slice containing the path to the
	// template files. The base must be first.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}

	// Parse the template files
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create an instance of a templateData struct
	// to hold the snippet data
	data := &templateData{
		Snippet: snippet,
	}

	// Execute the template files, passing in the snippet
	// data as the final parameter. The data passed to the
	// template will contain a models.Snippet struct
	// inside a templateData struct.
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Check for only a POST request
	// and return error if not a POST request
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "Is this a poem?"
	content := "Roses are red\n violets are blue\n -Rob\n"
	expires := 14

	// Pass data to SnippetModel.Insert() method.
	// The ID is returned
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	
	// Redirect user to new snippet page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}