package utils

import (
	"github.com/sirtalin/raspart/internal/model"
	"github.com/sirupsen/logrus"
)

func UniqueArtistsSlice(artistSlice []*model.Artist) []*model.Artist {
	keys := make(map[string]bool)
	artists := []*model.Artist{}
	for _, artist := range artistSlice {
		if _, value := keys[artist.Name]; !value {
			keys[artist.Name] = true
			artists = append(artists, artist)
		}
	}

	logrus.Infof("Removed %d repeated artists from a total of %d. There are %d artists available", len(artistSlice)-len(artists), len(artistSlice), len(artists))

	return artists
}
