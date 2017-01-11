package models

import ("time"
)


type Request struct {
	RequestUuid	string	 `db:"request_uuid" json:"request_uuid"`
	UserUuid	string		`db:"user_uuid" json:"user_uuid"`
	TypeRequest	TypeRequest		`db:"type_request" json:"type_request"`
	CreatedAt	time.Time	`db:"created_at" json:"created_at"`
	Status		Status		`db:"status" json:"status"`
	Params		[]byte		`db:"params" json:"params"`
}