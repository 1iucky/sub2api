# Design: Replay model marketplace and admin model catalog

## Architecture and Boundaries

The port creates a **self-contained vertical slice** (DB → repository → service → handlers → routes → frontend) that touches upstream-owned code only at additive registration seams.

```
┌─ feature-owned (new files, zero upstream edits) ─────────────────────┐
│ backend:                                                               │
│   service/model_catalog.go, model_catalog_errors.go,                   │
│   service/model_catalog_pricing_sync.go  (rewritten: own fetcher)      │
│   repository/model_catalog_repo.go       (raw SQL, USD-only)           │
│   handler/model_catalog_handler.go       (public List/Vendors)         │
│   handler/admin/model_catalog_handler.go (CRUD + sync + vendors)       │
│   handler/dto/model_catalog.go                                         │
│   ent/schema/model_catalog.go, model_vendor.go                         │
│   server/routes/public.go                (new /api/v1/public group)    │
│   migrations/239_model_catalog.sql, 240_model_vendor_soft_delete.sql   │
│ frontend:                                                              │
│   views/user/ModelMarketplaceView.vue (+spec, monitor helper +spec)    │
│   views/admin/ModelCatalogView.vue                                     │
│   api/models.ts, api/admin/models.ts, api/publicChannelMonitor.ts      │
│   components/home/PublicTopNav.vue, PublicFooter.vue                   │
│   i18n locales {en,zh}/models.ts, {en,zh}/admin/models.ts              │
└───────────────────────────────────────────────────────────────────────┘
┌─ additive seams (upstream-owned files, small blocks only) ───────────┐
│ routes/admin.go (+9 routes)  router.go (+1 call)                       │
│ handler/handler.go (+2 fields)  handler/wire.go, service/wire.go,      │
│ repository/wire.go (+providers)  cmd/server/wire_gen.go (generated)    │
│ frontend router/index.ts (+3 routes)  AppSidebar.vue (+2 entries)      │
│ HomeView.vue (+2 CTA links)  api/admin/index.ts (+models)              │
│ locales index/common/landing (+imports, keys)                          │
└───────────────────────────────────────────────────────────────────────┘
```

## Key Design Decisions

### D1. Pricing sync: catalog-owned fetcher (REF option B2)

REF coupled sync to upstream `PricingService.pricingData` and 12 deleted `LiteLLMModelPricing` fields. The port instead:

- Adds a catalog-owned fetch/parse component (in `service/model_catalog_pricing_sync.go` or a sibling new file) that:
  - GETs `cfg.Pricing.RemoteURL` (default `https://raw.githubusercontent.com/Wei-Shaw/model-price-repo/main/model_prices_and_context_window.json`, `config.go:2292`) with timeout; on failure falls back to `cfg.Pricing.FallbackFile` (`config.go:2295`).
  - Parses the raw JSON into a **catalog-local** struct that keeps `supports_*` capability keys, `max_*_tokens`, costs, `litellm_provider`, `mode` — the full facet set REF exposed.
- Reuses REF's pure mapping functions verbatim: `capabilitiesFromLiteLLM` (re-based on the local struct), `pricingMapFromLiteLLM`, `vendorNameForProvider`, `iconKeyForProvider`, `providerToPlatform`, `NormalizeModelCatalogPlatform`, `vendorSortOrder`, `endpointsForMode`, `modelTagsFromPricing`.
- `ModelCatalogService` constructor: `NewModelCatalogService(repo, cfg, httpClient-ish deps)` — no `PricingService`, no `GroupRepository`.
- Startup warm: `ProvideModelCatalogService` keeps REF's `go svc.SyncFromPricing(ctx)` (2-min timeout, log-only failure). Idempotent upsert → safe.

