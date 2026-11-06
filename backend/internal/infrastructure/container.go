package infrastructure

import (
	"log"

	"github.com/MAlifiahM/stockflow/backend/internal/config"
	"github.com/MAlifiahM/stockflow/backend/internal/domain"
	"github.com/MAlifiahM/stockflow/backend/internal/product"
	"github.com/MAlifiahM/stockflow/backend/internal/sales_order"
	"github.com/MAlifiahM/stockflow/backend/internal/sales_order_item"
	"github.com/MAlifiahM/stockflow/backend/internal/stock_movement"
	"github.com/MAlifiahM/stockflow/backend/pkg/xlogger"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

var (
	cfg config.Config

	productRepository        domain.ProductRepository
	stockMovementRepository  domain.StockMovementRepository
	salesOrderRepository     domain.SalesOrderRepository
	salesOrderItemRepository domain.SalesOrderItemRepository

	productService    domain.ProductService
	salesOrderService domain.SalesOrderService

	apiLogger *zerolog.Logger
)

func init() {
	// Load environment configuration
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] .env not found -> using system environment")
	}

	// Parse environment variables to config struct
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parsing env: %v", err)
	}

	// Initialize logging
	xlogger.Setup(cfg)
	apiLogger = xlogger.NewFileLogger("logs/api.log", cfg.IsDevelopment)

	dbSetup()
	txManager := NewTransactionManager(db)

	productRepository = product.NewPgsqlProductRepository(db)
	stockMovementRepository = stock_movement.NewPgsqlStockMovementRepository(db)
	salesOrderRepository = sales_order.NewPgsqlSalesOrderRepository(db)
	salesOrderItemRepository = sales_order_item.NewPgsqlSalesOrderItemRepository(db)

	productService = product.NewProductService(
		productRepository,
		stockMovementRepository,
		apiLogger,
		&cfg,
		txManager,
	)

	salesOrderService = sales_order.NewSalesOrderService(
		productRepository,
		salesOrderRepository,
		salesOrderItemRepository,
		stockMovementRepository,
		apiLogger,
		&cfg,
		txManager,
	)
}
