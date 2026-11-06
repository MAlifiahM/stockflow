package infrastructure

import (
	"fmt"

	"github.com/MAlifiahM/stockflow/backend/internal/product"
	"github.com/MAlifiahM/stockflow/backend/internal/sales_order"
	"github.com/MAlifiahM/stockflow/backend/pkg/xlogger"
	"github.com/gofiber/contrib/v3/zerolog"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/etag"
	recover2 "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func Run() {
	logger := xlogger.Logger

	app := fiber.New(fiber.Config{
		ProxyHeader:  cfg.ProxyHeader,
		ErrorHandler: defaultErrorHandler,
	})

	app.Use(zerolog.New(zerolog.Config{
		Logger: logger,
		Fields: cfg.LogFields,
	}))
	app.Use(recover2.New())
	app.Use(etag.New())
	app.Use(requestid.New())

	api := app.Group("/api")
	product.NewHttpHandler(api.Group("/products"), productService, &cfg, apiLogger)
	sales_order.NewHttpHandler(api.Group("/orders"), salesOrderService, &cfg, apiLogger)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Info().Msgf("Server is running on address: %s", addr)
	if err := app.Listen(addr); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
