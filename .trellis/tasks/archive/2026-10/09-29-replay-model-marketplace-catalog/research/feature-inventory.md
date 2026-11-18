# Research: Feature inventory — Model Marketplace + Model Catalog (reference branch)

- **Query**: complete file inventory of the marketplace/catalog features in `codex/custom-theme-upstream-main-20260717`
- **Scope**: internal
- **Date**: 2026-09-29
- **Reference (REF)**: `codex/custom-theme-upstream-main-20260717`
- **HEAD**: `codex/upstream-main-personalized-replay-20260917` (2bf02200c)

Legend: `in HEAD = NO` means `git cat-file -e HEAD:<path>` fails (file does not exist on current branch).

## Backend — feature-owned files (all absent from HEAD)

| File (REF) | Lines | in HEAD | Description |
|---|---|---|---|
| `backend/ent/schema/model_catalog.go` | 103 | NO | Ent schema `ModelCatalog`, table `model_catalogs` |
| `backend/ent/schema/model_vendor.go` | 66 | NO | Ent schema `ModelVendor`, table `model_vendors` |
| `backend/internal/service/model_catalog.go` | 420 | NO | `ModelCatalogService`, `ModelCatalogRepository` interface, domain structs, normalization, `SyncFromPricing` |
| `backend/internal/service/model_catalog_errors.go` | 9 | NO | `ErrModelCatalogNotFound`, `ErrModelCatalogExists`, `ErrModelVendorNotFound` |
| `backend/internal/service/model_catalog_pricing_sync.go` | 274 | NO | `PricingService.ListCatalogEntries()` (method on upstream-owned PricingService), LiteLLM → catalog mapping, vendor/platform/icon mapping helpers |
| `backend/internal/repository/model_catalog_repo.go` | 664 | NO | `modelCatalogRepository` — raw `*sql.DB` (NOT Ent), incl. raw SQL against upstream tables `channel_model_pricing`, `channels`, `channel_pricing_intervals`, `channel_groups`, `groups` |
| `backend/internal/handler/model_catalog_handler.go` | 72 | NO | User/public `ModelCatalogHandler` (`List`, `Vendors`) |
| `backend/internal/handler/admin/model_catalog_handler.go` | 240 | NO | Admin `ModelCatalogHandler` (CRUD + `SyncFromPricing` + vendor CRUD) |
| `backend/internal/handler/dto/model_catalog.go` | 113 | NO | `ModelCatalogResponse`, `ModelVendorResponse`, mappers |
| `backend/migrations/152_model_catalog.sql` | 45 | NO | Creates `model_vendors` + `model_catalogs` tables/indexes |
| `backend/migrations/154_model_catalog_vendor_soft_delete_platform_cleanup.sql` | 42 | NO | Adds `model_vendors.deleted_at`, platform normalization + dedupe |

Note: the Ent schemas exist in REF (and generated `backend/ent/modelcatalog*.go`, `modelvendor*.go` were committed), but the repository implementation does NOT use Ent — it uses raw SQL. No production code under `backend/internal` or `backend/cmd` references `ent.ModelCatalog`/`ent.ModelVendor` (verified via `git grep`). The Ent schemas are effectively dead weight kept in sync with the SQL migrations.

## Backend — upstream-owned files touched by the feature (merge surface)

