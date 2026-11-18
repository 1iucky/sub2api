# Research: Coupling analysis — reference catalog/marketplace code vs upstream-owned code

- **Query**: every dependency of the REF catalog implementation on other (upstream-owned) code, and how that code changed on HEAD
- **Scope**: internal
- **Date**: 2026-09-29
- **Sources read in full**: REF `backend/internal/service/model_catalog.go`, `model_catalog_pricing_sync.go`, `model_catalog_errors.go`, `group_models_list.go`, `domain/models_list_config.go`, `repository/model_catalog_repo.go`, `handler/model_catalog_handler.go`, `handler/admin/model_catalog_handler.go`, `handler/dto/model_catalog.go`, `ent/schema/model_catalog.go`, `ent/schema/model_vendor.go`

## 1. Dependency matrix (REF symbol → HEAD status)

### 1.1 Go package imports

| REF import | Used by | Status on HEAD |
|---|---|---|
| `internal/pkg/errors` (`infraerrors.BadRequest/NotFound/Conflict/InternalServer`) | service, handlers | UNCHANGED |
| `internal/pkg/pagination` (`PaginationParams{Page,PageSize,SortBy,SortOrder}`, `.Limit()`, `.NormalizedSortOrder(default)`) | service, repo, handlers | UNCHANGED (`pagination.go:7-75`) |
| `internal/pkg/response` (`ParsePagination`, `ErrorFrom`, `Success`, `Created`, `Paginated`) | handlers | UNCHANGED (standard handler kit) |
| `internal/service` (from repository) | repo | layering-consistent |
| `github.com/gin-gonic/gin`, `database/sql`, `encoding/json` | repo/handlers | — |

### 1.2 service.Package symbols the catalog code uses

| Symbol | Defined at (REF) | On HEAD | Impact |
|---|---|---|---|
| `PricingService` struct + `s.mu sync.RWMutex` + `s.pricingData map[string]*LiteLLMModelPricing` | `pricing_service.go` | EXISTS, same shape (`pricing_service.go:197-211`) | `ListCatalogEntries()` would compile if added |
| `LiteLLMModelPricing` capability fields: `MaxInputTokens`, `MaxOutputTokens`, `MaxTokens`, `SupportsAssistantPrefill`, `SupportsComputerUse`, `SupportsFunctionCalling`, `SupportsPDFInput`, `SupportsReasoning`, `SupportsResponseSchema`, `SupportsToolChoice`, `SupportsVision`, `SupportsWebSearch` | REF-added to upstream struct | **REMOVED on HEAD** (struct now ends at `CacheReadInputImageTokenCost` + `TokenPricingAbsent`) | **`model_catalog_pricing_sync.go` does NOT compile on HEAD** — `capabilitiesFromLiteLLM` and the `Metadata` map reference 12 missing fields |
| `PricingService.ListCatalogEntries()` | defined IN `model_catalog_pricing_sync.go` (feature-owned method on upstream type) | absent | self-contained once fields issue fixed |
| `NormalizeDisplayCurrency()` + `DisplayCurrencyUSD/CNY` | REF-added in upstream-owned `channel.go:44-55` | **ABSENT on HEAD** (no `currency` anywhere in `backend/internal` or migrations) | `model_catalog_repo.go` `CountPricingAssociations` fails to compile AND its SQL selects `display_currency` from `channel_model_pricing` — **column does not exist on HEAD** (REF migration `155_pricing_display_currency.sql` never ported) |
| `GroupRepository` (injected into `ModelCatalogService`) | service interface | EXISTS on HEAD | **dead dependency** — `s.groupRepo` is stored but never called anywhere in `model_catalog.go`; can be dropped from the constructor |
| `StatusActive` | `domain_constants.go:11` | UNCHANGED | used only in `enrichRelatedGroups` |
| `PlatformAnthropic/OpenAI/Gemini/Antigravity` | `domain_constants.go:41-44` | UNCHANGED | used by `providerToPlatform` / `NormalizeModelCatalogPlatform` |
| `Group.ModelsListConfig` / `GroupModelsListConfig` / `Group.CustomModelsListEnabled()` | REF-added (`domain/models_list_config.go`, `group_models_list.go`, ent group schema) | **ABSENT** — HEAD uses `model_allowlist` (`domain.GroupModelAllowlist{Enabled, Models}`, `ent/schema/group.go:265`, middleware `server/middleware/group_model_allowlist.go`) | the custom /v1/models list feature (optional) must be re-based or dropped |

