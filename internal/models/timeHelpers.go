package models

import "time"

// dbTimeFormat defines the format used to convert
// date and time to a SQLite-friendly datetime.
const dbTimeFormat = "2006-01-02 15:04:05"

// Create string variables to hold the record's
// created and expires fields for conversion to
// Go's time.Time format
var createdTime string
var expiredTime string

/*
stringToTime function takes in a string defining the
time format and a time string from SQLite. It returns
a GO time.Time format.
*/
func stringToTime(stringToConvert string) (time.Time) {
	res, _ := time.Parse(dbTimeFormat, stringToConvert)
	return res
}