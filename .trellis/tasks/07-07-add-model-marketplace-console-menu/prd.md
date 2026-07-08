# Add Model Marketplace Console Menu

## Goal

Add the model marketplace as a regular left-sidebar entry in the normal user console. Clicking the menu item should render the marketplace in the console content area while reusing the existing marketplace page content.

## What I Already Know

* The public model marketplace route exists at `/models` and renders `frontend/src/views/user/ModelMarketplaceView.vue`.
* The console uses Vue Router plus the shared `AppLayout`/`AppSidebar` navigation pattern.
* The user confirmed the implementation direction: reuse the existing marketplace content and hide only the public/home top navigation when shown inside the console.

## Requirements

* Add a normal user console route/menu item for the model marketplace.
* The menu item appears in the same left sidebar pattern as other normal user console menu items.
* Clicking the menu item displays the marketplace inside the right-side console content area.
* The marketplace content must be reused from the existing model marketplace implementation.
* Only the public marketplace top navigation/header should be omitted in the console version.
* No backend behavior changes are required.

## Acceptance Criteria

* [ ] A regular user can navigate to the model marketplace from the console sidebar.
* [ ] The route renders inside the console layout, not as the public standalone page.
* [ ] Marketplace cards, filters, search, pricing display, loading state, and empty/error states remain shared with the public marketplace page.
* [ ] The console rendering does not show the public/home top navigation bar.
* [ ] Existing public `/models` behavior remains unchanged.

## Definition of Done

* Frontend route and sidebar changes are implemented following existing patterns.
* Relevant i18n labels are present for English and Simplified Chinese.
* Targeted tests or existing frontend checks are run.
* No unrelated dirty files are modified or reverted.

## Out of Scope

* Changing marketplace API behavior.
* Changing pricing logic, currencies, or billing calculations.
* Redesigning the marketplace page content.
* Adding admin console entry points.

## Technical Notes

* Likely files: `frontend/src/router/index.ts`, `frontend/src/components/layout/AppSidebar.vue`, `frontend/src/views/user/ModelMarketplaceView.vue`, locale files.
* The existing public route should remain available at `/models`.
