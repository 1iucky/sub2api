# Research: Decoupling options — cutting runtime coupling to upstream-owned code

- **Query**: for each coupled point found in coupling-analysis.md, propose a HEAD-compatible re-implementation with zero (or minimal) runtime imports of volatile upstream code; note user-visible behavior changes
- **Scope**: internal
- **Date**: 2026-09-29

Coupled points, ordered by severity. "Upstream-owned" = files upstream/main also modifies (pricing_service.go, channel.go, gateway_handler.go, group repos/dto, utils/pricing.ts).

---

## Option A — Pricing-association read path (`CountPricingAssociations`, `WithPricingOnly`, dedupe ranking)

REF couples to upstream tables (`channel_model_pricing.display_currency`, `channels`, `channel_pricing_intervals`, `channel_groups`, `groups`) AND to `service.NormalizeDisplayCurrency` (REF-added to upstream-owned `channel.go`).

- **A1 (recommended): keep raw-SQL reads, drop currency entirely.**
  Remove `display_currency` from the SELECT and delete the `NormalizeDisplayCurrency` call (or define a tiny catalog-local `normalizeDisplayCurrency` — but since HEAD has no `display_currency` column anywhere, the value would always be USD). The catalog repository already owns its SQL; it imports nothing from upstream services at runtime — only shared, stable schema knowledge of upstream tables.
  *Behavior change vs REF*: prices always display in `$`; no CNY toggle in the marketplace/admin catalog. Everything else identical.
- **A2: also port display-currency** (migration + column + `NormalizeDisplayCurrency` in a catalog-owned file, e.g. extend `model_catalog_errors.go`-style new file). More work, touches upstream-owned `channel_repo_pricing.go` write paths and channel admin UI to be meaningful end-to-end. Only worth it if the CNY display feature is explicitly wanted later.
- **A3 (optional hardening): stop reading interval multiplier columns.** HEAD added `cache_write_1h_price` and `*_multiplier` columns to both pricing tables; REF's SELECT subset still works. Optionally extend `ModelCatalogPriceRange`/`Entry` with `cache_write_1h_price` to match HEAD's current pricing model.
  *Behavior change*: none if skipped; marketplace misses the 1h-cache price line if kept skipped.

Zero-import result: `model_catalog_repo.go` then depends only on `service` package types + stdlib + stable table names — same coupling class as any upstream repository.

---

## Option B — LiteLLM pricing sync (`SyncFromPricing` / `ListCatalogEntries`)

REF reads `PricingService.pricingData` (in-memory map of upstream-owned `LiteLLMModelPricing`) and 12 struct fields HEAD deleted.

- **B1 (recommended): sync only from fields HEAD still has.** Rewrite `capabilitiesFromLiteLLM`/`pricingMapFromLiteLLM`/metadata to use HEAD's struct: costs (input/output/cache/priority/1hr/image), `LongContextInputTokenThreshold` + multipliers, `SupportsServiceTier`, `SupportsPromptCaching`, `LiteLLMProvider`, `Mode`. Drop the 12 missing capability fields (`vision`, `reasoning`, `function_calling`, …) and `max_*_tokens` from capabilities/metadata.
  *Behavior change vs REF*: capability facet filters in the marketplace sidebar shrink to what HEAD's pricing data supports (cache, image-output, long-context, service-tier); context-window ranges (REF view's `contextBands`) will have no data and the filter becomes inert unless an alternate source (below) is used. Tags/cards otherwise unchanged.
- **B2: manual/CSV admin import instead of live struct coupling.** Catalog sync reads the same remote LiteLLM JSON directly (via its own fetch/parse in the catalog package) instead of `PricingService.pricingData`. Zero runtime import of upstream pricing internals; cost = duplicating the fetch/parsing and running on demand (admin clicks "Sync").
  *Behavior change*: sync no longer shares the PricingService cache/refresh cycle; startup warm-goroutine removed or made explicit.
- **B3: drop sync entirely; catalog is manual-only** (admin CRUD + vendor management). Smallest port.
  *Behavior change*: "Sync Remote Pricing" button removed from admin view; catalog starts empty; marketplace shows nothing until admin populates it.

All three remove the need to re-add fields to upstream-owned `pricing_service.go`. B1 keeps the exact REF UX minus a few facets; B2 is the only option that restores full capability facets without touching upstream files (it can parse the raw LiteLLM JSON which still contains `supports_*` keys).

---

## Option C — Group custom models list (`models_list_config`)

