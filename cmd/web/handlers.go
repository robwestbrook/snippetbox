package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/robwestbrook/snippetbox/internal/models"
	"github.com/robwestbrook/snippetbox/internal/validator"
)

/*
	Home function is the handler for the home page
*/
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Get the latest snippets
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Call the newTemplateData()helper to get a
	// templateData struct containing the default
	// data and add the snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Render the page
	app.render(w, http.StatusOK, "home.tmpl", data)
}

/*
	snippetView function handles a single snippet page
	determined by snippet ID, ID is stored in the request
	context. Retrieve ID using the ParseFromContext(),
	which returns a slice of parameter names and values.
*/
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Get request paramters
	params := httprouter.ParamsFromContext(r.Context())

	// Get id from URL parameters query string
	// and validate the id as an integer
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Use SnippetModel's GET method to retrieve data
	// for a record by ID. Return 404 if not found.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Call the newTemplateData()helper to get a
	// templateData struct containing the default
	// data and add the snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippet = snippet

	// Render the page
	app.render(w, http.StatusOK, "view.tmpl", data)
}

/*
	snippetCreate function responds to a GET function,
	processes the form template, and present it to
	the user.
*/
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Create a new template set
	data := app.newTemplateData(r)

	// Initialize a new createSnippetForm instance and
	// pass it to the template. This can be used to set
	// any 'initial' values for the form.
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	// Render the template
	app.render(w, http.StatusOK, "create.tmpl", data)
}

/*
	Define a snippetCreateForm to represent the form data
	and inherit all the fields and methods of the
	Validator type. All fields are exported, so they 
	are capitalized. 

	The struct includes struct tags which which tell the
	form decoder how to map HTML form values into struct
	fields.
*/
type snippetCreateForm struct {
	Title								string	`form:"title"`
	Content 						string	`form:"content"`
	Expires 						int			`form:"expires"`
	validator.Validator					`form:"-"`
}

/*
	snippetCreate function handles creating a new
	snippet
*/
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// Declare a new instance of the snippetCreateForm
	// struct
	var form snippetCreateForm

	// Decode the form data using the decodePostForm
	// helper function
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// BEGIN VALIDATION

	// Check for blank title
	 form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")

	// Check title for max length
	form.CheckField(
		validator.MaxChars(form.Title, 100),
		"title",
		"This field cannot be more than 100 characters long")
	
	// Check for blank content
	form.CheckField(
		validator.NotBlank(form.Content),
		"content",
		"This field cannot be blank")

	form.CheckField(
		validator.PermittedInt(form.Expires, 1, 7, 365),
		"expires",
		"This field must equal 1, 7, or 365")

	// Use Valid() method to check for any validation
	// fails. If so, re-render the template, passing in
	// the form
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}
 
	// ADD TO DATABASE

	// Pass data to SnippetModel.Insert() method.
	// The ID is returned
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	
	// Redirect user to new snippet page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}