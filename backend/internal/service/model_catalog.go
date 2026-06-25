package service

import (
	"context"
	"sort"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

const (
	ModelCatalogStatusActive   = "active"
	ModelCatalogStatusDisabled = "disabled"

	ModelCatalogVisibilityPublic = "public"
	ModelCatalogVisibilityAdmin  = "admin"

	ModelCatalogSourceManual  = "manual"
	ModelCatalogSourceLiteLLM = "litellm"
)

// ModelVendor represents provider/vendor metadata used by model catalog rows.
type ModelVendor struct {
	ID          int64
	Name        string
	ProviderKey string
	IconKey     string
	Description string
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ModelCatalog represents platform-wide model metadata.
type ModelCatalog struct {
	ID                int64
	ModelID           string
	NormalizedModelID string
	DisplayName       string
	Platform          string
	Provider          string
	VendorID          *int64
	Vendor            *ModelVendor
	Mode              string
	Description       string
	Tags              []string
	Capabilities      map[string]any
	Endpoints         []string
	Pricing           map[string]any
	Metadata          map[string]any
	Status            string
	Visibility        string
	Source            string
	IconKey           string
	LastSyncedAt      *time.Time
	RelatedPricing    ModelCatalogPricingAssociation
	RelatedGroups     []ModelCatalogGroupAssociation
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ModelCatalogPricingAssociation summarizes channel-pricing rows that reference a catalog model.
type ModelCatalogPricingAssociation struct {
	ChannelCount int64                             `json:"channel_count"`
	Channels     []string                          `json:"channels"`
	Entries      []ModelCatalogPricingEntry        `json:"entries"`
	Groups       []ModelCatalogPricingGroupSummary `json:"groups"`
	HasIntervals bool                              `json:"has_intervals"`
}

type ModelCatalogPricingEntry struct {
	ChannelID        int64                             `json:"channel_id"`
	ChannelName      string                            `json:"channel_name"`
	Platform         string                            `json:"platform"`
	Models           []string                          `json:"models"`
	BillingMode      string                            `json:"billing_mode"`
	InputPrice       *float64                          `json:"input_price"`
	OutputPrice      *float64                          `json:"output_price"`
	CacheWritePrice  *float64                          `json:"cache_write_price"`
	CacheReadPrice   *float64                          `json:"cache_read_price"`
	ImageOutputPrice *float64                          `json:"image_output_price"`
	PerRequestPrice  *float64                          `json:"per_request_price"`
	Intervals        []ModelCatalogPriceRange          `json:"intervals"`
	Groups           []ModelCatalogPricingGroupSummary `json:"groups"`
}

type ModelCatalogPriceRange struct {
	MinTokens       int      `json:"min_tokens"`
	MaxTokens       *int     `json:"max_tokens"`
	TierLabel       string   `json:"tier_label"`
	InputPrice      *float64 `json:"input_price"`
	OutputPrice     *float64 `json:"output_price"`
	CacheWritePrice *float64 `json:"cache_write_price"`
	CacheReadPrice  *float64 `json:"cache_read_price"`
	PerRequestPrice *float64 `json:"per_request_price"`
	SortOrder       int      `json:"sort_order"`
}

type ModelCatalogPricingGroupSummary struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Platform       string  `json:"platform"`
	RateMultiplier float64 `json:"rate_multiplier"`
}

// ModelCatalogGroupAssociation is the public-safe group metadata attached to a model.
type ModelCatalogGroupAssociation struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Platform       string  `json:"platform"`
	RateMultiplier float64 `json:"rate_multiplier"`
	IsExclusive    bool    `json:"is_exclusive"`
	Status         string  `json:"status"`
}

type ModelCatalogListFilters struct {
	Search          string
	Platform        string
	Provider        string
	VendorID        *int64
	Status          string
	Visibility      string
	PublicOnly      bool
	WithPricingOnly bool
}

type ModelVendorListFilters struct {
	Search string
}

type ModelCatalogUpsert struct {
	ModelID      string
	DisplayName  string
	Platform     string
	Provider     string
	VendorID     *int64
	Mode         string
	Description  string
	Tags         []string
	Capabilities map[string]any
	Endpoints    []string
	Pricing      map[string]any
	Metadata     map[string]any
	Status       string
	Visibility   string
	Source       string
	IconKey      string
	LastSyncedAt *time.Time
}

type ModelVendorUpsert struct {
	Name        string
	ProviderKey string
	IconKey     string
	Description string
	SortOrder   int
}

type ModelCatalogRepository interface {
	ListModels(ctx context.Context, params pagination.PaginationParams, filters ModelCatalogListFilters) ([]ModelCatalog, *pagination.PaginationResult, error)
	GetModel(ctx context.Context, id int64) (*ModelCatalog, error)
	CreateModel(ctx context.Context, input ModelCatalogUpsert) (*ModelCatalog, error)
	UpdateModel(ctx context.Context, id int64, input ModelCatalogUpsert) (*ModelCatalog, error)
	DeleteModel(ctx context.Context, id int64) error
	FindModelByPlatformAndName(ctx context.Context, platform, modelID string) (*ModelCatalog, error)
	UpsertModelFromSync(ctx context.Context, input ModelCatalogUpsert) (*ModelCatalog, bool, error)
	ListVendors(ctx context.Context, filters ModelVendorListFilters) ([]ModelVendor, error)
	GetVendor(ctx context.Context, id int64) (*ModelVendor, error)
	UpsertVendor(ctx context.Context, input ModelVendorUpsert) (*ModelVendor, error)
	DeleteVendor(ctx context.Context, id int64) error
	CountPricingAssociations(ctx context.Context, models []ModelCatalog) (map[int64]ModelCatalogPricingAssociation, error)
}

type ModelCatalogService struct {
	repo           ModelCatalogRepository
	pricingService *PricingService
	groupRepo      GroupRepository
}

func NewModelCatalogService(repo ModelCatalogRepository, pricingService *PricingService, groupRepo GroupRepository) *ModelCatalogService {
	return &ModelCatalogService{repo: repo, pricingService: pricingService, groupRepo: groupRepo}
}

func NormalizeModelCatalogID(modelID string) string {
	return strings.ToLower(strings.TrimSpace(modelID))
}

func (s *ModelCatalogService) ListModels(ctx context.Context, params pagination.PaginationParams, filters ModelCatalogListFilters) ([]ModelCatalog, *pagination.PaginationResult, error) {
	models, pag, err := s.repo.ListModels(ctx, params, normalizeModelCatalogFilters(filters))
	if err != nil {
		return nil, nil, err
	}
	if len(models) > 0 {
		assoc, err := s.repo.CountPricingAssociations(ctx, models)
		if err != nil {
			return nil, nil, err
		}
		for i := range models {
			models[i].RelatedPricing = assoc[models[i].ID]
		}
		if err := s.enrichRelatedGroups(ctx, models); err != nil {
			return nil, nil, err
		}
	}
	return models, pag, nil
}

func (s *ModelCatalogService) GetModel(ctx context.Context, id int64) (*ModelCatalog, error) {
	model, err := s.repo.GetModel(ctx, id)
	if err != nil {
		return nil, err
	}
	assoc, err := s.repo.CountPricingAssociations(ctx, []ModelCatalog{*model})
	if err != nil {
		return nil, err
	}
	model.RelatedPricing = assoc[model.ID]
	models := []ModelCatalog{*model}
	if err := s.enrichRelatedGroups(ctx, models); err != nil {
		return nil, err
	}
	model.RelatedGroups = models[0].RelatedGroups
	return model, nil
}

func (s *ModelCatalogService) CreateModel(ctx context.Context, input ModelCatalogUpsert) (*ModelCatalog, error) {
	input, err := normalizeModelCatalogInput(input, false)
	if err != nil {
		return nil, err
	}
	return s.repo.CreateModel(ctx, input)
}

func (s *ModelCatalogService) UpdateModel(ctx context.Context, id int64, input ModelCatalogUpsert) (*ModelCatalog, error) {
	input, err := normalizeModelCatalogInput(input, true)
	if err != nil {
		return nil, err
	}
	return s.repo.UpdateModel(ctx, id, input)
}

func (s *ModelCatalogService) DeleteModel(ctx context.Context, id int64) error {
	return s.repo.DeleteModel(ctx, id)
}

func (s *ModelCatalogService) ListVendors(ctx context.Context, filters ModelVendorListFilters) ([]ModelVendor, error) {
	return s.repo.ListVendors(ctx, filters)
}

func (s *ModelCatalogService) UpsertVendor(ctx context.Context, input ModelVendorUpsert) (*ModelVendor, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.ProviderKey = strings.ToLower(strings.TrimSpace(input.ProviderKey))
	input.IconKey = strings.TrimSpace(input.IconKey)
	input.Description = strings.TrimSpace(input.Description)
	if input.Name == "" {
		return nil, infraerrors.BadRequest("MODEL_VENDOR_NAME_REQUIRED", "vendor name is required")
	}
	return s.repo.UpsertVendor(ctx, input)
}

func (s *ModelCatalogService) DeleteVendor(ctx context.Context, id int64) error {
	if id <= 0 {
		return infraerrors.BadRequest("MODEL_VENDOR_ID_INVALID", "invalid vendor id")
	}
	return s.repo.DeleteVendor(ctx, id)
}

func (s *ModelCatalogService) SyncFromPricing(ctx context.Context) (ModelCatalogSyncResult, error) {
	if s.pricingService == nil {
		return ModelCatalogSyncResult{}, infraerrors.InternalServer("MODEL_PRICING_UNAVAILABLE", "pricing service is unavailable")
	}
	entries := s.pricingService.ListCatalogEntries()
	sort.Slice(entries, func(i, j int) bool { return entries[i].ModelID < entries[j].ModelID })

	result := ModelCatalogSyncResult{Total: len(entries)}
	vendorCache := map[string]*ModelVendor{}
	now := time.Now()
	for _, entry := range entries {
		vendorInput := defaultVendorForProvider(entry.Provider)
		vendor, ok := vendorCache[vendorInput.ProviderKey]
		if !ok {
			var err error
			vendor, err = s.UpsertVendor(ctx, vendorInput)
			if err != nil {
				return result, err
			}
			vendorCache[vendorInput.ProviderKey] = vendor
		}
		input := modelCatalogInputFromPricing(entry, vendor.ID, now)
		_, created, err := s.repo.UpsertModelFromSync(ctx, input)
		if err != nil {
			return result, err
		}
		if created {
			result.Created++
		} else {
			result.Updated++
		}
	}
	return result, nil
}

type ModelCatalogSyncResult struct {
	Total   int `json:"total"`
	Created int `json:"created"`
	Updated int `json:"updated"`
}

func normalizeModelCatalogFilters(filters ModelCatalogListFilters) ModelCatalogListFilters {
	filters.Search = strings.TrimSpace(filters.Search)
	filters.Platform = strings.ToLower(strings.TrimSpace(filters.Platform))
	filters.Provider = strings.ToLower(strings.TrimSpace(filters.Provider))
	filters.Status = strings.TrimSpace(filters.Status)
	filters.Visibility = strings.TrimSpace(filters.Visibility)
	return filters
}

func normalizeModelCatalogInput(input ModelCatalogUpsert, allowPartial bool) (ModelCatalogUpsert, error) {
	input.ModelID = strings.TrimSpace(input.ModelID)
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	input.Platform = strings.ToLower(strings.TrimSpace(input.Platform))
	input.Provider = strings.ToLower(strings.TrimSpace(input.Provider))
	input.Mode = strings.TrimSpace(input.Mode)
	input.Description = strings.TrimSpace(input.Description)
	input.Status = strings.TrimSpace(input.Status)
	input.Visibility = strings.TrimSpace(input.Visibility)
	input.Source = strings.TrimSpace(input.Source)
	input.IconKey = strings.TrimSpace(input.IconKey)
	if input.ModelID == "" && !allowPartial {
		return input, infraerrors.BadRequest("MODEL_ID_REQUIRED", "model id is required")
	}
	if input.ModelID == "" {
		return input, infraerrors.BadRequest("MODEL_ID_REQUIRED", "model id is required")
	}
	if input.Platform == "" {
		input.Platform = providerToPlatform(input.Provider)
	}
	if input.Platform == "" {
		return input, infraerrors.BadRequest("MODEL_PLATFORM_REQUIRED", "platform is required")
	}
	if input.Provider == "" {
		input.Provider = input.Platform
	}
	if input.DisplayName == "" {
		input.DisplayName = input.ModelID
	}
	if input.Status == "" {
		input.Status = ModelCatalogStatusActive
	}
	if input.Status != ModelCatalogStatusActive && input.Status != ModelCatalogStatusDisabled {
		return input, infraerrors.BadRequest("MODEL_STATUS_INVALID", "status must be active or disabled")
	}
	if input.Visibility == "" {
		input.Visibility = ModelCatalogVisibilityPublic
	}
	if input.Visibility != ModelCatalogVisibilityPublic && input.Visibility != ModelCatalogVisibilityAdmin {
		return input, infraerrors.BadRequest("MODEL_VISIBILITY_INVALID", "visibility must be public or admin")
	}
	if input.Source == "" {
		input.Source = ModelCatalogSourceManual
	}
	if input.Tags == nil {
		input.Tags = []string{}
	}
	if input.Capabilities == nil {
		input.Capabilities = map[string]any{}
	}
	if input.Endpoints == nil {
		input.Endpoints = []string{}
	}
	if input.Pricing == nil {
		input.Pricing = map[string]any{}
	}
	if input.Metadata == nil {
		input.Metadata = map[string]any{}
	}
	return input, nil
}

func (s *ModelCatalogService) enrichRelatedGroups(ctx context.Context, models []ModelCatalog) error {
	if len(models) == 0 {
		return nil
	}
	for i := range models {
		groups := make([]ModelCatalogGroupAssociation, 0, len(models[i].RelatedPricing.Groups))
		for _, group := range models[i].RelatedPricing.Groups {
			groups = append(groups, ModelCatalogGroupAssociation{
				ID:             group.ID,
				Name:           group.Name,
				Platform:       group.Platform,
				RateMultiplier: group.RateMultiplier,
				Status:         StatusActive,
			})
		}
		models[i].RelatedGroups = groups
	}
	_ = ctx
	return nil
}
