package models

import (
	"time"
)

type Customer struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar;not null"`
	Email       string    `json:"email" gorm:"type:varchar;unique;not null"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar;not null"`
	Gender      string    `json:"gender" gorm:"type:varchar;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:null"`
}
