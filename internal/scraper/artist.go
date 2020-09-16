package scraper

import (
	"fmt"
	"strings"
	"time"

	"github.com/gertd/go-pluralize"
	"github.com/gocolly/colly"
	"github.com/sirtalin/raspart/internal/model"
	"github.com/sirtalin/raspart/pkg/utils"
	"github.com/sirupsen/logrus"
)

// GetArtistURLs returns a list with the urls of the authors in that art movement
func GetArtistURLs(artMovementURL string) []string {
	var artistsURLs []string
	var wikiartURL string = "https://wikiart.org"

	c := colly.NewCollector()

	c.OnHTML(".artist-name a", func(e *colly.HTMLElement) {
		artistURL := fmt.Sprintf("%s%s", wikiartURL, e.Attr("href"))
		artistsURLs = append(artistsURLs, artistURL)
	})

	c.OnRequest(func(r *colly.Request) {
		logrus.Debugf("Visiting %s", r.URL.String())
	})

	c.Visit(artMovementURL)

	logrus.Infof("Obtained %d art movements URLs from %s", len(artistsURLs), artMovementURL)
	return artistsURLs
}

// GetArtist scrap the artist information
func GetArtist(artistURL string) *model.Artist {
	var artist *model.Artist = new(model.Artist)
	var nationalities []string
	var artMovements []string
	var paintingSchools []string
	var layout string = "January 2, 2006"
	pluralize := pluralize.NewClient()

	c := colly.NewCollector()

	c.OnHTML(".wiki-layout-artist-info", func(e *colly.HTMLElement) {
		var err error
		var birthDate time.Time
		var deathDate time.Time
		artist.Name = utils.PrepareString(e.ChildText("h1"))
		artist.OriginalName = e.ChildText("h2[itemprop=additionalName]")
		if artist.OriginalName == "" {
			artist.OriginalName = e.ChildText("h1")
		}

		birthDate, err = time.Parse(layout, e.ChildText("span[itemprop=birthDate]"))
		if err != nil {
			logrus.Warningf("Error while parsing birth date. %s", err)
		}
		artist.BirthDate = model.Date{birthDate}

		deathDate, err = time.Parse(layout, e.ChildText("span[itemprop=deathDate]"))
		if err != nil {
			logrus.Warningf("Error while parsing death date. %s", err)
		}
		artist.DeathDate = model.Date{deathDate}

		e.ForEach("span[itemprop=nationality]", func(_ int, el *colly.HTMLElement) {
			var nationality string = pluralize.Singular(strings.ToLower(el.Text))
			nationalities = append(nationalities, nationality)
		})
		artist.Nationalities = model.StringList{List: nationalities}

		e.ForEach("li.dictionary-values", func(_ int, el *colly.HTMLElement) {
			var text string = utils.TrimAllSpaces(el.Text)
			var textArray []string = strings.Split(text, ": ")
			if strings.Contains(text, "Art Movement") {
				for _, artMovementName := range strings.Split(textArray[1], ", ") {
					var artMovement string = utils.PrepareString(artMovementName)
					artMovements = append(artMovements, artMovement)
				}
				artist.ArtMovements = model.StringList{List: artMovements}
			}
			if strings.Contains(text, "Painting School") {
				for _, paintingSchoolName := range strings.Split(textArray[1], ", ") {
					var paintingSchool string = utils.PrepareString(paintingSchoolName)
					paintingSchools = append(paintingSchools, paintingSchool)
				}
				artist.PaintingSchools = model.StringList{List: paintingSchools}
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		logrus.Debugf("Visiting %s", r.URL.String())
	})

	c.Visit(artistURL)

	return artist
}
