package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	var tests = []struct {
		desc string
		name string

		expected string
	}{
		{
			"Given the keyword is found in the body response",
			"success",

			"Yup he's probably said something stupid",
		},
		{
			"Given the keyword is not found in the body response",
			"failure",

			"Na nothing today",
		},
	}

	for _, test := range tests {
		// build a new go vcr recorder
		rec, err := recorder.New("fixtures/" + test.name)
		if err != nil {
			log.Fatal(err)
		}
		defer rec.Stop()

		// setup a test server with our go-vcr transport recorder
		c := &http.Client{Transport: rec}
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
	}
}
