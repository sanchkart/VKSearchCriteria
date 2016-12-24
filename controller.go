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
	"strconv"
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
		var groups []string = r.URL.Query()["groups"]
		var memberMin, err = strconv.ParseInt(r.URL.Query().Get("member_min"), 10, 32)
		if(err != nil) {
			w.WriteHeader(404)
		}

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
		vkutils.MathGroups(groups,int(memberMin), 1000)
	}
	fmt.Fprintln(w, "MembersIntersect")
}

func GetResult(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Todo show:", 1)
}
