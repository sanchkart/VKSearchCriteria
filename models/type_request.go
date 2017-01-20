package models

import (
	"database/sql/driver"
	"fmt"
)

type TypeRequest string

const (
	MEMBERS_INTERSECT TypeRequest = "INTERSECTION"
)

func (g *TypeRequest) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", src)
	}
	*g = TypeRequest(string(b))
	return nil
}

func (u TypeRequest) Value() (driver.Value, error) { return string(u), nil }
