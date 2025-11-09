package domain

import "time"

type Customer struct {
	MaxID    string       `json:"max_id" db:"max_id"`
	Name     string       `json:"name" db:"name"`
	About    string       `json:"about" db:"about"`
	Type     CustomerType `json:"type" db:"type"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}