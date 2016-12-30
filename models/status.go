package models

type Status int

const (
	PROCESSING  = 1 + iota
	DONE
	CANCELLED
)