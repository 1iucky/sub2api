# Research: Current Frontend Inventory (Sub2API reskin scoping)

- **Query**: Inventory the current frontend (Vue 3 + Vite + Tailwind v3.4) to scope a visual reskin
- **Scope**: internal
- **Date**: 2026-06-18

## 1. Route Table

Source: `frontend/src/router/index.ts` (924 lines). Routes lazy-load `@/views/...`. `meta.requiresAuth` defaults to `true` when not set; `meta.requiresAdmin === true` gates admin area. Other gates: `requiresPayment`, `requiresRiskControl`.

### Auth / Setup area (public, `requiresAuth:false`)

| Path | View file | Notes |
|---|---|---|
| `/setup` | `views/setup/SetupWizardView.vue` | First-run wizard |
| `/login` | `views/auth/LoginView.vue` | Uses `AuthLayout` |
| `/register` | `views/auth/RegisterView.vue` | `AuthLayout` |
| `/forgot-password` | `views/auth/ForgotPasswordView.vue` | `AuthLayout` |
| `/reset-password` | `views/auth/ResetPasswordView.vue` | `AuthLayout` |
| `/email-verify` | `views/auth/EmailVerifyView.vue` | `AuthLayout` |
| `/auth/callback` (alias `/auth/oauth/callback`) | `views/auth/OAuthCallbackView.vue` | `AuthLayout` |
| `/auth/linuxdo/callback` | `views/auth/LinuxDoCallbackView.vue` | `AuthLayout` |
| `/auth/wechat/callback` | `views/auth/WechatCallbackView.vue` | `AuthLayout` |
| `/auth/wechat/payment/callback` | `views/auth/WechatPaymentCallbackView.vue` | |
| `/auth/dingtalk/callback` | `views/auth/DingTalkCallbackView.vue` | `AuthLayout` |
| `/auth/dingtalk/email-completion` | `views/auth/DingTalkEmailCompletionView.vue` | `AuthLayout` |
| `/auth/oidc/callback` | `views/auth/OidcCallbackView.vue` | `AuthLayout` |

### Public / Marketing (no auth)

| Path | View file | Notes |
|---|---|---|
| `/home` | `views/HomeView.vue` | **Homepage/landing**. Has iframe/`v-html` custom mode + default gradient hero. No `AppLayout`. |
| `/key-usage` | `views/KeyUsageView.vue` | Public API-key usage lookup |
| `/legal/:documentId` | `views/public/LegalDocumentView.vue` | Legal docs |
| `/` | redirect → `/home` | |

### User console (`requiresAuth:true`, `requiresAdmin:false`)

| Path | View file |
|---|---|
| `/dashboard` | `views/user/DashboardView.vue` |
| `/keys` | `views/user/KeysView.vue` |
| `/usage` | `views/user/UsageView.vue` |
| `/redeem` | `views/user/RedeemView.vue` |
| `/affiliate` | `views/user/AffiliateView.vue` |
| `/available-channels` | `views/user/AvailableChannelsView.vue` |
| `/profile` | `views/user/ProfileView.vue` |
| `/subscriptions` | `views/user/SubscriptionsView.vue` |
| `/purchase` | `views/user/PaymentView.vue` (`requiresPayment`) |
| `/orders` | `views/user/UserOrdersView.vue` (`requiresPayment`) |
| `/payment/qrcode` | `views/user/PaymentQRCodeView.vue` |
| `/payment/result` | `views/user/PaymentResultView.vue` (public) |
| `/payment/stripe` | `views/user/StripePaymentView.vue` (public) |
| `/payment/stripe-popup` | `views/user/StripePopupView.vue` (public) |
| `/payment/airwallex` | `views/user/AirwallexPaymentView.vue` (public) |
| `/custom/:id` | `views/user/CustomPageView.vue` | Admin-defined custom pages |
| `/monitor` | `views/user/ChannelStatusView.vue` | Public channel status |

