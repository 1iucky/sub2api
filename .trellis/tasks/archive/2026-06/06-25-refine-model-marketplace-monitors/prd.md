# Refine Model Marketplace and Catalog Monitors

## Goal

Improve model catalog, model marketplace, and public status behavior so production data stays stable after restarts, marketplace model availability is linked by model ID only, and the UI matches the console's polished form/list patterns.

## What I Already Know

- The public model marketplace consumes `/public/models`, `/public/models/vendors`, and public channel monitor data.
- Model monitor matching already exists in `frontend/src/views/user/modelMarketplaceMonitor.ts` and matches primary and extra monitor models by normalized model ID.
- Model catalog sync currently creates or updates vendors from pricing providers during `SyncFromPricing`, which can recreate vendors that an admin deleted.
- Model catalog rows are unique by `(platform, normalized_model_id)`, so the same model ID can appear more than once in public lists if synced through multiple providers/platform values.
- Pricing associations are matched through channel pricing model names and are required for public marketplace visibility.
- Existing unrelated dirty files are `deploy/docker-compose.yml` and `sub2api.tar`; they are outside this task.

## Requirements

- Deleted model vendors must not be automatically restored after production Docker Compose restarts or pricing sync runs.
- Status and model marketplace lists must show loading feedback while data is being loaded or refreshed.
- Model marketplace search must move from the left sidebar into the list's top filter area and use the console's focused input styling.
- Model marketplace left filters must be denser, vendor filter height must be taller, and vendors must have an independent search box.
- Model marketplace cards must be smaller, render three per row on wide screens, preserve content, and float subtly on hover.
- Model marketplace cards must show only the monitor timeline for availability, without the textual default availability or availability percentage copy.
- Model management model platform selection must only expose the account protocol platforms: `openai`, `anthropic`, `gemini`, `antigravity`.
- Existing model catalog platform values outside those four platform values must be cleaned/normalized.
- Public marketplace models must be deduplicated by model ID.
- Marketplace monitor timeline association must depend only on model ID and must match either a channel monitor primary model or any extra model.

## Acceptance Criteria

- [ ] Deleting a vendor keeps it hidden/deleted after service restart and pricing sync; manual creation can intentionally add a vendor again.
- [ ] Public model marketplace does not show duplicate cards for the same normalized model ID.
- [ ] Marketplace cards for models configured as primary or extra models in channel monitors show the matching monitor timeline.
- [ ] Marketplace list and status list show loading state on first load and refresh.
- [ ] Marketplace search/filter/card layout matches the requested UX changes and remains responsive.
- [ ] Add/edit model platform selection only offers `openai`, `anthropic`, `gemini`, and `antigravity`.
- [ ] Existing invalid model catalog platform values are migrated to an allowed value.
- [ ] Existing frontend typecheck/lint and relevant backend tests pass or any blockers are documented.

## Out of Scope

- Changing channel pricing business rules beyond model catalog association and marketplace visibility.
- Changing channel monitor execution/checking behavior.
- Reworking the full visual design of the home page, console navigation, or legal pages.

## Technical Notes

- Backend files likely impacted:
  - `backend/internal/service/model_catalog.go`
  - `backend/internal/repository/model_catalog_repo.go`
  - `backend/migrations/*`
- Frontend files likely impacted:
  - `frontend/src/views/user/ModelMarketplaceView.vue`
  - `frontend/src/views/StatusView.vue`
  - `frontend/src/views/admin/ModelCatalogView.vue`
  - `frontend/src/views/user/modelMarketplaceMonitor.ts`
- Cross-layer contract: model catalog platform is a protocol/platform field with only four allowed values. Marketplace association key is normalized model ID, not provider, monitor platform, or channel platform.
