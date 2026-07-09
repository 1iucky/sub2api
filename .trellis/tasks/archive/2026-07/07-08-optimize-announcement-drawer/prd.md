# Optimize Announcement Drawer Display

## Goal

Improve the console announcement experience so clicking the top-bar notification icon opens a right-side drawer that matches the refactored console UI. Announcements should be categorized (for example: announcements, model updates, changelog) and displayed as category-specific timelines while reusing the existing announcement visibility, read-status, and targeting logic.

## What I already know

* The current top-bar notification entry is `frontend/src/components/common/AnnouncementBell.vue`.
* The current UI uses a centered modal for the list and a second centered modal for details.
* User announcement data is loaded through `frontend/src/stores/announcements.ts` and `frontend/src/api/announcements.ts`.
* User permissions/visibility are handled by backend `AnnouncementService.ListForUser`; it filters active announcements by schedule, targeting, subscription groups, balance, and read state.
* Admin announcement management lives in `frontend/src/views/admin/AnnouncementsView.vue`.
* Current announcement model fields include title, content, status, notify_mode, targeting, starts_at, ends_at, read_at, created_at, updated_at.
* There is no existing category field in `backend/internal/domain/announcement.go`, `backend/ent/schema/announcement.go`, DTOs, or frontend `Announcement`/`UserAnnouncement` types.
* The attached screenshot shows a right drawer with title "通知", category tabs with counts, and timeline-style announcement entries.

## Requirements

* Replace the click behavior of the top-bar notification icon from a centered announcement list modal to a right-side drawer.
* Keep the top-bar notification icon and unread indicator behavior.
* Add announcement categories with at least:
  * Announcement / 公告
  * Model updates / 模型动态
  * Changelog / 更新日志
* Admin create/edit announcement forms must allow selecting a category.
* Existing announcements should remain valid after migration and default to the announcement category.
* User announcement API should return the category so the frontend can group data.
* Admin announcement API should return and accept the category so admins can manage it.
* The drawer must group/filter announcements by category tabs.
* The drawer must not include an "All" tab.
* Each category tab must display the number of currently visible announcements in that category.
* Category tab active styling should use the system orange accent instead of the previous blue accent.
* The default selected category should be the first category with data; if all categories are empty, default to Announcement / 公告.
* Each category list should render as a timeline ordered by existing announcement ordering, newest first.
* Timeline entries should place the timestamp beside the timeline dot, with the announcement title/content below it.
* Timeline dots should use the system orange solid fill for read entries.
* Announcement titles inside the drawer should be visually prominent with a larger, bold font.
* Long Markdown announcement content in the drawer should show a partial preview by default and provide a manual expand/collapse control to read the full content.
* Drawer open/close animation should feel smoother and less abrupt than the old modal transition.
* Category tabs should show counts per category.
* Existing permissions, targeting, schedule, read state, unread count, popup notify mode, and mark-read behavior must be reused and not redesigned.
* The implementation should stay within the current Vue/Tailwind/component patterns and avoid introducing a new UI dependency.

## Acceptance Criteria

* [ ] Clicking the top-bar notification icon opens a right-side drawer rather than a centered list modal.
* [ ] The drawer can be closed by the close button, overlay click, and Escape key.
* [ ] The drawer shows only the requested category tabs: announcement, model updates, and changelog.
* [ ] Each category tab shows a count based on currently visible announcement data.
* [ ] Active category tab uses the system orange border/background/count style.
* [ ] Announcements are displayed as timeline entries with timestamp, title, content preview/content, and unread/read visual state.
* [ ] Timeline timestamps appear beside the timeline dot, and read dots use system orange solid fill.
* [ ] Drawer announcement titles are larger and bolder than body text.
* [ ] Long Markdown announcement content is collapsed by default, can be expanded manually, and shows full rendered content after expansion.
* [ ] Drawer transition is slower/smoother than the first implementation.
* [ ] Counts in category tabs reflect the currently visible announcement data after backend permission filtering.
* [ ] Selecting a category changes the timeline list without refetching or bypassing existing permissions.
* [ ] Opening/reading an announcement continues to mark it as read through the existing API path.
* [ ] Admin create/edit supports selecting category, and update payload only sends changed category when editing.
* [ ] Backend stores category with a safe default for existing data.
* [ ] Existing `notify_mode=popup` behavior continues to work.
* [ ] Frontend i18n is updated for zh/en category labels and drawer copy.
* [ ] Focused tests cover category field plumbing and drawer/timeline behavior.

## Definition of Done

* Tests added/updated for frontend drawer/category behavior and backend category contract.
* Frontend `typecheck` and `lint:check` pass.
* Relevant backend unit tests pass.
* Migration is added for the announcement category column.
* No unrelated dirty files are staged or reverted.

## Out of Scope

* Changing announcement targeting/permission rules.
* Changing existing read-status semantics.
* Adding rich category management/custom categories.
* Replacing the separate popup notify-mode reminder flow unless required by the drawer refactor.
* Adding pagination/infinite scroll to the drawer unless the current 20-item store limit becomes a blocker.

## Technical Approach

* Add a backend category enum/constant set to the announcement domain/service layer.
* Add a non-null announcement category column with default `announcement` through Ent schema and SQL migration.
* Add category to repository mapping, DTOs, admin create/update requests, and frontend types/API payloads.
* Add a category selector to `AnnouncementsView.vue`.
* Refactor `AnnouncementBell.vue`:
  * Keep store usage and mark-read actions.
  * Replace the list modal with a right drawer.
  * Render tabs and timeline groups derived from the already-loaded `announcements`.
  * Preserve popup notify-mode flow via `AnnouncementPopup.vue` / store.
* Add or update tests near existing announcement/admin tests.

## Open Questions

* Resolved: do not include an "All / 全部" tab. Show counts on each category tab.

## Technical Notes

* Likely impacted frontend files:
  * `frontend/src/components/common/AnnouncementBell.vue`
  * `frontend/src/stores/announcements.ts`
  * `frontend/src/views/admin/AnnouncementsView.vue`
  * `frontend/src/types/index.ts`
  * `frontend/src/api/admin/announcements.ts`
  * `frontend/src/i18n/locales/zh.ts`
  * `frontend/src/i18n/locales/en.ts`
* Likely impacted backend files:
  * `backend/internal/domain/announcement.go`
  * `backend/internal/service/announcement_service.go`
  * `backend/internal/repository/announcement_repo.go`
  * `backend/internal/handler/dto/announcement.go`
  * admin/user announcement handlers
  * `backend/ent/schema/announcement.go`
  * `backend/migrations/*.sql`
