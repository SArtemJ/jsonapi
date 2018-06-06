package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestGetData(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/values?page[number]=3&page[size]=5", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	//go PrepareData(101, 201)
	//go PrepareData(201, 301)
	handler := http.HandlerFunc(GetData)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}
