package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func getHTTP(ts *httptest.Server, url string, headers map[string]string) (*http.Response, []byte) {
	req, _ := http.NewRequest("GET", ts.URL+url, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return res, body
}

// get all the pages in all languages
// the goal is to detect basic template error by checking the http result
func TestInstallHandlers(t *testing.T) {
	mux := http.NewServeMux()
	installHandlers(mux, true)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	testCases := []struct {
		url        string
		statusCode int
	}{
		{"/", 200},
		{"/index", 200},
		{"/index.html", 200},
		{"/portfolio", 200},
		{"/portfolio.html", 200},
		{"/xxxxxx", 404},
	}

	languages := []string{"fr", "en"}
	htmlLanguageArgument := regexp.MustCompile(`<html lang="(\S+)">`)

	for _, tc := range testCases {
		for _, lang := range languages {
			t.Run(fmt.Sprintf("test-url-status_code/%s%s", lang, tc.url), func(t *testing.T) {
				headers := map[string]string{
					"Accept-Language": lang,
				}

				res, body := getHTTP(ts, tc.url, headers)

				// check the status code is the one expected
				if res.StatusCode != tc.statusCode {
					t.Errorf("wrong status code, got %v expected %v", res.StatusCode, tc.statusCode)
					fmt.Printf("%s", body)
				}

				if res.StatusCode == 200 {
					// check the language argument is matching the Accept-Language header
					match := htmlLanguageArgument.FindSubmatch(body)
					if match == nil {
						t.Errorf("html lang argument not found")
					} else if string(match[1][:]) != lang {
						t.Errorf("wrong html lang argument, got %v expected %v", string(match[1][:]), lang)
					}
				}
			})
		}
	}
}