### 1.3 Database tables the repository queries with raw SQL

The repo bypasses Ent and issues raw SQL (`model_catalog_repo.go`). Cross-feature table dependencies:

| Table | Query site (REF repo) | On HEAD |
|---|---|---|
| `model_catalogs`, `model_vendors` | everywhere | feature-owned; created by REF migrations 152/154 |
| `channel_model_pricing` | dedupe subquery in `ListModels` (:31-56), `WithPricingOnly` filter (:423-438), `CountPricingAssociations` (:232-300) | EXISTS on HEAD but **schema drifted**: HEAD columns (`channel_repo_pricing.go:19`) = `id, channel_id, platform, models, billing_mode, input_price, output_price, cache_write_price, cache_write_1h_price, cache_read_price, fast_multiplier, flex_multiplier, max_reasoning_effort_multiplier, image_input_price, image_output_price, per_request_price, time_pricing, ...` — NO `display_currency` (REF selects it → runtime SQL error); HEAD adds `cache_write_1h_price`, multipliers, `image_input_price`, `time_pricing` that REF didn't know |
| `channels` (`c.status='active'`, `c.name`) | joins | EXISTS (raw-SQL managed on HEAD; no Ent schema) |
| `channel_pricing_intervals` (`pricing_id, min_tokens, max_tokens, tier_label, input/output/cache_write/cache_read/per_request_price, sort_order`) | `listCatalogPricingIntervals` (:303-345) | EXISTS; HEAD adds extra columns (`cache_write_1h_price`, `input_multiplier`, `output_multiplier`, `cache_write_multiplier`, `cache_read_multiplier` — `channel_repo_pricing.go:296`) — REF's SELECT subset still works |
| `channel_groups` + `groups` (`g.status='active'`, `g.is_exclusive=false`, `g.sort_order`, `g.rate_multiplier`) | `listCatalogPricingGroups` (:348-376) | EXISTS on HEAD (migration `081_create_channels.sql`); columns unchanged |

### 1.4 Repository helpers from the shared repository package

| Helper | Defined | HEAD |
|---|---|---|
| `isUniqueViolation(err)` | `channel_repo_pricing.go:319` | UNCHANGED |
| `escapeLike(s)` | `channel_repo_pricing.go:328` | UNCHANGED |
| `*sql.DB` injection (`NewModelCatalogRepository(db)`) | wire | UNCHANGED pattern |

### 1.5 Frontend runtime dependencies of the two views

| Dependency | REF | HEAD |
|---|---|---|
| `@/api/client` (`apiClient`) | `api/models.ts`, `api/admin/models.ts` | UNCHANGED |
| `@/utils/pricing` `currencySymbol`, `DisplayCurrency` | `ModelMarketplaceView.vue:449`, `api/models.ts:2` | **HEAD `utils/pricing.ts` has a different, incompatible `formatScaled` and no `DisplayCurrency`** (see collision-surface.md) |
| `@/api/publicChannelMonitor` → `/public/channel-monitors` | marketplace monitor badges | file + public route absent on HEAD; `ChannelMonitorUserHandler.List/GetStatus` DO exist on HEAD (`routes/user.go:141-142` registers them under auth) |
| `@/api/channelMonitor` types (`UserMonitorView`) | marketplace | EXISTS on HEAD (`api/channelMonitor.ts:24`) |
| `AppLayout`, `Icon`, `ModelIcon`, `Select`, `Pagination`, `usePersistedPageSize`, `extractApiErrorMessage`, `useAppStore/useAuthStore` | both views | ALL EXIST on HEAD |
| `PublicTopNav`, `PublicFooter`, `useTheme`, `useBrandHover` | public (non-embedded) marketplace chrome | `PublicTopNav`/`PublicFooter` absent; `useBrandHover` exists on HEAD; `useTheme` absent in HEAD but present as untracked working-tree file |
| `modelMarketplaceMonitor.ts` helpers | marketplace | feature-owned, portable as-is |

## 2. Self-contained vs coupled parts

### Self-contained (port with zero upstream coupling)

