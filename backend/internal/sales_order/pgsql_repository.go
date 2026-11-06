package sales_order

import (
	"github.com/MAlifiahM/stockflow/backend/internal/database"
	"github.com/MAlifiahM/stockflow/backend/internal/domain"
	"gorm.io/gorm"
)

type pgsqlSalesOrderRepository struct {
	db *gorm.DB
}

func NewPgsqlSalesOrderRepository(db *gorm.DB) domain.SalesOrderRepository {
	return &pgsqlSalesOrderRepository{db: db}
}

func (r *pgsqlSalesOrderRepository) WithTx(tx *gorm.DB) domain.SalesOrderRepository {
	return NewPgsqlSalesOrderRepository(tx)
}

func (r *pgsqlSalesOrderRepository) Create(payload *domain.SalesOrder) error {
	return r.db.Create(payload).Error
}

func (r *pgsqlSalesOrderRepository) FindByID(id int32) (*domain.SalesOrder, error) {
	var salesOrder domain.SalesOrder

	query := database.BuildBaseQuery(r.db, &domain.SalesOrder{}, []string{"Items", "Items.Product"})

	if err := query.First(&salesOrder, id).Error; err != nil {
		return nil, err
	}
	return &salesOrder, nil
}

func (r *pgsqlSalesOrderRepository) Updates(id int32, updates map[string]interface{}) error {
	return r.db.Model(&domain.SalesOrder{}).Where("id = ?", id).Updates(updates).Error
}
