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

### How migrations run

- `backend/migrations/*.sql` is an embed FS; every `*.sql` file auto-loads at boot (`migrations/migrations.go`).
- **Checksum validation**: applied files are hashed. NEVER edit a migration that has shipped — add a new file instead.
- Migration identity is the full filename (`schema_migrations.filename`), not the numeric prefix. Multiple distinct filenames may share a prefix. New migrations should sort after the current maximum; published migrations must keep their names and contents.

### Convention: personalized-feature migrations on this fork

This branch replays personalized features onto `upstream/main`. When porting a feature whose migrations occupy slots already used upstream:

1. Find HEAD's max prefix: `ls backend/migrations | tail`.
2. For migrations that have never shipped on THIS fork, assign filenames after the maximum (e.g. REF `152_model_catalog.sql` + `154_...sql` became `239_` / `240_` when HEAD was at 238).
3. Merge multi-step REF migrations only when the intermediate state never shipped on THIS branch (final-state DDL is enough).

**Why**: checksum validation applies to the full filename. Distinct upstream/fork filenames with the same prefix do not collide. Renaming a published migration instead makes it appear unapplied and replays its SQL, which can rerun data cleanup. Preserve published fork `239_model_catalog.sql` / `240_model_catalog_vendor_soft_delete_platform_cleanup.sql` alongside upstream's distinct `239_` / `240_` files.

### Convention: Ent schema + raw-SQL repository coexistence

`ent/schema/*.go` is the schema source of truth (CLAUDE.md), but a repository may still query with raw `*sql.DB` (e.g. `repository/model_catalog_repo.go` reads `channel_model_pricing` with hand-written SQL). Rules:

- Always add the Ent schema and run `make generate`; commit generated code.
- Raw-SQL repositories must only read columns that exist on HEAD's migrations — when porting, re-verify every selected column against the current `channel_repo_pricing.go` column list (upstream adds columns like `cache_write_1h_price` over time).
- Edges to upstream-owned ent schemas are merge-conflict risks; keep feature schemas edge-free or edge-internal only (catalog→vendor with `OnDelete(SetNull)` is the accepted example).

---

## Naming Conventions

<!-- Table names, column names, index names -->

(To be filled by the team)

---

## Common Mistakes

<!-- Database-related mistakes your team has made -->

(To be filled by the team)
