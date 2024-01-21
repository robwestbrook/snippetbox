package models

import "errors"

// ErrNoRecord generates a new error to use when
// a database query is made and no records are found
var ErrNoRecord = errors.New("models: no matching record found")

// ErrInvalidCredentials generates a new error when 
// a user tries to login with an incorrect email 
// or password
var ErrInvalidCredentials = errors.New("models: invalid credentials")

// ErrDuplicateEmail generates a new error when a user
// tries to signup with an email address that is
// already in use
var ErrDuplicateEmail = errors.New("models: duplicate email")

