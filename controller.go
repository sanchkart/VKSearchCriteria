package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"./data_access"
	"./models"
	"gopkg.in/pg.v5"
	"time"
	"strconv"
	"encoding/json"
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
		decoder := json.NewDecoder(r.Body)
		var query models.QueryMembers;
		err := decoder.Decode(&query);
		if err != nil {
			panic(err)
		}

		db := pg.Connect(&pg.Options{
			User: "postgres",
			Password: "411207",
		})

		request := &models.Request{
			UserUuid:	"",
			TypeRequest:	"Create",
			CreatedAt: time.Now(),
		}

		data_access.InsertRequest(db, request);
		//runtime.GOMAXPROCS(utils.LoadConfiguration().CountGoroutine)
		//log.Println(len(vk_utils.MathGroups(groups, int(memberMin),1000,utils.LoadConfiguration().CountGoroutine)))
	}
	fmt.Fprintln(w, "MembersIntersect")
}

func GetResult(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.RawQuery) > 0 {
		var id, err = strconv.ParseInt(r.URL.Query().Get("request_id"), 10, 32)
		if err != nil {
			w.WriteHeader(404)
		}

		db := pg.Connect(&pg.Options{
			User:     "postgres",
			Password: "411207",
		})

		var data= data_access.ReadResult(db, id);
		fmt.Fprintln(w, data.Id)
	}
}
