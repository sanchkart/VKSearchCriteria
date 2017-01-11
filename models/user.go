package models


type User struct {
	UserUuid	string  `db:"user_uuid" json:"user_uuid"`
	KeyAuth		string  `db:"keyauth" json:"keyauth"`
	Name		string  `db:"name" json:"name"`
}
