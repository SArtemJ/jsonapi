package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"sync"
	"runtime"
	"strconv"
	"fmt"
	"github.com/google/jsonapi"
	"bytes"
)

func TestGetData(t *testing.T) {
	URL := fmt.Sprintf("/api/v1/values?page[number]=%d&page[size]=%d", 1, 3)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetData)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}

func TestPostData(t *testing.T) {

		testData := CreateToInsert(9978)
		requestBody := bytes.NewBuffer(nil)
		jsonapi.MarshalOnePayloadEmbedded(requestBody, testData)

		req, err := http.NewRequest("POST", "/api/v1/values", requestBody)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", jsonapi.MediaType)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(PostData)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusOK)
		}

}

func TestConcurrent(t *testing.T) {
	wg := &sync.WaitGroup{}
	for j := 0; j < 10; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 500; i++ {
				TestPostData(t)
			}
			wg.Done()
		}()
		runtime.Gosched()
	}
	for j := 0; j < 10; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 500; i++ {
				TestGetData(t)
			}
			wg.Done()
		}()
		runtime.Gosched()
	}
	wg.Wait()
}

func CreateToInsert(id int) *NestedData {
	var k NestedData
	k.ID = id
	k.Name = "test" + strconv.Itoa(id)
	k.Value = "test" + strconv.Itoa(id)
	return &k
}
