package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	// Get the time right now for database record
	// created field
	now := time.Now()
	
	// Create a bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password), 12,
	)
	if err != nil {
		return err
	}

	// Create the SQL statement to insert user into db
	stmt := `
		INSERT INTO users (name, email, hashed_password, created)
		VALUES(?, ?, ?, ?)
	`

	// Use the Exec() method to insert user data and
	// hashed password into users table
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword), now.Format(dbTimeFormat))

	// If there is an error, process the error.
	// If the error string contains "UNIQUE" and "users.email"
	// it is an SQLite error for duplicate emails.
	// Return ErrDuplicateEmail. Else return err.
	if err != nil {
		errString := err.Error()
		if strings.Contains(errString, "UNIQUE") && 
				strings.Contains(errString, "users.email") {
					return ErrDuplicateEmail
				}
		return err
	}
	return nil
}

/*
	Authenticate verifies a user exists using the 
	provided email and password. The user's ID will be
	returned.
*/
func (m *UserModel) Authenticate(email, password string) (int, error) {
	// Get the ID and hashed password associated with the
	// given email. If no email matches, return 
	// ErrInvalidCredentials error
	var id int
	var hashedPassword []byte

	// Create SQL statement to retrieve user
	stmt := `
		SELECT id, hashed_password FROM users
		WHERE email = ?
		`
	
	// Query the database and scan in id and hashed password
	// to the variables declared above
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check hashed password against password user entered.
	// If they don't match, return ErrInvalidCredentials
	// error
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// If no errors, password is correct. Return user ID
	return id, nil
}

/*
	Exists is used to check if a user exists with
	a specific ID
*/
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}