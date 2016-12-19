package main

import (
	"gopkg.in/pg.v5"
	"time"
	"../models"
)

func createSchema(db *pg.DB) error {
	queries := []string{
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

func Insert(db *pg.DB, result *models.Result) {
	err := db.Insert(&result)
	if(err != nil) {
		panic(err)
	}
}

func Read(db *pg.DB, id int64)  {
	result := models.Result{Id: id}
	err := db.Select(&result)
	if err != nil {
		panic(err)
	}
}

func Update(db *pg.DB, result *models.Result) {
	err := db.Update(result)
	if(err != nil) {
		panic(err)
	}
}

func Delete(db *pg.DB, result *models.Result) {
	err := db.Delete(result)
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
	createSchema(db)

	result1 := &models.Result{
		ResultId:	1,
		RequestUuid:	1,
		Id:	1,
		AddedAt: time.Now(),
	}

	Insert(db, result1)

	//Insert(db, result)
}


