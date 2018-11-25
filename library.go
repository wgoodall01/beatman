package main

import "path/filepath"

type Library struct {
	Songs map[string]*Song
}

func LoadLibrary(path string) (library *Library, err error) {
	library = &Library{
		Songs: make(map[string]*Song),
	}

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

		library.Songs[song.ID] = song
	}

	return library, err
}
