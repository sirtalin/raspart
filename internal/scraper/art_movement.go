package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

// GetArtMovementURLs returns a list with the URLs of the art art movements of wikiart
func GetArtMovementURLs(artistByArtMovementURL string) []string {
	var artMovementURLs []string
	var wikiartURL string = "https://wikiart.org"
	var scraperListBeginning int = 0

	c := colly.NewCollector()

	c.OnHTML("ul[class=dictionaries-list] li[class=dottedItem] a", func(e *colly.HTMLElement) {
		var href string = e.Attr("href")
		if href != "" {
			artMovementURL := fmt.Sprintf("%s%s", wikiartURL, href)
			artMovementURLs = append(artMovementURLs, artMovementURL)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		logrus.Debugf("Visiting %s", r.URL.String())
	})

	c.Visit(artistByArtMovementURL)

	artMovementURLs = artMovementURLs[scraperListBeginning:]
	logrus.Infof("Obtained %d art movements URLs from %s", len(artMovementURLs), artistByArtMovementURL)

	return artMovementURLs
}

func GetArtMovementNameFromURL(artMovementURL string) string {
	var artMovementURLSplit []string = strings.Split(artMovementURL, "/")
	var artMovementName string = artMovementURLSplit[len(artMovementURLSplit)-1]
	return strings.Split(artMovementName, "#")[0]
}
