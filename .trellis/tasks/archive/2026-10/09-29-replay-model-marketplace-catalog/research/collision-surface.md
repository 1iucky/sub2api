# Research: Collision surface — reference feature vs current HEAD

- **Query**: route/naming/settings/migration conflicts between REF feature code and HEAD (incl. upstream Model Plaza)
- **Scope**: internal
- **Date**: 2026-09-29

## 1. Frontend route path conflicts

Checked HEAD `frontend/src/router/index.ts`:

| REF route | HEAD status |
|---|---|
| `/marketplace` (name `UserModelMarketplace`, auth) | FREE — no match |
| `/models` (name `ModelMarketplace`, public) | FREE — HEAD has only `/model-plaza` (:189-191, name `ModelPlaza`) and a redirect guard at :844 (`if (to.path === '/model-plaza')`) |
| `/admin/models` (name `AdminModels`, admin) | FREE — HEAD admin routes span :416-716; closest are `/admin/channels/pricing` (:480) and `/admin/channels/monitor` (:492) |

No path collisions. Naming similarity with `ModelPlaza` route/view is cosmetic but worth avoiding in route NAMES (REF uses `ModelMarketplace`/`UserModelMarketplace`/`AdminModels` — all free).

## 2. Backend route conflicts

- HEAD has **no** `backend/internal/server/routes/public.go` and no `/api/v1/public` group (only `/payment/public` in `routes/payment.go:53` and `/settings/public` in `routes/auth.go:246`). REF's `RegisterPublicRoutes` (new file + call in `router.go`) is additive — no conflict.
- HEAD `routes/admin.go` has **no** `admin.Group("/models")` (verified by grep). REF `registerModelCatalogRoutes` is additive. One callout: HEAD admin.go registers other `/models`-suffixed endpoints (`groups/:id/models-list-candidates` does NOT exist on HEAD; `accounts/:id/models`, `dashboard/models` do — different prefixes, no clash).
- HEAD `routes/user.go` has no authenticated `/models` group; REF's user.go addition (:87-91) is additive but appears **dead** (REF frontend always calls `/public/models`).
- HEAD `router.go:128-130` registers User/ModelPlaza/Admin routes; inserting `RegisterPublicRoutes(v1, h)` mirrors REF `router.go:119`.
- Gateway `/v1/models` etc. (`routes/gateway.go:211,377`) live outside `/api/v1` — no clash with `/api/v1/public/models`.

## 3. Naming conflicts

### Backend

| REF name | HEAD collision? |
|---|---|
| Ent schemas `ModelCatalog`, `ModelVendor`; tables `model_catalogs`, `model_vendors` | none (no such files/tables on HEAD) |
| `service.ModelCatalogService`, `ModelCatalogRepository` (interface), `ModelCatalog`, `ModelVendor` structs | none — HEAD has `ModelPlazaService`/`modelPlazaService`, distinct names |
| `handler.ModelCatalogHandler`, `admin.ModelCatalogHandler` | none — HEAD has `ModelPlazaHandler` only |
| `AdminHandlers.ModelCatalog`, `Handlers.ModelCatalog` fields | none — HEAD `handler.go:36-70` has `ModelPlaza` field but no `ModelCatalog` |
| `ProvideModelCatalogService`, `NewModelCatalogRepository`, `NewModelCatalogHandler` | none |
| `ErrModelCatalogNotFound/Exists`, `ErrModelVendorNotFound` | none |
| `NormalizeModelCatalogID/Platform`, `PricingCatalogEntry` | none |
| **CONFLICT**: `PricingService.ListCatalogEntries()` — method on upstream type, must be re-added to HEAD's upstream-owned `pricing_service.go` area (new file in `service` package is fine since same package) | compile-level, resolved by porting the file |
| **CONFLICT**: REF-added `LiteLLMModelPricing` fields (12 capability/token-limit fields) | HEAD struct lacks them (see coupling-analysis §1.2) |
| **CONFLICT**: `NormalizeDisplayCurrency`, `DisplayCurrencyUSD/CNY` + `display_currency` column | absent on HEAD entirely |

### Frontend

