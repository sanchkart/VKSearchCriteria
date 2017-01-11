package models

import (
	"time"
)

type Result struct {
	ResultId int `db:"result_id" json:"result_id"`
	RequestUuid	string `db:"request_uuid" json:"request_uuid"`
	Id		string `db:"id" json:"id"`
	AddedAt	time.Time `db:"added_at" json:"added_at"`
}