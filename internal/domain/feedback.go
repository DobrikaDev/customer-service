package domain

import "time"

type Feedback struct {
	ID         string    `json:"id" db:"id"`
	CustomerID string    `json:"customer_id" db:"customer_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Rating     int       `json:"rating" db:"rating"`
	Comment    string    `json:"comment" db:"comment"`
	TaskID     string    `json:"task_id" db:"task_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
