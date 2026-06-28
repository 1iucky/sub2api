# Region blocked page and Nginx GeoIP guide

## Goal

Add a public unsupported-region page for deployments that block Chinese mainland source IPs at Nginx, and document a concrete Nginx + GeoIP2 configuration that redirects blocked users to that page.

## What I already know

* The frontend is Vue 3 + Vue Router.
* Public routes live in `frontend/src/router/index.ts`.
* Public pages use the app store for `siteName`, `siteLogo`, and cached public settings.
* The user wants the page to show logo, site name, source IP geolocation, unsupported-region copy, and a retry button.
* The retry button should redirect to the site home address so Nginx re-checks IP geolocation.
* Nginx can provide source IP / country / region details to the static page via query params on redirect.
* Existing deployment docs include Caddy and Docker docs, but no GeoIP2 region-blocking recipe.

## Assumptions

* New route path: `/unsupported-region`.
* The frontend page reads query parameters such as `ip`, `country`, `region`, and `path`.
* Nginx performs the authoritative IP geolocation decision; frontend only displays the received information.
* Query values are treated as untrusted display-only strings and never as HTML.
* Retry button navigates to `/home` with a cache-busting query parameter to force a fresh Nginx request path.
* Documentation should live under `docs/deployment/nginx-geoip2-region-block.md`.

## Requirements

* Add a public unsupported-region route.
* Page must show site logo and configured site name.
* Page must show an unsupported-region icon/visual similar to the provided screenshot.
* Page must show the visitor source IP attribution from query params.
* Page must explain that the region is not supported.
* Page must include a retry button that navigates to `/home` so Nginx checks IP geolocation again.
* Page must be responsive and match the current public frontend visual language.
* Add i18n copy in Chinese and English.
* Add deployment documentation with a concrete Nginx + GeoIP2 example.
* Nginx config must avoid redirect loops for `/unsupported-region` and static assets.
* Nginx config must explain real client IP handling behind CDN/proxies.

## Acceptance Criteria

* [ ] Visiting `/unsupported-region?ip=1.2.3.4&country=China&region=Guangdong&path=/usage` renders a branded unsupported-region page.
* [ ] Page displays IP/location text without unsafe HTML rendering.
* [ ] Retry button redirects to `/home`.
* [ ] Public route does not require authentication.
* [ ] Nginx documentation includes GeoIP2 database/module prerequisites, `map` logic for CN blocking, redirect location, and proxy headers.
* [ ] Frontend typecheck, lint, and targeted tests pass.

## Definition of Done

* Implementation follows existing public Vue page patterns.
* No backend API change is required.
* Deployment recipe is concrete enough to copy and adapt.
* Existing unrelated dirty files remain untouched.

## Out of Scope

* Backend-side GeoIP lookup.
* Automatically downloading MaxMind/IP2Location databases.
* Admin UI for region policy management.
* Blocking non-China regions.

## Technical Notes

* Relevant frontend files: `frontend/src/router/index.ts`, `frontend/src/stores/app.ts`, `frontend/src/i18n/locales/zh.ts`, `frontend/src/i18n/locales/en.ts`.
* Existing public brand fallback: `siteName || 'SiliconBase'`, `siteLogo || '/logo.svg'`.
* Nginx should use `$realip_remote_addr` or `$remote_addr` after `real_ip_header` is configured for trusted proxies.
