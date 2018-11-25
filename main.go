package main

import (
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.Formatter = &logrus.TextFormatter{}
	log.Level = logrus.InfoLevel
}

func main() {
	log.Info(" --- Start Beatman ---")

	library, err := NewLibrary("_beatsaber")
	if err != nil {
		log.WithError(err).Fatal("Could not open library")
	}
	library.StartSync()

	schema := graphql.MustParseSchema(SchemaText, &Resolver{library})
	http.Handle("/graphql", &relay.Handler{Schema: schema})
	log.Info("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
