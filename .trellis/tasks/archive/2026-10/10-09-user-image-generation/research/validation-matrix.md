# Final validation matrix

Worktree: `/Users/liuliang/.codex/worktrees/user-image-generation/sub2api`, branch `codex/user-image-generation`, starting from `f59ded0eb`. No production deployment or real upstream image call was performed. Original master checkout and unrelated homepage/guideline/tests/image archive remain preserved.

| Layer | Evidence | Result / boundary |
| --- | --- | --- |
| Frontend regressions | Root full focused Vitest command including original KeysView, shared form, studio view/state/API/controls/history/download/settings/i18n | **71 passed**, 10 test files |
| Frontend static | Owner final targeted ESLint, TypeScript; root build's vue-tsc | Passed |
| Frontend production build | Root final `pnpm run build`, `/tmp/sub2api-image-studio-final-build.log` | Passed; latest backend-code aliases found in built studio chunk, MIT notices bundled |
| Backend final feature / old async / lifecycle | Owner `go test -tags=unit ./internal/service ./internal/handler ./internal/handler/admin ./internal/server/routes ./cmd/server -run 'TestImageStudio|TestAsyncImage|TestProvideCleanup' -count=1` | Passed; routes/admin compiled, service10.067s, handler7.836s, cmd8.473s |
| Backend build / generation | Final Wire generation and `go build ./cmd/server`; storage Ent generation | Passed; generated sources kept with standalone schemas |
| Existing auth / allowlist / billing regression | Root original targeted command plus final whole middleware unit package | Passed; all middleware tests2.465s |
| Existing subscription billing / client request ID / image tier | Root service regression command | Passed6.530s; subscription command multiplier preserved, separate from Studio subscription end-to-end |
| Mock full generation | Both providers, image/token modes, real APIKeyAuth and allowlist, original Images handler and Billing.Apply with mutable balance stub | Passed; one upstream/usage/debit, duplicate409/reuse and JWT-only owner reads; upstream/repositories are mocks |
| Download HTTP boundary | Studio content handler test | Exact bytes/MIME/length/disposition for supported formats, JWT401/owner404/deleted410; real encoded decoding covered separately |
| Managed image storage | Real image limits/SSRF/root/confined file/late lease/readonly startup regressions | Passed;16 parallel initializers repeated20 times |
| PostgreSQL history/claims | Real PostgreSQL15 testcontainers integration | Passed; concurrent unique submission, owner isolation, cleanup retries/leases/crash recovery, exact original fee snapshot and metadata preservation |
| Complete cleanup service against DB/files | Additional root-requested real PostgreSQL integration | Passed on real PostgreSQL15: concurrent leader exclusion/release, disabled eligible asset preserved, exactly1 image/75 bytes removed and all history metadata retained |
| Independent Trellis source check | `research/check-review.md`, Go vet, formatting, diff check | No remaining blocking source findings; golangci-lint unavailable |
| Browser rendered UI | Native IAB + isolated API fixture | Both providers, no automatic Key creation, shared form preselection/cancel, previews/history, cleaned state, Chinese/English, themes and390px/desktop layout observed |
| Browser saved download file | JWT content fetched, Blob/anchor download requested | Native IAB download event timed out; filesystem completion unconfirmed; separate HTTP/component tests passed |
| Real upstream / production billing | No live credentials used | Unverified; real PostgreSQL Studio subscription-debit end-to-end also not run |

## Acceptance mapping

- AC1/AC2/AC4/AC5/AC6/AC7: frontend/API/gateway/provider/form tests and browser fixture, with the browser-download boundary above.
- AC3: original middleware/handler/Billing.Apply exercised under standard balance mode for both providers and pricing modes; original subscription multiplier regression passed. Do not promote mocked debit evidence into production settlement evidence.
- AC8/AC10/AC11: real PostgreSQL owner/history/cost/asset claims, confined real files, exact content endpoint tests and cleanup-state UI. Expiration changes file availability and asset cleanup state only, not original historical information.
- AC9: explicit implementation authorization and low-conflict requirement recorded before worktree implementation. Context manifests validate10 implement entries and12 check entries.
- AC12: scoped submission uniqueness/conflict, no frontend POST replay, unknown recovery/lease tests, per-asset storage failure and original usage cost linkage.

## Merge footprint and preserved work

Public synchronous/Redis-async Images handlers and old uploader remain unchanged. Existing gateway/user/admin route files each add one feature registration; aggregate handlers add3 lines and provider wiring7 lines; cmd/server wire adds3 lines. Feature settings avoid expanding the upstream settings DTO. Standalone Ent schema edges do not touch upstream user/group/Key entities. Regenerate Ent/Wire when merging instead of resolving generated conflicts manually. The necessary shared Key-form extraction remains the largest handwritten existing-file change, covered by the original26 Key-view regressions.

Unrelated copied `.trellis/spec/frontend/component-guidelines.md` and local frontend dependency symlink are excluded from feature staging. See `docs/IMAGE_STUDIO.md` for enabling, storage, retention, rollback and upstream maintenance.

## Nonblocking follow-up

Uncertain submission with no matching record keeps the page's generation button disabled until read-only recovery succeeds or the user leaves/reopens the page. This conservative gate intentionally avoids implicit replay; a separately acknowledged new-submit control could improve usability later.
