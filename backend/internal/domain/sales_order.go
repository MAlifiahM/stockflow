package domain

import (
	"time"

	"gorm.io/gorm"
)

type SalesOrder struct {
	ID          int32      `json:"id" gorm:"primaryKey;autoIncrement"`
	Customer    string     `json:"customer" gorm:"type:text;not null"`
	Status      string     `json:"status" gorm:"type:text;not null;default:CREATED"` // CREATED or CANCELLED
	Total       int32      `json:"total" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"not null;autoCreateTime"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`

	Items []SalesOrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

type SalesOrderStoreRequest struct {
	Customer string                       `json:"customer" validate:"required"`
	Items    []SalesOrderItemStoreRequest `json:"items" validate:"required,min=1,dive"`
}

type SalesOrderItemStoreRequest struct {
	SKU string `json:"sku" validate:"required"`
	Qty int32  `json:"qty" validate:"required,gt=0"`
}

type SalesOrderStoreResponse struct {
	ID        int32                         `json:"id"`
	Customer  string                        `json:"customer"`
	Status    string                        `json:"status"`
	Total     int32                         `json:"total"`
	CreatedAt time.Time                     `json:"created_at"`
	Items     []SalesOrderItemStoreResponse `json:"items"`
}

type SalesOrderItemStoreResponse struct {
	SKU       string `json:"sku"`
	Qty       int32  `json:"qty"`
	UnitPrice int32  `json:"unit_price"`
}

type SalesOrderCancelledResponse struct {
	ID          int32      `json:"id"`
	Status      string     `json:"status"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
}

type SalesOrderPaginationRequest struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
}

type SalesOrderPaginationResponse struct {
	ID          int32                         `json:"id"`
	Customer    string                        `json:"customer"`
	Status      string                        `json:"status"`
	Total       int32                         `json:"total"`
	CreatedAt   time.Time                     `json:"created_at"`
	CancelledAt *time.Time                    `json:"cancelled_at,omitempty"`
	Items       []SalesOrderItemStoreResponse `json:"items,omitempty"`
}

type SalesOrderRepository interface {
	WithTx(tx *gorm.DB) SalesOrderRepository
	Create(payload *SalesOrder) error
	FindByID(id int32) (*SalesOrder, error)
	Updates(id int32, updates map[string]interface{}) error
	Find(payload *SalesOrderPaginationRequest) ([]SalesOrder, int64, error)
}

type SalesOrderService interface {
	Store(payload *SalesOrderStoreRequest) (*SalesOrderStoreResponse, error)
	Cancel(id int32) (*SalesOrderCancelledResponse, error)
	FindByID(id int32) (*SalesOrderStoreResponse, error)
	Find(payload *SalesOrderPaginationRequest) ([]SalesOrderPaginationResponse, int64, error)
}
