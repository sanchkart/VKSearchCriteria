package models

type QueryGetResult struct {
	Auth      string `json:"auth"`
	RequestId string `json:"request_id"`
	Offset    int    `json:"offset"`
}
