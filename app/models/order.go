package models

import (
	"time"
)

type Order struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	CustomerID    uint       `json:"customer_id" gorm:"constraint:OnDelete:CASCADE;not null"`
	ProductName   string     `json:"product_name" gorm:"type:varchar;not null"`
	Quantity      int        `json:"quantity" gorm:"not null;check:quantity >= 1"`
	TotalPrice    float64    `json:"total_price" gorm:"not null;check:total_price >= 0"`
	PaymentStatus string     `json:"payment_status" gorm:"type:varchar"`
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAt     time.Time  `json:"created_at" gorm:"default:null"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"default:null"`
}
