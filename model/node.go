package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type JSONB map[string]string

type Node struct {
	ID         uuid.UUID `json:"id"`
	NodeID     int       `json:"nodeid" gorm:"column:nodeid"`
	Model      string    `json:"model"`
	Enabled    bool      `json:"enabled"`
	Attributes JSONB     `json:"attributes"`
}

func (j *JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, j)
}
