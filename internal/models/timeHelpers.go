package models

import "time"

/*
stringToTime function takes in a string defining the
time format and a time string from SQLite. It returns
a GO time.Time format.
*/
func stringToTime(timeLayout string, stringToConvert string) (time.Time) {
	res, _ := time.Parse(timeLayout, stringToConvert)
	return res
}