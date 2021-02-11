package main

import (
	"io"

	"golang.org/x/net/html"
)

// tokenize gather the links, title, doctype, heading from input body
func tokenize(body io.Reader) *SiteInfo {
	tokenizer := html.NewTokenizer(body)

	sInfo := NewSiteInfo()

	for {
		tt := tokenizer.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return sInfo
		case tt == html.DoctypeToken:
			sInfo.HtmlVer = tokenizer.Token().String()

		case tt == html.StartTagToken:
			t := tokenizer.Token()

			// Check the opening "a" tag get the href attribute
			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						link := a.Val
						sInfo.links = append(sInfo.links, link)

						break
					}
				}
			}

			// get site title
			isTitle := t.Data == "title"
			if isTitle {
				tt := tokenizer.Next()

				if tt == html.TextToken {
					sInfo.Title = tokenizer.Token().Data
				}
			}

			// Count site heading
			_, ok := sInfo.Heading[t.Data]
			if ok {
				sInfo.Heading[t.Data]++
			}

		}
	}
}
