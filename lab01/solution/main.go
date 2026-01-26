package main

import (
	"time"

	"github.com/sirupsen/logrus"
)

// main is the entry point of the program.
// Go programs start execution in the main function of the main package.
func main() {
	// Print a welcome message using logrus
	logrus.Info("Welcome to Go!")
	
	// Get the current time and print it in a readable format
	// time.Now() returns the current local time
	// Format() uses Go's reference time: Mon Jan 2 15:04:05 MST 2006
	// This reference time is used as a template for formatting
	currentTime := time.Now()
	logrus.Infof("Current timestamp: %s", currentTime.Format("2006-01-02 15:04:05 MST"))
	
	// Alternative: Print in RFC3339 format (ISO 8601)
	logrus.Infof("ISO 8601 format: %s", currentTime.Format(time.RFC3339))
}
