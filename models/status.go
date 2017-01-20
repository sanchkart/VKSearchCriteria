package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type Status string

const (
	PROCESSING Status = "PROCESSING"
	DONE       Status = "DONE"
	CANCELLED  Status = "CANCELLED"
	ERROR      Status = "ERROR"
)

func (g *Status) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", src)
	}
	*g = Status(strings.ToLower(string(b)))
	return nil
}

func (u Status) Value() (driver.Value, error) { return string(u), nil }
