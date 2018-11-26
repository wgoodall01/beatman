package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type SongID struct {
	Name        string // From directory title
	ID, Version int    // From beatsaver.com key, format
}

var beatsaverIDRegex = regexp.MustCompile(`[0-9]+-[0-9]+`)

func ParseIDURL(escaped string) (id SongID, err error) {
	decoded, err := url.PathUnescape(escaped)
	if err != nil {
		return SongID{}, err
	}
	return ParseID(decoded), nil
}

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

func (sid *SongID) URLEncode() string {
	return url.PathEscape(sid.String())
}
