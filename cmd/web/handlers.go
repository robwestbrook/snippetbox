package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/robwestbrook/snippetbox/internal/models"
)

/*
	Home function is the handler for the home page
*/
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Get the latest snippets
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Call the newTemplateData()helper to get a
	// templateData struct containing the default
	// data and add the snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Render the page
	app.render(w, http.StatusOK, "home.tmpl", data)
}

/*
	snippetView function handles a single snippet page
	determined by snippet ID, ID is stored in the request
	context. Retrieve ID using the ParseFromContext(),
	which returns a slice of parameter names and values.
*/
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Get request paramters
	params := httprouter.ParamsFromContext(r.Context())

	// Get id from URL parameters query string
	// and validate the id as an integer
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
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

	// Call the newTemplateData()helper to get a
	// templateData struct containing the default
	// data and add the snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippet = snippet

	// Render the page
	app.render(w, http.StatusOK, "view.tmpl", data)
}

/**/
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display snippet creation form..."))
}

/*
	snippetCreate function handles creating a new
	snippet
*/
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

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
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}