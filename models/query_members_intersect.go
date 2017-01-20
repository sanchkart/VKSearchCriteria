package models

type QueryMembers struct {
	Auth      string   `json:"auth"`
	Groups    []string `json:"groups"`
	MemberMin int      `json:"member_min"`
}
