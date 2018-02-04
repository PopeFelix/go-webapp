package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(handler)

	hf.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code from handler: %v <> %v",
			status, http.StatusOK)
	}

	expected := "fnord\n"
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("Unexpected body from handler: %v <> %v", actual, expected)
	}
}

func TestRouter(t *testing.T) {

	r := newRouter()
	mock := httptest.NewServer(r)
	resp, err := http.Get(mock.URL + "/hello")

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code from handler: %v <> %v",
			resp.StatusCode, http.StatusOK)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}
	expected := "fnord\n"
	respString := string(b)
	if respString != expected {
		t.Errorf("Unexpected body from handler: %v <> %v", respString, expected)
	}
}

func TestRouterForNonExistentRoute(t *testing.T) {
	r := newRouter()
	mock := httptest.NewServer(r)
	resp, err := http.Post(mock.URL+"/hello", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected status code from handler: %v <> %v",
			resp.StatusCode, http.StatusMethodNotAllowed)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}
	expected := ""
	respString := string(b)
	if respString != expected {
		t.Errorf("Unexpected body from handler: %v <> %v", respString, expected)
	}
}

func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	mock := httptest.NewServer(r)
	resp, err := http.Get(mock.URL + "/assets/")

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code from handler: %v <> %v",
			resp.StatusCode, http.StatusOK)
	}

	contentType := resp.Header.Get("Content-Type")
	expected := "text/html; charset=utf-8"

	if expected != contentType {
		t.Errorf("Unexpected content type: %v <> %v", contentType, expected)
	}
}
