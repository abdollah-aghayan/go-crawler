package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type UrlStat struct {
	Internal     int `json:"internal"`
	External     int `json:"external"`
	InAccessable int `json:"inAccessable"`
}

type SiteInfo struct {
	HtmlVer string         `json:"htmlVersion"`
	Title   string         `json:"Title"`
	Heading map[string]int `json:"heading"`
	UrlStat
	links     []string ``
	LoginPage bool     `json:"isLoginPage"`
}

func NewSiteInfo() *SiteInfo {
	return &SiteInfo{
		Heading: map[string]int{"h1": 0, "h2": 0, "h3": 0, "h4": 0, "h5": 0, "h6": 0},
	}
}

func run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/fetch", fetchInfo)

	fmt.Println(http.ListenAndServe(":8090", mux))
}

func fetchInfo(w http.ResponseWriter, r *http.Request) {

	// get url from query string
	urls, ok := r.URL.Query()["url"]
	if !ok {
		JSONError(w, errors.New("url not defined"), 400)
		return
	}

	// validate the received url
	url, err := url.Parse(urls[0])
	if !ok {
		JSONError(w, errors.New("Malform url"), 400)
		return
	}

	// fetch html body
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		JSONError(w, errors.New("Can not make request to the specified url"), 404)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		JSONError(w, errors.New("Can not make request to the specified url"), 404)
		return
	}

	rc := res.Body
	defer rc.Close()

	sInfo := tokenize(rc)

	// get fetched urls statistic
	urlStat := getUrlStats(url, sInfo.links)

	links := checkLinks(url, sInfo.links)

	c := Client{}
	rotten := getRotten(c, links)

	urlStat.InAccessable = len(rotten)

	sInfo.UrlStat = urlStat

	// response
	JSON(w, sInfo, 200)
	return
}
