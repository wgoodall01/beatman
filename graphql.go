package main

import graphql "github.com/graph-gophers/graphql-go"

type Resolver struct {
	// The song library used in queries
	Library *Library
}

func (r *Resolver) Songs(params struct{ Query *string }) (srs []*SongResolver, err error) {

	// Search the songs, if enabled.
	var songs []*Song
	if params.Query != nil {
		songs = r.Library.Query(*params.Query)
	} else {
		songs = r.Library.Songs()
	}

	for _, song := range songs {
		srs = append(srs, &SongResolver{song})
	}

	return srs, nil
}

func (r *Resolver) Song(params struct{ ID graphql.ID }) (sr *SongResolver, err error) {
	id := ParseID(string(params.ID))
	song := r.Library.Get(id)
	if song == nil {
		return nil, nil
	}

	return &SongResolver{song}, nil
}

const SchemaText = `
schema {
	query: Query
}

type Query {
	songs(query:String): [Song!]!
	song(id: ID!): Song
}

type Song {
	id: ID!
	name: String!
	subName: String
	authorName: String
	beatsPerMinute: Float

	# URIs to gett the track cover and audio
	coverUri: String,
	audioUri: String,
}

`
