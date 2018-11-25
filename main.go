package main

import (
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

func main() {
	log.Println("Starting server")
	schema := graphql.MustParseSchema(SchemaText, &Resolver{"_beatsaber"})
	http.Handle("/graphql", &relay.Handler{Schema: schema})
	log.Println("Loaded graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
