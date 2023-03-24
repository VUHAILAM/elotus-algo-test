package models

import "time"

type Image struct {
	Id          int64     `json:"id" gorm:"primaryKey;auto_increment"`
	Username    string    `json:"username"`
	Name        string    `json:"name"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
