package models

import (
	"database/sql"
	"time"
)

// Snippet defines a type to hold data for an
// individual snippet. The fields of the struct
// correspond to the fields in SQLite snippets
// table.
type Snippet struct {
	ID				int
	Title			string
	Contect		string
	Created		time.Time
	Expires		time.Time
}

// SnippetModel defines a type to wrap an
// sql.DB connection pool.
type SnippetModel struct {
	DB	*sql.DB
}

const dbTimeFormat = "2006-01-02 15:04:05"

// Insert function inserts a new snippet into
// the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	
	// Get the time right now for database record
	// created field
	now := time.Now()

	// Create an expires date by adding the expires time
	// to the current time
	exp := now.AddDate(0, 0, expires)
	
	// SQL statement to execute.
	stmt := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES(?, ?, ?, ?)
	`
	// Execute the SQL statement
	result, err := m.DB.Exec(stmt, title, content, now.Format(dbTimeFormat), exp.Format(dbTimeFormat))
	if err != nil {
		return 0, err
	}

	// Use lastInsertId on the result to get the
	// ID of new record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Convert ID (int64) to int and return
	return int(id), nil
}

// Get function returns a specific snippet
// based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}