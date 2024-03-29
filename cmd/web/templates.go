package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/robwestbrook/snippetbox/internal/models"
)

// templateData struct is used to pass information
// to the templates. This data includes:
//	1. CurrentYear - the current year
//	2. Snippet - holds the information of a single snippet
//	3. Snippets - a slice of many snippets
//	4. Form - hold any submitted form data
//	5. Flash - holds any flash message generated
// 	6. IsAuthenticated - holds true or false for authenticated users
//	7. CSRFToken - Adds a CSRFToken
type templateData struct {
	CurrentYear			int
	Snippet					*models.Snippet
	Snippets				[]*models.Snippet
	Form						any
	Flash						string
	IsAuthenticated	bool
	CSRFToken				string
}

/*
	humanDate function returns a formatted string of
	a time.Time object
*/
func humanDate(t time.Time) string {
	// Return an empty string if time has zero value
	if t.IsZero() {
		return ""
	}
	// Convert the time to UTC before formatting
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

/*
	Initialize a template.FuncMap object and store in
	a global variable. This is a string-keyed map which
	acts as a lookup between the names of custom template
	functions and the functions themselves. These functions
	can accept any number of parameters but the MUST
	return only one value, except when returning an error.
*/
var functions = template.FuncMap{
	"humanDate": humanDate,
}

/*
	newTemplateCache function creates a map of all app
	pages, partials, and templates. The map will function
	as an in-memory cache.
*/
func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as a cache
	cache := map[string]*template.Template{}

	// Use the Glob() function to get a slice of all
	// file paths that match the pattern 
	// "./ui/html/pages/*.tmpl". This gives a slice
	// of all file paths for page templates.
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Loop through page file paths
	for _, page := range pages {

		// Extract the file name and assign it to name
		name := filepath.Base(page)

		// Register the template.FuncMap with the template
		// set before the ParseFiles() is used. An empty
		// template set has to be created, registering the
		// FuncMap using the Funcs() method.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Use ParseGlob() on this template set to add
		// any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Parse the files in this template set to add
		// the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache map, using the
		// name of the page as the key
		cache[name] = ts
	}
	// Return the map.
	return cache, nil
}