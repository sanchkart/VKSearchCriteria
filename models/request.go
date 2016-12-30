package models

import ("time"
	"github.com/satori/go.uuid"
)


type Request struct {
	RequestUuid	uuid.UUID	 `db:"request_uuid" json:"request_uuid"`
	UserUuid	uuid.UUID		`db:"user_uuid" json:"user_uuid"`
	TypeRequest	TypeRequest		`db:"type_request" json:"type_request"`
	CreatedAt	time.Time	`db:"created_at" json:"created_at"`
	Status		Status		`db:"status" json:"status"`
	Params		[]byte		`db:"params" json:"params"`
}