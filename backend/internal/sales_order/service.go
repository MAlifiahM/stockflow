package sales_order

import (
	"fmt"
	"strings"
	"time"

	"github.com/MAlifiahM/stockflow/backend/internal/config"
	"github.com/MAlifiahM/stockflow/backend/internal/domain"
	"github.com/MAlifiahM/stockflow/backend/internal/utilities"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type salesOrderService struct {
	productRepo        domain.ProductRepository
	salesOrderRepo     domain.SalesOrderRepository
	salesOrderItemRepo domain.SalesOrderItemRepository
	stockMovementRepo  domain.StockMovementRepository
	logger             *zerolog.Logger
	cfg                *config.Config
	txManager          domain.DBTransactionManager
}

func NewSalesOrderService(
	productRepo domain.ProductRepository,
	salesOrderRepo domain.SalesOrderRepository,
	salesOrderItemRepo domain.SalesOrderItemRepository,
	stockMovementRepo domain.StockMovementRepository,
	logger *zerolog.Logger,
	cfg *config.Config,
	txManager domain.DBTransactionManager,
) domain.SalesOrderService {
	return &salesOrderService{
		productRepo:        productRepo,
		salesOrderRepo:     salesOrderRepo,
		salesOrderItemRepo: salesOrderItemRepo,
		stockMovementRepo:  stockMovementRepo,
		logger:             logger,
		cfg:                cfg,
		txManager:          txManager,
	}
}

func (s *salesOrderService) Store(payload *domain.SalesOrderStoreRequest) (*domain.SalesOrderStoreResponse, error) {
	skus := s.extractSKUs(payload.Items)

	products, err := s.productRepo.FindProductsBySKUs(skus)
	if err != nil {
		utilities.LogError(
			s.logger,
			"SalesOrderService",
			"Store",
			nil,
			err,
			"Failed to fetch products by SKUs",
		)

		return nil, err
	}

	productMap := s.mapProductsBySKU(products)

	if err := s.validateProducts(payload, productMap); err != nil {
		return nil, err
	}

	total, orderItems := s.countTotalOrder(payload.Items, productMap)

	resp := &domain.SalesOrderStoreResponse{
		Customer: payload.Customer,
		Status:   "CREATED",
		Total:    total,
		Items:    orderItems,
	}

	err = s.txManager.Execute(func(tx *gorm.DB) error {
		return s.createSalesOrderTransaction(
			tx,
			payload,
			skus,
			total,
			orderItems,
			resp,
		)
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *salesOrderService) extractSKUs(
	items []domain.SalesOrderItemStoreRequest,
) []string {
	skus := make([]string, 0, len(items))

	for _, item := range items {
		skus = append(
			skus,
			strings.ToUpper(strings.TrimSpace(item.SKU)),
		)
	}

	return skus
}

func (s *salesOrderService) mapProductsBySKU(
	products []domain.Product,
) map[string]domain.Product {
	productMap := make(map[string]domain.Product, len(products))

	for _, product := range products {
		productMap[product.SKU] = product
	}

	return productMap
}

func (s *salesOrderService) validateProducts(payload *domain.SalesOrderStoreRequest, productMap map[string]domain.Product) error {
	for _, item := range payload.Items {
		sku := strings.ToUpper(strings.TrimSpace(item.SKU))
		product, exists := productMap[sku]

		if !exists {
			return fmt.Errorf("product not found - %s", sku)
		}

		if product.IsActive == false {
			return fmt.Errorf("product is not active - %s", sku)
		}

		if product.Stock < item.Qty {
			return fmt.Errorf("insufficient stock for product - %s - %d", sku, product.Stock)
		}
	}

	return nil
}

func (s *salesOrderService) countTotalOrder(
	items []domain.SalesOrderItemStoreRequest,
	productMap map[string]domain.Product,
) (int32, []domain.SalesOrderItemStoreResponse) {
	var total int32
	orderItems := make([]domain.SalesOrderItemStoreResponse, 0, len(items))

	for _, item := range items {
		sku := strings.ToUpper(strings.TrimSpace(item.SKU))
		product := productMap[sku]

		// count sub-total
		subtotal := product.Price * item.Qty

		// total accumulated
		total += subtotal

		// sales order items
		orderItems = append(orderItems, domain.SalesOrderItemStoreResponse{
			SKU:       product.SKU,
			Qty:       item.Qty,
			UnitPrice: product.Price,
		})
	}

	return total, orderItems
}

func (s *salesOrderService) createSalesOrderTransaction(
	tx *gorm.DB,
	payload *domain.SalesOrderStoreRequest,
	skus []string,
	total int32,
	orderItems []domain.SalesOrderItemStoreResponse,
	resp *domain.SalesOrderStoreResponse,
) error {
	utilities.LogInfo(
		s.logger,
		"SalesOrderService",
		"Store",
		nil,
		"Starting transaction for creating sales order",
	)

	lockedProducts, err := s.productRepo.WithTx(tx).
		FindProductsBySKUsForUpdate(skus)
	if err != nil {
		return err
	}

	lockedProductMap := s.mapProductsBySKU(lockedProducts)

	if err := s.validateProducts(payload, lockedProductMap); err != nil {
		return err
	}

	return s.createSalesOrderWithItems(
		tx,
		payload,
		total,
		orderItems,
		lockedProductMap,
		resp,
	)
}

func (s *salesOrderService) createSalesOrderWithItems(
	tx *gorm.DB,
	payload *domain.SalesOrderStoreRequest,
	total int32,
	orderItems []domain.SalesOrderItemStoreResponse,
	productMap map[string]domain.Product,
	resp *domain.SalesOrderStoreResponse,
) error {
	salesOrder := &domain.SalesOrder{
		Customer: payload.Customer,
		Total:    total,
	}

	if err := s.salesOrderRepo.WithTx(tx).Create(salesOrder); err != nil {
		utilities.LogError(
			s.logger,
			"SalesOrderService",
			"Store",
			nil,
			err,
			"Failed to create sales order",
		)
		return err
	}

	resp.ID = salesOrder.ID
	resp.CreatedAt = salesOrder.CreatedAt

	for _, item := range orderItems {
		if err := s.createSalesOrderItem(tx, salesOrder.ID, item, productMap, payload.Customer); err != nil {
			return err
		}
	}

	return nil
}

func (s *salesOrderService) createSalesOrderItem(
	tx *gorm.DB,
	orderID int32,
	item domain.SalesOrderItemStoreResponse,
	productMap map[string]domain.Product,
	customer string,
) error {
	product := productMap[item.SKU]

	salesOrderItem := &domain.SalesOrderItem{
		OrderID:   orderID,
		SKU:       item.SKU,
		Qty:       item.Qty,
		UnitPrice: item.UnitPrice,
	}

	if err := s.salesOrderItemRepo.WithTx(tx).Create(salesOrderItem); err != nil {
		utilities.LogError(
			s.logger,
			"SalesOrderService",
			"Store",
			nil,
			err,
			"Failed to create sales order item",
		)
		return err
	}

	if err := s.updateProductStock(tx, product, item.Qty); err != nil {
		return err
	}

	return s.createStockMovement(tx, product.ID, customer, item.Qty)
}

func (s *salesOrderService) updateProductStock(
	tx *gorm.DB,
	product domain.Product,
	qty int32,
) error {
	productUpdate := map[string]interface{}{
		"stock": product.Stock - qty,
	}

	if err := s.productRepo.WithTx(tx).Updates(product.ID, productUpdate); err != nil {
		utilities.LogError(
			s.logger,
			"SalesOrderService",
			"Store",
			nil,
			err,
			"Failed to update product stock",
		)
		return err
	}

	return nil
}

func (s *salesOrderService) createStockMovement(
	tx *gorm.DB,
	productID int32,
	customer string,
	qty int32,
) error {
	note := fmt.Sprintf("Sales Order %s", customer)
	stockMovement := &domain.StockMovement{
		ProductID: productID,
		Qty:       qty,
		Type:      "out",
		Note:      &note,
	}

	if err := s.stockMovementRepo.WithTx(tx).Create(stockMovement); err != nil {
		utilities.LogError(
			s.logger,
			"SalesOrderService",
			"Store",
			nil,
			err,
			"Failed to create stock movement",
		)
		return err
	}

	return nil
}

func (s *salesOrderService) Cancel(id int32) (*domain.SalesOrderCancelledResponse, error) {
	salesOrder, err := s.salesOrderRepo.FindByID(id)
	if err != nil {
		utilities.LogError(
			s.logger,
			"SalesOrderService",
			"Cancel",
			nil,
			err,
			"Failed to find sales order by ID",
		)
		return nil, err
	}

	cancelledAt := time.Now()
	updatesSalesOrder := map[string]interface{}{
		"status":       "CANCELLED",
		"cancelled_at": cancelledAt,
	}

	err = s.txManager.Execute(func(tx *gorm.DB) error {
		utilities.LogInfo(
			s.logger,
			"SalesOrderService",
			"Cancel",
			nil,
			"Starting transaction for cancelling sales order",
		)

		err := s.salesOrderRepo.WithTx(tx).Updates(id, updatesSalesOrder)
		if err != nil {
			utilities.LogError(
				s.logger,
				"SalesOrderService",
				"Cancel",
				nil,
				err,
				"Failed to update sales order status to CANCELLED",
			)
			return err
		}

		for _, item := range salesOrder.Items {
			product, err := s.productRepo.WithTx(tx).FindByIDForUpdate(item.Product.ID)
			if err != nil {
				utilities.LogError(
					s.logger,
					"SalesOrderService",
					"Cancel",
					nil,
					err,
					"Failed to find product by ID",
				)
				return err
			}

			productUpdate := map[string]interface{}{
				"stock": product.Stock + item.Qty,
			}

			if err := s.productRepo.WithTx(tx).Updates(product.ID, productUpdate); err != nil {
				utilities.LogError(
					s.logger,
					"SalesOrderService",
					"Cancel",
					nil,
					err,
					"Failed to update product stock",
				)
				return err
			}

			movementCancelNote := fmt.Sprintf("Sales Order %s Cancellation", salesOrder.Customer)
			movementStock := &domain.StockMovement{
				ProductID: product.ID,
				Qty:       item.Qty,
				Type:      "in",
				Note:      &movementCancelNote,
			}

			if err := s.stockMovementRepo.WithTx(tx).Create(movementStock); err != nil {
				utilities.LogError(
					s.logger,
					"SalesOrderService",
					"Cancel",
					nil,
					err,
					"Failed to create stock movement for cancellation",
				)
				return err
			}
		}

		salesOrder.Status = "CANCELLED"
		salesOrder.CancelledAt = &cancelledAt
		return nil
	})

	resp := &domain.SalesOrderCancelledResponse{
		ID:          salesOrder.ID,
		Status:      salesOrder.Status,
		CancelledAt: salesOrder.CancelledAt,
	}

	return resp, nil
}

func (s *salesOrderService) FindByID(id int32) (*domain.SalesOrderStoreResponse, error) {
	orderData, err := s.salesOrderRepo.FindByID(id)
	if err != nil {
		utilities.LogError(
			s.logger,
			"SalesOrderService",
			"FindByID",
			nil,
			err,
			"Failed to find sales order by ID",
		)
		return nil, err
	}

	resp := &domain.SalesOrderStoreResponse{
		ID:        orderData.ID,
		Customer:  orderData.Customer,
		Status:    orderData.Status,
		Total:     orderData.Total,
		CreatedAt: orderData.CreatedAt,
	}

	if len(orderData.Items) > 0 {
		items := make([]domain.SalesOrderItemStoreResponse, 0, len(orderData.Items))
		for _, item := range orderData.Items {
			items = append(items, domain.SalesOrderItemStoreResponse{
				SKU:       item.SKU,
				Qty:       item.Qty,
				UnitPrice: item.UnitPrice,
			})
		}
		resp.Items = items
	}

	return resp, nil
}
