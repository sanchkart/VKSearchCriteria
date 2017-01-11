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
		`CREATE TABLE users (user_uuid  text, keyAuth text, name  text)`,
		`CREATE TABLE results (result_id serial, request_uuid text, id text, added_at timestamp)`,
		`CREATE TYPE status_type AS ENUM ('PROCESSING', 'DONE', 'CANCELLED', 'ERROR')`,
		`CREATE TYPE request_type AS ENUM ('INTERSECTION')`,
		`CREATE TABLE requests (request_uuid uuid, user_uuid uuid, type_request request_type, created_at timestamp, status status_type, params text)`,
		`INSERT INTO users (user_uuid, keyAuth, name) VALUES ('123', '46259032-320e-4a42-b610-270661feb008', 'admin')`,
	}

	for _, q := range queries {
		_, err := sess.Exec(q)
		if err != nil {
			println(err.Error())
		}
	}
}