- Catalog/vendor CRUD: ent schemas, migrations 152/154 (renumbered), `modelCatalogRepository` CRUD methods, `ModelCatalogService` CRUD + normalization + validation, both handlers, dto, all frontend views/api/i18n.
- Vendor/platform/icon mapping tables (`vendorNameForProvider`, `iconKeyForProvider`, `providerToPlatform`, `NormalizeModelCatalogPlatform`, `vendorSortOrder`, `endpointsForMode`, `modelTagsFromPricing`) — pure functions, depend only on `Platform*` consts.
- `enrichRelatedGroups` — only reshapes pricing-group summaries; `groupRepo` injection is dead and can be removed.

### Coupled (needs decision)

1. **Pricing-association read path** (`CountPricingAssociations`, `WithPricingOnly`, dedupe ranking) — reads `channel_model_pricing` (+`display_currency`), `channels`, `channel_pricing_intervals`, `channel_groups`, `groups`. Compile break on `NormalizeDisplayCurrency`; runtime break on `display_currency` column. HEAD table otherwise compatible.
2. **LiteLLM pricing sync** (`SyncFromPricing` + `ListCatalogEntries` + `capabilitiesFromLiteLLM`) — reads upstream `PricingService.pricingData` and 12 `LiteLLMModelPricing` fields deleted on HEAD. Startup warm goroutine in `ProvideModelCatalogService` (`service/wire.go:35-51`) also triggers it.
3. **Group custom models list** (`models_list_config`) — an entire sub-feature entangled with upstream-owned group entity/repos/handlers and `gateway_handler.go` /v1/models flow; HEAD replaced the concept space with `group_model_allowlist` + `middleware/group_model_allowlist.go` + `pkg/requestmodel`.

## 3. Ent schema details (merge-conflict risk assessment)

### `ModelCatalog` (REF `ent/schema/model_catalog.go`)

- Table: `model_catalogs` (`entsql.Annotation`)
- Mixin: `mixins.TimeMixin{}` (created_at/updated_at)
- Fields: `model_id`(str,200,NotEmpty), `normalized_model_id`(str,200,NotEmpty), `display_name`(str,200,default ''), `platform`(str,50), `provider`(str,100), `vendor_id`(int64,Optional,Nillable), `mode`(str,50), `description`(str), JSON `tags []string`, `capabilities map`, `endpoints []string`, `pricing map`, `metadata map`, Enum `status`(active|disabled), Enum `visibility`(public|admin), `source`(str,50,default manual), `icon_key`(str,80), `last_synced_at`(time,Optional,Nillable)
- Edges: **only** `edge.To("vendor", ModelVendor.Type).Field("vendor_id").Unique()` with `OnDelete(SetNull)` — **no edges to upstream-owned schemas → no ent merge risk**
- Indexes: unique `(platform, normalized_model_id)`; plus `platform`, `provider`, `vendor_id`, `(status, visibility)`, `source`

### `ModelVendor` (REF `ent/schema/model_vendor.go`)

- Table: `model_vendors`
- Mixin: `TimeMixin`
- Fields: `name`(str,100,NotEmpty,**Unique**), `provider_key`(str,80), `icon_key`(str,80), `description`(str), `sort_order`(int), `deleted_at`(timestamptz,Optional,Nillable)
- Edges: `edge.From("models", ModelCatalog.Type).Ref("vendor")` — internal only
- Indexes: `provider_key`, `deleted_at`

**Note**: the only upstream-owned ent schema the REF feature touches is `group.go` (`models_list_config` field at :215) — that belongs to the optional group-models-list sub-feature, not the catalog proper.

## 4. Unrelated REF↔HEAD drift encountered (NOT catalog coupling)

- `backend/internal/service/gateway_model_availability.go`: REF rewrote `DiagnoseModelAvailabilityForPlatform` to use `accountRepo.ListModelAvailabilityCandidates` (persistent eligibility); HEAD uses `s.listSchedulableAccounts`. No catalog code references this file.
- HEAD-only upstream additions absent in REF: `Plugin` handler, `CNProvider` handler/api, `ChannelMonitorV2`, `Passkey`, `ModelPlaza` (`handler/model_plaza_handler.go`, `service/model_plaza_service.go`), `composite_model_route`, `pkg/requestmodel`, locales `channelMonitorV2.ts`/`batchImage.ts`/`plugins.ts`.

## Caveats / Not Found

- REF catalog has no caching layer and no settings keys; nothing in `domain_constants.go` settings section is catalog-specific.
- `ProvideModelCatalogService`'s startup warm sync runs with a 2-minute timeout and only logs failures.