### Admin console (`requiresAuth:true`, `requiresAdmin:true`)

| Path | View file |
|---|---|
| `/admin` → `/admin/dashboard` | redirect |
| `/admin/dashboard` | `views/admin/DashboardView.vue` |
| `/admin/ops` | `views/admin/ops/OpsDashboard.vue` |
| `/admin/users` | `views/admin/UsersView.vue` |
| `/admin/groups` | `views/admin/GroupsView.vue` |
| `/admin/channels` → `/admin/channels/pricing` | redirect |
| `/admin/channels/pricing` | `views/admin/ChannelsView.vue` |
| `/admin/channels/monitor` | `views/admin/ChannelMonitorView.vue` |
| `/admin/subscriptions` | `views/admin/SubscriptionsView.vue` |
| `/admin/accounts` | `views/admin/AccountsView.vue` |
| `/admin/announcements` | `views/admin/AnnouncementsView.vue` |
| `/admin/proxies` | `views/admin/ProxiesView.vue` |
| `/admin/redeem` | `views/admin/RedeemView.vue` |
| `/admin/promo-codes` | `views/admin/PromoCodesView.vue` |
| `/admin/settings` | `views/admin/SettingsView.vue` |
| `/admin/risk-control` | `views/admin/RiskControlView.vue` (`requiresRiskControl`) |
| `/admin/usage` | `views/admin/UsageView.vue` |
| `/admin/affiliates` → `/admin/affiliates/invites` | redirect |
| `/admin/affiliates/invites` | `views/admin/affiliates/AdminAffiliateInvitesView.vue` |
| `/admin/affiliates/rebates` | `views/admin/affiliates/AdminAffiliateRebatesView.vue` |
| `/admin/affiliates/transfers` | `views/admin/affiliates/AdminAffiliateTransfersView.vue` |
| `/admin/orders/dashboard` | `views/admin/orders/AdminPaymentDashboardView.vue` (`requiresPayment`) |
| `/admin/orders` | `views/admin/orders/AdminOrdersView.vue` (`requiresPayment`) |
| `/admin/orders/plans` | `views/admin/orders/AdminPaymentPlansView.vue` (`requiresPayment`) |

### 404

| `/:pathMatch(.*)*` | `views/NotFoundView.vue` | public, standalone gradient page |

**Route totals:** ~13 auth/setup, 3 public/marketing (home, key-usage, legal), 1 home-redirect, 1 404, 17 user (some payment routes are public but logically grouped), 22 admin.

## 2. Layout Components

All layout components live in `frontend/src/components/layout/` (with `index.ts` barrel).

| File | Lines | Role |
|---|---|---|
| `AppLayout.vue` | 52 | **Main shell for console pages (admin + user).** Renders `AppSidebar`, `AppHeader`, and `<main class="p-4 md:p-6 lg:p-8">` slot. Background: `bg-gray-50 dark:bg-dark-950` with `bg-mesh-gradient` overlay. Responsive margin that tracks `appStore.sidebarCollapsed` (72px collapsed / 256px expanded on `lg:`). |
| `AppSidebar.vue` | 1022 | **Sidebar nav + brand + theme toggle.** Logo+siteName+VersionBadge header, dual nav (admin section then personal section for admin users), collapsible groups, mobile slide-in. Theme toggle: toggles `document.documentElement.classList` `dark`, persists to `localStorage.theme`. |
| `AppHeader.vue` | 343 | **Top header.** Sticky glass, page title + description, announcement bell, docs link, locale switcher, subscriptions widget, balance widget, user dropdown. |
| `AuthLayout.vue` | 88 | **Shell for auth/setup pages.** Centered card on a gradient + grid pattern backdrop with logo/brand. Uses `appStore.siteName`, `siteLogo`, `site_subtitle`. |
| `TablePageLayout.vue` | 112 | **Reusable table-page wrapper** used inside `AppLayout`. Sticky actions/filters, scrollable table body, fixed pagination. Has scoped CSS for table styling. |

