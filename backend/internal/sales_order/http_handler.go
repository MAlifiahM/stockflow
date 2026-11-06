package sales_order

import (
	"strconv"
	"strings"

	"github.com/MAlifiahM/stockflow/backend/internal/config"
	"github.com/MAlifiahM/stockflow/backend/internal/domain"
	"github.com/MAlifiahM/stockflow/backend/internal/middleware/validation"
	"github.com/MAlifiahM/stockflow/backend/internal/utilities"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

type HttpSalesOrderHandler struct {
	salesOrderSvc domain.SalesOrderService
	cfg           *config.Config
	logger        *zerolog.Logger
}

func NewHttpHandler(
	r fiber.Router,
	salesOrderSvc domain.SalesOrderService,
	cfg *config.Config,
	logger *zerolog.Logger,
) {
	handler := &HttpSalesOrderHandler{
		salesOrderSvc: salesOrderSvc,
		cfg:           cfg,
		logger:        logger,
	}

	r.Get("/", handler.FindSalesOrders)
	r.Post("/", validation.New[domain.SalesOrderStoreRequest](), handler.StoreSalesOrder)
	r.Post("/:id/cancel", handler.CancelSalesOrder)
	r.Get("/:id", handler.FindSalesOrder)
}

func (h *HttpSalesOrderHandler) FindSalesOrders(c fiber.Ctx) error {
	payload := new(domain.SalesOrderPaginationRequest)
	if err := c.Bind().Query(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "INVALID_QUERY_PARAMETERS",
			},
		})
	}

	salesOrders, count, err := h.salesOrderSvc.Find(payload)
	if err != nil {
		utilities.LogError(
			h.logger,
			"HttpSalesOrderHandler",
			"FindSalesOrders",
			nil,
			err,
			"failed to find sales orders",
		)

		return c.Status(fiber.StatusInternalServerError).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "FAILED_TO_FIND_SALES_ORDERS",
			},
		})
	}

	totalPages := count / int64(payload.Limit)
	if count%int64(payload.Limit) != 0 {
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(domain.PaginationResponse{
		Data:       salesOrders,
		Page:       int64(payload.Page),
		Limit:      int64(payload.Limit),
		TotalPages: totalPages,
		TotalItems: count,
	})
}

func (h *HttpSalesOrderHandler) StoreSalesOrder(c fiber.Ctx) error {
	payload := utilities.ExtractStructFromValidator[domain.SalesOrderStoreRequest](c)

	salesOrder, err := h.salesOrderSvc.Store(payload)
	if err != nil {
		utilities.LogError(
			h.logger,
			"HttpSalesOrderHandler",
			"StoreSalesOrder",
			nil,
			err,
			"failed to store sales order",
		)

		parts := strings.SplitN(err.Error(), " - ", 3)

		if len(parts) >= 2 {
			errorCode := ""

			errorDetail := &domain.ErrorCodeDetail{}

			switch parts[0] {
			case "product not found":
				errorCode = "PRODUCT_NOT_FOUND"
				errorDetail.SKU = &parts[1]

			case "product is not active":
				errorCode = "PRODUCT_NOT_ACTIVE"
				errorDetail.SKU = &parts[1]

			case "insufficient stock for product":
				if len(parts) == 3 {
					errorCode = "INSUFFICIENT_STOCK"
					errorDetail.SKU = &parts[1]

					available, parseErr := strconv.ParseInt(parts[2], 10, 32)
					if parseErr == nil {
						stock := int32(available)
						errorDetail.Available = &stock
					}
				}
			}

			if errorCode != "" {
				errorDetail.Code = errorCode

				return c.Status(fiber.StatusBadRequest).JSON(&domain.ErrorDetail{
					Error: errorDetail,
				})
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "FAILED_TO_STORE_SALES_ORDER",
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Response{
		Data: salesOrder,
	})
}

func (h *HttpSalesOrderHandler) CancelSalesOrder(c fiber.Ctx) error {
	salesOrderID := c.Params("id")

	id, err := strconv.ParseInt(salesOrderID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "SALES_ORDER_NOT_FOUND",
			},
		})
	}

	productIDInt := int32(id)

	cancelledOrder, err := h.salesOrderSvc.Cancel(productIDInt)
	if err != nil {
		utilities.LogError(
			h.logger,
			"HttpSalesOrderHandler",
			"CancelSalesOrder",
			nil,
			err,
			"failed to cancel sales order",
		)

		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(&domain.ErrorDetail{
				Error: &domain.ErrorCodeDetail{
					Code: "SALES_ORDER_NOT_FOUND",
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "FAILED_TO_CANCEL_SALES_ORDER",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.Response{
		Data: cancelledOrder,
	})
}

func (h *HttpSalesOrderHandler) FindSalesOrder(c fiber.Ctx) error {
	salesOrderID := c.Params("id")

	id, err := strconv.ParseInt(salesOrderID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "SALES_ORDER_NOT_FOUND",
			},
		})
	}

	productIDInt := int32(id)

	salesOrder, err := h.salesOrderSvc.FindByID(productIDInt)
	if err != nil {
		utilities.LogError(
			h.logger,
			"HttpSalesOrderHandler",
			"FindSalesOrder",
			nil,
			err,
			"failed to find sales order",
		)

		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(&domain.ErrorDetail{
				Error: &domain.ErrorCodeDetail{
					Code: "SALES_ORDER_NOT_FOUND",
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "FAILED_TO_FIND_SALES_ORDER",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.Response{
		Data: salesOrder,
	})
}
