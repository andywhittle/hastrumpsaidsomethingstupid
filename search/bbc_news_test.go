package search

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/andrewstuart/goq"
	"github.com/stretchr/testify/assert"
)

type fakeClient struct {
	Resp *http.Response
	Err  error
}

func (f *fakeClient) Get(url string) (resp *http.Response, err error) {
	return f.Resp, f.Err
}

type fakeDecoder struct {
}

func (fd fakeDecoder) Decode(interface{}) error {
	return errors.New("decode failed")
}

func TestBBCNewsSearchHeadlines(t *testing.T) {
	var tests = []struct {
		desc     string
		keyword  string
		given    string
		givenErr error
		decoder  func(io.Reader) Decodeable

		expected    []string
		expectedErr error
	}{
		{
			"Given the keyword is found in the body response",
			"trump",
			`<h3 class="gs-c-promo-heading__title gel-pica-bold nw-o-link-split__text">trump said something ridiculous</h3>`,
			nil,
			nil,

			[]string{"trump said something ridiculous"},
			nil,
		},
		{
			"Given the keyword is not preset in the body response",
			"trump",
			`<h3 class="gs-c-promo-heading__title gel-pica-bold nw-o-link-split__text">some other headline</h3>`,
			nil,
			nil,

			[]string{},
			nil,
		},
		{
			"Given an error is returned",
			"trump",
			"",
			errors.New("some error"),
			nil,

			[]string{},
			errors.New("failed to fetch BBC news search: some error"),
		},
		{
			"Given a decoder error is returned",
			"trump",
			`<h3 class="gs-c-promo-heading__title gel-pica-bold nw-o-link-split__text">trump said something ridiculous</h3>`,
			nil,
			func(b io.Reader) Decodeable { return fakeDecoder{} },

			nil,
			errors.New("failed to decode BBC news search body: decode failed"),
		},
	}

	for _, test := range tests {
		bb := ioutil.NopCloser(bytes.NewBufferString(test.given))
		resp := http.Response{Body: bb}
		client := &fakeClient{Resp: &resp, Err: test.givenErr}
		s := BBCNews{
			Client:  client,
			Keyword: test.keyword,
			Decoder: test.decoder,
		}

		// subject
		hl, err := s.Headlines()

		if err != nil && assert.NotNil(t, test.expectedErr, test.desc) {
			assert.EqualError(t, err, test.expectedErr.Error(), test.desc)
		} else {
			assert.Equal(t, test.expected, hl, test.desc)
			assert.Nil(t, err, test.desc)
		}
	}
}

func TestDecodeBody(t *testing.T) {
	fd := func(b io.Reader) Decodeable { return fakeDecoder{} }
	var tests = []struct {
		desc   string
		decode func(b io.Reader) Decodeable

		expected interface{}
	}{
		{
			"Given no decoder is defined",
			nil,

			&goq.Decoder{},
		},
		{
			"Given a decoder is defined",
			fd,

			fakeDecoder{},
		},
	}

	for _, test := range tests {
		s := BBCNews{Decoder: test.decode}
		r := bytes.NewBufferString("test")
		assert.IsType(t, test.expected, s.DecodeBody(r))
	}
}
