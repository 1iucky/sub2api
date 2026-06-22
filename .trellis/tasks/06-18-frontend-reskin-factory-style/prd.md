# Frontend Reskin → 硅基链 / SiliconBase (Factory.ai style)

## Goal

Completely rebrand and restyle the entire Sub2API frontend — homepage, admin console, and user console — referencing https://factory.ai/ as the design north star, and relaunch it under a new brand **硅基链 / SiliconBase** (`siliconbase.link`). The result must be visually unrecognizable from the original open-source project, while (a) keeping all backend/business logic untouched and (b) staying rebasable on upstream `Wei-Shaw/sub2api` via a tiered, documented touch-point architecture.

**Delivery is phased.** This task = **Phase 1 (MVP)**: design-system foundation + homepage + app shell + full rebrand (+ authored SVG illustrations + key-surface i18n copy). Phase 2 (user console pages) and Phase 3 (admin console pages) are follow-up tasks that inherit Phase 1's cascade + shell.

## Requirements

**Brand**
* New brand: **硅基链** (zh) / **SiliconBase** (en), domain `siliconbase.link`.
* New logo: hand-authored geometric SVG mark in factory aesthetic (warm-dark + vermillion), evoking silicon / lattice / chain motifs.
* Apply brand to: default siteName fallback, default logo asset, favicon, document title, homepage wordmark/footer, loading page, empty states, 404 page.

**Illustration (authored SVG only)**
* Hand-authored **SVG vector illustrations** in factory geometric/technical style (line-art / isometric schematic / pattern motifs) for: homepage hero decoration, empty states, brand graphics.
* Constraint: SVG vector only — **no raster / hand-painted / AI-generated art** (no image-gen tool available).

**i18n copy (key surfaces, en/zh)**
* Brand-voice copy rewrite (en/zh) on key surfaces: homepage hero/features, login/register, empty states, nav, footer — tone: technical, restrained, SiliconBase voice.
* Console table/field/tooltip copy left as-is (full-locale rewrite is out of scope).

**Design system (factory.ai aesthetic)**
* Geist + Geist Mono fonts (npm `geist`).
* Warm near-black dark-first theme (`#020202` bg / `#0a0908` surface, warm-gray text scale) + working light theme (`#fafafa`).
* Vermillion accent (`#EF6F2E` / `#ee6018` / `#d15010`).
* Low radius (3–6px buttons/inputs, 8–12px cards), no shadows, hairline warm-gray borders, dashed dividers.
* Mono uppercase 12px eyebrow labels; tight negative tracking on `font-normal` headings (do NOT bold H1/H2).
* Color-inversion hover buttons (150ms, ease-in-out); subtle fade-up motion (respect `prefers-reduced-motion`).
* Signature motif: faint line-grid blueprint background (login / empty states / homepage, NOT behind data-dense tables).

**Scope (Phase 1)**
* Token + component-class overhaul in `tailwind.config.js` + `src/style.css` (cascades to ~80% of the app, both themes).
* New `src/styles/theme-override.css` (Geist @font-face wiring, line-grid bg, dashed dividers, mono eyebrow, focus-visible outline) + one import line in `main.ts`.
* App shell restyle: `AppLayout`, `AppSidebar`, `AppHeader`, `AuthLayout`.
* Homepage restyle: `HomeView` default mode — keep section flow (nav → hero+terminal → features → providers → footer), transform visual language + SVG hero illustration + homepage copy. Admin custom-content (iframe / v-html) mode untouched.
* Brand polish + content: empty states (with SVG illustrations), loading page, favicon, 404; key-surface i18n copy.

## Acceptance Criteria

* [ ] Geist + Geist Mono loaded; `fontFamily` wired in `tailwind.config.js`.
* [ ] `tailwind.config.js` tokens reflect factory palette (warm dark + vermillion + low radius + no-shadow) for dark & light.
* [ ] `src/style.css` component classes rewritten to factory aesthetic (inversion-hover buttons, hairline cards, dashed table dividers, mono labels) for both themes.
* [ ] `theme-override.css` created and imported in `main.ts`; line-grid + dashed-divider + mono-eyebrow utilities available.
* [ ] App shell (sidebar/header/auth layout) visually restyled; sidebar shows new brand + mono section labels; theme toggle still works.
* [ ] Homepage unrecognizable vs original; factory feel; new brand wordmark + domain in footer.
* [ ] Authored SVG illustrations present on homepage hero + empty states (+ brand graphics).
* [ ] i18n (en/zh) copy rewritten on homepage, auth, empty states, nav, footer in SiliconBase brand voice.
* [ ] New logo mark (SVG) + favicon + default siteName = 硅基链/SiliconBase applied.
* [ ] Empty states, loading page, 404 restyled to new system.
* [ ] Spot-QA: cascade produces a coherent look on ≥1 sample admin page + ≥1 sample user page, in dark & light.
* [ ] `pnpm run lint:check`, `pnpm run typecheck`, `pnpm run build` all green.
* [ ] Touch-point conflict surface documented in `CLAUDE.md` (or `REBRAND.md`) for future rebases.

## Definition of Done

* All AC met; lint/typecheck/build green.
* Manual visual QA (homepage + shell + sample console pages, both themes).
* No backend changes; backend build unaffected.
* Touch-point architecture documented.

## Decision (ADR-lite)

**Context**: Full reskin + rebrand of homepage + ~16 user + ~22 admin pages must be "unrecognizable" yet stay rebasable on upstream. Research showed (a) styling is centralized in `style.css` + `tailwind.config.js` (cascade covers ~80% of the app), (b) the shell (sidebar/header) + homepage carry bespoke structure that *defines* the original's look, so true unrecognizability requires editing them directly — which conflicts with the earlier "minimal upstream touch" goal.

