package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type modelCatalogRepository struct {
	db *sql.DB
}

func NewModelCatalogRepository(db *sql.DB) service.ModelCatalogRepository {
	return &modelCatalogRepository{db: db}
}

func (r *modelCatalogRepository) ListModels(ctx context.Context, params pagination.PaginationParams, filters service.ModelCatalogListFilters) ([]service.ModelCatalog, *pagination.PaginationResult, error) {
	where, args := buildModelCatalogWhere(filters)
	argIdx := len(args) + 1
	whereClause := strings.Join(where, " AND ")
	fromClause := `model_catalogs mc`
	if filters.DeduplicateByID {
		fromClause = fmt.Sprintf(`(
			SELECT *
			FROM (
				SELECT mc.*,
				       ROW_NUMBER() OVER (
				         PARTITION BY mc.normalized_model_id
				         ORDER BY
				           (
				             SELECT COUNT(DISTINCT c.id)
				             FROM channel_model_pricing cmp
				             JOIN channels c ON c.id = cmp.channel_id
				             WHERE c.status = 'active'
				               AND EXISTS (
				                 SELECT 1
				                 FROM jsonb_array_elements_text(cmp.models::jsonb) AS m(name)
				                 WHERE LOWER(TRIM(m.name)) = mc.normalized_model_id
				               )
				           ) DESC,
				           CASE mc.source WHEN 'manual' THEN 0 ELSE 1 END ASC,
				           mc.updated_at DESC,
				           mc.id DESC
				       ) AS model_catalog_rank
				FROM model_catalogs mc
				LEFT JOIN model_vendors mv ON mv.id = mc.vendor_id AND mv.deleted_at IS NULL
				WHERE %s
			) ranked_model_catalogs
			WHERE model_catalog_rank = 1
		) mc`, whereClause)
		whereClause = "1=1"
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM %s LEFT JOIN model_vendors mv ON mv.id = mc.vendor_id AND mv.deleted_at IS NULL WHERE %s`, fromClause, whereClause), args...).Scan(&total); err != nil {
		return nil, nil, fmt.Errorf("count model catalog: %w", err)
	}

	pageSize := params.Limit()
	page := params.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	orderBy := modelCatalogOrderBy(params)
	query := fmt.Sprintf(`
		SELECT mc.id, mc.model_id, mc.normalized_model_id, mc.display_name, mc.platform, mc.provider, mc.vendor_id,
		       mc.mode, mc.description, mc.tags, mc.capabilities, mc.endpoints, mc.pricing, mc.metadata,
		       mc.status, mc.visibility, mc.source, mc.icon_key, mc.last_synced_at, mc.created_at, mc.updated_at,
		       mv.id, mv.name, mv.provider_key, mv.icon_key, mv.description, mv.sort_order, mv.created_at, mv.updated_at
		FROM %s
		LEFT JOIN model_vendors mv ON mv.id = mc.vendor_id AND mv.deleted_at IS NULL
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`, fromClause, whereClause, orderBy, argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("query model catalog: %w", err)
	}
	defer func() { _ = rows.Close() }()

	models, err := scanModelCatalogRows(rows)
	if err != nil {
		return nil, nil, err
	}
	pages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if pages < 1 {
		pages = 1
	}
	return models, &pagination.PaginationResult{Total: total, Page: page, PageSize: pageSize, Pages: pages}, nil
}

func (r *modelCatalogRepository) GetModel(ctx context.Context, id int64) (*service.ModelCatalog, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT mc.id, mc.model_id, mc.normalized_model_id, mc.display_name, mc.platform, mc.provider, mc.vendor_id,
		       mc.mode, mc.description, mc.tags, mc.capabilities, mc.endpoints, mc.pricing, mc.metadata,
		       mc.status, mc.visibility, mc.source, mc.icon_key, mc.last_synced_at, mc.created_at, mc.updated_at,
		       mv.id, mv.name, mv.provider_key, mv.icon_key, mv.description, mv.sort_order, mv.created_at, mv.updated_at
		FROM model_catalogs mc
		LEFT JOIN model_vendors mv ON mv.id = mc.vendor_id AND mv.deleted_at IS NULL
		WHERE mc.id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("get model catalog: %w", err)
	}
	defer func() { _ = rows.Close() }()
	models, err := scanModelCatalogRows(rows)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, service.ErrModelCatalogNotFound
	}
	return &models[0], nil
}

func (r *modelCatalogRepository) FindModelByPlatformAndName(ctx context.Context, platform, modelID string) (*service.ModelCatalog, error) {
	filters := service.ModelCatalogListFilters{
		Platform: platform,
		Search:   service.NormalizeModelCatalogID(modelID),
	}
	models, _, err := r.ListModels(ctx, pagination.PaginationParams{Page: 1, PageSize: 1}, filters)
	if err != nil {
		return nil, err
	}
	for i := range models {
		if models[i].Platform == strings.ToLower(strings.TrimSpace(platform)) &&
			models[i].NormalizedModelID == service.NormalizeModelCatalogID(modelID) {
			return &models[i], nil
		}
	}
	return nil, service.ErrModelCatalogNotFound
}

func (r *modelCatalogRepository) CreateModel(ctx context.Context, input service.ModelCatalogUpsert) (*service.ModelCatalog, error) {
	id, err := r.insertModel(ctx, input)
	if err != nil {
		return nil, err
	}
	return r.GetModel(ctx, id)
}

func (r *modelCatalogRepository) UpdateModel(ctx context.Context, id int64, input service.ModelCatalogUpsert) (*service.ModelCatalog, error) {
	if err := r.updateModel(ctx, id, input); err != nil {
		return nil, err
	}
	return r.GetModel(ctx, id)
}

func (r *modelCatalogRepository) DeleteModel(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM model_catalogs WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete model catalog: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return service.ErrModelCatalogNotFound
	}
	return nil
}

func (r *modelCatalogRepository) DeleteVendor(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE model_vendors SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return fmt.Errorf("delete model vendor: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return service.ErrModelVendorNotFound
	}
	return nil
}

func (r *modelCatalogRepository) UpsertModelFromSync(ctx context.Context, input service.ModelCatalogUpsert) (*service.ModelCatalog, bool, error) {
	existing, err := r.FindModelByPlatformAndName(ctx, input.Platform, input.ModelID)
	if err == nil && existing != nil {
		updated, err := r.UpdateModel(ctx, existing.ID, input)
		return updated, false, err
	}
	if err != nil && err != service.ErrModelCatalogNotFound {
		return nil, false, err
	}
	created, err := r.CreateModel(ctx, input)
	return created, true, err
}

func (r *modelCatalogRepository) ListVendors(ctx context.Context, filters service.ModelVendorListFilters) ([]service.ModelVendor, error) {
	where := []string{"1=1"}
	args := []any{}
	if strings.TrimSpace(filters.Search) != "" {
		args = append(args, "%"+escapeLike(strings.TrimSpace(filters.Search))+"%")
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR provider_key ILIKE $%d)", len(args), len(args)))
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT id, name, provider_key, icon_key, description, sort_order, deleted_at, created_at, updated_at
		FROM model_vendors
		WHERE deleted_at IS NULL AND %s
		ORDER BY sort_order ASC, name ASC`, strings.Join(where, " AND ")), args...)
	if err != nil {
		return nil, fmt.Errorf("list model vendors: %w", err)
	}
	defer func() { _ = rows.Close() }()
	return scanModelVendorRows(rows)
}

