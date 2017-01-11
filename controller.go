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
	"os"
	"container/list"

	_"github.com/lib/pq"
)

type IntersectId struct {
	RequestId   uuid.UUID      `json:"request_id"`
}

type ResultResponse struct {
	Ids []string `json:"ids"`
	Status models.Status `json:"status"`
}

var sess sqlbuilder.Database
var num_goroutines int
var wp *workerpool.WorkerPool

func main() {
	file, err := os.OpenFile("log.txt", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file")
	}

	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	configuration := utils.LoadConfiguration()
	num_goroutines = configuration.CountGoroutine

	var settings,errs = postgresql.ParseURL(os.Getenv("DATABASE_URL"))

	if errs != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	wp = workerpool.New(num_goroutines)
	sess, err = postgresql.Open(settings)

	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	var requestsRemain []models.Request
	results := sess.Collection("results")
	requests := sess.Collection("requests")

	requests.Find(db.Cond{"status" : "PROCESSING"}).All(&requestsRemain)

	for _,element := range requestsRemain {
		wp.Submit(func() {
			var param  models.Param
			json.Unmarshal(element.Params, param)
			var list,statusVK = vkutils.MathGroups(param.Groups,param.MembersMin,"", element.RequestUuid)
			var w http.ResponseWriter
			w = UpdateData(requests, element.RequestUuid, err, w , list, results, statusVK)
		})
	}
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index)
	router.HandleFunc("/tsa/members_intersect", MembersIntersect)
	router.HandleFunc("/tsa/get_result", GetResult)

	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Backend VK App!")
}


func MembersIntersect(w http.ResponseWriter, r *http.Request) {
	log.Println("Start members intersect")
	if r.ContentLength == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		log.Println("Bad Request")
		return
	}

	var query models.QueryMembers

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&query)

	jsonParams, err := json.Marshal(models.Param{query.Groups,query.MemberMin})
	if err != nil {
		fmt.Println("error:", err)
	}

	var auth = query.Auth

	if auth == ""{
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w,"Unauthorized")
		log.Println("Unauthorized")
		return
	}else{
		var user models.User

		users := sess.Collection("users")
		users.Find("keyauth", auth).One(&user)

		if user.KeyAuth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w,"Unauthorized")
			log.Println("Unauthorized")
			return
		}
	}

	newUuid := uuid.NewV4()
	userUuid := uuid.NewV4()

	request := &models.Request {
		RequestUuid:  newUuid.String(),
		UserUuid:	userUuid.String(),
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
			var uuidString = newUuid.String()
			var list,statusVK = vkutils.MathGroups(query.Groups,query.MemberMin,query.Auth,uuidString)
			w = UpdateData(requests, uuidString, err, w, list, results, statusVK)
		}()
	})

	resNewUuid := IntersectId{
		RequestId: newUuid,
	}
	result, _ := json.Marshal(resNewUuid)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(result))
	log.Println("End members intersect")
}

func UpdateData(requests db.Collection, uuidString string, err error, w http.ResponseWriter, list *list.List, results db.Collection, status models.Status) http.ResponseWriter {
	var request_query models.Request
	requests.Find("request_uuid", uuidString).One(&request_query)
	query_update := sess.Update("requests").Set("status", status).Where("request_uuid = ?", request_query.RequestUuid)
	_, err = query_update.Exec()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	for e := list.Front(); e != nil; e = e.Next() {
		results.Insert(e.Value)
	}
	return w
}

func GetResult(w http.ResponseWriter, r *http.Request) {
	log.Println("Start get result")

	var auth  = r.URL.Query().Get("auth")
	if auth == ""{
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w,"Unauthorized")
		log.Println("Unauthorized")
		return
	}else{
		users := sess.Collection("users")

		var userKey models.User
		users.Find("keyauth", auth).One(&userKey)

		if userKey.KeyAuth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w,"Unauthorized")
			log.Println("Unauthorized")
			return
		}
	}

	var id = r.URL.Query().Get("request_id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		log.Println("Bad Request")
		return
	}else{
		requests := sess.Collection("requests")

		var requestKey models.Request
		requests.Find("request_uuid", id).One(&requestKey)

		if requestKey.RequestUuid == "" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w,"Not Found")
			log.Println("Not Found")
			return
		}
	}

	var resultMas []models.Result
	results := sess.Collection("results")
	results.Find("request_uuid", id).All(&resultMas)
	var ids [] string
	ids = make([]string, 0)
	for _,element := range resultMas {
		ids = append(ids, "id" + element.Id)
	}

	var request_query models.Request
	requests := sess.Collection("requests")
	requests.Find("request_uuid", id).One(&request_query)
	resGet := ResultResponse{
		Ids: ids,
		Status : request_query.Status,
	}
	result, _ := json.Marshal(resGet)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(result))
	log.Println("End get result")
}
