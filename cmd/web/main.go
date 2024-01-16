// Package main
package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robwestbrook/snippetbox/internal/models"
)

// Define an application struct to hold all application
// wide dependencies for the application. The route
// handlers will become methods against this
// application struct.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

// Open DB function
// Wraps sql.Open() and returns a sql.DB connection pool
// for the DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// Main function - entry point to app.
// Main function responsibilities:
//  1. Parsing runtime configuration settings for app
//  2. Establishing the dependencies for the handlers
//  3. Running the HTTP server
func main() {
	 // Define command line flags
	// "addr"	: 	http PORT (default: 8000)
	// "dsn"	:		database DSN string (database name)
	// Then parse the command line flags.
	// Read the command line flags and assign to variable
	addr := flag.String("addr", ":8000", "HTTP network address")
	dsn := flag.String("dsn", "./snippetbox.db", "SQLite data source file name")
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

	// Create a database connection pool, using the
	// openDB() function. Pass openDB() the DSN from
	// the command line flag.
	// Defer a call to db.Close(), so the connection
	// pool is closed before the main() function exits
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new instance of the application struct,
	// containing the dependencies. This makes all
	// dependencies available to the handlers
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	// Initialize a new http.Server struct
	// Set the Addr and Handler fields so the server uses
	// the addr flag and app.routes, and set the
	// ErrorLog field so the server uses the custom
	// errorLog logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Start server using defined loggers
	infoLog.Printf("Starting server on port %s", *addr)
	errSrv := srv.ListenAndServe()
	errorLog.Fatal(errSrv)
}
