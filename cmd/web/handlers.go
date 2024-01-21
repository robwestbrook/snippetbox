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

	// Create a session value for a flash message to user
	app.sessionManager.Put(
		r.Context(),
		"flash",
		"Snippet successfully created!",
	)
	
	// Redirect user to new snippet page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// userSignupForm struct takes in the form values from
// the user signup form.
type userSignupForm struct {
	Name						string	`form:"name"`
	Email						string	`form:"email"`
	Password 				string	`form:"password"`
	validator.Validator			`form:"-"`
}

/*
	userSignup function displays the form for a
	new user signup.
*/
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

/*
	userSignupPost processes the signup form and signs
	up a new user.
*/
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Create an instance of userSignupForm struct
	var form userSignupForm

	// Parse the form data into the struct
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the name field is not blank
	form.CheckField(
		validator.NotBlank(form.Name),
		"name",
		"This field cannot be blank",
	)

	form.CheckField(
		validator.NotBlank(form.Email),
		"email",
		"This field cannot be blank",
	)

	form.CheckField(
		validator.Matches(form.Email, validator.EmailRX),
		"email",
		"This field must be a valid email address",
	)

	form.CheckField(
		validator.NotBlank(form.Password),
		"password",
		"This field cannot be blank",
	)

	form.CheckField(
		validator.MinChars(form.Password, 8),
		"password",
		"This field must be at least 8 characters long",
	)

	// Check for form errors. If any redisplay the form
	// with a 422 status code.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	// Create a new user in the database and
	// check for errors
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// If no errors, add confirmation flash message
	// to session confirming signup.
	app.sessionManager.Put(
		r.Context(),
		"flash",
		"Signup was successful. Please log in",
	)

	// Redirect to login page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

// Create a userLoginForm struct
type userLoginForm struct {
	Email						string 	`form:"email"`
	Password				string 	`form:"password"`
	validator.Validator 		`form:"-"`
}

/*
	userLogin displays the form for a user to log in.
*/
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}

/*
	userLoginPost processes the login form and logs in
	a user.
*/
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	// Create a form variable of type userLoginForm and
	// decode the form data into variable
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate input
	form.CheckField(
		validator.NotBlank(form.Email),
		"email",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.Matches(form.Email, validator.EmailRX),
		"email",
		"This field must be a valid email address",
	)
	form.CheckField(
		validator.NotBlank(form.Password),
		"password",
		"This field cannot be blank",
	)

	// Check for any validation errors
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	// Check if the credentials are valid. If not, add a
	// generic non-field error message and re-diplay
	// the login page.
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use RenewToken() method on current session to 
	// change the session ID. RenewToken() changes the
	// ID of the current user's session but retain any
	// data associated with the session.
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Add the ID of the user to the session, so they
	// are now logged in.
	app.sessionManager.Put(r.Context(), "authenticatedID", id)

	// Redirect the user to the create snippet page
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

/*
	userLogoutPost logs a user out.
*/
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Use RenewToken() method on current session to
	// change the session ID
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Remove the authenticatedUserID from the session
	// data so user is logged out
	app.sessionManager.Remove(r.Context(), "authenticatedID")

	// Add flash message to session to confirm to user
	// they are logged out
	app.sessionManager.Put(
		r.Context(),
		"flash",
		"You have been logged out successfully",
	)

	// Redirect user to application home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}