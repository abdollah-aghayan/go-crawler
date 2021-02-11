package main

import (
	"strings"
	"testing"
)

func TestTokenizer(t *testing.T) {
	tt := []struct {
		name string
		body string
		out  SiteInfo
	}{
		{
			name: "No body",
			body: "",
			out:  SiteInfo{},
		},
		{
			name: "body with title",
			body: `<title>Google</title>`,
			out:  SiteInfo{Title: "Google"},
		},
		{
			name: "body with links",
			body: `<a class="gb_g" data-pid="23" href="https://mail.google.com/mail/" target="_top">Gmail</a>`,
			out:  SiteInfo{links: []string{"https://mail.google.com/mail/"}},
		},
		{
			name: "body with Doctype",
			body: `<!doctype html>`,
			out:  SiteInfo{HtmlVer: "<!DOCTYPE html>"},
		},
		{
			name: "body with Headings",
			body: `<h1>H1 title</h1><h1>H1 title</h1><h4>H4 title</h4>`,
			out:  SiteInfo{Heading: map[string]int{"h1": 2, "h4": 1}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			si := tokenize(strings.NewReader(tc.body))
			if tc.out.Title != si.Title {
				t.Fatalf("Expected title to be %s got %s", tc.out.Title, si.Title)
			}

			if tc.out.HtmlVer != si.HtmlVer {
				t.Fatalf("Expected Html version to be %s got %s", tc.out.HtmlVer, si.HtmlVer)
			}

			if len(tc.out.links) != len(si.links) {
				t.Fatalf("Expected link length to be %d got %d", len(tc.out.links), len(si.links))
			}

			for i, v := range tc.out.links {
				if v != si.links[i] {
					t.Fatalf("Expected index %d to be %s got %s", i, v, si.links[i])
				}
			}

			for key, val := range tc.out.Heading {
				if si.Heading[key] != val {
					t.Fatalf("Expected index %s's value to be %d got %d", key, val, si.Heading[key])
				}
			}

		})
	}
}
