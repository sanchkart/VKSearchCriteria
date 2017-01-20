package models

type User struct {
	UserUuid string `db:"user_uuid" json:"user_uuid"`
	Auth     string `db:"auth" json:"auth"`
	Name     string `db:"name" json:"name"`
}