**Decision**:
* **Phased delivery.** Phase 1 (this task, MVP) = design-system foundation + homepage + app shell + full rebrand (+ SVG illustrations + key-surface i18n). Phase 2 = user console pages. Phase 3 = admin console pages (separate follow-up tasks).
* **Tiered touch-points** (documented conflict surface, not a broad scatter) — see Technical Approach.
* **Defaults**: Geist + Geist Mono (npm `geist`); warm near-black dark-first with a working light theme; vermillion `#EF6F2E` accent; brand 硅基链/SiliconBase.
* **Homepage**: keep current section flow, transform visual language (not a factory IA transplant).
* **Rebrand**: full — name + logo + favicon + empty states + loading + 404. Logo + illustrations = authored geometric SVG (not AI raster art).
* **i18n**: key-surface brand-voice rewrite (en/zh); not a full-locale sweep.

**Consequences**: Convincing "unrecognizable" result, at the cost of a documented conflict surface (rebase effort concentrates there, not everywhere). Phase 2/3 remain follow-ups; Phase 1 alone re-frames every console page via the cascade + new shell.

## Technical Approach

**Tiered touch-points (rebasable conflict surface):**
* **Tier 0 (engine):** `frontend/tailwind.config.js`, `frontend/src/style.css`.
* **Tier 1 (new, zero-conflict):** `frontend/src/styles/theme-override.css` + one import in `frontend/src/main.ts`.
* **Tier 2 (shell + key pages):** `frontend/src/components/layout/{AppLayout,AppSidebar,AppHeader,AuthLayout}.vue`, `frontend/src/views/HomeView.vue`, `frontend/src/views/NotFoundView.vue`.
* **Brand assets:** `frontend/public/logo.png` (+ new SVG source), favicon, default `siteName` in `frontend/src/stores/app.ts`, doc-title fallback in `main.ts`.
* **i18n:** `frontend/src/i18n/` locale files (en/zh) — key-surface copy edits (homepage/auth/empty/nav/footer).
* **Tier 3 (Phase 2/3 only):** individual user/admin views — cascade-first; hand-edit only where the cascade falls short; huge files (`SettingsView` 10k lines) restyled via tokens, never line-by-line.

**Sync model:** `main` tracks `upstream/main` (ff-only); `custom/theme` long-lived branch rebased onto upstream. (Remotes `origin`=fork, `upstream`=source already configured.)

## Implementation Plan (small PRs)

* **PR1 — Design-system foundation:** add Geist fonts; `tailwind.config.js` token overhaul (warm dark + vermillion + low radius + no-shadow + Geist stacks + line-grid utility); create `theme-override.css` + import in `main.ts`; rewrite `style.css` `@layer base/components` to factory aesthetic for both themes. → whole app shifts to factory palette via cascade.
* **PR2 — Brand + app shell:** new logo mark (SVG) + favicon + default siteName; restyle `AppLayout` / `AppSidebar` / `AppHeader` / `AuthLayout` (brand block, mono section labels, hairline nav, inversion hover, warm-dark surface).
* **PR3 — Homepage:** restyle `HomeView` default mode (warm-dark hero + line-grid + tight `font-normal` H1 + mono eyebrow + vermillion CTA + factory-style terminal panel + **authored SVG hero illustration**; hairline feature cards; providers; footer) + apply brand wordmark/domain + **homepage en/zh copy (brand voice)**.
* **PR4 — Polish + content + QA:** empty states (**with SVG illustrations**), loading page, 404; **i18n (en/zh) brand-voice copy on auth/login/register, nav, footer, empty states**; spot-QA ≥1 admin + ≥1 user page (dark & light); fix token gaps; `lint/typecheck/build` green; document touch-points in `CLAUDE.md`/`REBRAND.md`.

## Out of Scope (explicit)

* Phase 2 (user console per-page polish) & Phase 3 (admin console per-page polish) — separate tasks (they inherit Phase 1's cascade + shell automatically).
* Raster / hand-painted / AI-generated illustration (only authored SVG is in scope).
* Full-locale copy rewrite of every console label/tooltip (only key-surface brand-voice copy is in scope; table/field copy unchanged).
* Backend, API, DB, business-logic, or route changes.
* New product features / new pages.
* Mobile-native / app builds.

## Research References

* [`research/factory-ai-design-language.md`](research/factory-ai-design-language.md) — factory.ai visual spec + Tailwind token starter + transferability notes (what to adopt vs. marketing-only).
* [`research/current-frontend-inventory.md`](research/current-frontend-inventory.md) — 57-route map, 4 layouts, styling touch-points, conflict-risk ranking, branding hooks.

## Technical Notes

* Stack: Vue 3 (Composition API) + Vite + Tailwind v3.4 (JS config) + Pinia + vue-i18n; pnpm (not npm).
* Build → `backend/internal/web/dist/` (embedded into Go binary via `-tags embed`).
* Branding runtime: `useAppStore` (`siteName`/`siteLogo`/`siteVersion`) from `PublicSettings` + `window.__APP_CONFIG__` pre-mount injection.
* Dark mode: `darkMode:'class'`, toggled in `AppSidebar`, persisted `localStorage.theme`; consumers read `document.documentElement.classList.contains('dark')`.
* Layout wiring: each view imports its own layout (no global layout map); public pages (home/404/key-usage) standalone.
