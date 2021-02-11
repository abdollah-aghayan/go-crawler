package main

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// HTTPClient is a interface has to be impelimented for Client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct{}

// Do send request
func (c Client) Do(req *http.Request) (*http.Response, error) {
	hc := http.Client{
		Timeout: 5 * time.Second,
	}

	return hc.Do(req)
}

// getUrlStats count the internal and external urls in a array without paying attention to repeated urls
func getUrlStats(url *url.URL, links []string) UrlStat {
	stat := UrlStat{}

	if len(links) == 0 {
		return stat
	}

	for _, v := range links {
		pLink, err := url.Parse(v)

		if err == nil {
			// check whether url is internal
			if isInternal(url, pLink) {
				stat.Internal++
				continue
			}
			stat.External++
		}
	}

	return stat
}

// isInternal check if a fetchedUrl is internal based on main page (pageUrl)
func isInternal(pageUrl *url.URL, fetchedUrl *url.URL) bool {
	return fetchedUrl.Host == pageUrl.Host || strings.Index(fetchedUrl.String(), "#") == 0 || len(fetchedUrl.Host) == 0
}

// getRotten check for dead links in a array
func getRotten(client HTTPClient, links []string) []string {
	wg := sync.WaitGroup{}
	wg.Add(len(links))

	// receive dead links through this chan
	rottenChan := make(chan string, len(links))

	for _, v := range links {
		// fmt.Println(v)
		// trigger a go routine
		go func(url string) {

			defer wg.Done()

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				rottenChan <- url
				return
			}

			// do the request and check the status code
			res, err := client.Do(req)
			if err != nil {
				rottenChan <- url
				return
			}

			defer res.Body.Close()

			if res.StatusCode >= http.StatusBadRequest {
				rottenChan <- url
				return
			}

		}(v)
	}

	wg.Wait()

	close(rottenChan)

	rotten := make([]string, 0, len(rottenChan))

	for v := range rottenChan {
		rotten = append(rotten, v)
	}

	return rotten
}

// checkLinks check for duplicate links and complite internal links with url Host
func checkLinks(url *url.URL, l []string) []string {
	if len(l) == 0 {
		return []string{}
	}

	// creat a map to check repeated links
	lMap := make(map[string]bool)

	res := make([]string, 0, len(l))

	for _, v := range l {

		// check repeated links
		if _, ok := lMap[v]; ok {
			continue
		}

		if strings.HasPrefix(v, "/") {
			res = append(res, url.Scheme+"://"+url.Host+v)
			continue
		}

		res = append(res, v)
	}

	return res
}
