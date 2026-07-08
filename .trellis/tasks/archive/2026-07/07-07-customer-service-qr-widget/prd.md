# Customer Service QR Widget

## Goal

Add a global customer service entry on the right side of the application. Users can click a customer-service robot icon to open a floating panel that shows a support group QR code, so they can scan it with a phone and join the after-sales group.

## Requirements

* Add admin-configurable customer service widget settings:
  * Enabled/disabled switch.
  * QR code image URL, using the existing image upload pattern.
* Expose the enabled flag and QR code URL through public settings.
* Render a global floating widget from the application root.
* Hide the widget when disabled or when no QR code image is configured.
* Clicking the robot icon opens a floating panel with the QR code and scan guidance.
* Clicking outside/close closes the panel.
* Add English and Simplified Chinese i18n copy.

## Acceptance Criteria

* [ ] Admin can configure the widget enabled state and QR code image in system settings.
* [ ] `/settings/public` returns the widget enabled state and QR code URL.
* [ ] The global widget appears on regular application pages when enabled and configured.
* [ ] The widget is hidden when disabled or missing QR image.
* [ ] The QR panel can be opened and closed.
* [ ] Existing settings, auth, and layout behavior remains unchanged.

## Definition of Done

* Backend and frontend types are updated consistently.
* Regression tests cover public settings exposure and widget visibility.
* Frontend typecheck, lint, and targeted tests pass.
* Backend targeted tests pass.

## Out of Scope

* Real-time chat or bot conversations.
* Per-user customer service routing.
* New storage/upload backend; use the existing image upload capability.