| File | REF anchor | What REF added |
|---|---|---|
| `backend/internal/server/routes/public.go` | whole file (25 lines) | NEW FILE in REF. `RegisterPublicRoutes(v1, h)` registers `/api/v1/public/models` (+`/vendors`) and `/api/v1/public/channel-monitors`. Does NOT exist on HEAD |
| `backend/internal/server/routes/admin.go` | :106 call, :703-715 `registerModelCatalogRoutes` | `admin.Group("/models")` with 9 routes (list/create/sync-pricing/vendors CRUD/get/update/delete) |
| `backend/internal/server/routes/user.go` | :87-91 | authenticated `/models` group (`GET ""`, `GET /vendors`) — appears UNUSED by REF frontend (frontend always calls `/public/models`) |
| `backend/internal/server/router.go` | :119 | calls `routes.RegisterPublicRoutes(v1, h)` |
| `backend/internal/handler/handler.go` | :36, :57 | `AdminHandlers.ModelCatalog *admin.ModelCatalogHandler`; `Handlers.ModelCatalog *ModelCatalogHandler` |
| `backend/internal/handler/wire.go` | :40, :79, :174, :197, :222, :260 | provider params + `NewModelCatalogHandler`, `admin.NewModelCatalogHandler` in ProviderSet |
| `backend/internal/service/wire.go` | :33-51, :761 | `ProvideModelCatalogService(repo, pricingService, groupRepo)` — warms catalog via `go svc.SyncFromPricing(ctx)` at startup; added to ProviderSet |
| `backend/internal/repository/wire.go` | :97 | `NewModelCatalogRepository` in ProviderSet |
| `backend/cmd/server/wire_gen.go` | :171-173, :247, :270 | generated DI wiring |
| `backend/internal/service/pricing_service.go` | `LiteLLMModelPricing` struct | REF ADDED fields to this upstream-owned struct: `MaxInputTokens`, `MaxOutputTokens`, `MaxTokens`, `SupportsAssistantPrefill`, `SupportsComputerUse`, `SupportsFunctionCalling`, `SupportsPDFInput`, `SupportsReasoning`, `SupportsResponseSchema`, `SupportsToolChoice`, `SupportsVision`, `SupportsWebSearch` — ABSENT on HEAD |
| `backend/internal/service/channel.go` | :44-55 | REF ADDED `DisplayCurrencyUSD/CNY` consts + `NormalizeDisplayCurrency()` — ABSENT on HEAD |

### Adjacent REF-only features that the catalog code touches (optional scope)

- **Group custom /v1/models list** (`models_list_config`): `backend/internal/domain/models_list_config.go` (7 lines, NEW), `backend/internal/service/group_models_list.go` (32 lines, NEW), migration `143_group_models_list_config.sql`, plus edits to upstream-owned `ent/schema/group.go` (:215), `repository/group_repo.go`, `repository/api_key_repo.go`, `service/admin_group.go`, `service/admin_service.go`, `service/admin_group_duplicate.go`, `handler/dto/types.go`, `handler/dto/mappers.go`, `handler/admin/group_handler.go`, `handler/gateway_handler.go` (:1017-1019 runtime filter, :1170-1186 `customModelsListSource`/`filterModelsByCustomList`). HEAD has NO `models-list-candidates` admin route; HEAD instead has `model_allowlist` (`domain.GroupModelAllowlist`, migration `235_group_model_allowlist.sql`).
- **Display currency** (separate task 07-06): migration `155_pricing_display_currency.sql`, `display_currency` column on `channel_model_pricing` + `usage_log`. The catalog repo (`CountPricingAssociations`) SELECTs `display_currency` and calls `service.NormalizeDisplayCurrency` — neither exists on HEAD.

## Frontend — feature-owned files (all absent from HEAD)

| File (REF) | Lines | in HEAD | Description |
|---|---|---|---|
| `frontend/src/views/user/ModelMarketplaceView.vue` | 1022 | NO | Marketplace page; `embedded` prop switches between `AppLayout` (authed `/marketplace`) and `PublicTopNav` (public `/models`) |
| `frontend/src/views/admin/ModelCatalogView.vue` | 1201 | NO | Admin catalog CRUD + vendor management + sync-pricing |
| `frontend/src/api/models.ts` | 199 | NO | `listModels`/`listVendors` → `/public/models`, `/public/models/vendors`; all catalog TS types |
| `frontend/src/api/admin/models.ts` | 100 | NO | `adminModelsAPI` → `/admin/models*` |
| `frontend/src/views/user/modelMarketplaceMonitor.ts` | 54 | NO | `dedupeModelsByModelId`, `matchMonitorsByModelId` helpers |
| `frontend/src/views/user/modelMarketplaceMonitor.spec.ts` | 41 | NO | helper tests |
| `frontend/src/views/user/__tests__/ModelMarketplaceView.spec.ts` | 238 | NO | view tests |
| `frontend/src/api/publicChannelMonitor.ts` | 17 | NO | `/public/channel-monitors` client (monitor badges in marketplace) |
| `frontend/src/components/home/PublicTopNav.vue` | 225 | NO | top nav for public pages; deps: `useBrandHover` (exists on HEAD), `useTheme` (NOT in HEAD, but untracked working-tree file `frontend/src/composables/useTheme.ts` exists), `LocaleSwitcher`, `Icon` |
| `frontend/src/components/home/PublicFooter.vue` | 79 | NO | public footer; :33 links `/models` via `home.footer2.links.modelMarketplace` |
| `frontend/src/i18n/locales/en/models.ts` | 93 | NO | top-level `models.*` keys |
| `frontend/src/i18n/locales/zh/models.ts` | 93 | NO | same, zh |
| `frontend/src/i18n/locales/en/admin/models.ts` | 59 | NO | top-level `models.*` merged into admin module |
| `frontend/src/i18n/locales/zh/admin/models.ts` | 59 | NO | same, zh |