func (r *modelCatalogRepository) GetVendor(ctx context.Context, id int64) (*service.ModelVendor, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, provider_key, icon_key, description, sort_order, deleted_at, created_at, updated_at
		FROM model_vendors WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return nil, fmt.Errorf("get model vendor: %w", err)
	}
	defer func() { _ = rows.Close() }()
	vendors, err := scanModelVendorRows(rows)
	if err != nil {
		return nil, err
	}
	if len(vendors) == 0 {
		return nil, service.ErrModelVendorNotFound
	}
	return &vendors[0], nil
}

func (r *modelCatalogRepository) UpsertVendor(ctx context.Context, input service.ModelVendorUpsert) (*service.ModelVendor, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO model_vendors (name, provider_key, icon_key, description, sort_order)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (name) DO UPDATE
		SET provider_key = EXCLUDED.provider_key,
		    icon_key = EXCLUDED.icon_key,
		    description = EXCLUDED.description,
		    sort_order = EXCLUDED.sort_order,
		    deleted_at = NULL,
		    updated_at = NOW()
		RETURNING id`, input.Name, input.ProviderKey, input.IconKey, input.Description, input.SortOrder).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("upsert model vendor: %w", err)
	}
	return r.GetVendor(ctx, id)
}

func (r *modelCatalogRepository) FindDeletedVendorByProviderKey(ctx context.Context, providerKey string) (*service.ModelVendor, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, provider_key, icon_key, description, sort_order, deleted_at, created_at, updated_at
		FROM model_vendors
		WHERE deleted_at IS NOT NULL AND provider_key = $1
		ORDER BY updated_at DESC
		LIMIT 1`, strings.ToLower(strings.TrimSpace(providerKey)))
	if err != nil {
		return nil, fmt.Errorf("find deleted model vendor: %w", err)
	}
	defer func() { _ = rows.Close() }()
	vendors, err := scanModelVendorRows(rows)
	if err != nil {
		return nil, err
	}
	if len(vendors) == 0 {
		return nil, service.ErrModelVendorNotFound
	}
	return &vendors[0], nil
}

