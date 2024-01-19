package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

/*
	routes function processes requests for the
	app's files and pages.

	Routes Table
	Meth 	|		Pattern					|	Handler						| 	Actiom
	-----------------------------------------------------
	GET		|	/									| home							| home page

	GET		|	/snippet/view/:id	| snippetView				| display
				|										|										| specific
				|										|										| snippet

	GET		|	/snippet/create		| snippetCreate			| Display form
				|										|										| to create
				|										|										| new snippet

	POST	|	/snippet/create		|	snippetCreatePost	| Create
				|										|										| new
				|										|										| snippet

	GET		| /static/*filepath	|	http.Fileserver		| serve
				|										|										|	static
				|										|										| file
*/
func (app *application) routes() http.Handler {
	// Initialize a router
	router := httprouter.New()

	// Create a handler function wrapping the notFound()
	// helper function. Assign it as the custom handler
	// for 404 not found responses.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Create a static file server and use the mux.Handle()
	// function to register the file server as the
	// handler for all URL paths that start with "/static/"
	// Strip the "/static" prefix for matching paths
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(
		http.MethodGet,
		"/static/*filepath",
		http.StripPrefix("/static/",
		fileServer),
	)


	// Create Routes with methods, patterns, handlers
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// Create a middleware chain containing the "standard"
	// middleware which will be sent for every request
	// the application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware followed
	// by the servermux.
	return standard.Then(router)
}