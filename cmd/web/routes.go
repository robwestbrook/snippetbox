package main

import (
	"net/http"

	"github.com/justinas/alice"
)

/*
	routes function processes requests for the
	app's files and pages.
*/
func (app *application) routes() http.Handler {
	// Create a new server mux
	mux := http.NewServeMux()

	// Create a static file server and use the mux.Handle()
	// function to register the file server as the
	// handler for all URL paths that start with "/static/"
	// Strip the "/static" prefix for matching paths
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))


	// Register routes with application methods
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Create a middleware chain containing the "standard"
	// middleware which will be sent for every request
	// the application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware followed
	// by the servermux.
	return standard.Then(mux)
}