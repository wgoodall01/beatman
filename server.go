package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/sirupsen/logrus"
)

type API struct {
	library *Library
	addr    string
}

// ServeAPI serves the API. This blocks until the server is terminated.
func (api *API) Serve() {
	// Set up mux
	r := mux.NewRouter()
	r.Use(api.logging)

	// Serve graphql
	schema := graphql.MustParseSchema(SchemaText, &Resolver{api.library})
	r.Handle("/graphql", &relay.Handler{Schema: schema})

	// Serve
	r.HandleFunc("/api/cover/{id}", api.handleCoverImages)
	r.HandleFunc("/api/audio/{id}", api.handleAudio)

	// Start listening
	srv := http.Server{
		Addr:         api.addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}
	log.web.WithField("addr", api.addr).Info("listening")
	log.web.Fatal(srv.ListenAndServe())
}

func (api *API) logging(next http.Handler) http.Handler {
	// Request logging
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.web.WithFields(logrus.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
			"time":   time.Since(start),
		}).Info("req")
	})
}

func sendSong404(w http.ResponseWriter, id string) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "song %s not found", id)
}

func (api *API) handleCoverImages(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		sendSong404(w, "")
		return
	}
	song := api.library.Get(ParseID(id))
	if song == nil {
		sendSong404(w, id)
		return
	}

	http.ServeFile(w, r, song.CoverPath)
	return
}

func (api *API) handleAudio(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		sendSong404(w, "")
		return
	}
	song := api.library.Get(ParseID(id))
	if song == nil {
		sendSong404(w, id)
		return
	}

	http.ServeFile(w, r, song.AudioPath)
}