func (r *modelCatalogRepository) CountPricingAssociations(ctx context.Context, models []service.ModelCatalog) (map[int64]service.ModelCatalogPricingAssociation, error) {
	result := make(map[int64]service.ModelCatalogPricingAssociation, len(models))
	for _, model := range models {
		rows, err := r.db.QueryContext(ctx, `
			SELECT cmp.id, cmp.channel_id, c.name, cmp.platform, cmp.models, cmp.billing_mode,
			       cmp.input_price, cmp.output_price, cmp.cache_write_price, cmp.cache_read_price,
			       cmp.image_output_price, cmp.per_request_price
			FROM channel_model_pricing cmp
			JOIN channels c ON c.id = cmp.channel_id
			WHERE c.status = 'active'
			  AND EXISTS (
			    SELECT 1
			    FROM jsonb_array_elements_text(cmp.models::jsonb) AS m(name)
			    WHERE LOWER(TRIM(m.name)) = $1
			  )
			ORDER BY c.name, cmp.id`, model.NormalizedModelID)
		if err != nil {
			return nil, fmt.Errorf("count pricing associations: %w", err)
		}
		assoc := service.ModelCatalogPricingAssociation{}
		channelSeen := map[string]struct{}{}
		groupSeen := map[int64]struct{}{}
		for rows.Next() {
			var entry service.ModelCatalogPricingEntry
			var pricingID int64
			var modelsJSON []byte
			var inputPrice, outputPrice, cacheWritePrice, cacheReadPrice, imageOutputPrice, perRequestPrice sql.NullFloat64
			if err := rows.Scan(
				&pricingID, &entry.ChannelID, &entry.ChannelName, &entry.Platform, &modelsJSON, &entry.BillingMode,
				&inputPrice, &outputPrice, &cacheWritePrice, &cacheReadPrice, &imageOutputPrice, &perRequestPrice,
			); err != nil {
				_ = rows.Close()
				return nil, fmt.Errorf("scan pricing association: %w", err)
			}
			_ = json.Unmarshal(modelsJSON, &entry.Models)
			entry.InputPrice = nullableFloat(inputPrice)
			entry.OutputPrice = nullableFloat(outputPrice)
			entry.CacheWritePrice = nullableFloat(cacheWritePrice)
			entry.CacheReadPrice = nullableFloat(cacheReadPrice)
			entry.ImageOutputPrice = nullableFloat(imageOutputPrice)
			entry.PerRequestPrice = nullableFloat(perRequestPrice)
			intervals, err := r.listCatalogPricingIntervals(ctx, pricingID)
			if err != nil {
				_ = rows.Close()
				return nil, err
			}
			entry.Intervals = intervals
			if len(intervals) > 0 {
				assoc.HasIntervals = true
			}
			groups, err := r.listCatalogPricingGroups(ctx, entry.ChannelID)
			if err != nil {
				_ = rows.Close()
				return nil, err
			}
			entry.Groups = groups
			for _, group := range groups {
				if _, ok := groupSeen[group.ID]; ok {
					continue
				}
				groupSeen[group.ID] = struct{}{}
				assoc.Groups = append(assoc.Groups, group)
			}
			if _, ok := channelSeen[entry.ChannelName]; !ok {
				channelSeen[entry.ChannelName] = struct{}{}
				assoc.Channels = append(assoc.Channels, entry.ChannelName)
			}
			assoc.Entries = append(assoc.Entries, entry)
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("iterate pricing associations: %w", err)
		}
		_ = rows.Close()
		assoc.ChannelCount = int64(len(assoc.Channels))
		result[model.ID] = assoc
	}
	return result, nil
}

