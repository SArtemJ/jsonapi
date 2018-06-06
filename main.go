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

	//ручное создание объекта но по сути это не обязательно post запрос
	//тонее post запрос без обработки Body
	//и таким образом можно только автогенерировать каждый раз объект с новыми  данными и зазполнять в слайс


	//ndEL := &NestedData{
	//		ID:    345,
	//		Name:  "name" + strconv.Itoa(345),
	//		Value: "value" + strconv.Itoa(345),
	//	}
	//
	//
	//

	//проблема при рьработке body
	//обработка вставки одного элемента из body запроса
	//индекс в слайсе создается но объект приходит nil
	//не разворачивается из json в структуру
	var tmp *NestedData
	jsonapi.UnmarshalPayload(r.Body, tmp)

	nD = append(nD, tmp)
	for k, v := range nD {
		fmt.Printf("%v - %v", k, v)
	}

}
