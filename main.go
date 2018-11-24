package main

import "fmt"

func main() {
	songs, err := LoadCustomSongs("_beatsaber/CustomSongs")
	if err != nil {
		fmt.Println(err)
	} else {
		for _, song := range songs {
			fmt.Println(song)
		}
	}
}
