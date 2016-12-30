package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"./models"
	"time"
	"encoding/json"
	"./vkutils"
	"./utils"
	"github.com/satori/go.uuid"
	"upper.io/db.v2"
	"upper.io/db.v2/postgresql"
	"upper.io/db.v2/lib/sqlbuilder"
	"github.com/gammazero/workerpool"
)

var sess sqlbuilder.Database
var num_goroutines int
var wp *workerpool.WorkerPool

func main() {
	configuration := utils.LoadConfiguration();
	num_goroutines = configuration.CountGoroutine
	var settings = postgresql.ConnectionURL{
		Host:     configuration.Host,
		User:     configuration.User,
		Password: configuration.Password,
	}

	wp = workerpool.New(num_goroutines)
	sess, _ = postgresql.Open(settings)

	var requests []models.Request
	sess.Collection("requests").Find(db.Cond{"status" : 1}).All(&requests)

	for _,element := range requests {
		wp.Submit(func() {
			var param  models.Param
			json.Unmarshal(element.Params, param)
			go vkutils.MathGroups(param.Groups,param.MembersMin,"", element.RequestUuid)
		})
	}
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index)
	router.HandleFunc("/tsa/members_intersect", MembersIntersect)
	router.HandleFunc("/tsa/get_result", GetResult)

	log.Fatal(http.ListenAndServe(":7000", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Backend VK App!")
}


func MembersIntersect(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var query models.QueryMembers
	decoder.Decode(&query)

	jsonParams, err := json.Marshal(models.Param{query.Groups,query.MemberMin})
	if err != nil {
		fmt.Println("error:", err)
	}

	newUuid := uuid.NewV4()
	userUuid := uuid.NewV4()

	request := &models.Request{
		RequestUuid:  newUuid,
		UserUuid:	userUuid,
		TypeRequest:	models.MEMBERS_INTERSECT,
		CreatedAt: time.Now(),
		Status: models.PROCESSING,
		Params: jsonParams,
	}

	requests := sess.Collection("requests")
	requests.Insert(request)

	results := sess.Collection("results")
	wp.Submit(func() {
		go func() {
			var list = vkutils.MathGroups(query.Groups,query.MemberMin,query.Auth,newUuid)
			
			var request_query models.Request
			res := requests.Find("request_uuid", newUuid.String())
			res.One(&request_query)
			request_query.Status = models.DONE
			res.Update(&request_query)

			for e := list.Front(); e != nil; e = e.Next() {
				results.Insert(e.Value)
			}
		}()
	})

	fmt.Fprintln(w, newUuid)
}

func GetResult(w http.ResponseWriter, r *http.Request) {
	results := sess.Collection("results")

	var id = r.URL.Query().Get("request_id")
	var res = results.Find("request_uuid", id)

	var result models.Result
	res.One(&result)
	fmt.Fprintln(w, result.Id)
}
