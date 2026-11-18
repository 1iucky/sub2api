# Implementation Plan: Repair Personalized Upstream Merge UI

## Ordered Checklist

1. Confirm the current branch remains unchanged except for the pre-existing working-tree files; record exact current/backup/upstream refs before edits.
2. Read the current frontend spec indexes and inspect the route, feature-flag, i18n, and theme contracts used by the homepage/header/sidebar.
3. Add focused failing regression assertions for:
   - the icon-only personalized locale trigger;
   - the header theme control and theme persistence/toggle behavior;
   - the absence of the sidebar-bottom theme item;
   - the default homepage's personalized structural markers/content.
4. Reconcile `HomeView.vue` with the personalized backup and current upstream. Restore only the intended personalized composition and copy the minimum related template/script/style/component/i18n dependencies needed for a valid current build.
5. Restore the icon-only locale switcher presentation in `LocaleSwitcher.vue`, keeping its existing locale selection and outside-click behavior.
6. Add `useTheme.ts` if needed and integrate it into `AppHeader.vue`; preserve all current header features and use existing localized theme labels.
7. Remove the sidebar theme menu, theme-specific icon declarations, local toggle function, and duplicate initialization from `AppSidebar.vue` without touching unrelated nav logic.
8. Run focused tests and fix only failures caused by this scope.
9. Run `pnpm run typecheck`, `pnpm run lint:check`, `pnpm run build`, and `git diff --check` from `frontend/`.
10. Perform browser/runtime checks if the dev server can start; inspect desktop/mobile homepage and console layouts and exercise locale/theme controls.
11. Review the final diff against `HEAD`, `upstream/main`, and the backup ref. Confirm no unrelated working-tree files were modified.

## Expected Files

- `frontend/src/views/HomeView.vue`
- `frontend/src/components/common/LocaleSwitcher.vue`
- `frontend/src/components/layout/AppHeader.vue`
- `frontend/src/components/layout/AppSidebar.vue`
- `frontend/src/composables/useTheme.ts` (new or reconciled)
- `frontend/src/views/__tests__/HomeView.compact.spec.ts` and/or a focused homepage test
- `frontend/src/components/layout/__tests__/AppHeader*.spec.ts`
- `frontend/src/components/layout/__tests__/AppSidebar.spec.ts`

Do not modify the pre-existing unrelated files `.gitattributes`, `deploy/docker-compose.yml`, `.agents/`, `.cursor/`, or unrelated backend files.

## Rollback Points

- Before production edits: only task artifacts are changed.
- After each UI surface: inspect `git diff` and run its focused test.
- If homepage reconciliation becomes incompatible with current upstream, preserve the current file and apply smaller section-level patches rather than restoring the backup wholesale.
