package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
)

type Song struct {
	// Song metadata
	ID         SongID
	Name       string `json:"songName"`
	SubName    string `json:"songSubName"`
	AuthorName string `json:"authorName"`

	// Track info
	EnvironmentName  string  `json:"environmentName"`
	BeatsPerMinute   float64 `json:"beatsPerMinute"`
	PreviewStartTime float64 `json:"previewStartTime"`
	PreviewDuration  float64 `json:"previewDuration"`
	Shuffle          float64 `json:"shuffle"`
	ShufflePeriod    float64 `json:"shufflePeriod"`
	OneSaber         bool    `json:"oneSaber"`

	// Attached information
	// note: make sure this is absolute, not relative to info.json
	CoverPath   string  `json:"coverImagePath"`
	AudioPath   string  `json:"audioPath"`
	AudioOffset float64 `json:"songTimeOffset"`

	DifficultyLevels []struct {
		Difficulty     string `json:"difficulty"`
		DifficultyRank int    `json:"difficultyRank"`
		JsonPath       string `json:"jsonPath"`

		// copy these properties to root. They're deprecated in SongLoader.
		AudioPath_ string  `json:"audioPath"`
		Offset_    float64 `json:"offset"`
	} `json:"difficultyLevels"`

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

				// Copy over deprecated props, top songs overriding others.
				for i := len(song.DifficultyLevels) - 1; i >= 0; i-- {
					dl := &song.DifficultyLevels[i]
					if dl.AudioPath_ != "" {
						song.AudioPath = dl.AudioPath_
					}
					if dl.Offset_ != 0 {
						song.AudioOffset = dl.Offset_
					}
				}

				// Make paths relative to the info.json file location
				dir := filepath.Dir(infoPath)
				song.CoverPath = filepath.Join(dir, song.CoverPath)
				song.AudioPath = filepath.Join(dir, song.AudioPath)

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

type SongResolver struct {
	song *Song
}

func (sr *SongResolver) ID() graphql.ID {
	return graphql.ID(sr.song.ID.String())
}

func (sr *SongResolver) Name() string {
	return sr.song.Name
}

func (sr *SongResolver) SubName() *string {
	return optStr(sr.song.SubName)
}

func (sr *SongResolver) AuthorName() *string {
	return optStr(sr.song.AuthorName)
}

func (sr *SongResolver) BeatsPerMinute() *float64 {
	return optFloat64(sr.song.BeatsPerMinute)
}

func (sr *SongResolver) CoverURI() *string {
	if sr.song.CoverPath != "" {
		return optStr(fmt.Sprintf("/api/cover/%s", sr.song.ID.String()))
	}
	return nil
}

func (sr *SongResolver) AudioURI() *string {
	if sr.song.AudioPath != "" {
		return optStr(fmt.Sprintf("/api/audio/%s", sr.song.ID.String()))
	}
	return nil
}
