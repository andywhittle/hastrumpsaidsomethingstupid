package search

import (
	"io"
	"net/http"
	"strings"

	"github.com/andrewstuart/goq"
	"github.com/pkg/errors"
)

const url = "http://www.bbc.co.uk/news"

// BBCNewsPage holds the query data for the BBC News landing page
type BBCNewsPage struct {
	Titles []string `goquery:"h3.gs-c-promo-heading__title"`
}

// BBCNews search struct
type BBCNews struct {
	Client  Gettable
	Decoder func(io.Reader) Decodeable
	Keyword string
}

// Gettable interface that can be satisfied by http client
type Gettable interface {
	Get(string) (*http.Response, error)
}

// ensure that http client adheres to gettable
var _ Gettable = &http.Client{}

// Decodeable is an interface to wrap decode on goq for testing
type Decodeable interface {
	Decode(interface{}) error
}

// ensure that the goq decoder is decodeable
var _ Decodeable = &goq.Decoder{}

// NewBBCNews initialise a new BBC News search
func NewBBCNews(client Gettable, key string) *BBCNews {
	return &BBCNews{
		Client:  client,
		Keyword: key,
		Decoder: func(b io.Reader) Decodeable { return goq.NewDecoder(b) },
	}
}

// Headlines returns all matching headlines to keyword
func (bns *BBCNews) Headlines() ([]string, error) {
	res, err := bns.Client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch BBC news search")
	}
	defer res.Body.Close()

	var page BBCNewsPage
	err = bns.Decoder(res.Body).Decode(&page)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode BBC news search body")
	}

	headlines := []string{}
	for _, title := range page.Titles {
		if strings.Contains(strings.ToLower(title), bns.Keyword) {
			headlines = append(headlines, title)
		}
	}

	return headlines, nil
}
