# Design System Spec — 硅基链 / SiliconBase (Factory.ai aesthetic)

Authoritative token + component-class spec for the reskin. **PR1 implements this** (tailwind config + style.css + theme-override.css); PR2–4 extend it to shell/homepage/brand. Exact values below — do not guess.

## Principles
1. **Keep token NAMES, repoint VALUES.** `primary`, `accent`, `dark`, `gray` keep their names so every existing `@apply`/utility cascades. Rename nothing.
2. **Dark-first.** Dark theme is the star (warm near-black). Light theme must still work (warm off-white).
3. **Warm, not cool.** Replace cool slate/zinc with warm neutral grays. Achieved by **overriding Tailwind `gray`** to a warm scale (warms the whole UI at once).
4. **No shadows, no glass glow.** Elevation = hairline borders + tone separation. Keep shadow *token names* (set to `none`) to avoid breaking `@apply`, but **remove shadow usages** in component classes.
5. **Low radius.** Do NOT override Tailwind's radius scale. Reduce radius inside component classes: buttons/inputs `rounded-sm`(3px)–`rounded`(4px); cards `rounded-md`(6px)–`rounded-lg`(8px).
6. **Color-inversion hover** for buttons (150ms color swap; no lift/scale/shadow).
7. **Mono uppercase eyebrows** (Geist Mono, 12px, tracking-wide) for section labels.
8. **Respect `prefers-reduced-motion`** (modals already do; keep).

## Color tokens (`tailwind.config.js` → `theme.extend.colors`)

### `primary` — vermillion / burnt-orange (signature). Anchor 500 ≈ `#ef6f2e`
```
50:#fdf4ef 100:#f9e1d3 200:#f3c3a6 300:#ec9e72 400:#e67f47
500:#ef6f2e 600:#d8581a 700:#b3430c 800:#8f3710 900:#743012 950:#411706
```
Primary actions, active states, links, focus, selection, key metrics.

### `gray` — OVERRIDE default cool gray → warm neutral (warms the whole UI)
```
50:#f7f5f3 100:#ece7e1 200:#d9d1c7 300:#bfb3a5 400:#a09486
500:#7d7264 600:#5e554a 700:#4a4339 800:#2b2722 900:#1a1714 950:#0c0a08
```

### `dark` — warm near-black surface scale (dark-theme surfaces)
```
50:#f7f5f3 100:#e8e3dd 200:#c9c0b5 300:#9a8f80 400:#6f655b
500:#4d4947 600:#3d3a39 700:#2e2c2b 800:#1a1816 900:#0a0908 950:#020202
```

### `accent` — repoint to warm neutral (same scale as `gray` above). Used in `.text-gradient`, sidebar, etc.

### Status colors — keep Tailwind defaults: `emerald`(success), `red`(danger), `amber`(warning), `blue`/`purple`(info). Payment brand hexes (`btn-stripe` #635bff, `btn-alipay` #00AEEF, `btn-wxpay` #2BB741, `btn-airwallex` #14171A/#7AF0C4) stay as-is.

## Fonts (`theme.extend.fontFamily`)
```
sans: ['Geist','Geist Fallback','ui-sans-serif','system-ui','-apple-system','BlinkMacSystemFont','Segoe UI','PingFang SC','Hiragino Sans GB','Microsoft YaHei','sans-serif']
mono: ['Geist Mono','Geist Mono Fallback','ui-monospace','SFMono-Regular','Menaco','Monaco','Consolas','monospace']
```
Load via `pnpm add geist`; import Geist CSS in `theme-override.css`.

## Shadows (`theme.extend.boxShadow`) — keep NAMES, set to `none`
```
glass:'none'  'glass-sm':'none'  glow:'none'  'glow-lg':'none'
card:'none'  'card-hover':'none'  'inner-glow':'none'
```

## Backgrounds (`theme.extend.backgroundImage`) — keep names, repoint
```
'gradient-primary': 'linear-gradient(135deg,#ef6f2e 0%,#d15010 100%)'
'gradient-dark':     'linear-gradient(135deg,#1a1714 0%,#020202 100%)'
'gradient-glass':    'linear-gradient(135deg,rgba(255,255,255,0.06) 0%,rgba(255,255,255,0.02) 100%)'
'mesh-gradient':     'radial-gradient(at 40% 20%,rgba(239,111,46,0.10) 0px,transparent 50%),radial-gradient(at 80% 0%,rgba(209,80,16,0.06) 0px,transparent 50%),radial-gradient(at 0% 50%,rgba(239,111,46,0.06) 0px,transparent 50%)'
'gradient-radial': (keep default)
```