## Frontend — upstream/shared files touched (entry points)

| File | REF anchor | What REF added | HEAD state |
|---|---|---|---|
| `frontend/src/router/index.ts` | :290-311, :527-536 | `/marketplace` (auth, `props:{embedded:true}`), `/models` (public, no auth), `/admin/models` (admin) | HEAD has only `/model-plaza` (:189-191, redirect guard :844); `/models`, `/marketplace`, `/admin/models` all FREE |
| `frontend/src/components/layout/AppSidebar.vue` | :697 user nav, :762 admin nav | `{ path: '/marketplace', label: t('nav.modelMarketplace'), icon: SparklesIcon }`; `{ path: '/admin/models', label: t('nav.modelManagement'), icon: SparklesIcon, hideInSimpleMode: true }` | HEAD user nav :696-706 has no marketplace entry; admin nav :764 area — insertion point between channels group and `/admin/subscriptions` |
| `frontend/src/views/HomeView.vue` | :83, :178 | two `to="/models"` CTA links | HEAD :53-54, :172-173, :269-270, :615 use `showModelPlazaEntry` → `/model-plaza` |
| `frontend/src/components/home/PublicFooter.vue` | :33 | footer link to `/models` | file absent on HEAD |
| `frontend/src/i18n/locales/en/index.ts` / `zh/index.ts` | :4, :12 | `import models from './models'` + `...models` spread | HEAD en/index.ts spreads landing/common/dashboard/channelMonitorV2/batchImage + `admin` + misc |
| `frontend/src/i18n/locales/en/admin/index.ts` / `zh/admin/index.ts` | :2, :13 | `import models from './models'` + spread | HEAD admin/index.ts has extra `plugins` module |
| `frontend/src/i18n/locales/{en,zh}/common.ts` | :193-194 | `nav.modelMarketplace`, `nav.modelManagement` | HEAD common.ts :183 has `nav.modelPlaza` — additive, no key clash |
| `frontend/src/i18n/locales/{en,zh}/landing.ts` | footer links | `home.footer2.links.modelMarketplace` | HEAD has `home.footer2.links.models` (= 'Model Marketplace', used for `/model-plaza` footer link at HomeView.vue:615) |
| `frontend/src/api/admin/index.ts` | — | `import adminModelsAPI from './models'` + `models:` entry + named export + type re-export | HEAD version additionally has `cnProvidersAPI` and `pluginsAPI` (upstream additions) — merge by addition only |
| `frontend/src/utils/pricing.ts` | whole file | REF REPLACED: `DisplayCurrency` type, `normalizeDisplayCurrency`, `currencySymbol`, 2-arg `formatScaled(value, scale, displayCurrency)` | HEAD version has 3-arg `formatScaled(value, scale, minFractionDigits)` + `resolveIntervalPrices`; used by `modelPlaza/PlazaModelPricingTable.vue`, `channels/PricingRow.vue`, `channels/SupportedModelChip.vue`, `utils/__tests__/intervalPrices.spec.ts` — DIRECT CONFLICT |

No Pinia store or composable is owned by the feature (uses `useAppStore`/`useAuthStore` from `@/stores`).

## Related specs / tasks in REF

- `.trellis/tasks/archive/2026-06/06-23-console-status-model-catalog/` — original catalog task (prd, research, implement.jsonl)
- `.trellis/tasks/archive/2026-06/06-25-refine-model-marketplace-monitors/prd.md` — monitor badges in marketplace
- `.trellis/tasks/archive/2026-06/06-26-optimize-model-status-monitor-selects/prd.md`
- `.trellis/tasks/07-07-add-model-marketplace-console-menu/prd.md` — sidebar entry
- `.trellis/tasks/archive/2026-07/07-06-pricing-currency-display/prd.md` — display-currency feature (coupling source)

## Caveats / Not Found

- No settings keys (`SettingKey*`) or feature flags gate the marketplace in REF — routes/nav are always on (admin catalog entry hidden only in simple mode).
- No backend unit tests for catalog service/repository/handler found in REF (`git grep ModelCatalog -- '*_test.go'` in backend returned nothing outside wire).
- REF frontend admin group tests `groupsModelsList*.spec.ts` cover the optional group models-list feature, not the catalog itself.
