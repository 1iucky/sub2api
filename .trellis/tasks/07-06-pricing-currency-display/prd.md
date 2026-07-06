# Add Pricing Currency Display for Channel Model Prices

## Goal

Allow admins to choose a display currency per channel model pricing entry so model marketplace prices and usage-record cost displays show either USD `$` or CNY `￥`. This setting is display-only for visible currency symbols and must not change billing, stored numeric price semantics, exchange rates, or cost calculations.

## What I Already Know

* User requirement: channel pricing management should let each model pricing entry choose USD or RMB; the choice only controls model marketplace price symbols.
* Current channel pricing data lives in `channel_model_pricing`.
* Current service/domain type is `service.ChannelModelPricing`.
* Current admin channel API request/response uses `channelModelPricingRequest` and `channelModelPricingResponse`.
* Current frontend admin pricing form uses `PricingFormEntry` and `PricingEntryCard.vue`.
* Current model marketplace uses `related_pricing.entries` from public model catalog responses.
* Marketplace price formatting is currently hard-coded to `$` in `formatUSDPerMillion`.
* User and admin usage record cost displays are also hard-coded to `$` in `frontend/src/views/user/UsageView.vue` and `frontend/src/components/admin/usage/UsageTable.vue`.
* Usage log DTOs currently return numeric costs but no currency/display-currency field.
* Usage logs already store `channel_id`, `model`, `requested_model`, and billing metadata, which can be used to resolve or snapshot a display currency.
* Existing billing code reads the numeric price fields from `ChannelModelPricing`; the new field must not be used there.

## Requirements

* Add a per-pricing-entry display currency with allowed values `USD` and `CNY`.
* Default existing and omitted values to `USD` for backward compatibility.
* Add the control to the channel pricing admin form for create and edit flows.
* Persist the value with the channel model pricing entry.
* Include the value in admin channel API responses and public model catalog `related_pricing.entries`.
* Update model marketplace price rendering so `USD` uses `$` and `CNY` uses `￥`.
* Include a display currency on user and admin usage record responses.
* Update user and admin usage record cost columns and cost tooltip amounts so `USD` uses `$` and `CNY` uses `￥`.
* Existing usage records without a display currency must render as USD `$`.
* Keep price numbers unchanged; do not convert between USD and CNY.
* Do not use display currency in billing, account stats pricing, gateway pricing resolution, or balance deduction.

## Acceptance Criteria

* [x] Admin can set each channel model pricing entry to USD or CNY.
* [x] Existing pricing entries with no value behave as USD.
* [x] Creating/updating a channel preserves the selected display currency.
* [x] Model marketplace cards show `$` for USD entries and `￥` for CNY entries.
* [x] Model marketplace drawer and interval price display use the same symbol as the selected entry.
* [x] User usage records show cost column and cost detail tooltip with the record's display currency symbol.
* [x] Admin usage records show cost column, account-billed row, and cost detail tooltip with the record's display currency symbol.
* [x] Usage records without display currency continue to show `$`.
* [x] Billing/cost calculation tests or code paths remain unaffected by the display currency field.
* [x] Tests cover API/type propagation and marketplace symbol formatting.

## Proposed Design

### Backend

* Add `display_currency VARCHAR(3) NOT NULL DEFAULT 'USD'` to `channel_model_pricing`.
* Add `DisplayCurrency string` to `service.ChannelModelPricing`.
* Accept and return `display_currency` in admin channel pricing request/response.
* Load/save the column in channel repository list/create/update/replace flows.
* Include `display_currency` in model catalog pricing association entries returned to the marketplace.
* Validate or normalize values to `USD` / `CNY`, defaulting blank values to `USD`.
* Add a usage-record display currency, preferably as a `usage_logs.display_currency` snapshot set when a usage log is written.
* Populate usage log `display_currency` from the matched channel model pricing entry when available; otherwise default to `USD`.
* Return usage log `display_currency` in both user and admin usage DTOs.
* Treat usage log display currency as an immutable request-time snapshot; later admin changes to channel model pricing currency must not alter historical usage record symbols.

### Frontend Admin

* Add `display_currency: 'USD' | 'CNY'` to channel pricing types and form entries.
* Add a compact select in `PricingEntryCard.vue` near billing mode.
* Replace static form unit labels with the selected symbol: `$/MTok` or `￥/MTok`, `$` or `￥`.
* Preserve the value through form-to-API and API-to-form transforms.

### Model Marketplace

* Add `display_currency` to `ModelPricingAssociationEntry`.
* Replace `formatUSDPerMillion` with a display-currency-aware formatter.
* Use the representative pricing entry currency for card prices.
* Attach the pricing entry currency to aggregated interval rows so drawer interval prices use the correct symbol.

### Usage Records

* Add `display_currency` to frontend `UsageLog` / `AdminUsageLog` types.
* Add shared frontend currency-symbol helpers so usage pages and model marketplace do not duplicate `$` / `￥` mapping.
* Replace hard-coded `$` prefixes in:
  * `frontend/src/views/user/UsageView.vue`
  * `frontend/src/components/admin/usage/UsageTable.vue`
  * `frontend/src/utils/usagePricing.ts`
* Keep dashboard/stat aggregate cards out of scope unless they are per-record rows; aggregated totals may contain mixed currencies and should continue existing display until a separate aggregate-currency design exists.

## Decision (ADR-lite)

**Context**: The field is explicitly display-only, but it originates in channel pricing and must be shown in the public model marketplace and per-record usage cost displays.

**Decision**: Store a separate `display_currency` field on `channel_model_pricing`, propagate it to marketplace association responses, snapshot it onto usage records at request time, and use it only in frontend formatting.

**Consequences**: Existing billing math remains stable. The database/API schema grows by one channel-pricing field and one usage-log display snapshot field. Historical usage records keep the symbol from the request time even if admins later change model pricing currency. Future real multi-currency billing would need a separate design because this field is intentionally not monetary semantics.

## Out of Scope

* Currency conversion.
* Changing stored price units or billing semantics.
* Changing payment currencies.
* Applying currency to account stats pricing unless it is needed only for shared type compatibility.
* Recomputing or converting aggregate usage/dashboard totals across mixed currencies.
* Supporting currencies beyond `USD` and `CNY`.

## Technical Notes

* Relevant backend files inspected:
  * `backend/internal/service/channel.go`
  * `backend/internal/repository/channel_repo_pricing.go`
  * `backend/internal/repository/model_catalog_repo.go`
  * `backend/internal/handler/admin/channel_handler.go`
  * `backend/internal/service/model_pricing_resolver.go`
  * `backend/internal/service/account_stats_pricing.go`
  * `backend/internal/repository/usage_log_repo.go`
  * `backend/internal/handler/dto/types.go`
  * `backend/internal/handler/dto/mappers.go`
  * `backend/internal/service/gateway_service.go`
* Relevant frontend files inspected:
  * `frontend/src/components/admin/channel/types.ts`
  * `frontend/src/components/admin/channel/PricingEntryCard.vue`
  * `frontend/src/views/admin/ChannelsView.vue`
  * `frontend/src/api/admin/channels.ts`
  * `frontend/src/api/models.ts`
  * `frontend/src/views/user/ModelMarketplaceView.vue`
  * `frontend/src/views/user/UsageView.vue`
  * `frontend/src/components/admin/usage/UsageTable.vue`
  * `frontend/src/utils/usagePricing.ts`