**Layout wiring pattern:** Each view explicitly imports and wraps its own layout — there is no global `<router-view>` → layout mapping. Views that need the console chrome do `<AppLayout>…</AppLayout>`; many additionally nest `<TablePageLayout>`. Auth views use `<AuthLayout>`. Public pages (home, 404, key-usage) are standalone (no layout wrapper). `App.vue` is just `<NavigationProgress /><RouterView /><Toast /><AnnouncementPopup /><AdminComplianceDialog />`.

## 3. Views per Area

### (a) Homepage / Landing — 1 file
- `views/HomeView.vue` (644 lines) — default gradient hero; optional `home_content` admin setting switches to iframe/`v-html` mode.

### (b) Admin console pages — 14 top-level + nested
Top-level (in `views/admin/`): `DashboardView`, `UsersView`, `GroupsView`, `ChannelsView`, `ChannelMonitorView`, `SubscriptionsView`, `AccountsView`, `AnnouncementsView`, `ProxiesView`, `RedeemView`, `PromoCodesView`, `SettingsView`, `RiskControlView`, `UsageView`, plus `BackupView.vue` (not routed, tooling?).
Nested:
- `views/admin/orders/` — `AdminOrdersView.vue`, `AdminPaymentDashboardView.vue`, `AdminPaymentPlansView.vue`, `PlanEditDialog.vue`
- `views/admin/affiliates/` — `AdminAffiliateInvitesView.vue`, `AdminAffiliateRebatesView.vue`, `AdminAffiliateTransfersView.vue`, `AdminAffiliateRecordsTable.vue`
- `views/admin/ops/` — `OpsDashboard.vue` + `components/` (OpsDashboardHeader 1627 lines, OpsSettingsDialog, OpsAlertEventsCard, OpsConcurrencyCard, OpsAlertRulesCard, OpsRuntimeSettingsCard, OpsSystemLogTable, OpsEmailNotificationCard) and `utils/`

### (c) User console pages — 16 in `views/user/`
`DashboardView`, `KeysView`, `UsageView`, `RedeemView`, `AffiliateView`, `AvailableChannelsView`, `ProfileView`, `SubscriptionsView`, `PaymentView`, `UserOrdersView`, `PaymentQRCodeView`, `PaymentResultView`, `StripePaymentView`, `StripePopupView`, `AirwallexPaymentView`, `ChannelStatusView`, `CustomPageView`.

### Auth/setup — 14 in `views/auth/` + `views/setup/`
`LoginView`, `RegisterView`, `ForgotPasswordView`, `ResetPasswordView`, `EmailVerifyView`, `OAuthCallbackView`, `LinuxDoCallbackView`, `WechatCallbackView`, `WechatPaymentCallbackView`, `DingTalkCallbackView`, `DingTalkEmailCompletionView`, `OidcCallbackView`, `SetupWizardView`.

### Public/other — 3
`HomeView`, `KeyUsageView`, `NotFoundView`, `views/public/LegalDocumentView.vue`.

## 4. Styling Touch-Points

### Global CSS entry
- `frontend/src/style.css` (765 lines) — the **single global stylesheet**, imported once in `main.ts` (line 7). Structure: `@tailwind base/components/utilities` then `@layer base`, `@layer components`, `@layer utilities`.
  - `@layer base` — scrollbar styling, selection color, default body bg/text, dark-mode overrides.
  - `@layer components` — ~80 component classes (`@apply`-based design system): all button variants (`btn`, `btn-primary`, `btn-secondary`, `btn-ghost`, `btn-danger`, `btn-success`, `btn-warning`, `btn-stripe`, `btn-airwallex`, `btn-alipay`, `btn-wxpay`, size mods `btn-sm/md/lg/icon`), inputs (`input`, `input-error`, `input-label`, `input-hint`), glass effects (`glass`, `glass-card`, `card-glass`), cards (`card`, `card-hover`, `card-header`, `card-body`, `card-footer`), stats (`stat-card`, `stat-icon-*`, `stat-value`, `stat-label`, `stat-trend-*`), tables (`.table`, `.table-container`, `.table-wrapper`), badges, dropdowns, modals (`.modal-*`, `.dialog-*`), toasts (`.toast-*`), sidebar classes (`.sidebar`, `.sidebar-header`, `.sidebar-link`, `.sidebar-link-active`, `.sidebar-section`, etc.), page header (`page-header`, `page-title`, `page-description`), empty states, spinners, skeletons, tabs, progress, switch, code blocks, tour helpers, `.text-gradient`, `.scrollbar-hide`, `.safe-top`, `.safe-bottom`.
  - `@layer utilities` — small utility additions.
