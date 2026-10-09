# Image console research

Research date: 2026-10-09 (Asia/Shanghai). Repository evidence was read locally; GitHub metadata, branches, trees, source and LICENSE files were read with `gh api`. No product code was changed and no runtime image generation or billing test was performed.

## Existing integration contracts

| Concern | Evidence | Design implication |
| --- | --- | --- |
| Image ingress | `backend/internal/server/routes/gateway.go:259`, `backend/internal/handler/openai_images.go:21` | Reuse the existing Images handler rather than implement a second upstream image client. |
| Image permission | `backend/internal/handler/openai_images.go:98` | Honor `allow_image_generation` on every generation. |
| Account scheduling, concurrency, eligibility and balance reservations | `backend/internal/handler/openai_images.go:106`, `:139`, `:155`, `:182` | A console adapter must enter this full gateway path, not just call an upstream forwarding method. |
| Pricing | `backend/internal/service/openai_gateway_usage.go:203`, `:602`, `:742`; `backend/internal/service/image_billing_multiplier.go:3` | When the resolved channel mode is token, existing token pricing applies. Otherwise the image branch uses the effective image multiplier. Independent image multiplier applies to the image branch, not all image-related token charges. |
| Cost records | `backend/internal/service/openai_gateway_usage.go:420` | UI may display the actual billing mode/rate but must not independently calculate or deduct the final charge. |
| Group permission | `backend/internal/service/api_key_service.go:478`, `:1045` | Public groups, user-restricted groups and subscription groups have different access predicates; picking a default must not bypass them. |
| Console authentication | `backend/internal/server/routes/user.go:21`, `backend/internal/server/middleware/jwt_auth.go:100` | The session user is available on authenticated console routes. |
| Gateway authentication | `backend/internal/server/middleware/api_key_auth.go:56`, `:177`, `:273` | Images needs API Key, user, group and subscription context. Simply setting JWT user context is insufficient. |
| API Key creation | `backend/internal/service/api_key_service.go:490` | Existing Key creation validates group authorization and per-user key limits. The selected credential must belong to the user and group. Reuse the real creation form and original restrictions; no automatic provisioning. |
| Async tasks | `backend/internal/handler/image_task_handler.go:33`, `:54`, `:164`, `:213` | Async Submit calls the same gateway; object storage gates submissions. Get uses owner `(user_id, api_key_id)` and returns results even if generation is disabled later. |
| Task retention | `backend/internal/service/image_task.go:25` | Default retention is 24 hours and execution timeout is 30 minutes. These are task limits, not a permanent works gallery. |
| Async platform support | `backend/internal/handler/image_task_handler.go:71` | Existing async handler accepts OpenAI and Grok groups; composite groups require additional support and cannot be assumed equivalent. |
| Models | `backend/internal/handler/gateway_handler.go:1148`, `backend/internal/service/gateway_service.go:1393` | Reuse group mapping/allowlist semantics, then restrict to image-capable models. The gateway list can fall back to defaults, so presence in `/v1/models` alone does not prove runtime availability. |
| Catalog | `backend/internal/handler/model_catalog_handler.go:20`, `backend/internal/handler/dto/model_catalog.go:17` | Catalog has capabilities/endpoints/pricing metadata useful for labels, but public catalog visibility is not generation authorization. |
| Image parameters | `backend/internal/service/openai_images.go:84`, `:271`, `:276`, `:283` | Existing parser covers size, quality, background, output format, compression and count. Parameter options must follow each supported model's contract. |
| Admin defaults | `backend/internal/service/settings_view.go:181`, `frontend/src/views/admin/SettingsView.vue:3935` | Found registration default subscriptions, not a dedicated console generation default group. Account auto-binding to `<platform>-default` at `backend/internal/service/admin_account.go:500` is a separate account onboarding behavior. |
| Existing batch UI | `frontend/src/router/index.ts:251`, `frontend/src/views/user/BatchImageGuideView.vue`, `frontend/src/composables/useBatchImageAccess.ts:16` | The view is more than a static guide: it manages batch tasks using user API Keys. Its access flag is Gemini batch-specific and should not gate the GPT image console. |
| Styling | `frontend/package.json`, `frontend/src/views/user/BatchImageGuideView.vue` | Vue 3 + TS + Tailwind + Pinia + vue-i18n, with shared AppLayout, Select and Icon components; retain the current console design system. |

