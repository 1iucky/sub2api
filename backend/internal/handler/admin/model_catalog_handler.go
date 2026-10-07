package admin

import (
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
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

type modelCatalogRequest struct {
	ModelID      string         `json:"model_id" binding:"required,max=200"`
	DisplayName  string         `json:"display_name" binding:"omitempty,max=200"`
	Platform     string         `json:"platform" binding:"omitempty,max=50"`
	Provider     string         `json:"provider" binding:"omitempty,max=100"`
	VendorID     *int64         `json:"vendor_id"`
	Mode         string         `json:"mode" binding:"omitempty,max=50"`
	Description  string         `json:"description"`
	Tags         []string       `json:"tags"`
	Capabilities map[string]any `json:"capabilities"`
	Endpoints    []string       `json:"endpoints"`
	Pricing      map[string]any `json:"pricing"`
	Metadata     map[string]any `json:"metadata"`
	Status       string         `json:"status" binding:"omitempty,oneof=active disabled"`
	Visibility   string         `json:"visibility" binding:"omitempty,oneof=public admin"`
	Source       string         `json:"source" binding:"omitempty,max=50"`
	IconKey      string         `json:"icon_key" binding:"omitempty,max=80"`
}

type modelVendorRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	ProviderKey string `json:"provider_key" binding:"omitempty,max=80"`
	IconKey     string `json:"icon_key" binding:"omitempty,max=80"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// List handles admin model catalog listing.
// GET /api/v1/admin/models
func (h *ModelCatalogHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	models, pag, err := h.modelCatalogService.ListModels(c.Request.Context(), pagination.PaginationParams{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    c.DefaultQuery("sort_by", ""),
		SortOrder: c.DefaultQuery("sort_order", "asc"),
	}, service.ModelCatalogListFilters{
		Search:     strings.TrimSpace(c.Query("search")),
		Platform:   strings.TrimSpace(c.Query("platform")),
		Provider:   strings.TrimSpace(c.Query("provider")),
		VendorID:   parseModelCatalogOptionalInt64(c.Query("vendor_id")),
		Status:     strings.TrimSpace(c.Query("status")),
		Visibility: strings.TrimSpace(c.Query("visibility")),
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, dto.ModelCatalogsToResponse(models), pag.Total, pag.Page, pag.PageSize)
}

// Get handles admin model catalog detail.
// GET /api/v1/admin/models/:id
func (h *ModelCatalogHandler) Get(c *gin.Context) {
	id, ok := parseModelCatalogIDParam(c, "id")
	if !ok {
		return
	}
	model, err := h.modelCatalogService.GetModel(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.ModelCatalogToResponse(model))
}

// Create handles admin model creation.
// POST /api/v1/admin/models
func (h *ModelCatalogHandler) Create(c *gin.Context) {
	var req modelCatalogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorFrom(c, infraerrors.BadRequest("VALIDATION_ERROR", err.Error()))
		return
	}
	model, err := h.modelCatalogService.CreateModel(c.Request.Context(), req.toService())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, dto.ModelCatalogToResponse(model))
}

// Update handles admin model update.
// PUT /api/v1/admin/models/:id
func (h *ModelCatalogHandler) Update(c *gin.Context) {
	id, ok := parseModelCatalogIDParam(c, "id")
	if !ok {
		return
	}
	var req modelCatalogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorFrom(c, infraerrors.BadRequest("VALIDATION_ERROR", err.Error()))
		return
	}
	model, err := h.modelCatalogService.UpdateModel(c.Request.Context(), id, req.toService())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.ModelCatalogToResponse(model))
}

// Delete handles admin model delete.
// DELETE /api/v1/admin/models/:id
func (h *ModelCatalogHandler) Delete(c *gin.Context) {
	id, ok := parseModelCatalogIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.modelCatalogService.DeleteModel(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

// SyncFromPricing syncs model catalog entries from LiteLLM pricing data.
// POST /api/v1/admin/models/sync-pricing
func (h *ModelCatalogHandler) SyncFromPricing(c *gin.Context) {
	result, err := h.modelCatalogService.SyncFromPricing(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}

// Vendors lists vendors.
// GET /api/v1/admin/models/vendors
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

// UpsertVendor creates or updates a vendor by name.
// POST /api/v1/admin/models/vendors
func (h *ModelCatalogHandler) UpsertVendor(c *gin.Context) {
	var req modelVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorFrom(c, infraerrors.BadRequest("VALIDATION_ERROR", err.Error()))
		return
	}
	vendor, err := h.modelCatalogService.UpsertVendor(c.Request.Context(), service.ModelVendorUpsert{
		Name:        req.Name,
		ProviderKey: req.ProviderKey,
		IconKey:     req.IconKey,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.ModelVendorToResponse(vendor))
}

// DeleteVendor deletes a vendor and leaves existing models without a vendor.
// DELETE /api/v1/admin/models/vendors/:id
func (h *ModelCatalogHandler) DeleteVendor(c *gin.Context) {
	id, ok := parseModelCatalogIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.modelCatalogService.DeleteVendor(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (req modelCatalogRequest) toService() service.ModelCatalogUpsert {
	return service.ModelCatalogUpsert{
		ModelID:      req.ModelID,
		DisplayName:  req.DisplayName,
		Platform:     req.Platform,
		Provider:     req.Provider,
		VendorID:     req.VendorID,
		Mode:         req.Mode,
		Description:  req.Description,
		Tags:         req.Tags,
		Capabilities: req.Capabilities,
		Endpoints:    req.Endpoints,
		Pricing:      req.Pricing,
		Metadata:     req.Metadata,
		Status:       req.Status,
		Visibility:   req.Visibility,
		Source:       req.Source,
		IconKey:      req.IconKey,
	}
}

func parseModelCatalogIDParam(c *gin.Context, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		response.ErrorFrom(c, infraerrors.BadRequest("INVALID_ID", "invalid id"))
		return 0, false
	}
	return id, true
}

func parseModelCatalogOptionalInt64(raw string) *int64 {
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