REF display-only filter for `GET /v1/models`, entangled with upstream-owned `ent/schema/group.go`, group repos, dto, admin group handler, and `gateway_handler.go`. HEAD replaced this concept space with enforcement-oriented `model_allowlist` (`domain.GroupModelAllowlist{Enabled, Models}` — structurally identical to REF's config — plus `middleware/group_model_allowlist.go` and `pkg/requestmodel`).

- **C1 (recommended): DROP the sub-feature.** The catalog/marketplace never reads `ModelsListConfig` (verified: no reference in any catalog file). Not porting it loses nothing for the two target features.
  *Behavior change vs REF*: admins cannot set a per-group curated display list for `/v1/models`; groups always show the full available list. REF admin UI section + `groupsModelsList*.spec.ts` tests not ported.
- **C2: re-implement on top of HEAD's allowlist read path** — treat a group's `model_allowlist` (when enabled) as the display filter for `/v1/models`. Reuses HEAD-owned fields; no schema changes; but semantics merge display with enforcement.
  *Behavior change*: curating the display list also blocks non-listed models from the gateway (allowlist is enforced upstream) — stricter than REF.
- **C3: port REF's separate `models_list_config` column verbatim** (migration 239+, edits to ~10 upstream-owned group files).
  Highest merge debt; preserves exact REF semantics (display-only, independent of allowlist).

---

## Option D — Dead `GroupRepository` dependency

REF `NewModelCatalogService(repo, pricingService, groupRepo)` stores `groupRepo` but never calls it (verified by reading all 420 lines of `model_catalog.go`).
- **D1 (recommended): remove the parameter** in the port (`NewModelCatalogService(repo, pricingService)`; `ProvideModelCatalogService` likewise). One less upstream-interface dependency; no behavior change. If B3 is chosen, `pricingService` goes too.

---

## Option E — Public routes + monitor badges

REF public marketplace needs `/api/v1/public/models(+vendors)` and (for monitor badges) `/api/v1/public/channel-monitors`. HEAD has no public group but has `ChannelMonitorUserHandler.List/GetStatus` (registered under auth at `routes/user.go:141-142`).

- **E1 (recommended): create `routes/public.go` on HEAD** registering models + channel-monitors, mirroring REF `public.go:8-25`, called from `router.go` next to `RegisterModelPlazaRoutes` (:129). Handlers already exist on HEAD. Additive; no upstream file edits beyond one line in `router.go`.
  *Behavior change*: monitor status (availability/latency badges) becomes publicly readable without login — same as REF. If that's undesirable, register only `/public/models` and drop the monitor badges + `publicChannelMonitor.ts` + `modelMarketplaceMonitor.ts` from the marketplace view (view edits required; badges are a self-contained section).
- Skip REF's authenticated `/models` group in `user.go` — dead code in REF (frontend always uses `/public/models`).

---

## Option F — Frontend `utils/pricing.ts` conflict

REF overwrote HEAD's `formatScaled` (incompatible 3rd arg) and deleted `resolveIntervalPrices`, which HEAD's `PlazaModelPricingTable.vue`, `PricingRow.vue`, `SupportedModelChip.vue` and `intervalPrices.spec.ts` still use.
- **F1 (recommended): do NOT touch HEAD's `utils/pricing.ts`.** Add `DisplayCurrency`/`currencySymbol`/`normalizeDisplayCurrency` either as appended exports (new names don't clash) or in a new `utils/currency.ts`, and change the marketplace view import accordingly.
  *Behavior change*: none. If option A1 is taken (USD-only), even `currencySymbol` can be inlined as `'$'` and the new module skipped.

---

## Option G — Ent schemas vs raw SQL

REF commits ent schemas + generated code for `model_catalogs`/`model_vendors` but never uses them (repository is raw SQL).
- **G1 (recommended): port schemas + regenerate** (`make generate`) to keep Ent the source of truth, as CLAUDE.md requires for `ent/schema/*.go`. No upstream schema edges involved (catalog↔vendor only) → zero merge risk. Generated files are new files → no conflicts.
- **G2: skip Ent entirely**, rely on migrations only. Less generated-code churn on future upstream merges, but diverges from repo convention. If chosen, do not add `ent/schema/model_catalog.go`/`model_vendor.go`.

---

## Option H — Startup warm sync

REF `ProvideModelCatalogService` launches `go svc.SyncFromPricing()` at boot (2-min timeout, log-only failure).
- **H1: keep** if B1/B2 chosen — harmless, idempotent upsert.
- **H2: drop** and rely on the admin "Sync" button — one less goroutine touching upstream pricing state at startup; catalog starts empty until first manual sync. Recommended if B2 (own fetcher) is chosen.

---

## Summary recommendation matrix

| Coupled point | Recommended option | Residual upstream coupling | User-visible delta vs REF |
|---|---|---|---|
| Pricing-association SQL | A1 (+A3 optional) | table names/columns only | USD-only display |
| LiteLLM sync | B1 (or B2 for full facets) | PricingService struct (stable core fields) / none | fewer capability facets (B1) |
| Group models list | C1 (drop) | none | feature absent |
| groupRepo injection | D1 (remove) | none | none |
| Public routes/monitors | E1 | none | public monitor badges |
| utils/pricing.ts | F1 | none | none |
| Ent schemas | G1 | none | none |
| Warm sync | H1/H2 | — | none / manual first sync |
