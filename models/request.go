package models

import ("time"
	"github.com/satori/go.uuid"
)


type Request struct {
	RequestUuid	uuid.UUID
	UserUuid	string
	TypeRequest	string
	CreatedAt	time.Time
	Status		string
	Params		Param
}