package product

import (
	"github.com/MAlifiahM/stockflow/backend/internal/database"
	"github.com/MAlifiahM/stockflow/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type pgsqlProductRepository struct {
	db *gorm.DB
}

func NewPgsqlProductRepository(db *gorm.DB) domain.ProductRepository {
	return &pgsqlProductRepository{db: db}
}

func (r *pgsqlProductRepository) WithTx(tx *gorm.DB) domain.ProductRepository {
	return NewPgsqlProductRepository(tx)
}

func (r *pgsqlProductRepository) FindPagination(payload *domain.ProductPaginationRequest) ([]domain.Product, int64, error) {
	var products []domain.Product
	var count int64

	query := database.BuildQueryWithJoins(
		r.db,
		&domain.Product{},
		[]string{
			`LEFT JOIN (
            SELECT
                product_id,
                MAX(created_at) AS last_movement_at
            FROM stock_movements
            GROUP BY product_id
        ) sm ON sm.product_id = products.id`,
		},
		nil,
	)

	if payload.Search != "" {
		search := "%" + payload.Search + "%"
		query = query.Where("products.sku ILIKE ? OR products.name ILIKE ?", search, search)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	query = query.Select(`
		products.*,
		sm.last_movement_at AS last_movement_at
	`)

	query = query.Order("products.id ASC")

	if payload.Page > 0 && payload.Limit > 0 {
		query = query.Limit(payload.Limit).Offset((payload.Page - 1) * payload.Limit)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (r *pgsqlProductRepository) Create(payload *domain.Product) error {
	return r.db.Create(payload).Error
}

func (r *pgsqlProductRepository) FindByID(id int32) (*domain.Product, error) {
	var product domain.Product

	query := database.BuildQueryWithJoins(
		r.db,
		&domain.Product{},
		[]string{
			`LEFT JOIN (
			SELECT
				product_id,
				MAX(created_at) AS last_movement_at
			FROM stock_movements
			GROUP BY product_id
		) sm ON sm.product_id = products.id`,
		},
		nil,
	)

	query = query.Select(`
		products.*,
		sm.last_movement_at AS last_movement_at
	`)

	if err := query.First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *pgsqlProductRepository) FindByIDForUpdate(id int32) (*domain.Product, error) {
	var product domain.Product

	query := database.BuildBaseQuery(
		r.db,
		&domain.Product{},
		nil,
	)

	if err := query.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *pgsqlProductRepository) Updates(id int32, payload map[string]interface{}) error {
	return r.db.Model(&domain.Product{}).Where("id = ?", id).Updates(payload).Error
}

func (r *pgsqlProductRepository) Delete(id int32) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *pgsqlProductRepository) FindProductsBySKUs(skus []string) ([]domain.Product, error) {
	var products []domain.Product

	if len(skus) == 0 {
		return products, nil
	}

	query := database.BuildQueryWithJoins(
		r.db,
		&domain.Product{},
		[]string{
			`LEFT JOIN (
			SELECT
				product_id,
				MAX(created_at) AS last_movement_at
			FROM stock_movements
			GROUP BY product_id
		) sm ON sm.product_id = products.id`,
		},
		nil,
	)

	query = query.Select(`
		products.*,
		sm.last_movement_at AS last_movement_at
	`)

	query = query.Where("sku IN ?", skus)

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *pgsqlProductRepository) FindProductsBySKUsForUpdate(skus []string) ([]domain.Product, error) {
	var products []domain.Product

	if len(skus) == 0 {
		return products, nil
	}

	query := database.BuildBaseQuery(
		r.db,
		&domain.Product{},
		nil,
	)

	query = query.Where("sku IN ?", skus)
	query = query.Order("sku ASC")

	if err := query.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
