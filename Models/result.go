package Models

import "time"

type Result struct {
	ResultId	int64
	RequestUuid	int64
	Id		int64
	AddedAt	time.Time
}