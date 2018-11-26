package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Library struct {
	sync.RWMutex

	// Where the library is located ("x/Beat Saber/")
	Path string

	songs  map[SongID]*Song
	ticker *time.Ticker
}

// NewLibrary creates a library and loads its songs from disk
func NewLibrary(path string) (library *Library, err error) {

	// Create the library
	library = &Library{
		songs: make(map[SongID]*Song),
		Path:  path,
	}

	return library, nil
}

// Load the library contents from disk, updating existing songs.
func (l *Library) Load() (new int, err error) {
	// using local locking so we don't block the Library for 20s

	log.lib.WithFields(logrus.Fields{
		"path": l.Path,
	}).Debug("Loading song library")
	customSongs := filepath.Join(l.Path, "CustomSongs")

	dir, err := openDir(customSongs)
	if err != nil {
		return 0, err
	}

	// Get the names of the song directories
	filenames, err := dir.Readdirnames(0)
	if err != nil {
		return 0, err
	}

	// keep track of songs we still have, so deletes work
	touched := make(map[SongID]bool)
	for id, _ := range l.songs {
		touched[id] = false
	}

FileLoop:
	for _, filename := range filenames {
		songPath := filepath.Join(customSongs, filename)
		if strings.HasPrefix(filename, ".") {
			log.lib.WithField("path", songPath).Warn("ignored dotfile")
			continue FileLoop
		}

		loadSong, err := LoadSong(songPath)
		if err == nil && loadSong == nil {
			log.lib.WithField("path", songPath).Warn("ignored song")
			continue FileLoop
		}
		if err != nil {
			return 0, err
		}

		songLog := log.lib.WithFields(logrus.Fields{
			"name":       loadSong.Name,
			"subName":    loadSong.SubName,
			"authorName": loadSong.AuthorName,
			"id":         loadSong.ID.String(),
			"path":       songPath,
		})

		touched[loadSong.ID] = true

		l.Lock()
		existingSong := l.songs[loadSong.ID]
		if existingSong == nil {
			// this is a new song, add it
			songLog.Info("discovered song")
			new++

			l.songs[loadSong.ID] = loadSong
		} else {
			// Update existing song
			*existingSong = *loadSong
		}
		l.Unlock()
	}

	for id, keep := range touched {
		if !keep {
			l.Lock()
			oldSong := l.songs[id]
			log.lib.WithFields(logrus.Fields{
				"name":       oldSong.Name,
				"subName":    oldSong.SubName,
				"authorName": oldSong.AuthorName,
				"id":         oldSong.ID.String(),
			}).Info("Forgot song")
			delete(l.songs, id)
			l.Unlock()
		}
	}

	return new, nil
}

// Song gets the song with the given ID
func (l *Library) Get(id SongID) *Song {
	l.RLock()
	defer l.RUnlock()

	return l.songs[id]
}

// Songs gets all songs in the library
func (l *Library) Songs() (songs []*Song) {
	l.RLock()
	defer l.RUnlock()

	songs = make([]*Song, 0, len(l.songs))
	for _, song := range l.songs {
		songs = append(songs, song)
	}

	return songs
}

// Query finds all songs in the library which match the filter.
func (l *Library) Query(filter string) (songs []*Song) {
	l.RLock()
	defer l.RUnlock()

	filter = strings.TrimSpace(filter)
	filter = strings.ToLower(filter)

FilterLoop:
	for _, song := range l.songs {
		text := fmt.Sprintf("%s %s %s", song.Name, song.SubName, song.AuthorName)
		text = strings.ToLower(text)

		tokens := strings.Split(filter, " ")
		for _, t := range tokens {
			if strings.Index(text, t) == -1 {
				continue FilterLoop
			}
		}

		songs = append(songs, song)
	}

	return songs
}

// StartSync begins to reload the Library from disk once every freq.
func (l *Library) StartSync(freq time.Duration) {
	l.ticker = time.NewTicker(freq)

	go func() {
		// run initial sync
		_, err := l.Load()
		if err != nil {
			log.lib.WithError(err).Fatal("failed to reload library in sync")
		}

		for _ = range l.ticker.C {
			new, err := l.Load()
			if err != nil {
				log.lib.WithError(err).Fatal("failed to reload library in sync")
			}
			if new != 0 {
				log.lib.WithFields(logrus.Fields{"count": new}).Info("loaded new songs")
			}
		}
	}()
}

// StopSync stops reloading the library
func (l *Library) StopSync() {
	l.ticker.Stop()
}
