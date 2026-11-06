package domain

import (
	"time"

	"gorm.io/gorm"
)

type StockMovement struct {
	ID        int32     `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID int32     `json:"product_id" gorm:"not null"`
	Type      string    `json:"type" gorm:"type:text;not null"` // in or out
	Qty       int32     `json:"qty" gorm:"not null"`
	Note      *string   `json:"note,omitempty" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;autoCreateTime"`

	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

type StockMovementRepository interface {
	WithTx(tx *gorm.DB) StockMovementRepository
	Create(payload *StockMovement) error
	FindByProductID(productID int32) ([]StockMovement, error)
}

type StockMovementService interface {
}
