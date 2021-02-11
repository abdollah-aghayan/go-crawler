package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestGetUrlStats(t *testing.T) {
	tt := []struct {
		name          string
		baseUrl       string
		links         []string
		internalCount int
		externalCount int
	}{
		{
			name:          "two internal",
			baseUrl:       "http://google.com",
			links:         []string{"http://google.com/info", "/gmail"},
			internalCount: 2,
			externalCount: 0,
		},
		{
			name:          "no internal with tree external link",
			baseUrl:       "http://amazon.com",
			links:         []string{"http://renewed.amazon.com/info", "/products", "/"},
			internalCount: 2,
			externalCount: 1,
		},
		{
			name:          "No links to check",
			baseUrl:       "http://amazon.com",
			links:         []string{},
			internalCount: 0,
			externalCount: 0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			url, _ := url.Parse(tc.baseUrl)
			stats := getUrlStats(url, tc.links)

			if tc.internalCount != stats.Internal {
				t.Fatalf("Expected %d internal link found %d", tc.internalCount, stats.Internal)
			}

			if tc.externalCount != stats.External {
				t.Fatalf("Expected %d internal link found %d", tc.externalCount, stats.External)
			}
		})
	}
}

func TestCheckLinks(t *testing.T) {
	tt := []struct {
		name    string
		baseUrl string
		links   []string
		out     []string
	}{
		{
			name:    "no links",
			baseUrl: "http://xyz.com",
			links:   []string{},
			out:     []string{},
		},
		{
			name:    "relative links",
			baseUrl: "http://xyz.com",
			links:   []string{"/sample", "/test"},
			out:     []string{"http://xyz.com/sample", "http://xyz.com/test"},
		},
		{
			name:    "absolute links",
			baseUrl: "http://xyz.com",
			links:   []string{"http://xyz.com/sample", "http://xyz.com/test"},
			out:     []string{"http://xyz.com/sample", "http://xyz.com/test"},
		},
		{
			name:    "absolute and relative links",
			baseUrl: "http://xyz.com",
			links:   []string{"http://xyz.com/sample", "/test"},
			out:     []string{"http://xyz.com/sample", "http://xyz.com/test"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			url, _ := url.Parse(tc.baseUrl)
			cl := checkLinks(url, tc.links)

			if len(cl) != len(tc.out) {
				t.Fatalf("Expected slice length to be %d, got %d", len(tc.out), len(cl))
			}

			for i, v := range tc.out {
				if v != cl[i] {
					t.Fatalf("expected index %d to be %s got %s", i, v, cl[i])
				}
			}
		})
	}
}

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func (mc MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func TestGetRotten(t *testing.T) {
	c := MockClient{}

	tt := []struct {
		name   string
		client func(req *http.Request) (*http.Response, error)
		links  []string
		out    []string
	}{
		{
			name:   "no url",
			client: nil,
			links:  []string{},
			out:    []string{},
		},
		{
			name: "no dead links",
			client: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("")),
				}, nil
			},
			links: []string{"http://xyz.net"},
			out:   []string{},
		},
		{
			name: "dead links",
			client: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 400,
					Body:       ioutil.NopCloser(strings.NewReader("")),
				}, nil
			},
			links: []string{"http://xyz.net"},
			out:   []string{"http://xyz.net"},
		},
		{
			name: "one dead and one healthy link",
			client: func(req *http.Request) (*http.Response, error) {

				if strings.HasPrefix(req.URL.String(), "http://deadlink.com") {
					return &http.Response{
						StatusCode: 404,
						Body:       ioutil.NopCloser(strings.NewReader("")),
					}, errors.New("page not found")
				}

				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("")),
				}, nil
			},
			links: []string{"http://xyz.net", "http://deadlink.com"},
			out:   []string{"http://deadlink.com"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			GetDoFunc = tc.client
			res := getRotten(c, tc.links)

			if len(res) != len(tc.out) {
				t.Fatalf("Expected slice length to be %d, got %d", len(tc.out), len(res))
			}

			for i, v := range tc.out {
				if v != res[i] {
					t.Fatalf("expected index %d to be %s got %s", i, v, res[i])
				}
			}
		})
	}
}
