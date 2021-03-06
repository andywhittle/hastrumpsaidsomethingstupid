package search

import (
	"log"
	"net/http"
	"strings"

	"github.com/andrewstuart/goq"
)

const url = "http://www.bbc.co.uk/news"

// BBCNewsPage holds the query data for the BBC News landing page
type BBCNewsPage struct {
	Titles []string `goquery:"h3.gs-c-promo-heading__title"`
}

// BBCNews search struct
type BBCNews struct {
	Keyword string
}

// Headlines returns all matching headlines to keyword
func (bns *BBCNews) Headlines() []string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var page BBCNewsPage
	err = goq.NewDecoder(res.Body).Decode(&page)
	if err != nil {
		log.Fatal(err)
	}

	headlines := []string{}
	for _, title := range page.Titles {
		if strings.Contains(strings.ToLower(title), bns.Keyword) {
			headlines = append(headlines, title)
		}
	}

	return headlines
}
