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
