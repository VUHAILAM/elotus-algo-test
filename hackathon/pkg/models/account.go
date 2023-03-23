package models

import "time"

type Account struct {
	Id        int64     `json:"id" gorm:"primaryKey;auto_increment"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password"`
	TokenHash string    `json:"token_hash"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
