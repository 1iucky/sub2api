# Replay model marketplace and admin model catalog

## Goal

Port the personalized Model Marketplace (用户侧模型集市) and Admin Model Catalog (管理员模型管理) from reference branch `codex/custom-theme-upstream-main-20260717` (REF) onto the current branch. The result must coexist with upstream's Model Plaza (`/model-plaza`) with zero route/API/table collisions, carry no runtime coupling to upstream-owned internals, and keep all edits to upstream-owned files as small additive blocks so future `upstream/main` merges stay conflict-minimal.

Parent task: `09-21-repair-personalized-upstream-merge-ui`.

## Background and Confirmed Facts

Evidence: `research/feature-inventory.md`, `research/coupling-analysis.md`, `research/collision-surface.md`, `research/decoupling-options.md` (2026-09-29).

### REF feature scope (all feature-owned files absent from HEAD — no file-level conflicts)

- **Marketplace (user/public)**: routes `/marketplace` (auth, embedded) + `/models` (public), view `frontend/src/views/user/ModelMarketplaceView.vue` (1022 lines), `api/models.ts` → `GET /api/v1/public/models(+vendors)`, monitor badges via `api/publicChannelMonitor.ts` → `/api/v1/public/channel-monitors`, `modelMarketplaceMonitor.ts` (+spec), `PublicTopNav.vue` / `PublicFooter.vue`, i18n `locales/{en,zh}/models.ts`, view spec.
- **Admin catalog**: route `/admin/models`, view `views/admin/ModelCatalogView.vue` (1201 lines), `api/admin/models.ts`, i18n `locales/{en,zh}/admin/models.ts`; backend `/admin/models` group (9 routes: CRUD + `sync-pricing` + vendor CRUD).
- **Backend feature-owned files (11)**: `service/model_catalog.go` / `_errors.go` / `_pricing_sync.go`, `repository/model_catalog_repo.go` (raw SQL, not Ent), `handler/model_catalog_handler.go`, `handler/admin/model_catalog_handler.go`, `handler/dto/model_catalog.go`, ent schemas `model_catalog.go` / `model_vendor.go`, migrations `152_model_catalog.sql` + `154_model_catalog_vendor_soft_delete_platform_cleanup.sql`.
- Ent schemas in REF are dead code (repo uses raw SQL); edges are catalog↔vendor only — zero ent merge risk. Port them per repo convention (`make generate`).
- Tables `model_catalogs` / `model_vendors` do not exist on HEAD. HEAD migrations occupy slots 152/154 and run to 238 → new migrations numbered **239+** (embed FS auto-load, checksum validation forbids editing applied files).

### Collision surface (verified, see collision-surface.md)

