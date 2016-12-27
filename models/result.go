package models

import (
	"time"
	"github.com/satori/go.uuid"
)

type Result struct {
	RequestUuid	uuid.UUID
	Id		int
	AddedAt	time.Time
}