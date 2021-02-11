package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchHandler(t *testing.T) {

	testSrv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			switch r.URL.Path {
			case "/main-page":
				res := `<!doctype html>
				<title>Google</title>
				<a href="/sub-page">Gmail</a>
				<h1>H1 Heading</h1>`

				w.WriteHeader(200)
				fmt.Fprint(w, res)

			case "/sub-page":
				w.WriteHeader(404)
				fmt.Fprint(w, "")
			}
		},
	))

	defer testSrv.Close()

	tt := []struct {
		name string
		url  string
		err  string
		res  string
	}{
		{name: "No Url", url: "", err: `{"error":"url not defined"}`},
		{name: "Malform Url", url: "http://localhost:8090/fetch?url=xyz.com", err: `{"error":"Can not make request to the specified url"}`},
		{name: "mocked url", url: "http://localhost:8090/fetch?url=" + testSrv.URL + "/main-page", res: `{"htmlVersion":"<!DOCTYPE html>","Title":"Google","heading":{"h1":1,"h2":0,"h3":0,"h4":0,"h5":0,"h6":0},"Internal":1,"External":0,"InAccessable":1,"isLoginPage":false}`},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			req, err := http.NewRequest("Get", tc.url, nil)
			if err != nil {
				t.Fatalf("Can not create request %v", err)
			}

			fetchInfo(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal("can not read response")
			}

			if tc.err != "" && tc.err != string(bytes.TrimSpace(b)) {
				t.Fatalf("Expect %s error got %s", tc.err, string(b))
			}

			if res.StatusCode == http.StatusOK {

				exp := SiteInfo{}

				_ = json.Unmarshal([]byte(tc.res), &exp)

				out := SiteInfo{}

				_ = json.Unmarshal([]byte(tc.res), &out)

				if exp.HtmlVer != out.HtmlVer {
					t.Fatalf("Expected %v got %v", exp.HtmlVer, out.HtmlVer)
				}

				if exp.Title != out.Title {
					t.Fatalf("Expected %v got %v", exp.Title, out.Title)
				}

				if exp.InAccessable != out.InAccessable {
					t.Fatalf("Expected %v got %v", exp.InAccessable, out.InAccessable)
				}

				if exp.Internal != out.Internal {
					t.Fatalf("Expected %v got %v", exp.Internal, out.Internal)
				}

				if exp.Heading["h1"] != out.Heading["h1"] {
					t.Fatalf("Expected %v got %v", exp.Heading["h1"], out.Heading["h1"])
				}
			}
		})
	}
}
