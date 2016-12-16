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

func Insert(user *User) {
	db := pg.Connect(&pg.Options{
		User: "postgres",
	})
	err := db.Insert(user)
	if(err != nil) {
		panic(err)
	}
}

func Read(id int)  {
	db := pg.Connect(&pg.Options{
		User: "postgres",
	})
	user := User{Id: id}
	err := db.Select(&user)
	if err != nil {
		panic(err)
	}
}
func Update(user *User) {
	db := pg.Connect(&pg.Options{
		User: "postgres",
	})
	err := db.Update(user)
	if(err != nil) {
		panic(err)
	}
}

func Delete(user *User) {
	db := pg.Connect(&pg.Options{
		User: "postgres",
	})
	err := db.Delete(user)
	if(err != nil) {
		panic(err)
	}
}


