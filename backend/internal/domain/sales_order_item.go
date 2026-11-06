package domain

import "gorm.io/gorm"

type SalesOrderItem struct {
	ID        int32  `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   int32  `json:"order_id" gorm:"not null"`
	SKU       string `json:"sku" gorm:"type:text;not null"`
	Qty       int32  `json:"qty" gorm:"not null"`
	UnitPrice int32  `json:"unit_price" gorm:"not null"`

	Product    Product    `json:"product,omitempty" gorm:"foreignKey:SKU;references:SKU"`
	SalesOrder SalesOrder `json:"-" gorm:"foreignKey:OrderID"`
}

type SalesOrderItemRepository interface {
	WithTx(tx *gorm.DB) SalesOrderItemRepository
	Create(payload *SalesOrderItem) error
}

type SalesOrderItemService interface {
}
