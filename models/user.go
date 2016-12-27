package models

import (
	"github.com/satori/go.uuid"
)

type User struct {
	UserUuid	uuid.UUID
	Key		string
	Name		string
}
