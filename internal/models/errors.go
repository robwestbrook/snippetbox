package models

import "errors"

// ErrNoRecord generates a new record to use when
// a database query is made and no records are found
var ErrNoRecord = errors.New("models: no matching record found")