package models

import "time"

type Result struct {
	RequestUuid	string
	Id		int
	AddedAt	time.Time
}