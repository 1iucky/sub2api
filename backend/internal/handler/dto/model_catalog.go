package dto

import "github.com/Wei-Shaw/sub2api/internal/service"

type ModelVendorResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ProviderKey string `json:"provider_key"`
	IconKey     string `json:"icon_key"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ModelCatalogResponse struct {
	ID             int64                                  `json:"id"`
	ModelID        string                                 `json:"model_id"`
	DisplayName    string                                 `json:"display_name"`
	Platform       string                                 `json:"platform"`
	Provider       string                                 `json:"provider"`
	VendorID       *int64                                 `json:"vendor_id"`
	Vendor         *ModelVendorResponse                   `json:"vendor"`
	Mode           string                                 `json:"mode"`
	Description    string                                 `json:"description"`
	Tags           []string                               `json:"tags"`
	Capabilities   map[string]any                         `json:"capabilities"`
	Endpoints      []string                               `json:"endpoints"`
	Pricing        map[string]any                         `json:"pricing"`
	Metadata       map[string]any                         `json:"metadata"`
	Status         string                                 `json:"status"`
	Visibility     string                                 `json:"visibility"`
	Source         string                                 `json:"source"`
	IconKey        string                                 `json:"icon_key"`
	LastSyncedAt   *string                                `json:"last_synced_at"`
	RelatedPricing service.ModelCatalogPricingAssociation `json:"related_pricing"`
	RelatedGroups  []service.ModelCatalogGroupAssociation `json:"related_groups"`
	CreatedAt      string                                 `json:"created_at"`
	UpdatedAt      string                                 `json:"updated_at"`
}

func ModelVendorsToResponse(vendors []service.ModelVendor) []ModelVendorResponse {
	out := make([]ModelVendorResponse, 0, len(vendors))
	for i := range vendors {
		out = append(out, ModelVendorToResponse(&vendors[i]))
	}
	return out
}

func ModelVendorToResponse(v *service.ModelVendor) ModelVendorResponse {
	if v == nil {
		return ModelVendorResponse{}
	}
	return ModelVendorResponse{
		ID:          v.ID,
		Name:        v.Name,
		ProviderKey: v.ProviderKey,
		IconKey:     v.IconKey,
		Description: v.Description,
		SortOrder:   v.SortOrder,
		CreatedAt:   v.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   v.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ModelCatalogsToResponse(models []service.ModelCatalog) []ModelCatalogResponse {
	out := make([]ModelCatalogResponse, 0, len(models))
	for i := range models {
		out = append(out, ModelCatalogToResponse(&models[i]))
	}
	return out
}

func ModelCatalogToResponse(m *service.ModelCatalog) ModelCatalogResponse {
	if m == nil {
		return ModelCatalogResponse{}
	}
	var synced *string
	if m.LastSyncedAt != nil {
		s := m.LastSyncedAt.Format("2006-01-02T15:04:05Z")
		synced = &s
	}
	var vendor *ModelVendorResponse
	if m.Vendor != nil {
		v := ModelVendorToResponse(m.Vendor)
		vendor = &v
	}
	return ModelCatalogResponse{
		ID:             m.ID,
		ModelID:        m.ModelID,
		DisplayName:    m.DisplayName,
		Platform:       m.Platform,
		Provider:       m.Provider,
		VendorID:       m.VendorID,
		Vendor:         vendor,
		Mode:           m.Mode,
		Description:    m.Description,
		Tags:           m.Tags,
		Capabilities:   m.Capabilities,
		Endpoints:      m.Endpoints,
		Pricing:        m.Pricing,
		Metadata:       m.Metadata,
		Status:         m.Status,
		Visibility:     m.Visibility,
		Source:         m.Source,
		IconKey:        m.IconKey,
		LastSyncedAt:   synced,
		RelatedPricing: m.RelatedPricing,
		RelatedGroups:  m.RelatedGroups,
		CreatedAt:      m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
