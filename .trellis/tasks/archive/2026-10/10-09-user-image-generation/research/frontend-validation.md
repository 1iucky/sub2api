# Frontend implementation validation

Frontend implementation is isolated in `/Users/liuliang/.codex/worktrees/user-image-generation/sub2api`.
No commit, staging, deployment, homepage edits, or credential persistence was performed by this worker.

## Delivered behavior

- Existing complete KeysView create/edit form is moved to `components/keys/ApiKeyFormDialog.vue` and consumed by both KeysView and ImageStudioView. Preserved custom Key, IP restrictions, quota/reset, rate windows/reset, expiration, provider/category selection, edit status rules and onboarding tour attributes.
- Studio config capabilities drive model fields. No page-level OpenAI/Grok branch. Dynamic enum/integer/boolean controls reset incompatible fields on model change.
- Manual group/Key selection; complete group-filtered Key pagination, no automatic creation or replacement. Cross-group creations do not get selected as current-group credentials.
- Generation uses one fetch POST with selected Key and UUID `X-Image-Studio-Submission-ID`, a 30-second bounded submit request, a UUID fallback for local HTTP consoles and no replay/retry interceptor. Lost responses recover via read-only JWT submission lookup. JWT history restores processing records after refresh.
- History filters/pagination/detail preserve owner snapshots, prompt/model/parameters/time and pending/confirmed fees. Asset tombstones never load/download. A cleanup-racing HTTP410 also removes download controls.
- JWT Blob previews/downloads use actual MIME filename extensions and release preview URLs. Logout/unmount clears in-memory credentials and stops timers/requests.
- Admin standalone settings choose backend-supported groups/default, daily/weekly/custom schedule, time zone, retention and cleanup toggle. Image-only irreversible deletion impact is shown; first-run absence is represented explicitly.
- English/Chinese independent locale modules plus small existing locale registration changes.

## Checks

Meaningful tests were added before implementation (initial red runs failed at the absent real component/composable/result module imports).

Focused Vitest suite includes original KeysView26 + new shared form6 + state7 + API5 + parameters3 + history1 + image result6 + admin settings2 + i18n3 = 59 tests.

Commands:

```sh
pnpm --dir frontend exec vitest run src/views/user/__tests__/KeysView.spec.ts src/components/keys/__tests__/ApiKeyFormDialog.spec.ts src/composables/__tests__/useImageStudio.spec.ts src/components/imageStudio/__tests__ src/components/admin/settings/__tests__/ImageStudioSettings.spec.ts src/api/__tests__/imageStudio.spec.ts src/i18n/__tests__/localeKeyCompleteness.spec.ts
pnpm --dir frontend typecheck
pnpm --dir frontend run build
```

Targeted ESLint (all touched frontend files, no global fix) and `git diff --check` were run.
Build output goes to the existing ignored `backend/internal/web/dist` target. Existing build warnings include Browserslist age, large chunks and app-store dynamic/static import overlap.

These are frontend unit/static/build checks, not browser CSS, live PostgreSQL ownership, mocked full gateway, real upstream output, or production billing validation. Root coordinates those acceptance layers separately.

## Actual third-party adaptation

Pinned sources and both complete root LICENSE files were retrieved with `gh api -X GET .../contents/...?... --jq .content | base64 --decode` (raw curl was unavailable).

1. Open Generative AI `packages/studio/src/utils/downloadImage.js`, commit `55e02f0b301fd971839588125096cdfbdd93cdbd`: MIME extension map, basename replacement and Blob/anchor download sequence adapted to `components/imageStudio/downloadImage.ts`, authenticated content and explicit failures.
2. Toonflow `packages/nodes/imageGenerationNode/src/components/generationSettings.vue`, commit `72a895c26aab3f54c5a914517615362208fa6008`: ratio thumbnail calculation and pressed-state ratio buttons adapted to native controls in `ImageParameterControls.vue`. No React/Element Plus/canvas dependencies.

`THIRD_PARTY_NOTICES.md` carries source paths, immutable commit IDs, actual copyright statements and both complete MIT license texts. Adapted source files carry attribution comments. A matching `frontend/public/THIRD_PARTY_NOTICES.md` is copied into the built web distribution so binary/web deployments carry the full licenses too.

## Upstream merge footprint

Unavoidable full KeysView form extraction is the largest existing-file delta, covered by the original26 regression tests and real shared-form tests. Other existing integration deltas are one sidebar entry, one route entry, one settings tab/component registration and independent locale imports/tab titles. New feature logic remains in dedicated files. No business/frontend dependency upgrades or unrelated refactoring.

## Final readability and bilingual error follow-up

Only new studio feature files received readability formatting; existing integration files were not reformatted. Formatting was checked against emitted JavaScript/Vue render ASTs before the subsequent authorized error-localization changes. No formatter dependency or configuration was added.

Stable default-group reasons and history/submission error categories now render English/Chinese messages. Unknown codes and arbitrary upstream diagnostics use a translated generic fallback, including inherited object names such as `constructor`. Only a bounded allowlist of local gateway parameter validation details survives. Exact backend admission codes `IMAGE_STUDIO_GROUP`, `IMAGE_STUDIO_MODEL`, and `IMAGE_STUDIO_STORAGE` are recognized after case normalization.

Follow-up runs, distinct from the earlier full 59-test suite:

- API5 + composable7 + initial page4 + i18n3: 19 tests passed.
- Expanded page6 + i18n3: 9 tests passed; all new feature/error files and both locale modules passed targeted ESLint; `vue-tsc --noEmit` passed.
- Final exact backend-code regressions: page12 + i18n3 = 15 tests passed. The error mapper and page test passed targeted ESLint, `vue-tsc --noEmit` passed, and `git diff --check -- frontend` passed (combined command exit0).

The test-only runtime i18n compiler emits its experimental-feature warning; application compilation was not changed. Root owns the final production build after these source changes. Native-browser saved-file download verification remains a separate acceptance limitation recorded by root.