Rationale: zero runtime imports of upstream pricing internals (user decision #2) while preserving the full capability-facet UX (B1 would shrink facets; B3 would empty the marketplace).

### D2. Pricing-association reads: raw SQL minus currency

- `model_catalog_repo.go` ported with `display_currency` removed from SELECTs and the `NormalizeDisplayCurrency` call deleted. All prices render as USD.
- HEAD's drifted pricing tables are compatible: REF's column subset exists unchanged on HEAD (`channel_model_pricing`, `channel_pricing_intervals`, `channels`, `channel_groups`, `groups`; coupling-analysis §1.3). Optionally surface `cache_write_1h_price` — **deferred**, not required for parity.
- Repository helpers `isUniqueViolation` / `escapeLike` exist in the shared repository package (unchanged) — reused, same package.

### D3. Ent schemas ported despite raw-SQL repository

CLAUDE.md mandates Ent as source of truth for DB models. Schemas ported exactly (coupling-analysis §3: fields/indexes, internal-only edge catalog→vendor with `OnDelete(SetNull)`), `make generate` run, generated code committed. The repository continues to use raw SQL (as in REF) — generated Ent code is for schema truth/future use, matching REF's state.

### D4. Routes

- New `backend/internal/server/routes/public.go` mirrors REF (`RegisterPublicRoutes(v1, h)`): `GET /api/v1/public/models`, `GET /api/v1/public/models/vendors`, `GET /api/v1/public/channel-monitors` (+status sub-route as REF). Called from `router.go` next to `RegisterModelPlazaRoutes`. Monitor handlers (`ChannelMonitorUserHandler.List/GetStatus`) already exist on HEAD.
- `routes/admin.go`: additive `registerModelCatalogRoutes` block (9 routes) + one call line.
- REF's authenticated `/models` group in `user.go`: **not ported** (dead code).

### D5. Frontend

- Routes added to `router/index.ts`: `/models` (public, name `ModelMarketplace`), `/marketplace` (auth, name `UserModelMarketplace`, `props:{embedded:true}`), `/admin/models` (admin, name `AdminModels`).
- Marketplace view ported verbatim except: USD display inlined (`'$'`), no import from `@/utils/pricing`; `DisplayCurrency` type removed.
- `PublicTopNav.vue` dependency chain: `useBrandHover` (HEAD ✓), `useTheme` (untracked working-tree file from parent task — verify presence at port time; if absent, port REF's `useTheme.ts`), `LocaleSwitcher`, `Icon` (HEAD ✓).
- Sidebar: user entry always rendered; admin entry `hideInSimpleMode: true`; both additive array elements.
- HomeView: two `to="/models"` CTAs ported into the parent's restored personalized HomeView structure; existing `showModelPlazaEntry`-gated plaza entries untouched.
- i18n: new modules `models.ts` + `admin/models.ts` wired into locale index files additively; `nav.modelMarketplace`/`nav.modelManagement` into `common.ts`; `home.footer2.links.modelMarketplace` into `landing.ts` (existing `links.models` stays, pointing to plaza).

### D6. Wire / DI

- `repository/wire.go`: `NewModelCatalogRepository(db)` added to ProviderSet.
- `service/wire.go`: `ProvideModelCatalogService(repo, cfg)` (warm goroutine inside) + `wire.Bind` for the repository interface; added to ProviderSet.
- `handler/wire.go`: both handlers + struct fields in `Handlers`/`AdminHandlers`.
- Regenerate `cmd/server/wire_gen.go` via `make generate`.

## Data Flow

**Catalog sync**: admin clicks sync (or boot warm) → catalog fetcher pulls remote LiteLLM JSON → parse to local structs → map (vendor/platform/icon/tags/capabilities/endpoints/pricing) → upsert by `(platform, normalized_model_id)` → set `last_synced_at`, `source='litellm'` (per REF) → vendors upserted by name.

**Public marketplace**: browser → `GET /api/v1/public/models?...` → service list (filters: platform/vendor/capability/search/price-range per REF) → repo raw SQL over `model_catalogs`/`model_vendors` + pricing-association enrichment from channel tables → dto → view. Monitor badges: `GET /api/v1/public/channel-monitors` → `matchMonitorsByModelId`.

**Admin catalog**: standard CRUD + vendor CRUD (soft delete via `deleted_at`) + `POST /admin/models/sync-pricing`.

## Compatibility and Migration Notes

- Migrations renumbered `239_model_catalog.sql` (= REF 152 content) and `240_model_catalog_vendor_soft_delete_platform_cleanup.sql` (= REF 154 content). Checksum validation: never edit applied files; new files only.
- Applies on fresh DB (full chain) and existing DBs at ≤238 (only 239/240 run).
- No settings keys, no feature flags, no config schema changes (reuses existing `pricing.remote_url`/`fallback_file`).
- Simple mode: admin entry hidden (as REF); public marketplace stays visible (REF behavior).

## Trade-offs

- USD-only prices (no CNY toggle) — accepted (user decision #5).
- No per-group curated `/v1/models` display list — accepted (user decision #5).
- Catalog-owned pricing fetch duplicates upstream PricingService's fetch/parse logic (~150 lines) — deliberate cost of zero-coupling.
- Startup warm goroutine performs network I/O at boot; failure is log-only (REF parity). Environments blocking egress simply start with an empty catalog until manual sync/fallback file.

## Rollback

- Feature is fully additive: rollback = revert the port commit(s). Tables `model_catalogs`/`model_vendors` are feature-owned; dropping them is optional (migrations have no down-files per repo convention — leave tables in place on rollback).
- No upstream-owned file receives semantic edits, so rollback cannot affect plaza/channels/gateway behavior.
