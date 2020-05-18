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

func TestIndexHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		indexHandler(w, r)
	}))
	defer ts.Close()

	// try to get url "/", shall return a success
	res, body := getHTTP(ts, "/")
	if res.StatusCode != 200 {
		t.Errorf("wrong status code, got %v expected %v", res.StatusCode, 200)
		fmt.Printf("%s", body)
	}

	// try to get url "/xxxxxx", shall return a 404 error
	res, body = getHTTP(ts, "/xxxxxx")
	if res.StatusCode != 404 {
		t.Errorf("wrong status code, got %v expected %v", res.StatusCode, 404)
		fmt.Printf("%s", body)
	}
}
