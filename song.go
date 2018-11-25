package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
)

type Song struct {
	// Song metadata
	ID         string // like "1183-814", from beatsaver, or "" if it can't be determined
	Name       string `json:"songName"`
	SubName    string `json:"songSubName"`
	AuthorName string `json:"authorName"`

	// Track info
	ImagePath        string  `json:"coverImagePath"`
	BeatsPerMinute   float64 `json:"beatsPerMinute"`
	PreviewStartTime int     `json:"previewStartTime"`
	PreviewDuration  int     `json:"previewDuration"`

	// Don't worry about the actual track data for now.
}

var beatsaverIDRegex = regexp.MustCompile("[0-9]+-[0-9]+")

func LoadSong(path string) (song *Song, err error) {
	song = &Song{}

	dir, err := openDir(path)
	if err != nil {
		return nil, err
	}

	// Get the directory contents
	filenames, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	// Look through files for info.json
	for _, filename := range filenames {
		if strings.ToLower(filename) == "info.json" {
			infoPath := filepath.Join(dir.Name(), filename)
			infoFile, err := os.Open(infoPath)
			if err != nil {
				return nil, err
			}

			// parse the info.json into the Song
			decoder := json.NewDecoder(infoFile)
			decoder.Decode(song)

			// set the song's ID based on the directory name
			// note: sometimes this will break, f.ex if the folders are named after the song.
			id := filepath.Base(dir.Name())
			if beatsaverIDRegex.MatchString(id) {
				song.ID = id
			} else {
				// TODO: internal IDs for non-beatsaver songs
				id := strings.Replace(song.Name, " ", "-", -1)
				id = strings.ToLower(id)
				song.ID = id
			}
		}
	}

	return song, nil
}

func (s *Song) String() string {
	return fmt.Sprintf("%s (%s) [by %s, %.1fbpm]", s.Name, s.SubName, s.AuthorName, s.BeatsPerMinute)
}

type SongResolver struct {
	song *Song
}

func (sr SongResolver) Id() *graphql.ID {
	if sr.song.ID == "" {
		return nil
	} else {
		id := graphql.ID(sr.song.ID)
		return &id
	}
}

func (sr SongResolver) Name() string {
	return sr.song.Name
}

func (sr SongResolver) SubName() *string {
	return optStr(sr.song.SubName)
}

func (sr SongResolver) AuthorName() *string {
	return optStr(sr.song.AuthorName)
}

func (sr SongResolver) BeatsPerMinute() *float64 {
	return optFloat64(sr.song.BeatsPerMinute)
}
