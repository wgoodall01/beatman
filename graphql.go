package main

type Resolver struct {
	// The song library used in queries
	Library *Library
}

func (r *Resolver) Songs(params struct{ Query *string }) (srs []SongResolver, err error) {

	// Search the songs, if enabled.
	var songs []*Song
	if params.Query != nil {
		songs = r.Library.Query(*params.Query)
	} else {
		songs = r.Library.Songs()
	}

	for _, song := range songs {
		srs = append(srs, SongResolver{song})
	}

	return srs, nil
}

const SchemaText = `
schema {
	query: Query
}

type Query {
	songs(query:String): [Song!]!
}

type Song {
	id: ID!
	name: String!
	subName: String
	authorName: String
	beatsPerMinute: Float
}

`