| REF name | HEAD collision? |
|---|---|
| `views/user/ModelMarketplaceView.vue`, `views/admin/ModelCatalogView.vue` | none (HEAD: `views/ModelPlazaView.vue`, `components/modelPlaza/*`) |
| `api/models.ts`, `api/admin/models.ts` | none (HEAD: `api/modelPlaza.ts`) |
| `adminAPI.models` / `adminModelsAPI` | none; HEAD `api/admin/index.ts` must be merged additively (HEAD has extra `cnProviders`/`plugins` entries REF lacks — keep HEAD's, add `models`) |
| i18n top-level key `models.*` (from `locales/{en,zh}/models.ts`) | **FREE** — no HEAD module exports top-level `models` (checked all `locales/en/*.ts`; `dashboard.ts` exports `modelPlaza`, `landing.ts` nests `models` under `home.*`) |
| i18n admin key `models.*` merged into admin module | **FREE** — admin modules export `overview/channels/accounts/resources/ops/settings/audit/promptAudit/plugins` |
| `nav.modelMarketplace`, `nav.modelManagement` (common.ts :193-194 in REF) | **FREE** — HEAD common.ts:183 has `nav.modelPlaza`; same file, additive keys |
| `home.footer2.links.modelMarketplace` (REF landing.ts) | HEAD landing.ts:246 already has `home.footer2.links.models` = 'Model Marketplace' currently used as the `/model-plaza` footer label (HomeView.vue:615). Adding REF's key alongside creates **duplicate-link UX** (two footer links to two model pages) — decide which entries stay |
| **CONFLICT**: `frontend/src/utils/pricing.ts` | REF version (DisplayCurrency + 3rd-arg `displayCurrency` in `formatScaled`) vs HEAD version (3rd-arg `minFractionDigits`, `resolveIntervalPrices`). HEAD callers: `components/modelPlaza/PlazaModelPricingTable.vue`, `components/channels/PricingRow.vue`, `components/channels/SupportedModelChip.vue`, `utils/__tests__/intervalPrices.spec.ts`. REF marketplace view needs `currencySymbol`. **Do not overwrite HEAD's file — add currency helpers under different names or in a new module** |
| `components/home/PublicTopNav.vue`, `PublicFooter.vue` | files absent on HEAD (HEAD `components/home/` has only BillingSparkline/GatewayDashboard/PoolingChart/RoutingDiagram). REF `PublicTopNav` needs `useTheme` (absent in HEAD, present as untracked working-tree file) and `useBrandHover` (present) |
| `SparklesIcon` reuse in AppSidebar for two entries (marketplace + modelManagement in REF; HEAD uses it elsewhere) | cosmetic only |

## 4. Settings keys

- REF catalog/marketplace introduces **no** settings keys and no frontend `FeatureFlags` entries (verified by grep over REF `domain_constants.go`, setting service/handler, and frontend constants).
- HEAD-only plaza settings `model_plaza_enabled` / `model_plaza_require_auth` (read in `HomeView.vue:687` and `ModelPlazaView`) are untouched.
- REF's *display-currency* sub-feature did add `pricing_display_currency` migration (155) — not part of the two target features but pulled in transitively by `CountPricingAssociations`.

## 5. Migration numbering

REF feature migrations and their HEAD slot occupants:

| REF file | HEAD file with same number |
|---|---|
| `152_model_catalog.sql` | `152_scheduler_outbox_dedup_key.sql` |
| `154_model_catalog_vendor_soft_delete_platform_cleanup.sql` | `154_account_spark_shadow.sql`, `154_add_ops_system_logs_api_key_id.sql` (+`154a_...`) |
| (optional) `143_group_models_list_config.sql` | HEAD 143 slot also taken by upstream content |
| (transitive) `155_pricing_display_currency.sql` | `155_add_ops_system_logs_api_key_id_index_notx.sql` |

Last numbered migrations on HEAD (tail): `231_add_usage_log_requested_reasoning_effort.sql`, `231_user_restrict_public_groups.sql`, `232_add_usage_log_upstream_request_id.sql`, `232_channel_cache_write_1h_pricing.sql`, `232_group_force_openai_fast.sql`, `232_group_reasoning_effort_over_limit.sql`, `233_add_usage_log_upstream_request_id_index_notx.sql`, `233_group_free_openai_fast.sql`, `234_channel_max_reasoning_effort_multiplier.sql`, `234_group_codex_models_manifest_config.sql`, `235_group_model_allowlist.sql`, `236_group_model_allowlist_repair.sql`, `237_add_minimax_platform.sql`, `238_opencode_go_platform.sql`, `238_purge_unlimited_user_platform_quotas.sql`.

→ New catalog migrations must be numbered **239+** (embed FS auto-loads `*.sql`, `migrations/migrations.go:33-34`; checksum validation means never edit applied files).

## 6. Coexistence with Model Plaza (product-level collision)

- HEAD upstream Model Plaza: public/user page `/model-plaza` (`views/ModelPlazaView.vue`), API `/api/v1/model-plaza` (`routes/model_plaza.go:23-28`, optional JWT, feature-gated by `model_plaza_enabled`), computes model list **live** from groups/channels (`service/model_plaza_service.go`), no DB tables.
- REF Marketplace: DB-backed catalog with vendor metadata, capability facets, pricing sync from LiteLLM; public `/models` + authed `/marketplace`.
- Overlap is **UX-level only** (both list models with pricing). Data paths are fully disjoint — no backend symbol or table overlap. Homepage/sidebar/footer will show both unless entries are deliberately deduplicated (HomeView.vue:53-54,172-173,269-270,615 for plaza; REF HomeView :83,:178 + sidebar :697 for marketplace).

## Caveats

- HEAD working tree contains uncommitted personalized changes to `AppSidebar.vue`, `HomeView.vue`, `landing.ts` locales — anchors above refer to committed HEAD; re-check line numbers at port time.
