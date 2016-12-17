package main
import "time"

type Result struct {
	request_uuid	int
	user_uuid	int
	typeRequest	string
	created_at	time.Time
	status		string
	params		string
}
