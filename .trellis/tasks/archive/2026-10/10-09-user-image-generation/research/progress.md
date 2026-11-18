# Implementation ledger — plan: .trellis/tasks/10-09-user-image-generation/implement.md

User approved implementation after final summary. Expired image files deleted only; all metadata retained. New requirement: minimize upstream/main merge conflict footprint.

Worktree created from current HEAD f59ded0eb on codex/user-image-generation; original checkout and dirty business files preserved. Only modified guideline copied as development context, must not be committed.

## Preflight and ownership
| Tasks | Shared output / input | Resolution |
| --- | --- | --- |
| Form + studio UI | created(ApiKey), current group | Single frontend owner; shared actual form |
| Settings + persistence | cleanup policy and provider DTO | Narrow policy interface; coordinate before use |
| Persistence + gateway | submission/history/result methods | Publish service types first, agree DTO in implementation-contracts |
| Gateway + UI | config/history endpoints, canonical Images payload | Shared contracts, route/key and JWT clients separate |
| All backend | Wire, routes, Ent generated | Backend owner Wire/routes; storage owner Ent generated; no concurrent generation |
| Cleanup + history | asset states and metadata retention | Storage owner both, deletes files only |
| Each implementation area | tests agree scope | Key no automatic creation; persistence server-owned; original billing only; local asset paths stable |

## Work queue
- Storage/history/local assets/cleanup: implemented; unit and real PostgreSQL evidence in storage-validation.md; readonly startup fix verified; full cleanup service PostgreSQL integration passed
- Backend settings/capabilities/gateway/JWT/admin integration: implemented; standard-mode auth + original billing debit covered; final serial feature/async/lifecycle tests and server build passed
- Frontend shared Key form/studio/history/settings/notices: implemented; 71 final focused tests passed; final lint/typecheck/build passed
- Review, meaningful tests, browser/integration, merge-maintenance docs: independent Trellis review complete, no blocking source findings; browser fixture evidence and deployment/merge docs written

## Baseline evidence
Original checkout KeysView.spec.ts: 26/26 passed. Original handler TestAsyncImage: passed. Docker daemon available. No real upstream credentials used. Browser fixture isolated at127.0.0.1:18080, Vite at127.0.0.1:13000; this environment only verifies UI against simulated API data.

## Implementation decisions
- Keep existing Images public-model admission: generic aliases such as portrait are currently rejected before upstream mapping. Support valid public image-family aliases and mappings without widening the shared parser, which reduces upstream conflict footprint.
- Price preview exposes existing resolved mode/effective multiplier as estimated; no fabricated total price. Actual costs remain original usage-derived snapshots.

## Review fixes and additional evidence
- Shared retention limit3650 used in both settings and cleanup, matching UI bounds.
- Root binding published atomically;16 parallel initializers x20 test repetitions pass.
- Registered crash files and genuinely revoked-lease late writes are cleaned without losing historical metadata.
- Safe existing readonly storage binding is reused without permission mutation; initialization failure degrades studio availability without aborting baseline server startup.
- Invalid history filters now return client400 and normalize submission UUID before any query.
- Unavailable default has a stable reason; UI maps stable error categories to English/Chinese without exposing raw upstream messages.
- Grok opaque mapped aliases supported where real gateway accepts them; conflicting/non-image mapped profiles fail closed. OpenAI keeps its original public-name parser constraints.
- Existing Handler aggregate field shortened to StudioGateway to avoid gofmt realigning all existing fields in two hot upstream files.
- Root original middleware/auth/allowlist/image-balance/key-billing regression command passed: middleware3.423s,handler4.309s.
- Standard-mode full Images tests now use real APIKeyAuth and original billing Apply; both providers and image/token modes debit the stub balance once, duplicate submission/history GET never debit again. Production billing still unverified.
- Browser fixture verified zero automatic Key creation, group-specific controls, both providers, JWT previews/history and cleanup tombstones,390px/light/dark layouts. Download bytes requested successfully, but native IAB saved-download event unavailable; do not claim browser filesystem completion.

## Final gate
Root full frontend suite71/71, production build, built latest code aliases and full notices passed. All middleware package and original subscription/client-ID/image-tier regressions passed. Independent source review, final backend feature/async/lifecycle tests/build, readonly/parallel initialization, four history/claim DB tests plus complete RunOnce leader/disabled/image-only DB/file test passed. Real upstream/production settlement and native IAB download filesystem completion remain explicitly unverified. Code/source frozen; feature staging excludes copied unrelated guideline and dependency symlink.
