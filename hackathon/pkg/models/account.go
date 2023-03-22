package models

import "time"

type Account struct {
	Id        int64
	Username  string
	Password  string
	TokenHash string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
