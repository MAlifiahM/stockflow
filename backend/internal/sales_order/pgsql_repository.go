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

func (r *pgsqlSalesOrderRepository) Find(payload *domain.SalesOrderPaginationRequest) ([]domain.SalesOrder, int64, error) {
	var salesOrders []domain.SalesOrder
	var count int64

	query := database.BuildBaseQuery(
		r.db,
		&domain.SalesOrder{},
		[]string{"Items", "Items.Product"},
	)

	if payload.Search != "" {
		search := "%" + payload.Search + "%"
		query = query.Where("customer ILIKE ?", search)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if payload.Page > 0 && payload.Limit > 0 {
		query = query.Limit(payload.Limit).Offset((payload.Page - 1) * payload.Limit)
	}

	if err := query.Find(&salesOrders).Error; err != nil {
		return nil, 0, err
	}

	return salesOrders, count, nil
}