- `frontend/src/styles/onboarding.css` (228 lines) — imported by `AppLayout.vue` (`import '@/styles/onboarding.css'`). Styles the onboarding-tour overlay (driver.js-style guided tours).

### Tailwind tokens
- `frontend/tailwind.config.js` (134 lines) — `darkMode: 'class'`, content globs. Extended tokens:
  - `colors.primary` — teal/cyan scale (50–950, `#14b8a6` is 500).
  - `colors.accent` and `colors.dark` — slate/`#1e293b` palette for dark surfaces.
  - `fontFamily.sans` — system-ui + CJK fallbacks (PingFang SC, Microsoft YaHei…). `fontFamily.mono` standard.
  - `boxShadow` — `glass`, `glass-sm`, `glow`, `glow-lg`, `card`, `card-hover`, `inner-glow`.
  - `backgroundImage` — `gradient-primary`, `gradient-dark`, `gradient-glass`, `mesh-gradient`.
  - `animation`/`keyframes` — `fade-in`, `slide-up`, `slide-down`, `slide-in-right`, `scale-in`, `pulse-slow`, `shimmer`, `glow`.
  - `backdropBlur.xs`, `borderRadius.4xl`.
  - `plugins: []` (no typography/forms plugin).

### Scoped `<style>` blocks in components
- ~51 `.vue` files (excluding tests) contain a `<style>` block (out of ~254 total `.vue`). Sampling: most are plain Tailwind-utility templates with no scoped CSS; the ones with scoped styles are typically larger components needing one-off tweaks (e.g. `AppHeader`, `AppSidebar`, `AuthLayout`, `TablePageLayout`, `HomeView`, `LoginView`, many `DataTable`/`Select`/`Toast` common components, all auth callback views, ops subcomponents).
- All inspected layout scoped blocks use `<style scoped>` (CSS-scoped to component) — verified for `AppHeader.vue:332`, `AuthLayout.vue:84`, `TablePageLayout.vue:46`, `AppSidebar.vue:895`, `HomeView.vue:483`, `LoginView.vue:555`.

### How widely scoped CSS vs utilities?
- **Dominant pattern:** Tailwind utility classes inline in templates, augmented by the global `@apply` component classes from `style.css`. Scoped CSS is the exception, used for layout primitives (`TablePageLayout`), one-off decorative gradients/orbs (`AuthLayout`, `HomeView`, `LoginView`), or third-party-style overrides (toast animations, table sticky columns). A reskin that changes tokens in `tailwind.config.js` + component classes in `style.css` will cascade across most of the app; the remaining deltas will be in the ~51 scoped blocks.

## 5. Conflict-Risk Ranking (largest view/component files)

These are the biggest rebase-conflict risks if edited in place. View files first, then component files.

