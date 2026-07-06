# Fix Chinese Labels for Subscription Request Limits

## Goal

Fix the admin order management -> subscription plans create/edit form so the three request limit labels render Simplified Chinese text instead of raw i18n keys.

## What I Already Know

* The affected component reads `payment.admin.dailyRequestLimit`, `payment.admin.weeklyRequestLimit`, and `payment.admin.monthlyRequestLimit` in `frontend/src/views/admin/orders/PlanEditDialog.vue`.
* English locale defines these keys under `payment.admin`.
* Simplified Chinese locale currently defines matching strings under `admin.groups`, not under `payment.admin`.
* This explains why Chinese mode displays raw keys like `payment.admin.dailyRequestLimit`.

## Requirements

* Add the missing Simplified Chinese `payment.admin` translations for the daily, weekly, and monthly request limit labels.
* Add the matching request limit placeholder in the same `payment.admin` scope.
* Keep existing `admin.groups` translations intact because they may be used by group management.
* Do not change backend behavior or subscription plan payload handling.

## Acceptance Criteria

* [ ] In Simplified Chinese, the subscription plan create/edit dialog labels show Chinese text for daily, weekly, and monthly request limits.
* [ ] The request limit placeholder in the same form shows Chinese text.
* [ ] English locale remains unchanged.
* [ ] A focused frontend test verifies the Chinese `payment.admin` keys exist and resolve to the expected strings.

## Definition of Done

* Focused test passes.
* Relevant frontend lint/typecheck or nearest available verification passes, or any existing unrelated failure is documented.
* No unrelated dirty files are included in this task.

## Out of Scope

* Redesigning the subscription plan form.
* Changing request limit validation or API payload behavior.
* Cleaning up broader duplicate locale structure.

## Technical Notes

* Root cause investigation: `rg` showed the component uses `payment.admin.*`; `frontend/src/i18n/locales/en.ts` has those keys; `frontend/src/i18n/locales/zh.ts` only had those strings under `admin.groups`.
