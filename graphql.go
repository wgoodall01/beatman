package main

import (
	"path/filepath"

	"github.com/blevesearch/bleve"
)

type Resolver struct {
	// The root path of the Beat Saber install directory
	LibraryRootPath string
}

func (r *Resolver) Songs(params struct{ Query *string }) (srs []SongResolver, err error) {
	// Load the CustomSongs folder from disk
	library, err := LoadLibrary(filepath.Join(r.LibraryRootPath, "CustomSongs"))
	if err != nil {
		return nil, err
	}

	// Search the songs, if enabled.
	if params.Query != nil {
		// Create search index
		mapping := bleve.NewIndexMapping()
		index, err := bleve.NewMemOnly(mapping)
		if err != nil {
			return nil, err
		}

		// Add songs to index
		for _, song := range library.Songs {
			err := index.Index(song.ID, song)
			if err != nil {
				return nil, err
			}
		}

		// Search the index
		query := bleve.NewMatchQuery(*params.Query)
		search := bleve.NewSearchRequest(query)
		searchResult, err := index.Search(search)
		if err != nil {
			return nil, err
		}

		results := searchResult.Hits
		for _, result := range results {
			song := library.Songs[result.ID]
			srs = append(srs, SongResolver{song})
		}
	} else {
		for _, song := range library.Songs {
			srs = append(srs, SongResolver{song})
		}
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
	id: ID
	name: String!
	subName: String
	authorName: String
	beatsPerMinute: Float
}

`
