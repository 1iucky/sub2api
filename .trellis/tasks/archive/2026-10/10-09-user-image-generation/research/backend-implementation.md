# Backend implementation evidence

Worktree: `/Users/liuliang/.codex/worktrees/user-image-generation/sub2api`. No commits or staging by the backend worker.

## Ownership and integration

Added feature-local settings/capability services, gateway/JWT/admin handlers, route helpers and provider sets. Existing gateway, user/admin routes and aggregate handlers have additive registrations only. Public Images synchronous handlers and the Redis async handler are unchanged. Studio reuses the latter's request validation, security audit, context copy and complete OpenAI/Grok execution dispatcher; the Studio sink is durable and does not call the old object uploader.

`cmd/server/image_studio.go` resolves `setup.GetDataDir()` above the repository layer (avoids setup/repository import cycle). If asset-store initialization fails it logs a warning and returns nil: config reports storage unavailable and Studio accepts no new execution; the existing gateway can still start. Wire and its existing cleanup-constructor regression are updated. Cleanup starts via the feature provider and stops in application cleanup.

## Canonical HTTP details

- JWT `/api/v1/user/image-studio/config`, `/history`, `/history/:id`, `/history/:id/images/:asset_id/content?download=1` enter the existing user JWT/audit/panel-rate-limit chain. History/content require only the owning user; no Key or current feature/balance gate.
- API Key `POST /v1/images/generations/studio` enters the existing body-limit/request-ID/ops/endpoint/APIKeyAuth/model-allowlist/composite/group chain. JSON only, 64 KiB request maximum, UUID `X-Image-Studio-Submission-ID`, dynamic field validation, nonstreaming output. Accepted response is raw HTTP202 with `history_id`, `submission_id`, `status`, `poll_url`.
- Admin GET/PUT `/api/v1/admin/settings/image-studio` retains the existing admin auth and audit middleware. Settings store one dedicated key, default disabled; default group must be nonzero/configured when enabled. Retention1–3650 days, five-field cron, valid time zone. GET additionally reports `supported_platforms` and `last_cleanup`.
- Config preserves `group_ids` order, rechecks the original user group authorization, current active/image permission and schedulable accounts. Configured defaults unavailable to this user report `default_unavailable_reason=default_group_unavailable`.
- Models have `binding_id=<group_id>:<model_id>`, capability fields `enum`/`integer`, GPT size/quality/background/format/compression and Grok resolution/ratio, with n1–4. Aliases resolve through account mapping and original group allowlist. Arbitrary OpenAI aliases remain excluded because the unchanged Images parser rejects them; Grok opaque image aliases are supported. Accounts with incompatible mapped parameter profiles for the same public model cause that model to be omitted.
- `billing_preview` contains `mode`, `effective_rate_multiplier`, `estimated`, `image_rate_independent`, `resolved_image_multiplier`, `effective_token_multiplier`, `observed_at`. Existing gateway channel pricing and peak/image multiplier functions own the decision; no guessed final price or separate debit is introduced.

## Execution and errors

A server-generated correlation UUID replaces the worker's `ctxkey.ClientRequestID`; history associates original owner+Key with `client:<uuid>`, exactly as original image billing. A hash of canonical request plus Key/group prevents dispatching the same submission twice or accepting changed content under that UUID. Worker URL and inbound endpoint are canonical `/v1/images/generations`.

The response recorder is capped128MiB; heartbeat every30s extends a2min lease. Generation runs with30min timeout and result persistence with5min timeout. Interruptions/panics/timeouts are never replayed; history reconciliation marks expired work unknown. Worker non2xx failures retain allowlisted actionable original code/type (balance, quota/rate, auth, permission, billing/concurrency), with safe fallback; upstream diagnostic strings are never copied into history. Boolean capability values are checked and unknown field kinds fail closed.

Download replies stream owner-scoped stored bytes with correct MIME/length, no-store, nosniff and MIME-aware disposition. Missing/foreign records404, cleaned410, temporarily unavailable503.

## Verification and boundaries

Tests intentionally failed before settings/providers existed; later regression REDs reproduced retention3651 acceptance, missing compression, missing default/preview/error contracts, Grok opaque alias exclusion and configured group order drift. Each was implemented and rerun.

`TestImageStudioGatewayFullImagesHandlerAndIdempotency` executes the full unchanged OpenAI/Grok gateway handlers against an actual `httptest.NewServer` HTTP upstream in **standard balance mode**, across image and explicit token pricing. It uses real APIKeyAuthMiddleware and GroupModelAllowlist, the original usage/cost path, and an atomic billing repository stub with actual mutable balance. It asserts a single upstream request, single original usage row, single matching BalanceCost debit, image multiplier7 versus token multiplier2, missing-Key401, duplicate submission reuse, conflicting-body409, JWT-only history GET without a Key and cross-user404. Three polling reads do not increase upstream/debit counts.

`TestImageStudioContentJWTAndCleanedStatus` checks content MIME/length/disposition and exact streamed bytes for PNG/JPEG/WebP metadata, missing JWT401, foreign user404, cleaned410 and zero additional store opens for denied/cleaned reads. Storage's separate tests validate real encoded image parsing.

Unit settings/config/provider tests cover disabled/invalid groups/defaults, authorization, empty scheduler, allowlist/mapping, Grok eligibility, compression/count, image/token preview and configured order. `TestImageStudioAssetStoreUnavailableDoesNotAbortStartup` checks degraded startup for an invalid root. The existing cleanup lifecycle regression includes the new service.

These mocks verify middleware/handler/cost/debit call semantics, not production account credentials or a real PostgreSQL billing transaction. Storage worker owns real PostgreSQL history/claims/storage integration evidence. No real upstream generation, production balance/subscription debit, or live browser acceptance is claimed by this report. Subscription-specific end-to-end debit remains an acceptance boundary; the original shared subscription gateway paths remain unchanged.

Final fresh verification (2026-10-09):

- `go generate ./cmd/server`: Wire regenerated successfully after optional-store provider, aggregate handler registration and cleanup lifecycle changes.
- `go test -tags=unit ./internal/service ./internal/handler ./internal/handler/admin ./internal/server/routes ./cmd/server -run 'TestImageStudio|TestAsyncImage|TestProvideCleanup' -count=1`: PASS. service10.067s, handler7.836s, cmd/server8.473s; admin/routes compile successfully (no dedicated matching tests). Includes configured-order/Grok-mapping/capability/error-code regression corrections and unchanged public async regression.
- `go build ./cmd/server`: PASS (same final shell exited0 after tests).
