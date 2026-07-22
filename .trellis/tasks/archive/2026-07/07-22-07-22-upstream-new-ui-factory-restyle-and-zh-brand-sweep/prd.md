# Upstream New UI Factory Restyle + zh Brand Sweep

## Goal

After merging upstream/main into the custom theme branch, upstream's newly added UI files (21 Vue files) arrived with zero adaptation to the custom factory aesthetic, upstream edits injected off-theme fragments into already-restyled pages, and the zh locale was only partially brand-swept. Bring all of these up to the custom standard: factory aesthetic (rounded-sm, primary/bordered language, eyebrow labels) has highest priority, and brand copy must read "SiliconBase" in both en and zh.

## What I already know (from merge review 2026-07-22)

* Global theming is intact: `frontend/src/style.css` tokens + `frontend/src/styles/theme-override.css` (Geist fonts, `.eyebrow`, `.bg-blueprint`) loaded in that order in `main.ts`. Semantic classes (`btn`, `input`, `card`, `input-label`) auto-theme every page.
* All 21 upstream-added Vue files are byte-identical to upstream/main (0 adaptation). Off-theme utility counts measured:
  * `views/user/BatchImageGuideView.vue` — 72 (user-facing, worst)
  * `features/prompt-audit/` — 7 Vue files, 4–18 each (PromptAuditView, EventWorkspace, RuntimeOverview, PolicyPanel, EndpointPool, EventDetailDialog, FilterDeleteDialog)
  * `views/admin/AuditLogView.vue` — 20
  * `components/admin/usage/UserTokenRanking.vue` — 10, no semantic classes
  * `components/account/`: OpenAIQuotaResetCell (10), UpstreamBillingRateCell (6), GrokQuotaProbeCell (4), HeaderOverrideEditor (4), GrokBaseUrlPresets (1), HeaderOverrideJsonTools
  * `components/auth/TotpStepUpDialog.vue` — 2 (mostly fine)
  * Clean already: IpGeoCell, IpGeoBatchToolbar, GrokFreeIcon, BulkEditUserModal (0–2)
* Off-theme fragments injected by upstream into restyled pages:
  * `views/admin/GroupsView.vue` — platform badges use pastel palette (bg-orange-100/emerald-100/purple-100/blue-100, ~line 140-149), rounded-lg dropdown (line ~76); 10 → 16 off-theme
  * `views/admin/DashboardView.vue` — quick-action cards with rounded-lg + sky/emerald icon blocks (~line 95, 206-227); 2 → 10 off-theme
  * `views/user/KeysView.vue` — 6 → 7, minor
* zh brand gaps:
  * Merge-introduced asymmetry (en adapted, zh not): `zh/dashboard.ts:176-178` (Grok group descriptions), `zh/admin/settings.ts:1028` (userIdsHint), `zh/admin/accounts.ts` (whole file raw upstream, incl. trustWarning)
  * Pre-existing leaks carried over from custom/theme: `zh/landing.ts:354-355` (setup wizard), `zh/misc.ts:158/160/167/279/281` (welcome onboarding), `zh/admin/settings.ts` placeholders (siteNamePlaceholder, fromNamePlaceholder, OAuth descriptions, easypay hint), `UseKeyModal.vue:792` (default name "Sub2API Grok"), `BatchImageGuideView.vue:1022` (example placeholder)
* Failing tests to fix (6 files):
  * `router/__tests__/unsupported-region-route.spec.ts` — vi.mock('@/router/title') missing `resolveRouteDocumentTitle` export
  * `i18n/__tests__/openaiFastPolicyLocales.spec.ts` — expects "Sub2API users", actual is "SiliconBase users" (brand)
  * `components/charts/__tests__/ModelDistributionChart.spec.ts` — expects #3b82f6, custom brand color is #ef6f2e
  * `views/admin/__tests__/groupsModelsListLayout.spec.ts` — expects upstream layout classes, GroupsView is factory-styled
  * `views/admin/__tests__/ChannelMonitorView.grok.spec.ts` — reads `.element.value` on a Select component (custom model-search Select replaced plain input); functionality intact, test must assert differently
  * `i18n/__tests__/localesMessageCompile.spec.ts` — needs `pnpm install` (@intlify/message-compiler not in node_modules)

## Requirements

* Restyle the upstream-added Vue files to the factory aesthetic, following `.trellis/spec/frontend/component-guidelines.md` and matching the patterns used in already-restyled pages (rounded-sm, primary-500 accent, bordered icon chips instead of pastel filled blocks, eyebrow labels where a section label exists, dark-mode `dark-*` tokens).
  * Priority 1: `BatchImageGuideView.vue`, all 7 `features/prompt-audit/` Vue files
  * Priority 2: `AuditLogView.vue`, `UserTokenRanking.vue`
  * Priority 3: account cell components (OpenAIQuotaResetCell, UpstreamBillingRateCell, GrokQuotaProbeCell, HeaderOverrideEditor, GrokBaseUrlPresets, HeaderOverrideJsonTools), TotpStepUpDialog touch-up
* Fix off-theme fragments in GroupsView.vue (platform badges, dropdown) and admin DashboardView.vue (quick-action cards, amber icon block) to match surrounding factory style.
* zh brand sweep — align with the en side which is already adapted:
  * `zh/dashboard.ts` Grok descriptions, `zh/admin/settings.ts` userIdsHint + placeholders + OAuth descriptions, `zh/admin/accounts.ts` (adapt whole file as en was)
  * Legacy: `zh/landing.ts` setup wizard, `zh/misc.ts` welcome content, `UseKeyModal.vue` default name, `BatchImageGuideView.vue` placeholder
  * Brand rule: user-visible product references read "SiliconBase"; keep technical identifiers (GitHub org, package names) untouched.
* Fix the 6 failing frontend test files so `pnpm test:run` is green:
  * Update assertions that encode upstream styles/brand to the custom expected values (brand color #ef6f2e, SiliconBase copy, factory layout classes)
  * Update `unsupported-region-route.spec.ts` mock to export `resolveRouteDocumentTitle`
  * Update `ChannelMonitorView.grok.spec.ts` to assert against the Select component (e.g. via component props/vm) instead of raw input `.element.value`
  * Run `pnpm install` first so `@intlify/message-compiler` resolves
* Do not change any functionality, API contracts, or upstream logic — this is a style/i18n/test adaptation only.

## Acceptance Criteria

* [ ] All 21 upstream-added Vue files contain no off-theme utilities (`rounded-lg/xl`, pastel `bg-*-100`/`text-*-700` blocks) except where a status/semantic color is genuinely required by convention used elsewhere in the custom theme
* [ ] GroupsView + admin DashboardView off-theme counts return to custom/theme baseline (10 and 2 respectively) or lower
* [ ] No user-visible "Sub2API" brand strings remain in `frontend/src/i18n/locales/zh/**` where the en counterpart says "SiliconBase"; welcome/setup/placeholder copy aligned
* [ ] `pnpm run typecheck` passes
* [ ] `pnpm test:run` passes 100% (previously 6 failed files / 6 failed tests)
* [ ] `pnpm run lint:check` passes on touched files
* [ ] Backend untouched and still builds (`go build ./...`)
