# Repair personalized upstream merge UI

## Goal

Restore the intentionally personalized public homepage and console controls that were lost while replaying the local customization branch onto the current `upstream/main`, while retaining newer upstream functionality and all unrelated current-branch work.

## Background and Confirmed Facts

- The working branch is `codex/upstream-main-personalized-replay-20260917`, based on `upstream/main` and currently ahead by 22 commits.
- The comparison backup is `codex/backup-main-before-upstream-20260917-eb3d604` at `eb3d6043c`.
- The working tree already contains unrelated user changes in `.gitattributes`, `deploy/docker-compose.yml`, `.agents/`, `.cursor/`, and `.trellis/`; these must remain intact.
- `frontend/src/views/HomeView.vue` is 468 lines on the current branch, versus 754 lines on the personalized backup and 744 lines on `upstream/main`. The current default home lost the backup's personalized v2/factory structure and related content.
- `frontend/src/components/common/LocaleSwitcher.vue` on the backup uses the personalized icon-only language control; the current branch still renders the upstream flag/code presentation.
- The personalized backup adds `frontend/src/composables/useTheme.ts` and moves the console theme toggle into `frontend/src/components/layout/AppHeader.vue`.
- The current `AppSidebar.vue` still owns the bottom theme-toggle UI and its local theme initialization. The requested personalized behavior is to remove that sidebar control while retaining theme behavior through the shared/header control.
- Theme initialization is already performed globally in `frontend/src/main.ts`; the repaired components must not reintroduce competing initialization logic.

## Requirements

### R1. Restore personalized homepage behavior and presentation

Reconcile `HomeView.vue` against the personalized backup and current upstream implementation. Restore the personalized default homepage structure, content sections, navigation presentation, responsive/mobile behavior, and associated styles/components/i18n keys that are absent from the current branch, while retaining current upstream route/feature-flag contracts and current custom-home/compact-home behavior.

Do not blindly restore the backup file. Preserve newer upstream homepage changes where they do not conflict with the requested personalization, and keep links feature-gated or route-compatible with the current router.

### R2. Restore personalized language-switcher styling on public and console surfaces

Make the shared `LocaleSwitcher` render the personalized control style used by both the homepage and console header, including its icon-only compact presentation and accessible label/title. Keep locale selection behavior, dropdown behavior, persistence, and all supported locale entries unchanged.

### R3. Move the console theme control to the top header

Add the personalized theme-color icon control to `AppHeader.vue`, backed by shared theme state and using the existing `nav.lightMode` / `nav.darkMode` translations. The control must toggle the global `dark` class and persist the selected theme in `localStorage`.

### R4. Remove the sidebar-bottom theme menu

Remove the theme-toggle menu item and its now-unused local implementation from `AppSidebar.vue`. The sidebar must not render a second theme control at its bottom; other navigation, collapse, and scroll behavior must remain unchanged.

## Constraints and Out of Scope

- Do not reset, checkout, cherry-pick, or wholesale copy branch files in a way that can overwrite current work.
- Do not alter backend behavior, deployment configuration, branding/payment changes, or unrelated frontend navigation.
- Do not remove legitimate upstream routes or feature flags merely because the backup predates them.
- Do not redesign the homepage beyond restoring the identified personalized version and reconciling it with current upstream contracts.

## Acceptance Criteria

- [ ] The default homepage visibly contains the personalized backup's intended structure/content and responsive navigation, while custom-home and compact-home branches still work.
- [ ] Current homepage links use routes and feature gates valid on the current branch; no broken `/status` or `/models` assumptions are introduced if the current router uses different paths.
- [ ] The homepage and console use the same personalized language-switcher appearance; the switcher remains keyboard/assistive-technology accessible and can change locale in both contexts.
- [ ] The console top header visibly contains one theme-color icon control; activating it changes the global theme and persists the choice.
- [ ] The bottom of the left sidebar contains no theme-switch menu/control, and no duplicate sidebar theme implementation remains.
- [ ] Existing targeted homepage/sidebar/header tests are updated or added to cover the restored contracts, and they fail before the production fix where practical.
- [ ] Frontend typecheck, relevant Vitest tests, lint check, and production build pass, with browser/runtime verification reported separately from static checks.
- [ ] Unrelated pre-existing working-tree changes remain untouched.

## Open Questions

None blocking. The implementation decision is to reconcile the personalized backup with current upstream rather than restore either branch wholesale.
