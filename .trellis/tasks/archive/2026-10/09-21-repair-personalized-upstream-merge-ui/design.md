# Technical Design: Repair Personalized Upstream Merge UI

## Change Boundary

The smallest behavior gap is in the frontend presentation layer: the current branch has an older/partial homepage composition, an upstream-style locale control, and a sidebar-owned theme control, while the personalized backup expects a factory-style homepage, icon-only locale control, and header-owned theme control.

The behavior lives in these boundaries:

- Public homepage composition and public navigation: `frontend/src/views/HomeView.vue`, its imported home visualization components, and the already-present landing locale messages/styles.
- Shared language control: `frontend/src/components/common/LocaleSwitcher.vue`.
- Console top navigation: `frontend/src/components/layout/AppHeader.vue`.
- Console sidebar: `frontend/src/components/layout/AppSidebar.vue`.
- Theme state synchronization: a shared `frontend/src/composables/useTheme.ts`, with initial class application remaining in `frontend/src/main.ts`.
- Regression coverage: the existing homepage and layout component tests under `frontend/src/views/__tests__` and `frontend/src/components/layout/__tests__`.

## Three-Way Reconciliation

| Surface | Current branch | Personalized backup | Merge decision |
| --- | --- | --- | --- |
| Default homepage | 468-line partial/factory composition; several personalized sections and content are absent | Full personalized v2/factory page with motion, hero dashboard, provider/tool marquees, defining cards, bento, CTA, and footer | Restore the personalized composition, but re-check route names, feature flags, and current upstream contracts before keeping each link/section |
| Compact homepage | Current behavior includes monitor and model-plaza entries | Backup intentionally omits monitor but keeps model-plaza | Keep current upstream entries unless they conflict with explicit personalization; preserve compact mode behavior and add only required shared-control styling |
| Locale control | Flag plus locale code in the trigger | Icon-only `languages` trigger with screen-reader label; no flags in menu rows | Adopt backup presentation and preserve current `setLocale`/dropdown behavior |
| Console header | Locale control but no theme icon | Theme icon in header, shared composable, personalized compact styling | Add header control using shared theme state and current translations |
| Sidebar footer | Theme toggle, local initialization, Sun/Moon icon definitions | No theme toggle and no local theme ownership | Remove only theme UI/state; preserve all sidebar nav and collapse behavior |
| Global theme bootstrap | `main.ts` initializes `dark` before mount | Same global bootstrap plus composable synchronization | Keep `main.ts` as the single bootstrap source; use the composable for component interaction |

## Data and Interaction Flow

1. `main.ts` applies the persisted/system theme class before Vue mounts.
2. `useTheme()` derives reactive `isDark` from the document class and exposes `toggleTheme()` / `setTheme()` for controls.
3. `HomeView` and `AppHeader` call the shared toggle; both update the same document class and `theme` storage key.
4. `LocaleSwitcher` remains a shared component mounted by both `HomeView` and `AppHeader`; it calls the existing i18n `setLocale` and does not own locale persistence beyond that API.
5. `AppSidebar` no longer owns theme state or renders a theme action.

## Compatibility and Risk Controls

- Preserve current route names and feature gates by comparing the personalized markup with `frontend/src/router/index.ts` and current feature-flag helpers before editing.
- Keep custom HTML/iframe homepage and compact homepage branches untouched except for shared controls that are explicitly in scope.
- Preserve current upstream `AppHeader` features such as support links, announcements, model plaza, subscriptions, balance, and user menu.
- Avoid copying the backup's unrelated sidebar navigation changes. Only remove the bottom theme action and unused theme code.
- Keep edits ASCII unless an existing file requires otherwise.

## Verification Design

- Add or update structural component tests for the locale trigger, header theme button, and sidebar absence. Where a test can mount the real interaction, assert the document class/storage behavior rather than only implementation details.
- Run the focused Vitest suites first, then frontend typecheck, lint check, and build.
- Start the frontend dev server if available and perform browser checks at desktop and mobile widths for homepage, homepage language switching, console header language/theme controls, and sidebar absence. Report browser verification separately if the environment cannot run it.
