package models

import (
	"time"
)

type Result struct {
	ResultId    int64     `db:"result_id,omitempty" json:"result_id"`
	RequestUuid string    `db:"request_uuid" json:"request_uuid"`
	Id          string    `db:"id" json:"id"`
	AddedAt     time.Time `db:"added_at" json:"added_at"`
}
