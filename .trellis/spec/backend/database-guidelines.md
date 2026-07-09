# Database Guidelines

> Database patterns and conventions for this project.

---

## Overview

<!--
Document your project's database conventions here.

Questions to answer:
- What ORM/query library do you use?
- How are migrations managed?
- What are the naming conventions for tables/columns?
- How do you handle transactions?
-->

(To be filled by the team)

---

## Query Patterns

<!-- How should queries be written? Batch operations? -->

(To be filled by the team)

---

## Migrations

<!-- How to create and run migrations -->

(To be filled by the team)

---

## Naming Conventions

<!-- Table names, column names, index names -->

(To be filled by the team)

---

## Common Mistakes

<!-- Database-related mistakes your team has made -->

(To be filled by the team)

---

## Scenario: Model Catalog Vendor Soft Delete And Public Deduplication

### 1. Scope / Trigger
- Trigger: model catalog work that changes `model_vendors`, `model_catalogs`, public marketplace listing, or pricing/monitor association behavior.
- This is cross-layer because persisted catalog rows are synchronized from pricing data, exposed through public APIs, and rendered by marketplace/status UI.

### 2. Signatures
- DB: `model_vendors.deleted_at TIMESTAMPTZ NULL` means the vendor was intentionally removed by an admin.
- DB: `model_catalogs.platform` is a protocol platform and must be one of `openai`, `anthropic`, `gemini`, `antigravity`.
- Service: `ModelCatalogListFilters.DeduplicateByID` requests public-list deduplication by `normalized_model_id`.

### 3. Contracts
- Admin vendor delete is soft delete: set `deleted_at`, keep the row for sync memory.
- Vendor list APIs must exclude `deleted_at IS NOT NULL`.
- Pricing sync must not automatically resurrect a vendor whose `provider_key` exists in a soft-deleted row.
- Manual vendor upsert by name is allowed to clear `deleted_at`; this is an intentional restore.
- Public marketplace association keys are model IDs only. Do not include provider/platform/channel platform in monitor matching logic.
- Public marketplace lists should return at most one row per `normalized_model_id`.

### 4. Validation & Error Matrix
- Delete unknown or already soft-deleted vendor -> `MODEL_VENDOR_NOT_FOUND`.
- Empty vendor name -> `MODEL_VENDOR_NAME_REQUIRED`.
- Empty model ID -> `MODEL_ID_REQUIRED`.
- Invalid/unknown catalog platform -> normalize to `openai` unless it maps to `anthropic`, `gemini`, or `antigravity`.

### 5. Good/Base/Bad Cases
- Good: deleting `DeepSeek` keeps its soft-deleted row; later pricing sync updates DeepSeek models with `vendor_id = NULL` and does not recreate `DeepSeek`.
- Base: public `/models` with two catalog rows for `glm-5.2` returns one marketplace row.
- Bad: hard-deleting `model_vendors` lets the next pricing sync recreate the same vendor, undoing the admin action.

### 6. Tests Required
- Unit or integration test that a deleted vendor is hidden from lists and sync does not upsert it automatically.
- Regression test that marketplace model deduplication keeps one row per normalized model ID.
- Type/lint checks for frontend consumers whenever public model response shape changes.

### 7. Wrong vs Correct

#### Wrong
```sql
DELETE FROM model_vendors WHERE id = $1;
```

This loses the admin's deletion intent, so automated sync can recreate the vendor on restart.

#### Correct
```sql
UPDATE model_vendors
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
```

Keep the row as a tombstone, filter it from user-facing lists, and let only manual upsert restore it.

---

## Scenario: Announcement Category Contract

### 1. Scope / Trigger
- Trigger: changes to announcement categories, announcement create/update APIs, user announcement responses, or announcement storage.
- This is cross-layer because category is persisted in `announcements`, validated by service/admin API, serialized to admin and user DTOs, and rendered by frontend admin forms and user notification drawers.

### 2. Signatures
- DB: `announcements.category VARCHAR(32) NOT NULL DEFAULT 'announcement'`.
- DB constraint: category must be one of `announcement`, `model_update`, `changelog`.
- Service constants: expose the same three category values from the announcement service/domain layer.
- Admin create request: optional `category` string.
- Admin update request: optional pointer `category` string.
- Admin/user response DTOs: include `category` string.

### 3. Contracts
- Empty category on create defaults to `announcement`.
- Existing rows must remain valid through the migration default.
- User announcement listing must not add new permission logic for category; it groups only the announcements already visible after existing targeting, schedule, and read-state filtering.
- Frontend should defensively treat missing or unknown response category as `announcement` during rollout compatibility.
- Category currently controls presentation grouping and labels only. Do not add billing, permission, or delivery semantics to this field without a new contract.

### 4. Validation & Error Matrix
- Create with empty category -> persist `announcement`.
- Create/update with `announcement`, `model_update`, or `changelog` -> accepted.
- Create/update with any other value -> `ErrAnnouncementInvalidCategory`.
- API binding with invalid category -> bad request before service mutation.
- Legacy response with missing category -> frontend displays under `announcement`.

### 5. Good/Base/Bad Cases
- Good: add a new category by updating constants, service validation, DB check constraint/migration, DTO tests, frontend union type, i18n labels, admin select options, and drawer grouping together.
- Base: old announcement rows created before the migration appear under the Announcement tab.
- Bad: adding only a frontend tab without backend validation and DB migration; admins could save values the service rejects or users could receive ungroupable data.

### 6. Tests Required
- Service create default test asserts omitted category becomes `announcement`.
- Service invalid category test asserts `ErrAnnouncementInvalidCategory`.
- Service update test asserts category can change to a valid value.
- DTO tests assert admin and user announcement responses include category.
- Frontend component test asserts drawer tabs show category counts based on already-loaded visible announcements and no `all` tab appears.
- Frontend type-check must pass after response/request type changes.

### 7. Wrong vs Correct

#### Wrong
```go
if input.Category != "" {
    a.Category = input.Category
}
```

This stores arbitrary category values and makes frontend grouping depend on unchecked strings.

#### Correct
```go
category := strings.TrimSpace(input.Category)
if category == "" {
    category = AnnouncementCategoryAnnouncement
}
if !isValidAnnouncementCategory(category) {
    return nil, ErrAnnouncementInvalidCategory
}
a.Category = category
```

Validate at the service boundary, keep the DB constraint aligned, and let the frontend group only known categories.
