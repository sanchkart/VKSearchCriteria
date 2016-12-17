package main

import (
	"fmt"
	"gopkg.in/pg.v5"
)

type User struct {
	Id     int64
	Name   string
	Emails []string
}

type Result struct {
	result_id	int
	request_uuid	int
	id		int
	added_at	int
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.Id, u.Name, u.Emails)
}

type Story struct {
	Id       int64
	Title    string
	AuthorId int64
	Author   *User
}

func (s Story) String() string {
	return fmt.Sprintf("Story<%d %s %s>", s.Id, s.Title, s.Author)
}

func createSchema(db *pg.DB) error {
	queries := []string{
		`CREATE TEMP TABLE request (request_uuid bigint, user_uuid bigint, type text, created_at datetime , status  text, params  jsonb)`,
		`CREATE TEMP TABLE result (result_id  bigint, request_uuid bigint, id bigint, added_at  dateTime)`,
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func Insert(db *pg.DB, result *Result) {
	err := db.Insert(&result)
	if(err != nil) {
		panic(err)
	}
}

//func Read(db *pg.DB, id int)  {
//	user := User{Id: id}
//	err := db.Select(&user)
//	if err != nil {
//		panic(err)
//	}
//}
func Update(db *pg.DB, result *Result) {
	err := db.Update(result)
	if(err != nil) {
		panic(err)
	}
}

func Delete(db *pg.DB, result *Result) {
	err := db.Delete(result)
	if(err != nil) {
		panic(err)
	}
}

func main()  {
	db := pg.Connect(&pg.Options{
		User: "postgres",
		Password: "411207",
	})

	result :=  &Result{
		result_id : 1,
		request_uuid : 1,
		id : 1,
		added_at : 1,
	}

	Insert(db, result)
}


