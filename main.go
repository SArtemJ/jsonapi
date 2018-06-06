package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/google/jsonapi"
	"strconv"
	"fmt"
	"os"
)

type NestedData struct {
	ID    int    `jsonapi:"primary,values"`
	Name  string `jsonapi:"attr,name"`
	Value string `jsonapi:"attr,value"`
}

var nD []*NestedData

func PrepareData() {
	for i := 0; i < 10; i++ {
		nD = append(nD, &NestedData{
			ID:    i,
			Name:  "name" + strconv.Itoa(i),
			Value: "value" + strconv.Itoa(i),
		})
	}
}

func init() {
	PrepareData()
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

	pn, _ := strconv.Atoi(r.URL.Query().Get("pagenumber"))
	//ps, _ := strconv.Atoi(r.URL.Query().Get("pagesize"))

	jsonapi.MarshalPayload(os.Stdout, nD)
	if pn < len(nD) {
		fmt.Printf("\n %v", nD[pn])
	}
}

func PostData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusCreated)

	var tmp []*NestedData
	jsonapi.UnmarshalPayload(r.Body, &tmp)
	for k, v := range tmp {
		fmt.Printf("%v - %v", k, v)
	}

}
