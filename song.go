package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

func openDir(path string) (dir *os.File, err error) {
	// Open the directory
	dir, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	// Fail if the file isn't a directory.
	dirInfo, err := dir.Stat()
	if err != nil {
		return nil, err
	}
	if !dirInfo.IsDir() {
		return nil, errors.New("openDir: not a directory")
	}

	return dir, err
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
			}
		}
	}

	return song, nil
}

func LoadCustomSongs(path string) (songs []*Song, err error) {
	dir, err := openDir(path)
	if err != nil {
		return nil, err
	}

	// Get the names of the song directories
	filenames, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	for _, filename := range filenames {
		song, err := LoadSong(filepath.Join(path, filename))
		if err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}

	return songs, err
}

func (s *Song) String() string {
	return fmt.Sprintf("%s (%s) [by %s, %.1fbpm]", s.Name, s.SubName, s.AuthorName, s.BeatsPerMinute)
}