### Views (top by line count, excluding tests)
| Lines | File |
|---|---|
| 10410 | `src/views/admin/SettingsView.vue` |
| 4351 | `src/views/admin/GroupsView.vue` |
| 2337 | `src/views/admin/RiskControlView.vue` |
| 2067 | `src/views/admin/ProxiesView.vue` |
| 1800 | `src/views/admin/UsersView.vue` |
| 1790 | `src/views/user/KeysView.vue` |
| 1722 | `src/views/admin/AccountsView.vue` |
| 1632 | `src/views/admin/ChannelsView.vue` |
| 1627 | `src/views/admin/ops/components/OpsDashboardHeader.vue` |
| 1421 | `src/views/admin/SubscriptionsView.vue` |
| 1189 | `src/views/admin/RedeemView.vue` |
| 1120 | `src/views/user/UsageView.vue` |
| 1102 | `src/views/auth/WechatCallbackView.vue` |
| 1077 | `src/views/user/PaymentView.vue` |
| 1001 | `src/views/KeyUsageView.vue` |

### Components (top by line count, excluding tests)
| Lines | File |
|---|---|
| 5530 | `src/components/account/CreateAccountModal.vue` |
| 4241 | `src/components/account/EditAccountModal.vue` |
| 1800 | `src/components/account/BulkEditAccountModal.vue` |
| 1269 | `src/components/account/AccountUsageCell.vue` |
| 1080 | `src/components/keys/UseKeyModal.vue` |
| 1022 | `src/components/layout/AppSidebar.vue` |
| 933 | `src/components/common/DataTable.vue` |
| 862 | `src/components/account/OAuthAuthorizationFlow.vue` |
| 749 | `src/components/account/AccountStatsModal.vue` |
| 713 | `src/components/admin/account/AccountStatsModal.vue` |
| 703 | `src/components/payment/PaymentProviderDialog.vue` |
| 684 | `src/components/admin/account/ScheduledTestsPanel.vue` |
| 660 | `src/components/user/profile/ProfileIdentityBindingsSection.vue` |
| 650 | `src/components/admin/ErrorPassthroughRulesModal.vue` |
| 625 | `src/components/admin/TLSFingerprintProfilesModal.vue` |

`SettingsView.vue` at 10.4k lines and the three account modals (5.5k/4.2k/1.8k) are the highest-risk single files for an in-place reskin. `AppSidebar.vue` (1k lines) is also high-risk because it carries the brand block, theme toggle, and nav-item rendering — any sidebar reskin will land here.

## 6. Entry / Import Chain & Dark-Mode Mechanism

**Global CSS chain (confirmed):**
1. `frontend/src/main.ts:7` — `import './style.css'` (the only global stylesheet import).
2. `style.css` `@tailwind` directives + `@layer base/components/utilities` produces the global design system.
3. `AppLayout.vue:26` — `import '@/styles/onboarding.css'` (extra, only loaded when `AppLayout` mounts, i.e., on every console page).
4. No CSS is imported in `App.vue` itself; per-view styling comes from Vue SFC `<style>` blocks compiled by Vite.

**Dark-mode mechanism (`darkMode: 'class'`):**
- Bootstrap: `main.ts` `initThemeClass()` runs **before** mount — reads `localStorage.theme`; falls back to `window.matchMedia('(prefers-color-scheme: dark)')`; toggles `dark` class on `document.documentElement`.
- Runtime toggle: `AppSidebar.vue` `toggleTheme()` (line ~795) — flips `document.documentElement.classList`, writes `localStorage.theme`.
- Consumers read dark state via `document.documentElement.classList.contains('dark')` directly (e.g. charts `TokenUsageTrend.vue:58`, `AccountStatsModal.vue:500`, Stripe appearance `StripePaymentInline.vue:119`).
- All global component classes in `style.css` use explicit `dark:` variants; the dark surface palette is `colors.dark` from `tailwind.config.js` (slate scale).

## 7. i18n / Branding Hooks

### Branding (site name / logo / subtitle)
- **Runtime store:** `frontend/src/stores/app.ts` — `useAppStore`:
  - `siteName` (default `'Sub2API'`), `siteLogo` (default `''`), `siteVersion`, `contactInfo`, `apiBaseUrl`, `docUrl`.
  - `applySettings(config)` (line ~290) sets these from `PublicSettings` (also stashes into `window.__APP_CONFIG__` to prevent flash).
  - `fetchPublicSettings(force)` (line ~308) — checks `window.__APP_CONFIG__` then API `/auth/public-settings` (`api/auth.ts`).
