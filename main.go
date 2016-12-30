package main

import (
	_"upper.io/db.v2/postgresql"
	"upper.io/db.v2/postgresql"
)

func main() {
	var settings = postgresql.ConnectionURL{
		Host:     "127.0.0.1:5432",
		User:     "postgres",
		Password: "411207",
	}
	sess, err := postgresql.Open(settings)
	if err != nil {

	}
	queries := []string{
		`CREATE TABLE users (user_uuid  uuid, key text, name  text)`,
		`CREATE TABLE results (result_id serial, request_uuid uuid, id text, added_at timestamp)`,
		`CREATE TABLE requests (request_uuid uuid, user_uuid uuid, type_request integer, created_at timestamp, status integer, params text)`,
	}

	for _, q := range queries {
		_, err := sess.Exec(q)
		if err != nil {

		}
	}
}
