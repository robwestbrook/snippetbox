package main

import "net/http"

func (app *application) routes() *http.ServeMux {
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

	return mux
}