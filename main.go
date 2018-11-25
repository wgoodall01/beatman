package main

import (
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

func main() {

	library, err := NewLibrary("_beatsaber")
	if err != nil {
		log.lib.WithError(err).Fatal("Could not open library")
	}
	library.StartSync()

	schema := graphql.MustParseSchema(SchemaText, &Resolver{library})
	http.Handle("/graphql", &relay.Handler{Schema: schema})
	log.web.Info("Listening on :8080")
	log.web.Fatal(http.ListenAndServe(":8080", nil))
}
