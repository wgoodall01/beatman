package main

import (
	"time"
)

func main() {
	// Set up beatsaber library
	library, err := NewLibrary("_beatsaber")
	if err != nil {
		log.lib.WithError(err).Fatal("Could not open library")
	}
	library.StartSync(time.Second * 30)

	// Start serving the API
	api := &API{library, ":8080"}
	api.Serve()
}
