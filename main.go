package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/google/jsonapi"
	"strconv"
	"io"
)

type NestedData struct {
	ID    int    `jsonapi:"primary,values"`
	Name  string `jsonapi:"attr,name"`
	Value string `jsonapi:"attr,value"`
}

var nD []*NestedData

func PrepareData(start int, end int) {
	for i := start; i < end; i++ {
		nD = append(nD, &NestedData{
			ID:    i,
			Name:  "name" + strconv.Itoa(i),
			Value: "value" + strconv.Itoa(i),
		})
	}
}

func init() {
	PrepareData(0, 10)
}

func main() {

	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/values", GetData).Methods("GET")
	s.HandleFunc("/values", PostData).Methods("POST")
	http.ListenAndServe(":8000", r)
}

func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonapi.MediaType)

	pn, _ := strconv.Atoi(r.URL.Query().Get("page[number]"))
	ps, _ := strconv.Atoi(r.URL.Query().Get("page[size]"))

	if validData, logic := validationPageSize(pn, ps); logic {
		jsonapi.MarshalPayload(w, validData)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Parameters under the limit")
	}
}

func PostData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonapi.MediaType)

	var k NestedData
	if err := jsonapi.UnmarshalPayload(r.Body, &k); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	nD = append(nD, &k)
}

func validationPageSize(number int, size int) ([]*NestedData, bool) {
	startFromSlice := number * size
	endFromSlice := size * (number + 1)
	if startFromSlice < len(nD) && (endFromSlice < len(nD) && endFromSlice > startFromSlice) {
		return nD[startFromSlice:endFromSlice], true
	}
	return nil, false
}
