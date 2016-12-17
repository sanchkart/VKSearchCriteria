package Models

import "time"

type Request struct {
	RequestUuid	int64
	UserUuid	int64
	TypeRequest	string
	CreatedAt	time.Time
	Status		string
	Params		string
}