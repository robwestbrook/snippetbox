// Package main
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold all application
// wide dependencies for the application. The route
// handlers will become methods against this
// application struct.
type application struct {
	errorLog 	*log.Logger
	infoLog		*log.Logger
}

// Main function - entry point to app.
// Main function responsibilities:
//	1. Parsing runtime configuration settings for app
//	2. Establishing the dependencies for the handlers
// 	3. Running the HTTP server
func main() {

	// Define command line flags
	// "addr": http PORT (default: 8000)
	// Then parse the command line flags.
	// Read the command line flags and assign to variable
	addr := flag.String("addr", ":8000", "HTTP network address")
	flag.Parse()

	// Create a logger for writing information  and 
	// error messages.
	// Takes 3 paramters:
	//	1. Destination to write logs to
	//	2. A prefix for message
	// 	3. Flags for additional information - flags are
	//			joined using the bitwise OR (|)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of the application struct,
	// containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	// Initialize a new http.Server struct
	// Set the Addr and Handler fields so the server uses
	// the addr flag and app.routes, and set the
	// ErrorLog field so the server uses the custom
	// errorLog logger
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	// Start server using defined loggers
	infoLog.Printf("Starting server on port %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}