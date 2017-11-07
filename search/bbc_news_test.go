package search

import (
	"log"
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func TestBBCNewsSearchHeadlines(t *testing.T) {
	var tests = []struct {
		desc    string
		keyword string

		expected []string
	}{
		{
			"Given the keyword is found in the body response",
			"trump",

			[]string{
				`'I'm channelling people's Trump frustration'`,
				"Magnate Trump v career communist Xi",
			},
		},
		{
			"Given the keyword is not present in the body response",
			"leeeeroy!",

			[]string{},
		},
	}

	for _, test := range tests {
		r, err := recorder.New("fixtures/bbc_news")
		if err != nil {
			log.Fatal(err)
		}
		defer r.Stop()

		c := &http.Client{Transport: r}
		s := BBCNews{Client: c, Keyword: test.keyword}

		// subject
		assert.Equal(t, test.expected, s.Headlines(), test.desc)
	}
}
