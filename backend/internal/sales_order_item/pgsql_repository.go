package sales_order_item

import (
	"github.com/MAlifiahM/stockflow/backend/internal/domain"
	"gorm.io/gorm"
)

type pgsqlSalesOrderItemRepository struct {
	db *gorm.DB
}

func NewPgsqlSalesOrderItemRepository(db *gorm.DB) domain.SalesOrderItemRepository {
	return &pgsqlSalesOrderItemRepository{db: db}
}

func (r *pgsqlSalesOrderItemRepository) WithTx(tx *gorm.DB) domain.SalesOrderItemRepository {
	return NewPgsqlSalesOrderItemRepository(tx)
}

func (r *pgsqlSalesOrderItemRepository) Create(payload *domain.SalesOrderItem) error {
	return r.db.Create(payload).Error
}
