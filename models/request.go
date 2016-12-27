package models

import "time"


type Request struct {
	RequestUuid	string
	UserUuid	string
	TypeRequest	string
	CreatedAt	time.Time
	Status		string
	Params		Param
}