// Package main
package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	// "github.com/alexedwards/scs"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robwestbrook/snippetbox/internal/models"
)

// Define an application struct to hold all application
// wide dependencies for the application. The route
// handlers will become methods against this
// application struct.
type application struct {
	errorLog 				*log.Logger
	infoLog  				*log.Logger
	snippets 				*models.SnippetModel
	templateCache		map[string]*template.Template
	formDecoder			*form.Decoder
	sessionManager	*scs.SessionManager
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

	// Initialize a new template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a form decoder instance
	formDecoder := form.NewDecoder()

	// Initialize a new session manager with scs.New().
	// The scs.New() function returns a pointer to a struct
	// which holds configuration settings for the sessions.
	// Configure to use SQLite as session store, setting
	// a lifetime of 12 hours for session. Set "Secure"
	// to ensure a cookie will only be sent using an
	// HTTPS connection.
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// Initialize a new instance of the application struct,
	// containing the dependencies. This makes all
	// dependencies available to the handlers
	//
	// Dependencies available:
	//	1. errorLog - error logger
	//	2. infoLog - information logger
	//	3. snippets - snippet model and methods
	//	4. templateCache - template in-memory cache
	// 	5. formDecoder - decodes all form input
	//	6. sessionManager - manages all user sessions
	app := &application{
		errorLog: 			errorLog,
		infoLog:  			infoLog,
		snippets: 			&models.SnippetModel{DB: db},
		templateCache: 	templateCache,
		formDecoder: 		formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize a tls.Config struct to hold non-default
	// TLS settings for the server. Here, change the curve
	// preference value, so that only elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Initialize a new http.Server struct using the
	// following parameters:
	//	1.	Addr: the TCP address the server listens on
	//	2.	ErrorLog: logger to use for errors
	//	3. 	Handler: the handler for routing
	//	4.	TLSConfig: provides optional TLS configuration
	//	5.	IdleTimeout: Max time to wait for next request when keep-alive is enabled
	//	6.	ReadTimeout: Max time for reading entire request
	//	7.	WriteTimeout: max time before timing out writes of the response
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		TLSConfig: tlsConfig,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server using defined loggers. Using the
	// ListenAndServeTLS() to start an HTTPS server,
	// passing in the paths to the TLS certificate and
	// corresponding private key as the two parameters.
	infoLog.Printf("Starting server on port %s", *addr)
	errSrv := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(errSrv)
}
