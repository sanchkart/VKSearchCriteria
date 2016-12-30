package models

import (
	"time"
	"github.com/satori/go.uuid"
)

type Result struct {
	RequestUuid	uuid.UUID `db:"request_uuid" json:"request_uuid"`
	Id		string `db:"id" json:"id"`
	AddedAt	time.Time `db:"added_at" json:"added_at"`
}