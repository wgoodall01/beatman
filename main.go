package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

func main() {
	// Set up beatsaber library
	library, err := NewLibrary("_beatsaber")
	if err != nil {
		log.lib.WithError(err).Fatal("Could not open library")
	}
	library.StartSync()

	// Set up mux
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		// Request logging
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.web.WithFields(logrus.Fields{
				"method": r.Method,
				"uri":    r.RequestURI,
				"time":   time.Since(start),
			}).Info()
		})
	})

	// Serve graphql
	schema := graphql.MustParseSchema(SchemaText, &Resolver{library})
	r.Handle("/graphql", &relay.Handler{Schema: schema})

	// Start listening
	srv := http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}
	log.web.Info("Listening on :8080")
	log.web.Fatal(srv.ListenAndServe())
}
