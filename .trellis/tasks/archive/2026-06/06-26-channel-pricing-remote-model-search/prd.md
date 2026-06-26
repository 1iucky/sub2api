# fix channel pricing remote model search

## Goal

Optimize the channel pricing model selector so model lookup is driven by remote keyword search instead of a one-time 100-item preload. This prevents models outside the initial page from being unselectable when configuring channel model pricing.

## What I already know

* The user reports the current dropdown requests once with a default of 100 models.
* Keyword matching currently happens only within those 100 local options.
* The desired behavior is: opening the dropdown does not request models; typing a keyword triggers remote fuzzy search with a controlled result size.
* The existing admin model catalog API supports pagination and `search` filters.

## Requirements

* Channel pricing model selection must not fetch model catalog data merely because the dropdown opens or the card mounts.
* Typing in the selector search box must trigger a remote model catalog query.
* The query must search by model keyword only, without platform/type filtering.
* The result count should be bounded to keep the dropdown lightweight.
* Existing manual model tag entry and default pricing autofill behavior must continue to work.

## Acceptance Criteria

* [ ] Opening a model selector with an empty search does not call `/admin/models`.
* [ ] Typing a keyword calls `/admin/models` with `search=<keyword>` and a bounded `page_size`.
* [ ] Search results are not limited to the first locally preloaded 100 models.
* [ ] Selecting a result adds the model to the pricing entry and retains existing price autofill behavior.
* [ ] Frontend typecheck/build validation passes for the touched code.

## Out of Scope

* Backend API changes.
* Changing the LiteLLM pricing sync behavior.
* Reworking channel pricing data persistence.

## Technical Notes

* Main files found: `frontend/src/components/admin/channel/PricingEntryCard.vue`, `frontend/src/components/common/Select.vue`.
* Admin model list API: `frontend/src/api/admin/models.ts` calls `/admin/models` with `page`, `page_size`, and filters.
