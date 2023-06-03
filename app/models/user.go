package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"type:varchar;unique;not null"`
	Password  string    `json:"password" gorm:"type:varchar;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"default:null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:null"`
}
