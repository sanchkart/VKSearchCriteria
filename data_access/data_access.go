package data_access

import (
	"gopkg.in/pg.v5"
	"time"
	"../models"
)
//TODO IN COMMON Generalize this class for working with T
func CreateSchema(db *pg.DB) error {
	queries := []string{
		`CREATE TABLE users (user_uuid  serial, key text, name  text`,
		`CREATE TABLE results (result_id serial, request_uuid serial, id serial, added_at timestamp)`,
		`CREATE TABLE requests (request_uuid serial, user_uuid serial, type_request text, created_at timestamp, status text, params text)`,
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertUser(db *pg.DB, result *models.Result) {
	err := db.Insert(&result)
	if(err != nil) {
		panic(err)
	}
}

func InsertResult(db *pg.DB, result *models.Result) {
	err := db.Insert(&result)
	if(err != nil) {
		panic(err)
	}
}

func InsertRequest(db *pg.DB, result *models.Request) {
	err := db.Insert(&result)
	if(err != nil) {
		panic(err)
	}
}

func ReadUser(db *pg.DB, id int64)  (result models.User) {
	result = models.User{UserUuid: id}
	err := db.Select(&result)
	if err != nil {
		panic(err)
	}
	return
}

func ReadResult(db *pg.DB, id int64) (result models.Result){
	result = models.Result{Id: id}
	err := db.Select(&result)
	if err != nil {
		panic(err)
	}
	return
}

func ReadRequest(db *pg.DB, id int64)  (result models.Request){
	result = models.Request{RequestUuid: id}
	err := db.Select(&result)
	if err != nil {
		panic(err)
	}

	return
}

func UpdateUser(db *pg.DB, result *models.User) {
	err := db.Update(result)
	if(err != nil) {
		panic(err)
	}
}

func UpdateResult(db *pg.DB, result *models.Result) {
	err := db.Update(result)
	if(err != nil) {
		panic(err)
	}
}

func UpdateRequest(db *pg.DB, result *models.Request) {
	err := db.Update(result)
	if(err != nil) {
		panic(err)
	}
}

func DeleteResult(db *pg.DB, result *models.Result) {
	err := db.Delete(result)
	if(err != nil) {
		panic(err)
	}
}

func DeleteRequest(db *pg.DB, request *models.Request) {
	err := db.Delete(request)
	if(err != nil) {
		panic(err)
	}
}

func main()  {
	//TODO Theese need to handle from main.go
	db := pg.Connect(&pg.Options{
		User: "postgres",
		Password: "411207",
	})
	CreateSchema(db)

	result1 := &models.Result{
		ResultId:	1,
		RequestUuid:	1,
		Id:	1,
		AddedAt: time.Now(),
	}

	InsertResult(db, result1)

	//Insert(db, result)
}


