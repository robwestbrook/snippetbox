package models

import (
	"database/sql"
	"errors"
	"time"
)

// Snippet defines a type to hold data for an
// individual snippet. The fields of the struct
// correspond to the fields in SQLite snippets
// table.
type Snippet struct {
	ID				int
	Title			string
	Content		string
	Created		time.Time
	Expires		time.Time
}

// SnippetModel defines a type to wrap an
// sql.DB connection pool.
type SnippetModel struct {
	DB	*sql.DB
}

/*
Insert function inserts a new snippet into
the database.
*/
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


/*
Get function returns a specific snippet
based on its id.
*/
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	
	// Get the time right now for database record
	// created field
	now := time.Now()

	// SQL statement to get snippet
	stmt:= `SELECT id, title, content, created, expires 
					FROM snippets WHERE expires > ? AND id = ?`
	
	// Use the QueryRow() method to get row
	row := m.DB.QueryRow(stmt, now.Format(dbTimeFormat), id)

	// Initialize a pointer to a new Snippet struct
	s := &Snippet{}

	// Scan the values of the row to the corresponding
	// fields in the Snippet struct. The arguments are
	// *pointers" to the copied data location. The
	// variables "createdTime" and "expiredTime" are
	// created at the top, to avoid duplicity.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &createdTime, &expiredTime)
	
	// Convert the record's time strings to Go's
	// time.Time format and add to snippet struct
	s.Created = stringToTime(createdTime)
	s.Expires = stringToTime(expiredTime)

	
	// If the query returns no rows, row.Scan() returns
	// a sql.ErrNoRows error. Check for error with the
	// errors.Is function
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	// Return the Snippet
	return s, nil
}

/*
Latest function gets the latest 10 unexpired
snippets.
*/
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	
	// Get the time right now for database record
	// created field
	now := time.Now()

	// SQL statement to execute
	stmt := `SELECT id, title, content, created, expires
					FROM snippets WHERE expires > ?
					ORDER BY id DESC LIMIT 10`

	// Use the Query() method, which returns a sql.Rows
	// result
	rows, err := m.DB.Query(stmt, now.Format(dbTimeFormat))
	if err != nil {
		return nil, err
	}

	// Defer rows.Close() after checking for errors
	defer rows.Close()

	// Initialize an empty slice to hold the snippet structs
	snippets := []*Snippet{}

	// Use rows.Next() to iterate through the results,
	// which prepares each row for the rows.Scan() method.
	for rows.Next() {

		// Create a pointer to a new zeroed Snippet struct
		s := &Snippet{}

		// Scan the values of the row to the corresponding
		// fields in the Snippet struct. The arguments are
		// *pointers" to the copied data location. The
		// variables "createdTime" and "expiredTime" are
		// created at the top, to avoid duplicity.
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &createdTime, &expiredTime)
		
		// Convert the record's time strings to Go's
		// time.Time format and add to snippet struct
		s.Created = stringToTime(createdTime)
		s.Expires = stringToTime(expiredTime)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// After loop, call rows.Err() to retrieve any error
	// encountered during iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return Snippets slice
	return snippets, nil
}