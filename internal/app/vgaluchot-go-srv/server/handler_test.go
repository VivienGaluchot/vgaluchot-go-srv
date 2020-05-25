package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getHTTP(ts *httptest.Server, url string) (*http.Response, []byte) {
	res, err := http.Get(ts.URL + url)
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
		{"/contact", 200},
		{"/contact.html", 200},
		{"/xxxxxx", 404},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("test-url-status_code%s", tc.url), func(t *testing.T) {
			res, body := getHTTP(ts, tc.url)
			if res.StatusCode != tc.statusCode {
				t.Errorf("wrong status code, got %v expected %v", res.StatusCode, tc.statusCode)
				fmt.Printf("%s", body)
			}
		})
	}
}
