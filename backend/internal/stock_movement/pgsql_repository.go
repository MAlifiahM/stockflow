package stock_movement

import (
	"github.com/MAlifiahM/stockflow/backend/internal/database"
	"github.com/MAlifiahM/stockflow/backend/internal/domain"
	"gorm.io/gorm"
)

type pgsqlStockMovementRepository struct {
	db *gorm.DB
}

func NewPgsqlStockMovementRepository(db *gorm.DB) domain.StockMovementRepository {
	return &pgsqlStockMovementRepository{db: db}
}

func (r *pgsqlStockMovementRepository) WithTx(tx *gorm.DB) domain.StockMovementRepository {
	return &pgsqlStockMovementRepository{db: tx}
}

func (r *pgsqlStockMovementRepository) Create(payload *domain.StockMovement) error {
	return r.db.Create(payload).Error
}

func (r *pgsqlStockMovementRepository) FindByProductID(productID int32) ([]domain.StockMovement, error) {
	var movements []domain.StockMovement

	query := database.BuildBaseQuery(r.db, &domain.StockMovement{}, nil)

	query = query.Order("created_at DESC")

	err := query.Where("product_id = ?", productID).Find(&movements).Error
	if err != nil {
		return nil, err
	}

	return movements, nil
}
