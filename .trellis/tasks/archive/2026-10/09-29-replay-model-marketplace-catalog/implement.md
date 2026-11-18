# Implement: Replay model marketplace and admin model catalog

Source of truth for REF code: `git show codex/custom-theme-upstream-main-20260717:<path>`. Never checkout/reset onto REF; copy file contents via `git show` only.

## Ordered Checklist

### 0. Preflight
1. Confirm working tree state: parent task's uncommitted frontend files (HomeView/AppSidebar/AppHeader/LocaleSwitcher/useTheme/locales) are present and must be preserved; record `git status` snapshot.
2. Verify `frontend/src/composables/useTheme.ts` exists in working tree (PublicTopNav dependency). If missing, plan to port REF's version.
3. Verify no migrations numbered 239/240 exist on HEAD (`ls backend/migrations | tail`).

### 1. Backend: schema + migrations
4. Port `backend/ent/schema/model_catalog.go` + `model_vendor.go` from REF verbatim.
5. Create `backend/migrations/239_model_catalog.sql` (= REF `152_model_catalog.sql`) and `240_model_catalog_vendor_soft_delete_platform_cleanup.sql` (= REF 154).
6. Run `cd backend && make generate` (ent + wire after later steps; ent now, wire at step 12).

### 2. Backend: feature-owned code
7. Port `service/model_catalog_errors.go` verbatim.
8. Port `service/model_catalog.go` with edits: drop `groupRepo` field/param; drop any display-currency remnants.
9. Rewrite `service/model_catalog_pricing_sync.go`: catalog-local LiteLLM struct (keep `supports_*`, `max_*_tokens`, costs, `litellm_provider`, `mode`); own fetch (`cfg.Pricing.RemoteURL` → fallback file); port all pure mapping functions; `SyncFromPricing` upsert flow per REF.
10. Port `repository/model_catalog_repo.go` with edits: remove `display_currency` from SELECTs; remove `NormalizeDisplayCurrency` call.
11. Port `handler/model_catalog_handler.go`, `handler/admin/model_catalog_handler.go`, `handler/dto/model_catalog.go` verbatim.

### 3. Backend: additive seams
12. `repository/wire.go`: add `NewModelCatalogRepository` to ProviderSet.
13. `service/wire.go`: add `ProvideModelCatalogService(repo, cfg)` with startup warm goroutine (2-min timeout, log-only) + `wire.Bind(new(service.ModelCatalogRepository), ...)`; add to ProviderSet.
14. `handler/handler.go`: add `Handlers.ModelCatalog` + `AdminHandlers.ModelCatalog` fields.
15. `handler/wire.go`: provider params + both constructors in ProviderSet.
16. New `routes/public.go` (models + vendors + channel-monitors); call from `router.go` beside `RegisterModelPlazaRoutes`.
17. `routes/admin.go`: additive `registerModelCatalogRoutes` (9 routes) + call.
18. `cd backend && make generate` (wire_gen.go) and `go build ./...`.

Validation gate A: `cd backend && make test-unit` + `golangci-lint run` (or `make test`).

### 4. Frontend: feature-owned code
19. Port `api/models.ts` (drop `DisplayCurrency` import; USD inline), `api/admin/models.ts`, `api/publicChannelMonitor.ts`.
20. Port `views/user/modelMarketplaceMonitor.ts` + spec, `views/user/ModelMarketplaceView.vue` (USD inline; no `@/utils/pricing` import) + spec, `views/admin/ModelCatalogView.vue`.
21. Port `components/home/PublicTopNav.vue` + `PublicFooter.vue`.
22. Port i18n `locales/{en,zh}/models.ts`, `locales/{en,zh}/admin/models.ts`.

### 5. Frontend: additive seams
23. `router/index.ts`: +3 routes (`/models`, `/marketplace`, `/admin/models`).
24. `AppSidebar.vue`: +user `/marketplace` entry, +admin `/admin/models` entry (`hideInSimpleMode`), +`SparklesIcon` import if needed.
25. `HomeView.vue`: +2 `to="/models"` CTAs (REF :83/:178 equivalents) in the parent's restored structure; +footer marketplace link.
26. Locale index files (`{en,zh}/index.ts`, `{en,zh}/admin/index.ts`): +imports/spreads. `common.ts`: +`nav.modelMarketplace`/`nav.modelManagement`. `landing.ts`: +`home.footer2.links.modelMarketplace`.
27. `api/admin/index.ts`: +`adminModelsAPI` (keep HEAD's `cnProvidersAPI`/`pluginsAPI`).

Validation gate B: `cd frontend && pnpm run typecheck && pnpm run lint:check && pnpm test:run && pnpm run build`.

### 6. End-to-end verification
28. Fresh-DB boot: run backend against empty Postgres → migrations 239/240 apply → startup warm sync populates catalog (or logs failure) → `GET /api/v1/public/models` returns data.
29. Browser checks (report separately): public `/models` (filters, badges, USD prices), `/marketplace` embedded, admin `/admin/models` CRUD + vendor CRUD + sync button, homepage CTAs, sidebar entries, `/model-plaza` unaffected (toggle off → plaza entries still hidden).
30. `git grep` audit: no `display_currency`/`NormalizeDisplayCurrency` in ported code; `git diff upstream/main -- backend/internal/service/pricing_service.go backend/internal/service/channel.go backend/ent/schema/group.go frontend/src/utils/pricing.ts` is empty.
31. `git status` final review: only expected files touched; parent-task and unrelated working-tree files intact.

## Validation Commands

```bash
cd backend && make generate && make test-unit && golangci-lint run
cd frontend && pnpm run typecheck && pnpm run lint:check && pnpm test:run && pnpm run build
git diff --check
```

## Risky Files / Rollback Points

- **After each validation gate**: inspect `git diff --stat`; unexpected upstream-file edits = stop and re-scope.
- `AppSidebar.vue` / `HomeView.vue` / locale files: parent task has uncommitted edits — port as additive hunks on top of the WORKING TREE, never `git checkout` from REF.
- `service/wire.go` warm goroutine: if depguard/lint objects to config injection, move warm call to an explicit `Start`-style method invoked from `cmd/server` additively.
- If the remote pricing URL is unreachable in the build/verify environment, verify sync via `pricing.fallback_file` instead.
- Rollback: revert commits; feature-owned tables may remain.

## Follow-up Checks Before Completion

- Migration checksum validation passes on boot (fresh + at-238 DB).
- Simple mode: admin catalog entry hidden; public marketplace still reachable (REF parity — flag in report).