- Frontend routes `/models`, `/marketplace`, `/admin/models` FREE on HEAD (HEAD has only `/model-plaza`, `router/index.ts:189`). Route names `ModelMarketplace`/`UserModelMarketplace`/`AdminModels` free.
- Backend: no `/api/v1/public` group and no `/admin/models` group on HEAD; gateway `/v1/models` lives outside `/api/v1` — no clash. REF's authenticated `/models` group in `user.go` is dead code (frontend always calls `/public/models`) — not ported.
- i18n top-level `models.*` (public + admin) FREE; `nav.modelMarketplace` / `nav.modelManagement` additive in `common.ts` (HEAD has `nav.modelPlaza` at :183).
- **Do not touch** HEAD `frontend/src/utils/pricing.ts` (used by plaza/channels: `formatScaled` 3-arg, `resolveIntervalPrices`). Marketplace inlines USD display (`'$'`); no currency helpers module.
- `home.footer2.links.modelMarketplace` added alongside HEAD's existing `home.footer2.links.models` (the plaza footer label, HomeView.vue:615) — both links intentional (two distinct pages).
- `api/admin/index.ts` merged additively (keep HEAD's `cnProvidersAPI`/`pluginsAPI`).
- No settings keys / feature flags introduced (REF had none; user decision #4).

### REF coupling defects fixed in this port (see coupling-analysis.md)

- `model_catalog_pricing_sync.go` reads 12 `LiteLLMModelPricing` capability fields HEAD deleted → compile break. Fixed by catalog-owned fetch/parse (R1).
- `CountPricingAssociations` selects `display_currency` + calls `service.NormalizeDisplayCurrency` — absent on HEAD → compile+runtime break. Fixed by dropping currency (USD-only).
- `GroupRepository` injected into `ModelCatalogService` but never called → removed from constructor.
- REF-added fields on upstream-owned `pricing_service.go` and consts on `channel.go` are NOT re-added.
- HEAD pricing tables drifted but compatible: `channel_model_pricing` adds `cache_write_1h_price`/multipliers/`image_input_price`/`time_pricing` (REF's SELECT subset still works); `channel_pricing_intervals` extra columns harmless.

### User decisions (all resolved)

1. Marketplace and plaza **fully coexist**; marketplace keeps independent routes/APIs/tables; upstream plaza untouched.
2. **Cut runtime coupling** to upstream internals; keep only data-sync entry points.
3. Child task of `09-21-repair-personalized-upstream-merge-ui`.
4. Marketplace entries **always visible, no feature flag** (mirrors REF; avoids opt-in trap and settings-key merge debt). No `SettingKey` additions.
5. Adjacent REF sub-features **dropped**: per-group custom `/v1/models` display list (`models_list_config`, ~10 upstream-owned files, never read by catalog/marketplace; HEAD's `model_allowlist` covers the concept) and display-currency/CNY pricing (07-06 task). User-visible deltas: prices display `$` only; no per-group curated `/v1/models` display list.

## Requirements

### R1. Backend: self-contained catalog domain

- Port service/errors/repository/dto/handlers as **new files only**. No struct-field additions to upstream-owned types.
- Repository keeps raw-SQL approach; pricing-association reads (`CountPricingAssociations`, `WithPricingOnly`, dedupe ranking) drop `display_currency` (USD-only); no `NormalizeDisplayCurrency`.
- Constructor drops dead `GroupRepository` param.
- **Pricing sync = catalog-owned fetcher**: fetch `cfg.Pricing.RemoteURL` (fallback `cfg.Pricing.FallbackFile`), parse raw LiteLLM-format JSON into catalog-local structs retaining `supports_*` capability keys; upsert catalog (vendor/platform/icon mapping pure functions ported as-is). Zero runtime import of upstream `PricingService` internals.
- Sync triggers: admin `POST /admin/models/sync-pricing` endpoint + startup warm goroutine (2-min timeout, log-only failure — REF behavior preserved).
- Ent schemas `model_catalog.go`/`model_vendor.go` ported (fields/indexes per coupling-analysis §3); `make generate` run; generated files committed.
- Migrations: REF 152+154 content merged/renumbered to `239_`/`240_` (final state: `model_vendors` incl. `deleted_at` + platform normalization, `model_catalogs` + indexes).

### R2. Backend: routes (additive)

- New `routes/public.go`: `/api/v1/public/models(+vendors)` and `/api/v1/public/channel-monitors` (monitor handlers already exist on HEAD); one additive call in `router.go` beside `RegisterModelPlazaRoutes`.
- `/admin/models` group added in `routes/admin.go` (9 routes).
- Handler/wire registries: additive fields/providers only (`handler.go`, `handler/wire.go`, `service/wire.go`, `repository/wire.go`, regenerate `wire_gen.go`).

### R3. Frontend: marketplace

- Port view/api/monitor helpers/specs, `PublicTopNav.vue`, `PublicFooter.vue`, i18n `models.ts` (en/zh), routes `/models` (public) + `/marketplace` (auth, `props:{embedded:true}`).
- `PublicTopNav.vue` needs `useTheme` — present as untracked working-tree file from the parent task; verify it lands with the parent task or co-deliver an equivalent.
- USD-only display: inline `'$'`; do not import from `utils/pricing.ts`.
- Homepage: port REF HomeView's two `to="/models"` CTA links (REF :83, :178) into current HomeView nav structure, always visible; footer marketplace link via `PublicFooter.vue` + `home.footer2.links.modelMarketplace`.
- User sidebar entry `{ path: '/marketplace', label: t('nav.modelMarketplace') }` (REF AppSidebar :697), always rendered.

### R4. Frontend: admin catalog

- Port `ModelCatalogView.vue`, `api/admin/models.ts`, admin i18n `models.ts` (en/zh), route `/admin/models` (name `AdminModels`, requiresAdmin), additive `api/admin/index.ts` merge.
- Admin sidebar entry `{ path: '/admin/models', label: t('nav.modelManagement'), icon: SparklesIcon, hideInSimpleMode: true }` inserted between channels group and `/admin/subscriptions` (REF AppSidebar :762).

### R5. Merge-debt isolation

- Every edit to an upstream-owned file is a small additive block: `router/index.ts`, `AppSidebar.vue`, `HomeView.vue`, locale index/common/landing files, `api/admin/index.ts`, `routes/admin.go`, `router.go`, `handler/handler.go`, three `wire.go` files.
- No edits to: `pricing_service.go`, `channel.go`, ent `group.go`, `utils/pricing.ts`, `gateway_handler.go`, any group repos/DTOs, plaza files, `domain_constants.go` settings section.

## Out of Scope

- REF per-group custom `/v1/models` display list (`models_list_config` sub-feature) — dropped (user decision #5).
- Display-currency/CNY pricing (07-06 task) — dropped; USD-only (user decision #5).
- REF authenticated `/models` backend group (dead code).
- Any change to upstream Model Plaza behavior/routes/settings/entries.
- Backend unit tests for catalog (REF had none); frontend specs are ported.

## Acceptance Criteria

- [ ] Public `/models` renders catalog from `GET /api/v1/public/models` with vendor/platform/capability filters and live monitor badges, no login required.
- [ ] Authenticated `/marketplace` renders the same marketplace embedded in the console layout; user sidebar shows the always-visible `/marketplace` entry.
- [ ] Homepage top nav shows an always-visible marketplace entry to `/models` (plaza entry remains flag-gated as before); footer shows both model-page links.
- [ ] Admin `/admin/models` supports catalog CRUD, vendor CRUD (soft delete), and "sync from remote pricing" which populates/updates entries incl. capability tags; admin sidebar shows `nav.modelManagement` (hidden in simple mode).
- [ ] Startup warm sync populates the catalog on fresh deploy without manual action; failure is log-only and does not block boot.
- [ ] Upstream Model Plaza (`/model-plaza`, `model_plaza_enabled`/`require_auth` settings, its homepage/header/footer entries) behaves exactly as before.
- [ ] `git grep` confirms: no `display_currency`, no `NormalizeDisplayCurrency`, no `GroupRepository` in catalog code; no diff in `pricing_service.go` / `channel.go` / ent `group.go` / `utils/pricing.ts`.
- [ ] Backend: `make generate` clean, `make test-unit` pass, `golangci-lint` pass (depguard layering respected).
- [ ] Frontend: `pnpm run typecheck`, `pnpm run lint:check`, `pnpm test:run` (incl. ported specs), `pnpm run build` pass.
- [ ] Migrations 239/240 apply cleanly on fresh Postgres 15 and on a DB already at 238.
- [ ] Browser verification (reported separately from static checks): public marketplace, console marketplace, admin catalog, homepage entries, plaza unaffected.
- [ ] Unrelated pre-existing working-tree changes (`.gitattributes`, `deploy/`, `.agents/`, `.cursor/`, `.trellis/`, parent-task frontend files) remain intact.
