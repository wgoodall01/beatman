package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
)

type SongID struct {
	Name        string // From directory title
	ID, Version int    // From beatsaver.com key, format
}

var beatsaverIDRegex = regexp.MustCompile(`[0-9]+-[0-9]+`)

// ParseID takes text and turns it into a usable song ID
func ParseID(text string) (id SongID) {
	getInt := func(s string) int {
		i64, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			log.lib.Fatal("couldn't parse valid ID--what? Please open an issue.")
		}
		return int(i64)
	}

	if beatsaverIDRegex.MatchString(text) {
		// Parse the bsaver id
		split := strings.Split(text, "-")
		id.ID = getInt(split[0])
		id.Version = getInt(split[1])
	} else {
		// Use the dir name
		id.Name = text
	}

	return id
}

func (sid *SongID) String() string {
	if sid.Name != "" {
		return sid.Name
	} else {
		return fmt.Sprintf("%d-%d", sid.ID, sid.Version)
	}
}

type Song struct {
	// Song metadata
	ID         SongID
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

// LoadSong attempts to load the given path as a song.
// If it succeeds, it returns the song and err=nil
// If it doesn't know how to load that path, it returns song=nil, err=nil
// If it fails somehow, it returns song=nil, err=why
func LoadSong(path string) (song *Song, err error) {
	song = &Song{}

	// Open the song file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Get file info--dir? zip?
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		// Load from zip file
		return nil, nil // can't do this yet
	} else {
		// Load from directory

		// Look through files for info.json
		err = filepath.Walk(path, func(infoPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			basename := filepath.Base(infoPath)

			if strings.ToLower(basename) == "info.json" {
				infoFile, err := os.Open(infoPath)
				if err != nil {
					return err
				}

				// parse the info.json into the Song
				decoder := json.NewDecoder(infoFile)
				decoder.Decode(song)

				// set the song's ID based on the directory name
				// note: sometimes this will break, f.ex if the folders are named after the song.
				song.ID = ParseID(filepath.Base(file.Name()))
			}

			return nil
		})
		if err != nil {
			return nil, err
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

func (sr SongResolver) Id() graphql.ID {
	return graphql.ID(sr.song.ID.String())
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
