package models

import "time"


type Request struct {
	RequestUuid	int
	UserUuid	int
	TypeRequest	string
	CreatedAt	time.Time
	Status		string
	Params		Param
}