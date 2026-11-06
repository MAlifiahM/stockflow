package product

import (
	"errors"
	"strconv"
	"strings"

	"github.com/MAlifiahM/stockflow/backend/internal/config"
	"github.com/MAlifiahM/stockflow/backend/internal/domain"
	"github.com/MAlifiahM/stockflow/backend/internal/middleware/validation"
	"github.com/MAlifiahM/stockflow/backend/internal/utilities"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type HttpProductHandler struct {
	productSvc domain.ProductService
	cfg        *config.Config
	logger     *zerolog.Logger
}

func NewHttpHandler(
	r fiber.Router,
	productSvc domain.ProductService,
	cfg *config.Config,
	logger *zerolog.Logger,
) {
	handler := &HttpProductHandler{
		productSvc: productSvc,
		cfg:        cfg,
		logger:     logger,
	}

	r.Get("/", handler.FindProducts)
	r.Post("/", validation.New[domain.ProductStoreRequest](), handler.StoreProduct)
	r.Put("/:id", validation.New[domain.ProductUpdateRequest](), handler.UpdateProduct)
	r.Delete("/:id", handler.DeleteProduct)
	r.Post("/:id/stock", validation.New[domain.ProductUpdateStockRequest](), handler.UpdateProductStock)
	r.Get("/:id/movements", handler.FindProductMovements)
}

func (h *HttpProductHandler) FindProducts(c fiber.Ctx) error {
	payload := new(domain.ProductPaginationRequest)
	if err := c.Bind().Query(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "INVALID_QUERY_PARAMETERS",
			},
		})
	}

	products, count, err := h.productSvc.Find(payload)
	if err != nil {
		utilities.LogError(
			h.logger,
			"HttpProductHandler",
			"FindProducts",
			nil,
			err,
			"Failed to retrieve products",
		)
		return c.Status(fiber.StatusInternalServerError).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "FAILED_TO_RETRIEVE_PRODUCTS",
			},
		})
	}

	totalPages := count / int64(payload.Limit)
	if count%int64(payload.Limit) != 0 {
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(domain.PaginationResponse{
		Data:       products,
		Page:       int64(payload.Page),
		Limit:      int64(payload.Limit),
		TotalPages: totalPages,
		TotalItems: count,
	})
}

func (h *HttpProductHandler) StoreProduct(c fiber.Ctx) error {
	payload := utilities.ExtractStructFromValidator[domain.ProductStoreRequest](c)

	product, err := h.productSvc.Store(payload)

	if err != nil {
		utilities.LogError(
			h.logger,
			"HttpProductHandler",
			"StoreProduct",
			nil,
			err,
			"Failed to store product",
		)

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(&domain.ErrorDetail{
				Error: &domain.ErrorCodeDetail{
					Code: "SKU_ALREADY_EXISTS",
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "FAILED_TO_STORE_PRODUCT",
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Response{
		Data: product,
	})
}

func (h *HttpProductHandler) UpdateProduct(c fiber.Ctx) error {
	productID := c.Params("id")

	id, err := strconv.ParseInt(productID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "PRODUCT_NOT_FOUND",
			},
		})
	}

	productIDInt := int32(id)

	payload := utilities.ExtractStructFromValidator[domain.ProductUpdateRequest](c)

	product, err := h.productSvc.Updates(productIDInt, payload)

	if err != nil {
		utilities.LogError(
			h.logger,
			"HttpProductHandler",
			"UpdateProduct",
			nil,
			err,
			"Failed to update product",
		)

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(&domain.ErrorDetail{
				Error: &domain.ErrorCodeDetail{
					Code: "SKU_ALREADY_EXISTS",
				},
			})
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&domain.ErrorDetail{
				Error: &domain.ErrorCodeDetail{
					Code: "PRODUCT_NOT_FOUND",
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "FAILED_TO_UPDATE_PRODUCT",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.Response{
		Data: product,
	})
}

func (h *HttpProductHandler) DeleteProduct(c fiber.Ctx) error {
	productID := c.Params("id")

	id, err := strconv.ParseInt(productID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "PRODUCT_NOT_FOUND",
			},
		})
	}

	productIDInt := int32(id)

	if err := h.productSvc.Delete(productIDInt); err != nil {
		utilities.LogError(
			h.logger,
			"HttpProductHandler",
			"DeleteProduct",
			nil,
			err,
			"Failed to delete product",
		)

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return c.Status(fiber.StatusConflict).JSON(&domain.ErrorDetail{
				Error: &domain.ErrorCodeDetail{
					Code: "PRODUCT_HAS_HISTORY",
				},
			})
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&domain.ErrorDetail{
				Error: &domain.ErrorCodeDetail{
					Code: "PRODUCT_NOT_FOUND",
				},
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "FAILED_TO_DELETE_PRODUCT",
			},
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(domain.Response{
		Data: nil,
	})
}

func (h *HttpProductHandler) UpdateProductStock(c fiber.Ctx) error {
	productID := c.Params("id")

	id, err := strconv.ParseInt(productID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "PRODUCT_NOT_FOUND",
			},
		})
	}

	productIDInt := int32(id)

	payload := utilities.ExtractStructFromValidator[domain.ProductUpdateStockRequest](c)

	product, err := h.productSvc.UpdateStock(productIDInt, payload)

	if err != nil {
		utilities.LogError(
			h.logger,
			"HttpProductHandler",
			"UpdateProductStock",
			nil,
			err,
			"Failed to update product stock",
		)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&domain.ErrorDetail{
				Error: &domain.ErrorCodeDetail{
					Code: "PRODUCT_NOT_FOUND",
				},
			})
		}

		parts := strings.SplitN(err.Error(), "-", 2)

		if len(parts) == 2 && parts[0] == "INSUFFICIENT_STOCK" {
			available, _ := strconv.ParseInt(parts[1], 10, 32)
			availableStock := int32(available)

			return c.Status(fiber.StatusConflict).JSON(&domain.ErrorDetail{
				Error: &domain.ErrorCodeDetail{
					Code:      "INSUFFICIENT_STOCK",
					Available: &availableStock,
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "FAILED_TO_UPDATE_PRODUCT_STOCK",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.Response{
		Data: product,
	})
}

func (h *HttpProductHandler) FindProductMovements(c fiber.Ctx) error {
	productID := c.Params("id")

	id, err := strconv.ParseInt(productID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "PRODUCT_NOT_FOUND",
			},
		})
	}

	productIDInt := int32(id)

	product, err := h.productSvc.FindMovementsByIDProduct(productIDInt)
	if err != nil {
		utilities.LogError(
			h.logger,
			"HttpProductHandler",
			"FindProductMovements",
			nil,
			err,
			"Failed to find product movements",
		)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&domain.ErrorDetail{
				Error: &domain.ErrorCodeDetail{
					Code: "PRODUCT_NOT_FOUND",
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(&domain.ErrorDetail{
			Error: &domain.ErrorCodeDetail{
				Code: "FAILED_TO_FIND_PRODUCT_MOVEMENTS",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.Response{
		Data: product,
	})
}
