# Storage/history/cleanup implementation evidence

Worktree: `/Users/liuliang/.codex/worktrees/user-image-generation/sub2api`. No staging or commits by storage worker. All image deletion tests use isolated temporary directories and testcontainers PostgreSQL.

## Implemented boundaries

- `service/image_studio_history.go`: persistent owner DTOs/repository interface, atomic submission orchestration, bounded result handling and per-asset register-before-save, conditionally completed execution, user-owned content. Store failures preserve successful upstream generation status; never rerun the gateway.
- `repository/image_studio_repo.go`: parameterized SQL against new tables; unique owner/submission plus digest conflict; conditional execution leases; filters; exact original usage owner/key/request cost reconciliation and immutable confirmed snapshot; SKIP LOCKED cleanup claims and retry tombstones.
- `service/image_studio_asset_store.go`, `repository/image_studio_storage_local.go`: fixed image-studio subdirectory, atomic fsynced root binding via nonreplace hardlink, UUID/MIME keys, root/file0700/0600, os.Root confinement and O_NOFOLLOW, same-directory temp+rename. Crash temporary names derive from registered keys and are removed alongside final file. No arbitrary URL/path input or directory scan.
- `service/image_result_reader.go`: reuse existing uploader decoder/download mechanics without changing public uploader behavior; bounded byte/base64/response limits, real PNG/JPEG/WebP decoding, outbound HTTPS/redirect checks and DNS-pinned public IP connections,60second timeout.
- `service/image_studio_cleanup.go`: dedicated PostgreSQL session advisory leader lock, pinned-session heartbeat, database asset leases, enabled/policy refresh each batch, bounded4minute/1000asset runs; cron timezone and periodic settings refresh; only files removed, all generation and asset metadata retained.
- `ent/schema/image_studio_{generation,asset}.go`, `migrations/242_image_studio_history_assets.sql`: standalone feature schemas with no edges to upstream entities; no foreign-key cascades tied to credentials/groups/users. Ent generated using `go generate ./ent`, never Wire in this worker.

## Verification

- Initial RED: targeted unit command failed for absent store/reader constructors; later cleanup RED failed for absent cleanup service.
- Targeted units passed: `go test -tags=unit ./internal/service ./internal/repository -run 'TestImageStudio(History|Cleanup|Storage)|TestImageResultReader|TestImageStorage' -count=1` before final root-binding/late-write regression additions. Final rerun PASS on latest implementation: service6.412seconds/repository4.195seconds, including full/truncated image bytes, registered crash temp removal and revoked-lease late-write cleanup.
- Late-write regression: persisted bytes with a revoked lease reproduced a missing cleanup before fix. Final suite confirms these bytes are removed when the execution lease is genuinely lost; a still-live lease preserves an acknowledged asset.
- Atomic initialization regression:16 simultaneous constructors reproduced empty binding failure before correction. `go test -tags=unit ./internal/repository -run TestImageStudioStorageConcurrentRootInitialization -count=20` PASS after atomic publish.
- `SUB2API_TEST_POSTGRES_IMAGE=postgres:15-alpine go test -tags=integration ./internal/repository -run 'TestImageStudio' -count=1 -timeout=240s -v`: PASS on real PostgreSQL15, four tests,12.738seconds (latest cost test run).
- Real database coverage:8 concurrent same submissions create exactly one record; different hash conflict; cross-user404 and persistent parameter/name snapshots; missing cost remains null/pending; active assets protected; concurrent leases claim once; deletion failure retry/stale lease rejection; interrupted registered-file save becomes unknown, file deleted and all history retained; original usage cost linked by user/key/request and confirmed cache survives hard Key deletion and asset cleanup. An unrelated Key's usage row does not confirm cost.
- Final service-level integration: `SUB2API_TEST_POSTGRES_IMAGE=postgres:15-alpine go test -tags=integration ./internal/repository -run '^TestImageStudioCleanupRunOnceConcurrentLeaderAndDisabledPolicy$' -count=1 -timeout=240s -v` PASS on real PostgreSQL15 (test0.08seconds, package23.234seconds including Docker startup). Uses a real PNG in an isolated local root and the production repository/service. A disabled Start/Reload/RunOnce/Stop leaves eligible image bytes and all metadata unchanged. Two RunOnce instances overlap while the first holds its actual advisory session lock; the second returns without starting/claiming/deleting. The leader then deletes exactly1 image/75bytes with valid summary timestamps, retains every generation field and asset identity/MIME/size/checksum/saved time, and preserves history count. A later invocation by the second instance confirms the leader lock was released and performs0 deletes. No production source changes were required.
- Public image_storage.go was not modified; existing ImageStorage tests are included in targeted command.

## Runtime limits

- Single-instance local storage or shared volume with identical `.root-id` binding is required across replicas; database leases cannot provide file sharing.
- Execution is process-bound, not a durable queue. Expired executions become unknown and are never automatically replayed.
- URLs only come from authenticated gateway results. Production constructor uses the guarded default client; injected HTTP clients are trusted test/provider dependencies.
- History cost snapshots are read from original logs only. If logs are purged before first successful snapshot, cost remains pending rather than guessed.
- Real upstream generation, production billing and browser validation are outside this worker's evidence.