func (r *modelCatalogRepository) listCatalogPricingIntervals(ctx context.Context, pricingID int64) ([]service.ModelCatalogPriceRange, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT min_tokens, max_tokens, COALESCE(tier_label, ''), input_price, output_price,
		       cache_write_price, cache_read_price, per_request_price, sort_order
		FROM channel_pricing_intervals
		WHERE pricing_id = $1
		ORDER BY sort_order ASC, min_tokens ASC`, pricingID)
	if err != nil {
		return nil, fmt.Errorf("list pricing intervals: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := []service.ModelCatalogPriceRange{}
	for rows.Next() {
		var item service.ModelCatalogPriceRange
		var maxTokens sql.NullInt64
		var inputPrice, outputPrice, cacheWritePrice, cacheReadPrice, perRequestPrice sql.NullFloat64
		if err := rows.Scan(
			&item.MinTokens, &maxTokens, &item.TierLabel, &inputPrice, &outputPrice,
			&cacheWritePrice, &cacheReadPrice, &perRequestPrice, &item.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("scan pricing interval: %w", err)
		}
		if maxTokens.Valid {
			max := int(maxTokens.Int64)
			item.MaxTokens = &max
		}
		item.InputPrice = nullableFloat(inputPrice)
		item.OutputPrice = nullableFloat(outputPrice)
		item.CacheWritePrice = nullableFloat(cacheWritePrice)
		item.CacheReadPrice = nullableFloat(cacheReadPrice)
		item.PerRequestPrice = nullableFloat(perRequestPrice)
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pricing intervals: %w", err)
	}
	return out, nil
}

func (r *modelCatalogRepository) listCatalogPricingGroups(ctx context.Context, channelID int64) ([]service.ModelCatalogPricingGroupSummary, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT g.id, g.name, g.platform, g.rate_multiplier
		FROM channel_groups cg
		JOIN groups g ON g.id = cg.group_id
		WHERE cg.channel_id = $1
		  AND g.status = 'active'
		  AND g.is_exclusive = FALSE
		ORDER BY g.sort_order ASC, g.name ASC`, channelID)
	if err != nil {
		return nil, fmt.Errorf("list pricing groups: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := []service.ModelCatalogPricingGroupSummary{}
	for rows.Next() {
		var group service.ModelCatalogPricingGroupSummary
		if err := rows.Scan(&group.ID, &group.Name, &group.Platform, &group.RateMultiplier); err != nil {
			return nil, fmt.Errorf("scan pricing group: %w", err)
		}
		out = append(out, group)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pricing groups: %w", err)
	}
	return out, nil
}

func nullableFloat(value sql.NullFloat64) *float64 {
	if !value.Valid {
		return nil
	}
	return &value.Float64
}

func (r *modelCatalogRepository) insertModel(ctx context.Context, input service.ModelCatalogUpsert) (int64, error) {
	tags, capabilities, endpoints, pricing, metadata, err := marshalModelCatalogJSON(input)
	if err != nil {
		return 0, err
	}
	var id int64
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO model_catalogs
		(model_id, normalized_model_id, display_name, platform, provider, vendor_id, mode, description,
		 tags, capabilities, endpoints, pricing, metadata, status, visibility, source, icon_key, last_synced_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb, $10::jsonb, $11::jsonb, $12::jsonb, $13::jsonb, $14, $15, $16, $17, $18)
		RETURNING id`,
		input.ModelID, service.NormalizeModelCatalogID(input.ModelID), input.DisplayName, input.Platform, input.Provider, input.VendorID,
		input.Mode, input.Description, tags, capabilities, endpoints, pricing, metadata, input.Status, input.Visibility, input.Source, input.IconKey, input.LastSyncedAt,
	).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return 0, service.ErrModelCatalogExists
		}
		return 0, fmt.Errorf("insert model catalog: %w", err)
	}
	return id, nil
}

func (r *modelCatalogRepository) updateModel(ctx context.Context, id int64, input service.ModelCatalogUpsert) error {
	tags, capabilities, endpoints, pricing, metadata, err := marshalModelCatalogJSON(input)
	if err != nil {
		return err
	}
	result, err := r.db.ExecContext(ctx, `
		UPDATE model_catalogs
		SET model_id = $1,
		    normalized_model_id = $2,
		    display_name = $3,
		    platform = $4,
		    provider = $5,
		    vendor_id = $6,
		    mode = $7,
		    description = $8,
		    tags = $9::jsonb,
		    capabilities = $10::jsonb,
		    endpoints = $11::jsonb,
		    pricing = $12::jsonb,
		    metadata = $13::jsonb,
		    status = $14,
		    visibility = $15,
		    source = $16,
		    icon_key = $17,
		    last_synced_at = $18,
		    updated_at = NOW()
		WHERE id = $19`,
		input.ModelID, service.NormalizeModelCatalogID(input.ModelID), input.DisplayName, input.Platform, input.Provider, input.VendorID,
		input.Mode, input.Description, tags, capabilities, endpoints, pricing, metadata, input.Status, input.Visibility, input.Source, input.IconKey, input.LastSyncedAt, id,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return service.ErrModelCatalogExists
		}
		return fmt.Errorf("update model catalog: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return service.ErrModelCatalogNotFound
	}
	return nil
}

func buildModelCatalogWhere(filters service.ModelCatalogListFilters) ([]string, []any) {
	where := []string{"1=1"}
	args := []any{}
	if filters.PublicOnly {
		where = append(where, "mc.status = 'active'", "mc.visibility = 'public'")
	}
	if filters.WithPricingOnly {
		where = append(where, `EXISTS (
			SELECT 1
			FROM channel_model_pricing cmp
			JOIN channels c ON c.id = cmp.channel_id
			WHERE c.status = 'active'
			  AND EXISTS (
			    SELECT 1
			    FROM jsonb_array_elements_text(cmp.models::jsonb) AS m(name)
			    WHERE LOWER(TRIM(m.name)) = mc.normalized_model_id
			  )
		)`)
	}
	if filters.Search != "" {
		args = append(args, "%"+escapeLike(filters.Search)+"%")
		where = append(where, fmt.Sprintf("(mc.model_id ILIKE $%d OR mc.display_name ILIKE $%d OR mc.provider ILIKE $%d OR mv.name ILIKE $%d)", len(args), len(args), len(args), len(args)))
	}
	if filters.Platform != "" {
		args = append(args, filters.Platform)
		where = append(where, fmt.Sprintf("LOWER(mc.platform) = $%d", len(args)))
	}
	if filters.Provider != "" {
		args = append(args, filters.Provider)
		where = append(where, fmt.Sprintf("LOWER(mc.provider) = $%d", len(args)))
	}
	if filters.VendorID != nil {
		args = append(args, *filters.VendorID)
		where = append(where, fmt.Sprintf("mc.vendor_id = $%d", len(args)))
	}
	if filters.Status != "" {
		args = append(args, filters.Status)
		where = append(where, fmt.Sprintf("mc.status = $%d", len(args)))
	}
	if filters.Visibility != "" {
		args = append(args, filters.Visibility)
		where = append(where, fmt.Sprintf("mc.visibility = $%d", len(args)))
	}
	return where, args
}

func modelCatalogOrderBy(params pagination.PaginationParams) string {
	order := params.NormalizedSortOrder(pagination.SortOrderAsc)
	switch params.SortBy {
	case "model_id":
		return "mc.model_id " + order
	case "platform":
		return "mc.platform " + order + ", mc.model_id ASC"
	case "provider":
		return "mc.provider " + order + ", mc.model_id ASC"
	case "updated_at":
		return "mc.updated_at " + order
	case "created_at":
		return "mc.created_at " + order
	default:
		return "mc.platform ASC, mv.sort_order ASC NULLS LAST, mc.model_id ASC"
	}
}

func marshalModelCatalogJSON(input service.ModelCatalogUpsert) (tags, capabilities, endpoints, pricing, metadata []byte, err error) {
	if tags, err = json.Marshal(input.Tags); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("marshal tags: %w", err)
	}
	if capabilities, err = json.Marshal(input.Capabilities); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("marshal capabilities: %w", err)
	}
	if endpoints, err = json.Marshal(input.Endpoints); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("marshal endpoints: %w", err)
	}
	if pricing, err = json.Marshal(input.Pricing); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("marshal pricing: %w", err)
	}
	if metadata, err = json.Marshal(input.Metadata); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("marshal metadata: %w", err)
	}
	return tags, capabilities, endpoints, pricing, metadata, nil
}

func scanModelCatalogRows(rows *sql.Rows) ([]service.ModelCatalog, error) {
	models := []service.ModelCatalog{}
	for rows.Next() {
		var model service.ModelCatalog
		var tagsJSON, capabilitiesJSON, endpointsJSON, pricingJSON, metadataJSON []byte
		var vendorID sql.NullInt64
		var lastSyncedAt sql.NullTime
		var vendor modelVendorNullable
		if err := rows.Scan(
			&model.ID, &model.ModelID, &model.NormalizedModelID, &model.DisplayName, &model.Platform, &model.Provider, &vendorID,
			&model.Mode, &model.Description, &tagsJSON, &capabilitiesJSON, &endpointsJSON, &pricingJSON, &metadataJSON,
			&model.Status, &model.Visibility, &model.Source, &model.IconKey, &lastSyncedAt, &model.CreatedAt, &model.UpdatedAt,
			&vendor.ID, &vendor.Name, &vendor.ProviderKey, &vendor.IconKey, &vendor.Description, &vendor.SortOrder, &vendor.CreatedAt, &vendor.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan model catalog: %w", err)
		}
		if vendorID.Valid {
			model.VendorID = &vendorID.Int64
		}
		if lastSyncedAt.Valid {
			model.LastSyncedAt = &lastSyncedAt.Time
		}
		_ = json.Unmarshal(tagsJSON, &model.Tags)
		_ = json.Unmarshal(capabilitiesJSON, &model.Capabilities)
		_ = json.Unmarshal(endpointsJSON, &model.Endpoints)
		_ = json.Unmarshal(pricingJSON, &model.Pricing)
		_ = json.Unmarshal(metadataJSON, &model.Metadata)
		if vendor.ID.Valid {
			model.Vendor = &service.ModelVendor{
				ID:          vendor.ID.Int64,
				Name:        vendor.Name.String,
				ProviderKey: vendor.ProviderKey.String,
				IconKey:     vendor.IconKey.String,
				Description: vendor.Description.String,
				SortOrder:   int(vendor.SortOrder.Int64),
				CreatedAt:   vendor.CreatedAt.Time,
				UpdatedAt:   vendor.UpdatedAt.Time,
			}
		}
		if model.Tags == nil {
			model.Tags = []string{}
		}
		if model.Capabilities == nil {
			model.Capabilities = map[string]any{}
		}
		if model.Endpoints == nil {
			model.Endpoints = []string{}
		}
		if model.Pricing == nil {
			model.Pricing = map[string]any{}
		}
		if model.Metadata == nil {
			model.Metadata = map[string]any{}
		}
		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate model catalog: %w", err)
	}
	return models, nil
}

type modelVendorNullable struct {
	ID          sql.NullInt64
	Name        sql.NullString
	ProviderKey sql.NullString
	IconKey     sql.NullString
	Description sql.NullString
	SortOrder   sql.NullInt64
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}

func scanModelVendorRows(rows *sql.Rows) ([]service.ModelVendor, error) {
	vendors := []service.ModelVendor{}
	for rows.Next() {
		var v service.ModelVendor
		var deletedAt sql.NullTime
		if err := rows.Scan(&v.ID, &v.Name, &v.ProviderKey, &v.IconKey, &v.Description, &v.SortOrder, &deletedAt, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan model vendor: %w", err)
		}
		if deletedAt.Valid {
			v.DeletedAt = &deletedAt.Time
		}
		vendors = append(vendors, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate model vendors: %w", err)
	}
	return vendors, nil
}