## Open-source candidates

Stars are the values returned by GitHub during this research, not a stable acceptance target. LICENSE contents were read, not inferred only from SPDX metadata.

| Project | Stars | License | Pinned commit | Fit |
| --- | ---: | --- | --- | --- |
| [Open Generative AI](https://github.com/Anil-matcha/Open-Generative-AI) | 29,895 | MIT, copyright 2026 Open Generative AI Contributors | `55e02f0b301fd971839588125096cdfbdd93cdbd` | Complete image studio with prompt composer, model parameters, gallery, progress/error handling and download. React/Muapi-specific; port selected UI behavior and plain JS helpers into the Vue console. |
| [Toonflow](https://github.com/HBAI-Ltd/Toonflow-app) | 16,806 | MIT, copyright 2026 HBAI-Ltd | `72a895c26aab3f54c5a914517615362208fa6008` | Vue generation settings component with resolution and aspect-ratio controls. The overall application depends on canvas scaffolding and Element Plus; selective component adaptation is viable. |
| [NextChat](https://github.com/ChatGPTNextWeb/NextChat) | 88,835 | MIT, copyright 2023-2025 NextChat | `839a24df5d655eb8c3477fb35e8da578dc206f39` | Highly starred, but the inspected independent image panel is Stable Diffusion/Stability oriented; less directly suited to a GPT Images console. |

### Source anchors actually inspected

- Open Generative AI [ImageStudio.jsx](https://github.com/Anil-matcha/Open-Generative-AI/blob/55e02f0b301fd971839588125096cdfbdd93cdbd/packages/studio/src/components/ImageStudio.jsx): parameter controls, prompt composer, history and preview/download gallery. Its API adapter is Muapi; it cannot be used as this project's billing client.
- Open Generative AI [imageSizing.js](https://github.com/Anil-matcha/Open-Generative-AI/blob/55e02f0b301fd971839588125096cdfbdd93cdbd/packages/studio/src/imageSizing.js): capability-driven aspect ratios and width/height constraints. GPT-specific dimensions still need this gateway's validated contract.
- Open Generative AI [downloadImage.js](https://github.com/Anil-matcha/Open-Generative-AI/blob/55e02f0b301fd971839588125096cdfbdd93cdbd/packages/studio/src/utils/downloadImage.js): blob download and MIME-aware extension mapping. Adaptation needs HTTP error checking, visible failure state and deliberate handling of URL/CORS behavior.
- Open Generative AI [studio package](https://github.com/Anil-matcha/Open-Generative-AI/blob/55e02f0b301fd971839588125096cdfbdd93cdbd/packages/studio/package.json): React peer dependencies and multiple local workspace dependencies make importing the complete studio expensive.
- Toonflow [generationSettings.vue](https://github.com/HBAI-Ltd/Toonflow-app/blob/72a895c26aab3f54c5a914517615362208fa6008/packages/nodes/imageGenerationNode/src/components/generationSettings.vue): Vue models, size buttons, ratio buttons and ratio thumbnail geometry; replace Element Plus controls with the existing system components/classes.
- Toonflow [node package](https://github.com/HBAI-Ltd/Toonflow-app/blob/72a895c26aab3f54c5a914517615362208fa6008/packages/nodes/imageGenerationNode/package.json): Vue Flow, Element Plus and internal scaffold dependencies support selective adaptation over importing the full node subsystem.
- NextChat [sd-panel.tsx](https://github.com/ChatGPTNextWeb/NextChat/blob/839a24df5d655eb8c3477fb35e8da578dc206f39/app/components/sd/sd-panel.tsx): SD3, Stability Ultra/Core options and provider-specific parameters.
- All three root LICENSE files were read. Any copied or substantially adapted code must carry the applicable copyright and MIT license notice in the delivered third-party notices.

## Current scope and integration decisions (2026-10-09 revision)

The latest user instruction explicitly replaces automatic Key creation with a quick-create button beside the group selector, reusing the existing API Key creation form. Persistent history querying/downloading and administrator-controlled stored-image cleanup are now in scope. Earlier automatic provisioning and session-only results recommendations are superseded; they must not guide implementation.

Adapt selected MIT modules/interactions from the pinned Open Generative AI and Toonflow sources above into the native Vue console. Actual source adaptation, attribution and license notices are implementation acceptance items, not work already completed.

## Manual Key form evidence

- `frontend/src/views/user/KeysView.vue:470`: existing create/edit form is embedded in a BaseDialog, not an already reusable component. It includes provider/group, custom Key, IP rules, quota, expiration and window limits. Extract the real component and keep both consumers on the same form.
- `KeysView.vue:1530`, `:1542`, `:1549`: provider selection and the open/groups-loaded watcher can clear group_id; group preselection must establish the correct provider before validation. Existing ordinary page defaults/edit behavior and tour attributes must survive.
- `KeysView.vue:1800`: handleSubmit uses existing validation and keysAPI; create needs to emit the actual DTO for the studio to refresh/select. `frontend/src/api/keys.ts` already provides group_id listing, getById and create returns ApiKey.
- `backend/internal/service/api_key_service.go:441`, `:490`: creation limits and original authorization apply to user-initiated form creation. No new provisioning lock, ensure endpoint or APIKeyRepository transaction change is required.
- Existing KeysView tests are relevant. The simplified fake in `components/__tests__/ApiKeyCreate.spec.ts` is not proof of real extracted form behavior.

## Persistent history and execution evidence

- `handler/image_task_handler.go:51`, `:211`, `:223`, `:262`: old async Submit validates/audits and copies the authenticated Gin context, uses context.WithoutCancel with a bounded timeout, then calls full OpenAI/Grok handlers through the existing dispatcher. This is an appropriate shared runner seam. Preserve all middleware-derived context, not only the Key pointer.
- Public old async stores owner(user_id,api_key_id) Redis records for 24h and is gated by current object storage. Those properties cannot satisfy persistent console history after Key deletion. Keep them compatible while a studio sink writes separate durable owner(user_id) history.
- `service/openai_images.go:530`: endpoint normalization recognizes /images/generations paths. A studio submission route under the current gateway chain must be tested against endpoint normalization, allowlist, auth eligibility and logging; worker execution should explicitly use canonical /v1/images/generations.
- `server/routes/gateway.go:183`: bodyLimit, clientRequestID, ops logger, endpoint normalization, APIKeyAuth, groupModelAllowlist, compositeTarget and group requirements compose the existing route chain. Studio ingress must keep this chain and then validate configured studio access too.
- `middleware/client_request_id.go:23`: a server-generated client_request_id is persisted in request context. `service/gateway_usage_billing.go:217` uses client:<id> for ordinary Images billing; `openai_gateway_usage.go:337` calls it. Browser submission_id/upstream response ID are not interchangeable with this billing ID.
- `repository/usage_log_repo_query.go:103`, `:107`, `:119`: filters include exact UserID, APIKeyID and RequestID. Existing record insertion deduplicates request+Key. Use owner+Key+request filters to associate actual cost; asynchronous usage writes require a pending/unconfirmed display rather than zero price.
- History capture must be server-side from parsed authorized generation and actual execution result, not a frontend result-upload API. A unique user+submission ID and request hash prevents duplicate dispatch without changing original billing semantics.
- Current process runner is not a durable queue. Reconcile timed-out leases to unknown; never replay a possibly billed upstream call automatically. Distinguish generated-but-save-failed from upstream failure.

## Storage and cleanup evidence

- `service/image_storage.go:22`, `:81`, `:114`, `:194`: old storage returns URLs; uploader rewrites base64/URLs and does not persist object key/provenance. Reader/limits are reusable, but saving only a presigned URL is not a durable identity.
- `repository/image_storage_s3.go:54`: current implementation only PutObject plus public/presigned URL (default24h). `repository/backup_s3_store.go` has GetObject/DeleteObject patterns for future adapters; switching current settings is not a migration of older objects.
- `internal/setup/setup.go:45`: GetDataDir resolves DATA_DIR, then writable /app/data, then current directory. Provide the resolved root to the asset-store constructor at wiring time. Do not couple image paths to pricing.data_dir.
- `deploy/docker-compose.local.yml:40`, standard and standalone compose: /app/data already has a persistent mount. A managed local store avoids a new required service; multi-replica use needs a shared root. This is a proposed technical choice for final review, not a statement of current implemented storage.
- `service/ops_cleanup_service.go:27`, `:90`, `:183`, `:363`: 5-field cron, Start/Stop, live Reload, Redis leader lock with DB fallback patterns. Separate image cleanup key/lease; settings refresh must work across instances, not only the admin request node.
- `service/batch_image_cleanup.go:224`: mark deleted after successful remote deletion; failed deletes remain retryable. This service concerns Gemini provider batch assets and cannot serve the new console files unchanged.
- Asset registration before file write, stable keys, database conditional claims and deletion tombstones/retry state prevent lost-object or false-cleaned states; cleanup must never select in-progress assets or billing logs.
- Current numeric migration maximum is241; both 241_add_payment_order_bonus_amount.sql and 241_add_typesafe_platform.sql exist. Migration identity is filename; add a new file after the actual implementation-time max, preserve published SQL and Ent source-of-truth conventions.

## Model and billing compatibility

- `handler/grok_media.go:21`: Grok Images already implemented; `routes/gateway.go:82` dispatches OpenAI/Grok by Key group. Composite image admission remains unsupported for this scope.
- `service/grok_media_image_geometry.go:13`, `:85`: ratio and current1k/2k definitions. Use model-specific resolution/aspect_ratio instead of universal GPT size/quality/background fields.
- Official xAI generation documentation was fetched previously on2026-10-09: https://docs.x.ai/developers/model-capabilities/images/generation . Base64 Output explicitly supports response_format=b64_json; default URLs are temporary.
- `handler/gateway_key_billing.go:33`: token scope only; image billing_preview must resolve existing channel/group/user image branch functions, not borrow token multiplier as the independent image rate.

## Workspace and confirmed retention decision

Initial unrelated dirty paths remain component-guidelines.md, GatewayDashboard.vue, home/__tests__/ and sub2api.tar. This revision modifies only the task artifacts; no product implementation, database migration, tests, image generation, storage deletion or billing verification has been performed.

On 2026-10-09 the user explicitly confirmed: "图片到期后，只删除图片，其余保留". Cleanup therefore deletes image files only; generation/asset metadata, prompt, model, parameters, time and cost information remain queryable with an image-cleaned state. No history-metadata deletion branch remains in the design. Requirements are converged; validate contexts and submit the latest final planning summary, then await implementation approval before task.py start.


## Implementation follow-through

The earlier research/approval wording above records planning-time evidence. User subsequently authorized implementation. Actual pinned MIT modules and full licenses were adapted and bundled; implementation now lives in the managed codex/user-image-generation worktree. Read validation-matrix.md and the owner/review/browser reports for current evidence and verification boundaries. No automatic Key creation or historical metadata deletion was introduced.
