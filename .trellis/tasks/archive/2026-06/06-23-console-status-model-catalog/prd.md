# Optimize Console Navigation, Channel Status, and Model Catalog

## Goal

Improve the product navigation and model discovery experience by moving the console theme switch into the top navigation, promoting channel status into a standalone page linked from the homepage, and adding platform-wide model information management/display that is associated with channel pricing model definitions.

## What I Already Know

* The homepage already has a top navigation theme button in `frontend/src/views/HomeView.vue`.
* The console theme switch currently lives at the bottom of `frontend/src/components/layout/AppSidebar.vue`.
* The console header already has a right-side action area in `frontend/src/components/layout/AppHeader.vue`, making it the natural target for the console theme button.
* User channel status currently exists at route `/monitor`, implemented by `frontend/src/views/user/ChannelStatusView.vue`.
* Channel status data is already fetched from `frontend/src/api/channelMonitor.ts` using `/channel-monitors` and `/channel-monitors/:id/status`.
* The homepage footer status link currently points to `githubUrl`; it should point to the new/standalone status page.
* Channel pricing model definitions currently live under channel management. The admin channel page stores `model_pricing` entries with model names grouped by platform.
* User-facing available channels already expose supported model names and pricing in `backend/internal/handler/available_channel_handler.go` and `frontend/src/components/channels/AvailableChannelsTable.vue`.
* CoSphere has a model management page with vendor tabs, model filters, add/edit/delete, missing-model discovery, upstream sync, model icons, vendor metadata, tags, endpoints, groups, quota types, and bound channels.

## Requirements

* Move the console theme color switch from the sidebar bottom area into the top console header action row.
* Keep homepage and console theme behavior consistent: same persisted `localStorage.theme`, same `document.documentElement.dark` toggling, same light/dark icon semantics.
* Remove the duplicated sidebar theme switch so console theme control is available from top navigation.
* Add a standalone channel status page entry linked from the homepage top navigation and homepage footer status link.
* Reuse existing channel monitor data APIs and formatting helpers for the standalone status page.
* Improve the standalone status page presentation based on the current channel status page: preserve monitor hero/card/detail behavior while making the page appropriate for homepage/public navigation.
* Add platform-wide model information management and display capability inspired by CoSphere's model management/model marketplace.
* Associate model information with channel pricing model definitions, so pricing-configured model names can resolve to richer model metadata where available.
* Keep channel pricing as the source of billing configuration; model catalog metadata should enrich and organize model display, not silently change billing behavior.
* Use a persistent model catalog implementation, not a pricing-only derived view.
* Model vendor icons should reference the same icon concept used by CoSphere via LobeHub icons (`https://icons.lobehub.com/components/lobe-hub`).
* Initial model catalog data should be seeded/synchronized from this project's existing remote model information/pricing data, including names, parameters/capabilities where available, and pricing fields.

## Acceptance Criteria

* [ ] Console top header includes a theme toggle button visible in normal console pages.
* [ ] Sidebar no longer contains the theme toggle control.
* [ ] Homepage top navigation includes a status link/button pointing to the standalone status page.
* [ ] Homepage footer status hyperlink points to the standalone status page.
* [ ] The standalone status page reuses current channel monitor list/detail data and still supports refresh, window selection, and detail drilldown.
* [ ] Existing `/monitor` console route remains usable or redirects cleanly to the standalone page without breaking sidebar navigation.
* [ ] Admin users can manage model information across platforms.
* [ ] Users can browse/display model information grouped or filtered by platform/vendor/category.
* [ ] Channel pricing model names visibly associate with catalog model records when model metadata exists, while still handling unmatched names gracefully.
* [ ] Lint, type-check, and relevant frontend/backend tests pass or known failures are documented.

## Technical Approach

### Navigation and Theme

Create or reuse a small shared theme composable if needed, then wire it into both homepage and console header. Move the console button to `AppHeader.vue` and remove sidebar theme state/toggle code from `AppSidebar.vue`.

### Channel Status

Keep current monitor data contracts and components. Add a standalone route, likely `/status`, with a page layout suitable for top-level navigation. The current `/monitor` route can continue using `AppLayout` for authenticated console users or be adapted to reuse the standalone page internals.

Homepage top nav and footer should use Vue Router links for internal status navigation instead of external anchors.

### Model Catalog

Implement model metadata as a first-class catalog. The recommended direction is a persistent backend entity plus admin CRUD APIs, user/public read APIs, frontend admin management page, and user-facing model marketplace/display page. Channel pricing model names should link by normalized platform + model name or matching rules, with explicit unmatched states.

## Feasible Approaches

### Approach A: Persistent Model Catalog (Chosen)

Create backend model/vendor entities and CRUD/read APIs. Channel pricing remains unchanged but resolves model names against the catalog for richer display and association.

Pros: durable, administrable, scalable, closest to CoSphere, supports future sync/import and marketplace experiences.

Cons: larger backend/frontend implementation and database surface.

### Approach B: Derived Catalog from Channel Pricing

Build model display purely by aggregating existing channel pricing model names, with optional frontend-only enrichment.

Pros: smaller, fastest path, minimal backend migration risk.

Cons: weak "management" capability, no real metadata ownership, harder to evolve into a model marketplace.

### Approach C: Hybrid MVP

Add a persistent catalog table for metadata and read-side association, but defer vendor CRUD, sync wizard, missing-model tooling, and bulk operations.

Pros: establishes the correct data model with smaller first release.

Cons: less feature-complete than CoSphere reference, follow-up work needed.

## Out of Scope

* Changing pricing calculation semantics.
* Replacing existing channel pricing UI wholesale.
* Upstream official model metadata sync unless included by explicit decision.
* Public exposure of private channel/account details beyond existing channel status/available-channel data contracts.

## Research References

* [`research/cosphere-model-pages.md`](research/cosphere-model-pages.md) — CoSphere model management and model marketplace patterns worth adapting.

## Technical Notes

* Dirty worktree existed before this task. Existing modified files include frontend layout files that this task may need to touch; changes must be reviewed carefully before editing.
* Relevant frontend files: `frontend/src/views/HomeView.vue`, `frontend/src/components/layout/AppHeader.vue`, `frontend/src/components/layout/AppSidebar.vue`, `frontend/src/views/user/ChannelStatusView.vue`, `frontend/src/router/index.ts`, i18n locale files.
* Relevant backend files: channel/available-channel handlers and services, channel pricing service/domain types, Ent schema/migration patterns.
* Relevant CoSphere files: `/Users/liuliang/work/openS/CoSphere/web/src/components/table/models/*`, `/Users/liuliang/work/openS/CoSphere/web/src/pages/Model/index.jsx`, `/Users/liuliang/work/openS/CoSphere/web/src/pages/Setting/Model/*`, `/Users/liuliang/work/openS/CoSphere/web/src/pages/Setting/Ratio/*`.

## Decision (ADR-lite)

**Context**: The feature needs true model information management and marketplace-style display, and must create a business association with existing channel pricing model definitions.

**Decision**: Implement Approach A, a persistent model catalog with backend storage, admin management APIs, user-facing display APIs, frontend management/display pages, and pricing association by platform + normalized model name. Vendor icons should use LobeHub icon identifiers/metadata compatible with the CoSphere reference.

**Consequences**: This adds database/API surface but gives the project an extensible model metadata system. Pricing calculation remains owned by channel pricing; catalog metadata enriches browsing, management, and association views.
