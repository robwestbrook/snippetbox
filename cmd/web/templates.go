package main

import (
	"html/template"
	"path/filepath"

	"github.com/robwestbrook/snippetbox/internal/models"
)

type templateData struct {
	Snippet		*models.Snippet 	// Hold a single snippet
	Snippets	[]*models.Snippet	// holds many snippets
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

		// Create a slice containing the file paths for the
		// base template, partials, and page
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}

		// Parse the files into a template set
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name
		// of the page as the key
		cache[name] = ts
	}
	// Return the map.
	return cache, nil
}