package models

import (
	"database/sql"
	"time"
)

// User defines a User type.
type User struct {
	ID							int
	Name 						string
	Email 					string
	HashedPassword	[]byte
	Created 				time.Time
}

// UserModel is a type that wraps a database connection
// pool
type UserModel struct {
	DB		*sql.DB
}

/*
	Insert adds a new record to the users table
*/
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

/*
	Authenticate verifies a user exists using the 
	provided email and password. The user's ID will be
	returned.
*/
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

/*
	Exists is used to check if a user exists with
	a specific ID
*/
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}