package models

import "time"

type Result struct {
	ResultId	int
	RequestUuid	string
	Id		int
	AddedAt	time.Time
}