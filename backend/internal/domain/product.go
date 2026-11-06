package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID             int32      `json:"id" gorm:"primaryKey"`
	SKU            string     `json:"sku" gorm:"unique;not null"`
	Name           string     `json:"name" gorm:"not null"`
	Price          int32      `json:"price" gorm:"not null"`
	Stock          int32      `json:"stock" gorm:"not null;default:0"`
	IsActive       bool       `json:"is_active" gorm:"not null;default:true"`
	CreatedAt      time.Time  `json:"created_at" gorm:"not null;autoCreateTime"`
	LastMovementAt *time.Time `json:"last_movement_at" gorm:"->"`
}

type ProductPaginationRequest struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
}

type ProductPaginationResponse struct {
	ID             int32      `json:"id"`
	SKU            string     `json:"sku"`
	Name           string     `json:"name"`
	Price          int32      `json:"price"`
	Stock          int32      `json:"stock"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	LastMovementAt *time.Time `json:"last_movement_at,omitempty"`
}

type ProductStoreRequest struct {
	SKU   string `json:"sku" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Price int32  `json:"price" validate:"required,gte=0"`
}

type ProductUpdateRequest struct {
	SKU      *string `json:"sku"`
	Name     *string `json:"name"`
	Price    *int32  `json:"price"`
	IsActive *bool   `json:"is_active"`
}

type ProductUpdateStockRequest struct {
	Type string  `json:"type" validate:"required,oneof=in out"`
	Qty  int32   `json:"qty" validate:"required,gt=0"`
	Note *string `json:"note"`
}

type ProductRepository interface {
	WithTx(tx *gorm.DB) ProductRepository
	FindPagination(payload *ProductPaginationRequest) ([]Product, int64, error)
	Create(payload *Product) error
	FindByID(id int32) (*Product, error)
	FindByIDForUpdate(id int32) (*Product, error)
	Updates(id int32, payload map[string]interface{}) error
	Delete(id int32) error
	FindProductsBySKUs(skus []string) ([]Product, error)
	FindProductsBySKUsForUpdate(skus []string) ([]Product, error)
}

type ProductService interface {
	Find(payload *ProductPaginationRequest) ([]Product, int64, error)
	Store(payload *ProductStoreRequest) (*Product, error)
	Updates(id int32, payload *ProductUpdateRequest) (*Product, error)
	Delete(id int32) error
	UpdateStock(id int32, payload *ProductUpdateStockRequest) (*Product, error)
	FindMovementsByIDProduct(id int32) ([]StockMovement, error)
}
