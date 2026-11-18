# Homepage light backdrop inversion

## Goal

Make the personalized homepage feel bright in light mode by reversing the light/dark treatment of its decorative backdrop. The user reports that the dense dots and dark upper shading look like an oppressive cloud and requests a visual trial.

## Confirmed evidence

- `frontend/src/views/HomeView.vue:120` uses `bg-factory-surface-grid` and `bg-blueprint-fade` for a first-viewport decorative layer.
- `frontend/src/styles/theme-override.css:57` defines the light backdrop: warm off-white base, overlapping dark edge/top gradients, brown radial shadows, dark dots, and mixed highlights/shadows. The top linear shadow alone has 18% opacity; dots have 11% opacity.
- `frontend/src/styles/theme-override.css:101` supplies a separate dark-mode background. Source inspection shows distinct theme definitions; it does not establish a runtime theme-toggle failure.
- `frontend/src/styles/theme-override.css:150` fades the backdrop out toward the first viewport bottom.
- `frontend/src/composables/useTheme.ts:10` toggles the document dark class and persists the preference.

## Requirements

- R1: In light mode, reverse the decorative backdrop's dark/light emphasis: dark top/edge shading becomes soft light, dark dots become light dots, and bright sheen becomes a very subtle neutral shade where needed to preserve texture.
- R2: Keep a warm off-white base and a restrained engineering texture; avoid trading the upper dark cloud for a new large dark patch elsewhere.
- R3: Preserve dot spacing, sparse grid geometry, and the gradual fade into the lower page.
- R4: Keep the existing dark appearance, typography, foreground colors, dashboard preview, branding, layout, and theme persistence behavior.

## Acceptance Criteria

- [ ] In light mode, the first viewport has no broad dark top/edge cloud; dots read as subtle light detail against a warm off-white surface (R1, R2).
- [ ] The background keeps its existing texture scale and fades smoothly into the page below (R3).
- [ ] Desktop and narrow-screen views preserve text contrast and show no background boundary or overflow (R2–R4).
- [ ] Switching light → dark → light restores each appearance; reloading preserves the chosen theme (R4).
- [ ] Before/after screenshots at matching viewport sizes are available for visual review; static checks alone do not establish visual acceptance.

## Out of scope

- Whole-page color inversion, dark-theme redesign, dashboard changes, backend changes, and the unrelated image-resolution investigation.

## Review status and risks

Lightweight PRD-only task. User approved implementation on 2026-10-22: “实现一版，我看下效果。只涉及浅色背景的调整，暗色的不要动”. Only light background styles may change; the dark CSS block must remain byte-for-byte unchanged. Exact highlight opacity needs browser visual tuning: white dots can disappear on an overly bright base. Final visual acceptance remains with the user after the trial.
