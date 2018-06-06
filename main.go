package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/google/jsonapi"
	"strconv"
	"fmt"
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
	w.WriteHeader(http.StatusCreated)

	pn, _ := strconv.Atoi(r.URL.Query().Get("page[number]"))
	ps, _ := strconv.Atoi(r.URL.Query().Get("page[size]"))

	jsonapi.MarshalPayload(w, nD)

	fmt.Println(pn)
	fmt.Println(ps)
}

func PostData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusCreated)

	var k NestedData

	if err := jsonapi.UnmarshalPayload(r.Body, &k); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	nD = append(nD, &k)
	for k, v := range nD {
		fmt.Printf("%v - %v \n", k, v)
	}
}
