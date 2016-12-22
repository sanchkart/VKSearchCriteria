package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"./vkutils"
	"./data_access"
	"./models"
	"gopkg.in/pg.v5"
	"time"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/tsa/members_intersect", MembersIntersect)
	router.HandleFunc("/tsa/get_result", GetResult)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Backend VK App!")
}

func MembersIntersect(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.RawQuery) > 0 {
		var groups []string = r.URL.Query().Get("groups")
		var memberMin = r.URL.Query().Get("member_min")
		db := pg.Connect(&pg.Options{
			User: "postgres",
			Password: "411207",
		})
		data_access.CreateSchema(db)

		request := &models.Request{
			RequestUuid:	1,
			UserUuid:	1,
			TypeRequest:	"Merge",
			CreatedAt: time.Now(),
		}

		data_access.InsertRequest(db, request);
		vkutils.MathGroups(groups,memberMin, 1000)
	}
	fmt.Fprintln(w, "MembersIntersect")
}

func GetResult(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Todo show:", 1)
}
