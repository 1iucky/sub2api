# Directory Structure

> How backend code is organized in this project.

---

## Overview

<!--
Document your project's backend directory structure here.

Questions to answer:
- How are modules/packages organized?
- Where does business logic live?
- Where are API endpoints defined?
- How are utilities and helpers organized?
-->

(To be filled by the team)

---

## Directory Layout

```
<!-- Replace with your actual structure -->
src/
├── ...
└── ...
```

---

## Module Organization

### Convention: personalized-feature isolation (upstream fork)

This branch carries personalized features on top of `upstream/main`. A personalized feature MUST be a **self-contained vertical slice** so future `upstream/main` merges stay conflict-minimal:

- **Feature-owned new files only** for: `service/<feature>.go`, `repository/<feature>_repo.go`, `handler/<feature>_handler.go`, `handler/admin/<feature>_handler.go`, `handler/dto/<feature>.go`, `ent/schema/<feature>*.go`, `server/routes/<feature>.go` (or a block in `admin.go`), frontend `views/**`, `api/**`, locale modules.
- **Additive seams** on upstream-owned files are limited to small blocks: wire ProviderSets, `handler/handler.go` struct fields, one call line in `router.go`, route array elements, sidebar nav items, locale index imports/spreads.
- **Never add fields or methods to upstream-owned structs** (e.g. `PricingService`, `LiteLLMModelPricing`, `channel.go` consts, ent `group.go`). When the feature needs data upstream also parses, give the feature its OWN struct + fetcher.

**Canonical example** (model catalog, task `09-29-replay-model-marketplace-catalog`):

```go
// Wrong: re-adding 12 capability fields to upstream's LiteLLMModelPricing
// so catalog sync can read them from PricingService.pricingData.

// Correct: catalog-owned struct + fetcher; zero runtime coupling.
type catalogLiteLLMPricing struct { /* supports_*, max_*_tokens, costs... */ }
func (s *ModelCatalogService) listCatalogEntries(ctx context.Context) { /* fetch cfg.Pricing.RemoteURL, fallback cfg.Pricing.FallbackFile */ }
```

- **Public read-only APIs** for standalone public pages live in `backend/internal/server/routes/public.go` under `/api/v1/public/*`, registered via `RegisterPublicRoutes(v1, h)` from `router.go`. No auth middleware; fail-soft on the frontend.
- **Audit before committing a port**: `git diff HEAD -- <upstream-owned files>` must show additive-only hunks; `git grep` must show no resurrected symbols the port deliberately dropped (e.g. `display_currency`, `NormalizeDisplayCurrency`).

---

## Naming Conventions

<!-- File and folder naming rules -->

(To be filled by the team)

---

## Examples

<!-- Link to well-organized modules as examples -->

(To be filled by the team)
