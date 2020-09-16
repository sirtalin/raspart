package main

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/sirtalin/raspart/internal/model"
	"github.com/sirtalin/raspart/internal/scraper"
	"github.com/sirtalin/raspart/pkg/utils"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
	})
}

func main() {
	var artMovementURLs []string = scraper.GetArtMovementURLs("https://www.wikiart.org/en/artists-by-art-movement")
	var artists []*model.Artist

	artistsFile, err := os.OpenFile("artists.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		logrus.Panic(err)
	}
	defer artistsFile.Close()

	err = gocsv.MarshalFile(&artists, artistsFile)

	if err := gocsv.UnmarshalFile(artistsFile, &artists); err != nil {
		logrus.Error(err)
	}
	if _, err := artistsFile.Seek(0, 0); err != nil {
		logrus.Panic(err)
	}

	for _, artMovementURL := range artMovementURLs {
		var artistsURLs []string = scraper.GetArtistURLs(artMovementURL)
		var artMovementName string = scraper.GetArtMovementNameFromURL(artMovementURL)

		for _, artistURL := range artistsURLs {
			var artist *model.Artist = scraper.GetArtist(artistURL)
			logrus.WithFields(logrus.Fields{
				"movement": artMovementName,
			}).Debug(artist)
			artists = append(artists, artist)
		}
	}

	artists = utils.UniqueArtistsSlice(artists)
	err = gocsv.MarshalFile(&artists, artistsFile)
	if err != nil {
		logrus.Panic(err)
	}
}
