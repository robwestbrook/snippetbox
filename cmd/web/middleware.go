package main

import (
	"fmt"
	"net/http"
)

/*
	secureHeaders function adds headers to increase
	app security.
*/
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts-gstatic.com")
		
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		// FLOW CONTROL: Any code here will execute on the 
		// way down the chain of control
		next.ServeHTTP(w, r)
		// FLOW CONTROL: Any code here will execute on the
		// way back up the chain of control
	})
}

/*
	logRequest function logs all requests as a method of
	the applicatio. Becuase it is a method against
	application, it has access to the handler dependencies
	including the information logger.
*/
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

/*
	recoverPanic function helps server gracefully recover
	from a panic situation.
*/
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function which will always be
		// run in the event of a panic as Go unwinds the stack
		defer func() {
			// Use the built in recover function to check if
			// there has been a panic.
			if err := recover(); err != nil {
				// Set a "Connection:Close" header on response.
				// This makes Go's HTTP server automatically
				// close the connection after a response is sent.
				// It alos informs the user the connection 
				// will be closed.
				w.Header().Set("Connection", "close")
				// Call the app serverError method to return a
				// 500 Internal Server response
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

/*
	requireAuthentication protects a page from an
	unauthorized user if it is only for athenticated users
*/
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect to
		// the login page. Return from the middleware
		// chain so no other middleware handlers are called
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// If the user is authenticated, set the
		// "Cache-Control: no-store" header so pages requiring
		// authentication are not stored in the browser cache
		w.Header().Add("Cache-Control", "no-store")

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}