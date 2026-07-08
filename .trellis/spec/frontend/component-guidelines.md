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

### Pattern: Remote Search Select

Use `frontend/src/components/common/Select.vue` for dropdowns that need the project-standard trigger, popover, option, loading, and empty states.

**Contract**:
- `searchable`: set to `true` when the search input must always be visible.
- `@search`: receives the raw search text whenever the dropdown search input changes.
- `loading` / `loadingText`: show an in-dropdown loading state while remote data is being fetched.
- `filter-options`: set to `false` when options already come from a remote query. Leave the default `true` for local-only selects.
- `emptyText`: provide an action-oriented empty state, especially before the first keyword is typed.

**Remote search behavior**:
- Do not fetch a large default option page on mount or on dropdown open just to support fuzzy search.
- Trigger remote lookup from `@search` after trimming and debouncing the keyword in the owning component.
- Clear options for empty keywords unless the product explicitly needs default suggestions.
- Abort or invalidate stale requests so older responses cannot overwrite newer search results.

**Example**:
```vue
<Select
  :model-value="null"
  :options="remoteOptions"
  searchable
  :filter-options="false"
  :loading="searching"
  :empty-text="searchEmptyText"
  @search="onRemoteSearch"
/>
```

### Pattern: Shared Public and Console Pages

When a public page must also appear inside the authenticated console layout,
reuse the same view component and pass a small route prop to control only the
outer page chrome. Do not duplicate the page content into a second component.

**Contract**:
- The public route keeps the default view behavior.
- The console route passes a typed prop such as `embedded: true`.
- The shared view uses that prop only for layout chrome such as a public top
  navigation bar and console `AppLayout`; filters, cards, API calls, loading states, and drawers stay
  shared.
- Add a focused component test that mounts the view with the console prop and
  asserts public-only chrome is hidden and the console layout slot is used.

**Example**:
```ts
{
  path: '/marketplace',
  component: () => import('@/views/user/ModelMarketplaceView.vue'),
  props: { embedded: true },
}
```

---

## Props Conventions

<!-- How props should be defined and typed -->

(To be filled by the team)

---

## Styling Patterns

<!-- How styles are applied (CSS modules, styled-components, Tailwind, etc.) -->

(To be filled by the team)

---

## Accessibility

<!-- A11y requirements and patterns -->

(To be filled by the team)

---

## Common Mistakes

<!-- Component-related mistakes your team has made -->

(To be filled by the team)
