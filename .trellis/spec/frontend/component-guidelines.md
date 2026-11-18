# Component Guidelines

> How components are built in this project.

---

## Overview

<!--
Document your project's component conventions here.

Questions to answer:
- What component patterns do you use?
- How are props defined?
- How do you handle composition?
- What accessibility standards apply?
-->

(To be filled by the team)

---

## Component Structure

<!-- Standard structure of a component file -->

(To be filled by the team)

---

## Props Conventions

<!-- How props should be defined and typed -->

(To be filled by the team)

---

## Styling Patterns

<!-- How styles are applied (CSS modules, styled-components, Tailwind, etc.) -->

### Public Homepage Factory Styles

The default public homepage uses template classes from both Tailwind and
`frontend/src/styles/theme-override.css`. When reconciling or replaying
`HomeView.vue`, verify the paired global classes exist before considering the
homepage restored. In particular, `bg-factory-surface-grid`,
`bg-blueprint-fade`, `sf-logo-pill`, `sf-brand-mark`, `sf-defining-card`,
`sf-model-radar-*`, and `sf-sdlc-*` are not Tailwind utilities; they must be
kept in `theme-override.css` and covered by a focused style regression test.

Homepage navigation entries that point to optional routes must keep the same
feature-gate semantics as the router/sidebar. Channel Monitor uses the
`FeatureFlags.channelMonitor` opt-out contract, while Model Plaza uses
`FeatureFlags.modelPlaza` plus `model_plaza_require_auth` visibility.

### Homepage Gateway Dashboard

`GatewayDashboard.vue` has a solid theme-aware panel (`fill-white
dark:fill-dark-900`), not an SVG grid backdrop. Keep the homepage's outer
`bg-factory-surface-grid` texture independent of this panel.

Preserve the personalized Claude/OpenAI/Gemini mini charts when replaying
upstream changes: distinct curve colors, `/min` SVG `tspan` labels, growth
percentages, time ticks, and focusable data points with tooltips and crosshairs.
The `sf-mini-chart-*` styles in `theme-override.css` provide hover and
`:focus-visible` states. Verify both states in a browser; mounted component
tests cannot validate CSS visibility. Cover chart data/geometry and the absence
of a panel grid in the focused GatewayDashboard regression suite.

---

## Accessibility

<!-- A11y requirements and patterns -->

(To be filled by the team)

---

## Common Mistakes

<!-- Component-related mistakes your team has made -->

(To be filled by the team)
