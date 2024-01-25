package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// Create a slice of anonymous structs containing
	// the test case name, input to the humanDate()
	// function (the tm field), and expected output,
	// (the want field).

	tests := []struct {
		name 	string
		tm 		time.Time
		want	string
	} {
		{
			name: "UTC",
			tm: time.Date(2024, 1, 25, 17, 30, 0, 0, time.UTC),
			want: "25 Jan 2024 at 17:30",
		},
		{
			name: "Empty",
			tm: time.Time{},
			want: "",
		},
		{
			name: "CST",
			tm: time.Date(2024, 1, 25, 17, 30, 0, 0, time.FixedZone("CST", -5*60*60)),
			want: "25 Jan 2024 at 22:30",
		},
	}

	// Loop over tesr cases
	for _, tt := range tests {
		// Use the t.Run() function to run a sub-test for
		// each test case.
		// Parameters:
		//	1. 	Name of test
		//	2. 	The anonymous function containing the
		//			actual test for each case.
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			if hd != tt.want {
				t.Errorf("got %q, want %q", hd, tt.want)
			}
		})
	}
}