- **Config injection:** `main.ts` calls `appStore.initFromInjectedConfig()` before mount — server-injected `window.__APP_CONFIG__` from backend.
- **Types:** `frontend/src/types/index.ts:204-206` — `site_name`, `site_logo`, `site_subtitle` are fields of `PublicSettings`.
- **Logo consumers** (hardcoded fallback `/logo.png` — static asset at `frontend/public/logo.png`):
  - `components/layout/AppSidebar.vue:13`, `components/layout/AuthLayout.vue:36`, `views/HomeView.vue:45`, `views/KeyUsageView.vue:8`, `views/public/LegalDocumentView.vue:7`. Pattern: `:src="siteLogo || '/logo.png'"`.
- **Title:** `router/index.ts` resolves document title via `resolveDocumentTitle(meta.title, siteName, titleKey)` per navigation; `App.vue` also updates title on mount. Custom pages (`/custom/:id`) use the menu item label.
- **Favicon:** `App.vue` `updateFavicon(siteLogo)` watches `appStore.siteLogo`.

### i18n
- `frontend/src/i18n/` — vue-i18n runtime + JIT (CSP-safe). `main.ts` awaits `initI18n()` before mount. Locales lazy-loaded (`en`, `zh`). Routes carry `titleKey` / `descriptionKey` for localized page titles.

### Theme tokens for rebrand
- **Primary color** lives in `tailwind.config.js` `theme.extend.colors.primary` (teal scale). Changing this single object re-themes every `primary-*` utility and the `@apply from-primary-*` component classes in `style.css` (buttons, focus rings, gradients, mesh background).
- **Dark surface palette** — `colors.dark` + `colors.accent`.
- **Shadows / gradients** — `theme.extend.boxShadow` (`glass`, `glow`, `card`) and `backgroundImage` (`gradient-primary`, `mesh-gradient`).
- **Brand text gradient** — `.text-gradient` in `style.css:739` + scoped `.text-gradient` in `AuthLayout.vue:85`.

### Where a rebrand touches
1. `tailwind.config.js` — color tokens, shadows, gradients, fonts.
2. `style.css` — `@layer base` (body bg, selection, scrollbar) and `@layer components` (the entire design system).
3. `AppSidebar.vue` (brand block + theme toggle) and `AppHeader.vue` (chrome).
4. `AuthLayout.vue` + `HomeView.vue` + `NotFoundView.vue` — bespoke gradient/orb backdrops that bypass the global classes.
5. `public/logo.png` — default logo asset.
6. ~51 component scoped `<style>` blocks — mostly small; exceptions: `AppSidebar.vue:895`, `HomeView.vue:483`, `LoginView.vue:555`, `TablePageLayout.vue:46`.

## Caveats / Not Found

- `views/admin/BackupView.vue` (639 lines) exists but is **not routed** in `router/index.ts` — likely orphaned or used via a non-route flow. Mentioned for completeness; not a reskin target unless reintroduced.
- `views/auth/` contains `README.md`, `USAGE_EXAMPLES.md`, `VISUAL_GUIDE.md` — documentation only, not Vue components.
- `views/auth/WechatPaymentCallbackView.vue` is listed in the route table but is not in the `views/auth/` `ls` output above — verify presence (may exist; route file references it). If missing, build would fail, so it must exist.
- Did not enumerate every scoped `<style>` block's purpose; the ~51-file list above is the complete set of files containing `<style>` blocks (tests excluded). Some contain tiny one-liners, others (sidebar/home/login/table-layout) carry substantial rules.
- Frontend tech confirmed: Vue 3 (Composition API), Pinia, vue-i18n, vue-router 4, Vite, Tailwind v3.4 (`tailwind.config.js` CommonJS-style export; no v4 `@theme`/`@config` — confirms v3).
