package product

import (
	"errors"
	"fmt"

	"github.com/MAlifiahM/stockflow/backend/internal/config"
	"github.com/MAlifiahM/stockflow/backend/internal/domain"
	"github.com/MAlifiahM/stockflow/backend/internal/utilities"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type productService struct {
	productRepo       domain.ProductRepository
	stockMovementRepo domain.StockMovementRepository
	logger            *zerolog.Logger
	cfg               *config.Config
	txManager         domain.DBTransactionManager
}

func NewProductService(
	productRepo domain.ProductRepository,
	stockMovementRepo domain.StockMovementRepository,
	logger *zerolog.Logger,
	cfg *config.Config,
	txManager domain.DBTransactionManager,
) domain.ProductService {
	return &productService{
		productRepo:       productRepo,
		stockMovementRepo: stockMovementRepo,
		logger:            logger,
		cfg:               cfg,
		txManager:         txManager,
	}
}

func (s *productService) Find(payload *domain.ProductPaginationRequest) ([]domain.Product, int64, error) {
	if payload.Page <= 0 {
		payload.Page = 1
	}

	if payload.Limit <= 0 {
		payload.Limit = 10
	}

	if len(payload.Search) > 150 {
		utilities.LogError(
			s.logger,
			"ProductService",
			"Find",
			nil,
			fmt.Errorf("search len is %d", len(payload.Search)),
			"search length exceeds 150 characters",
		)
		return nil, 0, errors.New("search length exceeds 150 characters")
	}

	products, count, err := s.productRepo.FindPagination(payload)
	if err != nil {
		utilities.LogError(
			s.logger,
			"ProductService",
			"Find",
			nil,
			err,
			"failed to find products",
		)
		return nil, 0, errors.New("failed to find products")
	}

	return products, count, nil
}

func (s *productService) Store(payload *domain.ProductStoreRequest) (*domain.Product, error) {
	if payload.Price < 0 {
		return nil, errors.New("price must be greater than or equal to 0")
	}

	product := &domain.Product{
		SKU:   payload.SKU,
		Name:  payload.Name,
		Price: payload.Price,
	}

	err := s.txManager.Execute(func(tx *gorm.DB) error {
		utilities.LogInfo(
			s.logger,
			"ProductService",
			"Store",
			nil,
			"creating product",
		)

		err := s.productRepo.WithTx(tx).Create(product)
		if err != nil {
			utilities.LogError(
				s.logger,
				"ProductService",
				"Store",
				nil,
				err,
				"failed to create product",
			)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) Updates(id int32, payload *domain.ProductUpdateRequest) (*domain.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		utilities.LogError(
			s.logger,
			"ProductService",
			"Updates",
			nil,
			err,
			"failed to find product by id",
		)

		return nil, err
	}

	updates := make(map[string]interface{})
	if payload.SKU != nil {
		updates["sku"] = *payload.SKU
		product.SKU = *payload.SKU
	}

	if payload.Name != nil {
		updates["name"] = *payload.Name
		product.Name = *payload.Name
	}

	if payload.Price != nil {
		updates["price"] = *payload.Price
		product.Price = *payload.Price
	}

	if payload.IsActive != nil {
		updates["is_active"] = *payload.IsActive
		product.IsActive = *payload.IsActive
	}

	err = s.txManager.Execute(func(tx *gorm.DB) error {
		utilities.LogInfo(
			s.logger,
			"ProductService",
			"Updates",
			nil,
			"updating product",
		)

		err := s.productRepo.WithTx(tx).Updates(id, updates)
		if err != nil {
			utilities.LogError(
				s.logger,
				"ProductService",
				"Updates",
				nil,
				err,
				"failed to update product",
			)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) Delete(id int32) error {
	_, err := s.productRepo.FindByID(id)
	if err != nil {
		utilities.LogError(
			s.logger,
			"ProductService",
			"Delete",
			nil,
			err,
			"failed to find product by id",
		)

		return err
	}

	err = s.txManager.Execute(func(tx *gorm.DB) error {
		utilities.LogInfo(
			s.logger,
			"ProductService",
			"Delete",
			nil,
			"deleting product",
		)

		err := s.productRepo.WithTx(tx).Delete(id)
		if err != nil {
			utilities.LogError(
				s.logger,
				"ProductService",
				"Delete",
				nil,
				err,
				"failed to delete product",
			)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *productService) UpdateStock(id int32, payload *domain.ProductUpdateStockRequest) (*domain.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		utilities.LogError(
			s.logger,
			"ProductService",
			"UpdateStock",
			nil,
			err,
			"failed to find product by id",
		)

		return nil, err
	}

	finalStock := product.Stock

	switch payload.Type {
	case "in":
		finalStock += payload.Qty

	case "out":
		if payload.Qty > product.Stock {
			return nil, fmt.Errorf("INSUFFICIENT_STOCK-%d", product.Stock)
		}

		finalStock -= payload.Qty
	}

	product.Stock = finalStock

	err = s.txManager.Execute(func(tx *gorm.DB) error {
		utilities.LogInfo(
			s.logger,
			"ProductService",
			"UpdateStock",
			nil,
			"updating product stock",
		)

		err := s.productRepo.WithTx(tx).Updates(id, map[string]interface{}{
			"stock": finalStock,
		})

		if err != nil {
			utilities.LogError(
				s.logger,
				"ProductService",
				"UpdateStock",
				nil,
				err,
				"failed to update product stock",
			)
			return err
		}

		stockMovement := &domain.StockMovement{
			ProductID: id,
			Type:      payload.Type,
			Qty:       payload.Qty,
			Note:      payload.Note,
		}

		err = s.stockMovementRepo.WithTx(tx).Create(stockMovement)
		if err != nil {
			utilities.LogError(
				s.logger,
				"ProductService",
				"UpdateStock",
				nil,
				err,
				"failed to create stock movement",
			)
			return err
		}

		product.LastMovementAt = &stockMovement.CreatedAt

		return nil
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) FindMovementsByIDProduct(id int32) ([]domain.StockMovement, error) {
	_, err := s.productRepo.FindByID(id)
	if err != nil {
		utilities.LogError(
			s.logger,
			"ProductService",
			"FindMovementsByIDProduct",
			nil,
			err,
			"failed to find product by id",
		)

		return nil, err
	}

	movements, err := s.stockMovementRepo.FindByProductID(id)
	if err != nil {
		utilities.LogError(
			s.logger,
			"ProductService",
			"FindMovementsByIDProduct",
			nil,
			err,
			"failed to find stock movements by product id",
		)

		return nil, err
	}

	return movements, nil
}
