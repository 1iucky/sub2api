package handler

import (
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type ModelCatalogHandler struct {
	modelCatalogService *service.ModelCatalogService
}

func NewModelCatalogHandler(modelCatalogService *service.ModelCatalogService) *ModelCatalogHandler {
	return &ModelCatalogHandler{modelCatalogService: modelCatalogService}
}

// List handles user-facing model marketplace browsing.
// GET /api/v1/models
func (h *ModelCatalogHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	vendorID := parseOptionalInt64(c.Query("vendor_id"))
	models, pag, err := h.modelCatalogService.ListModels(c.Request.Context(), pagination.PaginationParams{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    c.DefaultQuery("sort_by", ""),
		SortOrder: c.DefaultQuery("sort_order", "asc"),
	}, service.ModelCatalogListFilters{
		Search:          strings.TrimSpace(c.Query("search")),
		Platform:        strings.TrimSpace(c.Query("platform")),
		Provider:        strings.TrimSpace(c.Query("provider")),
		VendorID:        vendorID,
		PublicOnly:      true,
		WithPricingOnly: true,
		DeduplicateByID: true,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, dto.ModelCatalogsToResponse(models), pag.Total, pag.Page, pag.PageSize)
}

// Vendors handles user-facing vendor list.
// GET /api/v1/models/vendors
func (h *ModelCatalogHandler) Vendors(c *gin.Context) {
	vendors, err := h.modelCatalogService.ListVendors(c.Request.Context(), service.ModelVendorListFilters{
		Search: strings.TrimSpace(c.Query("search")),
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.ModelVendorsToResponse(vendors))
}

func parseOptionalInt64(raw string) *int64 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return nil
	}
	return &id
}
