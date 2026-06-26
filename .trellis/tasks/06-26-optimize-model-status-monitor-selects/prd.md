# Optimize model marketplace, public status, and monitor model selects

## Goal

Improve three frontend workflows: make the public model marketplace load progressively, simplify the public status page grouping, and align channel monitor model selection with the existing channel pricing remote-search behavior.

## What I already know

* `/models` is implemented in `frontend/src/views/user/ModelMarketplaceView.vue`.
* The marketplace currently fetches 48 models per page and renders pagination controls.
* The marketplace left filters use `space-y-*` / grid gaps inside each filter category.
* The initial marketplace loading skeleton only appears when `loading && models.length === 0`; initial timing needs to reliably show a loading state.
* `/status` is implemented in `frontend/src/views/StatusView.vue`.
* The status page currently groups by provider first, then by monitor group, and shows provider/platform count and group/monitor summaries.
* Channel monitor create/edit is implemented in `frontend/src/components/admin/monitor/MonitorFormDialog.vue`.
* Channel monitor currently preloads up to 100 catalog models when the dialog opens.
* Channel pricing already has a remote-search model selector in `frontend/src/components/admin/channel/PricingEntryCard.vue`; it does not load default model options, debounces keyword search, aborts stale requests, and uses `Select` with `filter-options=false`.

## Assumptions

* Marketplace lazy loading should use the existing `/public/models` pagination API with page size 20.
* Marketplace infinite loading should append pages 1, 2, 3... as the user scrolls near the end of the list.
* Search and platform filters remain server-side; vendor/capability/price/context filters remain client-side against the loaded pages.
* Public status should still allow text search, but platform category filtering and platform grouping should be removed.
* Channel monitor model search should query active model catalog entries by keyword and still allow manual model entry through the existing creatable behavior.

## Requirements

* `/models` left filter category item spacing is tightened by removing per-item gaps inside each filter category.
* `/models` shows a first-entry loading state while the first 20 models are being fetched.
* `/models` initially loads 20 models.
* `/models` loads the next 20 models when the user scrolls near the end of the rendered list, until no pages remain.
* `/models` avoids stale responses overwriting newer filter/search results.
* `/status` removes platform/category filter UI.
* `/status` displays monitor cards grouped directly by group name, not by provider/platform.
* `/status` removes right-side monitor quantity description text from grouping headers.
* `/status` reduces card padding, card gaps, group gaps, and supporting metric spacing so more cards fit in one viewport.
* Channel monitor primary and extra model selects do not fetch catalog data on dialog open.
* Channel monitor primary and extra model selects perform debounced fuzzy remote lookup only after the user types a keyword.
* Channel monitor model selects mirror channel pricing's loading, empty, abort, and stale-response handling pattern.

## Acceptance Criteria

* [ ] Opening `/models` triggers a visible loading state and requests page 1 with `page_size=20`.
* [ ] Scrolling near the bottom of `/models` appends page 2 with another 20 models without replacing page 1.
* [ ] `/models` no longer shows old pagination buttons for normal browsing.
* [ ] Changing `/models` search or platform resets loaded models and starts over from page 1.
* [ ] `/status` has no platform select/filter and no provider-first sections.
* [ ] `/status` shows one section per group, with more compact monitor cards.
* [ ] Channel monitor dialog opening does not call the model catalog list API solely to populate model dropdowns.
* [ ] Typing into channel monitor primary or extra model select searches the admin model catalog remotely and displays results.
* [ ] Existing manual model entry behavior remains available.
* [ ] Frontend lint/type-check or targeted tests pass.

## Definition of Done

* Implementation follows existing Vue Composition API and project `Select` patterns.
* Stale async requests are aborted or sequence-guarded.
* No unrelated dirty files are modified.
* Quality verification is run and results are recorded.

## Out of Scope

* Backend API changes.
* Reworking model marketplace filter semantics beyond the requested lazy loading.
* Changing monitor availability calculations or public monitor API shape.
* Changing channel pricing behavior.

## Technical Notes

* Relevant spec: `.trellis/spec/frontend/component-guidelines.md`, especially the remote search select contract.
* `Select.vue` emits `search` whenever dropdown search text changes and clears search text on close.
* Pricing entry remote search uses `adminModelsAPI.list(1, 30, { search, status: 'active', visibility: 'public', sort_by: 'model_id', sort_order: 'asc' }, { signal })`.
