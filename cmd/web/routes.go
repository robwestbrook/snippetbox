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

	GET		| /user/signup			| userSignup				| Display form
				|										|										| for signing up
				|										|										| new user

	POST	| /user/signup			| userSignupPost		| Create a
				|										|										| new user

	GET		| /user/login				| userLogin					| display form
				|										|										| for logging in
				|										|										| a user

	POST	| /user/login				|										| Authenticate
				|										|										| and login
				|										|										| a user

	POST	| /user/logout			| userLogoutPost		| Logout a
				|										|										| user

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

	// Create a new middleware chain containing middleware
	// specific to dynamic application routes. Alice
	// manages middleware chains.
	// Includes:
	//	1. LoadAndSave session middleware
	//	2. noSurf CSRF preventing middleware
	dynamic := alice.New(
		app.sessionManager.LoadAndSave, 
		noSurf,
		app.authenticate,
	)

	// UNPROTECTED ROUTES - Open to all app users

	// Create routes with methods, patterns, 
	// handlers. Wrap the unprotextedhandlers with the 
	// DYNAMIC middleware for session control.
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// PROTECTED ROUTES- Only available to authenticated user

	// Create a chain where the DYNAMIC middleware is
	// appended with the REQUIREAUTHENTICATION middleware
	protected := dynamic.Append(app.requireAuthentication)

	// Create routes with methods, patterns, 
	// handlers. Wrap the unprotextedhandlers with the 
	// PROTECTED middleware for authenticated session control.
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
	
	// Create a middleware chain containing the "standard"
	// middleware which will be sent for every request
	// the application receives. Alice manages middleware 
	// chains.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware followed
	// by the servermux.
	return standard.Then(router)
}