## Radius — do NOT override Tailwind scale; reduce radius inside component classes. Keep existing `borderRadius.4xl`.

## Component-class translation (`src/style.css` — KEEP ALL CLASS NAMES & SELECTORS)
General: `rounded-xl/2xl`→`rounded-sm`/`rounded`/`rounded-md`; remove `shadow-*` usages; warm surfaces via new `gray`/`dark`; both light & dark variants on every class.

- `.btn`: `rounded-sm` ~h-9 px-4, `text-sm font-medium`, `transition-colors duration-150`, `focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500`. Remove `active:scale`, remove `focus:ring`.
- `.btn-primary`: `bg-primary-500 text-white` hover `bg-primary-600`. No gradient/shadow. (Inversion option: hover `bg-gray-950 text-primary-400`.)
- `.btn-secondary`: light `bg-white text-gray-800 border border-gray-300` / dark `dark:bg-dark-900 dark:text-gray-200 dark:border-dark-700`; hover color-swap, no shadow.
- `.btn-ghost`: transparent, muted text; hover `bg-gray-100 dark:bg-dark-800`.
- `.btn-danger/success/warning`: solid status fills, no gradient/shadow, darker on hover.
- `.btn-sm/md/lg/icon`: keep size deltas, `rounded-sm`.
- `.input`: `rounded-sm`, hairline `border border-gray-300 dark:border-dark-700`, `bg-white dark:bg-dark-900`, focus `border-primary-500` + outline. No ring.
- `.card`: `rounded-md border border-gray-200 dark:border-dark-700 bg-white dark:bg-dark-900`. NO shadow.
- `.card-hover`: hover `border-gray-300 dark:border-dark-600` (no translate/shadow).
- `.glass`/`.glass-card`/`.card-glass`: translucent warm surface + `backdrop-blur` + hairline border, NO `shadow-glass`.
- `.table`: `th` `font-mono font-medium uppercase tracking-wide text-xs bg-gray-50 dark:bg-dark-900` + `border-b border-dashed border-gray-200 dark:border-dark-700`; `td` hairline bottom; row hover tone.
- `.sidebar`: `bg-white dark:bg-dark-950 border-r border-gray-200 dark:border-dark-800`. `.sidebar-link` `rounded-sm` + hover tone; `.sidebar-link-active` vermillion text + `bg-primary-50 dark:bg-primary-500/10` (NO fill gradient); `.sidebar-section-title` mono eyebrow `font-mono text-xs uppercase tracking-wide text-gray-400 dark:text-dark-400`.
- `.badge`: `rounded-sm` or `rounded-full`, mono `font-medium`.
- `.modal-content`/`.dialog-container`/`.dropdown`/`.toast`: `rounded-lg`/`rounded-sm`, hairline border, `bg-white dark:bg-dark-900`, no shadow (or barely-there).
- `.page-title`: `text-2xl font-normal tracking-tight` (NOT bold).
- `.code`/`.code-block`: Geist Mono; inline vermillion; block dark surface.
- `.tabs`/`.tab`: hairline; active = tone swap, no shadow.
- body: `bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-gray-100` (warm via gray override); `::selection bg-primary-500/20`.
- scrollbars: keep, warm tones.

## Utilities to add in `theme-override.css`
- Geist `@import` (from `geist` package CSS).
- `.bg-blueprint` — faint warm line-grid:
  `background-image:linear-gradient(rgba(125,114,100,0.06) 1px,transparent 1px),linear-gradient(90deg,rgba(125,114,100,0.06) 1px,transparent 1px);background-size:32px 32px;`
- `.divider-dashed` — dashed hairline divider.
- `.eyebrow` — `font-mono text-xs uppercase tracking-[0.12em] text-gray-500 dark:text-dark-400`.

## Import order — `main.ts`: `import './style.css'` then `import './styles/theme-override.css'` (override AFTER base).
