package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	var tests = []struct {
		desc        string
		name        string
		useRecorder bool

		expected string
		code     int
	}{
		{
			"Given the keyword is found in the body response",
			"success",
			true,

			"Yup he's probably said something stupid",
			http.StatusOK,
		},
		{
			"Given the keyword is not found in the body response",
			"failure",
			true,

			"Na nothing today",
			http.StatusOK,
		},
		{
			"Given ...",
			"missing",
			false,

			"",
			http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		// build a new go vcr recorder
		rec, err := recorder.New("fixtures/" + test.name)
		if err != nil {
			log.Fatal(err)
		}
		defer rec.Stop()
		c := &http.Client{Timeout: 10 * time.Nanosecond}
		if test.useRecorder {
			c.Transport = rec
		}

		// setup a test server with our go-vcr transport recorder
		r := Router(c)
		ts := httptest.NewServer(r)
		defer ts.Close()

		// subject
		res, err := http.Get(ts.URL)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		// assertions
		assert.Contains(t, string(body), test.expected, test.desc)
		assert.Equal(t, test.code, res.StatusCode)
	}
}
