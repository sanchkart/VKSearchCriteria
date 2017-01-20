package main

import (
	"encoding/json"
	"fmt"
	"github.com/gammazero/workerpool"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"os"
	"time"
	"upper.io/db.v2"
	"upper.io/db.v2/lib/sqlbuilder"
	"upper.io/db.v2/postgresql"

	"vkbackend/models"
	"vkbackend/vk"

	_ "github.com/lib/pq"
)

type IntersectId struct {
	RequestId string `json:"request_id"`
}

type ResultResponse struct {
	Ids    []string      `json:"ids"`
	Status models.Status `json:"status"`
}

var sess sqlbuilder.Database
var wp *workerpool.WorkerPool

func handler(requestId string) {
	log.Printf("processing request %v", requestId)

	requests := sess.Collection("requests")
	var rq models.Request
	res := requests.Find("request_uuid", requestId)
	err := res.One(&rq)
	if err != nil {
		log.Fatalf("unable to find request by id: %v", err)
	}

	param := models.Param{}
	err = json.Unmarshal(rq.Params, &param)
	if err != nil {
		log.Fatalf("unable to decode request params: %v", err)
	}

	var list, status = vk.MathGroups(param.Groups, param.MembersMin, "", requestId)

	log.Printf("storing results for request %v", requestId)

	results := sess.Collection("results")

	for e := list.Front(); e != nil; e = e.Next() {
		_, err = results.Insert(e.Value)
		if err != nil {
			log.Fatalf("unable to save result: %v", err)
		}
	}

	rq.Status = status
	err = res.Update(rq)
	if err != nil {
		log.Fatalf("unable to update request: %v", err)
	}

	log.Printf("request %v proceed", requestId)
}

func main() {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file")
	}

	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	config := LoadConfig("config.json")

	var settings = postgresql.ConnectionURL{
		Host:     config.Host,
		Database: config.Db,
		User:     config.User,
		Password: config.Password,
	}

	wp = workerpool.New(config.CountGoroutine)
	sess, err = postgresql.Open(settings)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	defer sess.Close()

	var notDoneRqs []models.Request
	requests := sess.Collection("requests")
	err = requests.Find(db.Cond{"status": "PROCESSING"}).All(&notDoneRqs)

	if err != nil {
		log.Fatalf("unable to find requests: %v", err)
	}

	log.Printf("not finished requests count: %v", len(notDoneRqs))

	for _, each := range notDoneRqs {
		requestId := each.RequestUuid
		wp.Submit(func() {
			handler(requestId)
		})
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/tsa/members_intersect", MembersIntersect)
	router.HandleFunc("/tsa/get_result", GetResult)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.ServicePort), router))
}

func MembersIntersect(w http.ResponseWriter, r *http.Request) {
	log.Println("members_intersect request")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not found.")
		log.Println("not a POST request")
		return
	}

	var query models.QueryMembers

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Can't decode JSON.")
		log.Printf("can't decode json %q", err)
		return
	}

	jsonParams, err := json.Marshal(models.Param{query.Groups, query.MemberMin})
	if err != nil {
		log.Fatalf("error:", err)
	}

	var auth = query.Auth

	if len(auth) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Not auth token passed.")
		log.Println("no auth token")
		return
	}

	var user models.User
	users := sess.Collection("users")
	err = users.Find("auth", auth).One(&user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error")
		log.Printf("error looking for user by auth token: %q", err)
		return
	}

	if len(user.Auth) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized")
		log.Println("no user found")
		return
	}

	request := &models.Request{
		RequestUuid: uuid.NewV4().String(),
		UserUuid:    user.UserUuid,
		TypeRequest: models.MEMBERS_INTERSECT,
		CreatedAt:   time.Now(),
		Status:      models.PROCESSING,
		Params:      jsonParams,
	}

	requests := sess.Collection("requests")
	_, err = requests.Insert(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error")
		log.Println("can't insert request: %q", err)
		return
	}

	requestId := request.RequestUuid
	wp.Submit(func() {
		handler(requestId)
	})

	resp := IntersectId{
		RequestId: requestId,
	}

	result, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(result))

	log.Println("request proceed")
}

func GetResult(w http.ResponseWriter, r *http.Request) {
	log.Println("get_result request")

	var auth = r.URL.Query().Get("auth")
	if len(auth) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No auth passed")
		log.Println("no auth token passed")
		return
	}

	users := sess.Collection("users")

	var userKey models.User
	err := users.Find("auth", auth).One(&userKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error")
		log.Printf("unable to find user: %v", err)
		return
	}

	if len(userKey.Auth) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized")
		log.Printf("user not found by key %v", auth)
		return
	}

	var requestId = r.URL.Query().Get("request_id")
	if len(requestId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No request_id passed")
		log.Println("no request_id passed")
		return
	}

	requests := sess.Collection("requests")

	var request models.Request
	err = requests.Find("request_uuid", requestId).One(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error")
		log.Printf("unable to find request: %v", err)
		return
	}

	if len(request.RequestUuid) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Request not found")
		log.Printf("request not found by id %v", requestId)
		return
	}

	var resultMas []models.Result
	results := sess.Collection("results")
	results.Find("request_uuid", requestId).All(&resultMas)
	var ids []string
	ids = make([]string, 0)
	for _, element := range resultMas {
		ids = append(ids, element.Id)
	}

	resGet := ResultResponse{
		Ids:    ids,
		Status: request.Status,
	}
	result, err := json.Marshal(resGet)
	if err != nil {
		log.Fatalf("unable to serialize json: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(result))

	log.Println("get_result done")
}
