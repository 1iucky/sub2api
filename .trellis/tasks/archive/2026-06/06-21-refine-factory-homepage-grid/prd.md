# Refine Homepage Factory.ai Background Grid

## Goal

Refresh the default homepage first-screen background so it matches the current Factory.ai visual language, while keeping console page surfaces quiet and readable with a pure theme-color background.

## Requirements

* Re-check the live Factory.ai homepage and base the change on current visual cues.
* Update the default homepage background/grid treatment.
* Keep console pages on a pure theme-color background rather than the point-grid treatment, so operational content stays easy to scan.
* Preserve the existing homepage layout, copy, dashboard hero, nav, marquee, bento cards, auth/settings behavior, and dark-mode support.
* Keep the grid subtle and low contrast; it should support the hero rather than compete with text or the dashboard frame.
* Do not use orange as a homepage background glow/color wash.
* Use a mixed dot-and-line treatment to create the black + white/gray metallic texture seen on Factory.ai.
* Increase the micro-dot density so the metallic texture reads more clearly.
* Add subtle brown-black edge blending plus black layered edge gradients on the left, right, and top edges, similar to the current lower background mixing.
* Restore the Gateway realtime/dashboard window interior to its plain panel background; it should not have an internal grid backdrop.
* Soften the console dark-theme shell color so navigation from the homepage dark hero to console pages feels smooth without making the console background pure black.
* Keep the background implemented locally in CSS, without depending on remote Factory.ai assets.

## Acceptance Criteria

* [ ] Homepage first screen uses a softer Factory.ai-inspired point-and-line background rather than the current uniform dense grid.
* [ ] Micro-dot density is visibly higher than the first point-line version.
* [ ] Top/left/right edges have subtle brown-black depth and black edge-to-center layering without becoming orange.
* [ ] Background works in light and dark modes.
* [ ] Homepage background has no orange glow/color wash.
* [ ] Console pages wrapped by the shared app layout use a pure theme-color surface with no point-grid background.
* [ ] Gateway realtime/dashboard window content area has no internal grid fill.
* [ ] The grid fades before the marquee/below-fold content.
* [ ] No new runtime dependency or remote image dependency is introduced.
* [ ] Frontend type-check/build succeeds.

## Definition of Done

* Lint/type-check/build or the closest practical frontend verification passes.
* Visual implementation is scoped to homepage background utilities and the shared app layout shell.
* Existing user changes outside this task are not reverted.

## Technical Approach

Use a Factory surface backdrop utility layered from CSS gradients for the homepage, and keep console layout backgrounds as pure theme-color surfaces:

* Base ambient wash matching Factory.ai's `#020202` / `#eeeeee` theme direction.
* Very faint large-scale line grid, closer to architectural drafting than a dense graph-paper grid.
* Micro-dot texture and diagonal/linear gray-white highlights to approximate the black + white/gray metallic finish.
* Mask fade so the homepage effect disappears before below-fold sections.
* Console shell uses a softer dark pure color (`dark-900`) instead of the homepage's near-black point-grid surface.

## Decision (ADR-lite)

**Context**: The existing homepage already references Factory.ai but the dense 18px grid reads too mechanical and high-frequency.

**Decision**: Replace the homepage-only `bg-blueprint-dense` usage with a new `bg-factory-home-grid` utility based on dot/line layers, tune shared blueprint utilities downward, and remove the SVG dashboard's internal grid fill.

**Consequences**: The homepage becomes closer to Factory.ai's quiet industrial surface while keeping the implementation lightweight. Future full-theme changes remain out of scope.

## Out of Scope

* Reworking homepage layout, content, dashboard SVG, or motion system.
* Re-theming admin/user app pages.
* Bundling or copying Factory.ai image assets.

## Technical Notes

* Main files: `frontend/src/views/HomeView.vue`, `frontend/src/components/layout/AppLayout.vue`, `frontend/src/styles/theme-override.css`.
* Factory.ai live page checked on 2026-06-21. Observed cues: dark default theme, body background around `#020202`, fixed blurred header surface, and very subtle point/line/metallic overlay treatment. The homepage background should not read as orange.
* Existing repo already has Factory palette tokens in `frontend/tailwind.config.js`.
