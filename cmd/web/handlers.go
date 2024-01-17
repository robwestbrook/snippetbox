package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/robwestbrook/snippetbox/internal/models"
)

/*
	Home function is the handler for the home page
*/
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check for exact path to "/"
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Get the latest snippets
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Render the page
	app.render(w, http.StatusOK, "home.tmpl", &templateData{
		Snippets: snippets,
	})
}

/*
	snippetView function handles a single snippet page
	determined by snippet ID
*/
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

	// Render the page
	app.render(w, http.StatusOK, "view.tmpl", &templateData{
		Snippet: snippet,
	})
}

/*
	snippetCreate function handles creating a new
	snippet
*/
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