package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSetAndCacheHandlers(t *testing.T) {
	// Setup handlers as in main.go
	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.FormValue("key")
		value := r.FormValue("value")
		store[key] = value
	})

	http.HandleFunc("/cache", func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("key")
		v, ok := store[key]
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(v))
	})

	// Test /set
	form := url.Values{}
	form.Set("key", "foo")
	form.Set("value", "bar")
	req := httptest.NewRequest("POST", "/set", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)

	if store["foo"] != "bar" {
		t.Fatalf("expected store[foo] = bar, got %q", store["foo"])
	}

	// Test /cache
	req2 := httptest.NewRequest("GET", "/cache?key=foo", nil)
	w2 := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w2, req2)

	if w2.Body.String() != "bar" {
		t.Fatalf("expected response 'bar', got %q", w2.Body.String())
	}
}