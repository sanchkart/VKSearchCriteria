package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/tsa/members_intersect", MembersIntersect)
	router.HandleFunc("/tsa/get_result", GetResult)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func MembersIntersect(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "MembersIntersect")
}

func GetResult(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Todo show:", 1)
}
