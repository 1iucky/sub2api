# Research: Factory.ai Visual Design Language

- **Query**: Produce an implementation-ready design spec (Tailwind token system) for a frontend reskin in the visual style of https://factory.ai/
- **Scope**: external (web/brand research)
- **Date**: 2026-06-18
- **Method**: PRIMARY — `curl` fetch of live production HTML + CSS bundles from factory.ai, then static analysis of Tailwind class usage, CSS custom properties (`--*`), hex/oklch color frequency, font-family declarations, radius/shadow/gradient rules, and keyframe definitions. **chrome-devtools MCP was NOT available** in this environment, so no rendered screenshots were captured; all values below are mined directly from the shipped source. A fallback `mcp__exa__web_search_exa` pass was not required because the raw assets gave exhaustive concrete data.

**Pages analyzed**:
- `/` (homepage, `factory_home.html`, ~479 KB HTML)
- `/pricing` (pricing page, `factory_pricing.html`, ~460 KB HTML)
- CSS bundle 1: `/_next/static/css/2f8a58de57cfbf1c.css` (~193 KB — main Tailwind v4 output)
- CSS bundle 2: `/_next/static/css/84edd78de243cde6.css` (~3.7 KB — `@font-face` for Geist)

**Tech stack detected**: Next.js App Router + React Server Components (`data-radix-collection-item`, RSC streaming markers), **Tailwind CSS v4** (signaled by `oklch()` color tokens, `--spacing:.25rem` multiplier, `--radius-*` / `--text-*` / `--container-*` custom-property scale, and `@theme`-style token names). Default theme is **dark** (`<html data-theme="dark">`); a `.light:` variant namespace exists for light mode.

---

## 1. Overall Aesthetic

**Keywords**: dark, monochrome-near-black, warm-tinted-neutral, monospaced-label, fine-line-grid, blueprint/engineering, high-contrast inversion, low-radius (sharp), editorial, restrained-accent-orange, terminal-cadence.

**Prose vibe**: Factory reads as a *blueprint for software* — a near-black canvas layered with a barely-visible fine line-grid (like graph paper / a CAD drawing), warm-tinted neutral grays instead of pure zinc, and a single restrained vivid-orange accent (#EF6F2E) used surgically. Headings are large, set in **Geist**, and noticeably *not bold* (`font-normal`) with aggressively tight negative letter-spacing (down to `-0.18rem`), giving a confident, engineered feel. Navigation, eyebrow labels, and metadata are set in **Geist Mono uppercase at 12px**, evoking a terminal/IDE chrome. Buttons are small (31px tall), sharp-cornered (3px radius), and invert foreground/background on hover rather than lifting. The whole thing feels like the UI of a developer tool that happens to be a marketing site — quiet, precise, technical, with a "cursor blink" animation motif. Default theme is **dark**; light theme is offered but secondary.

---

## 2. Color Palette

All values below are extracted from the shipped CSS custom properties. Factory uses **oklch()** for its neutral scale (Tailwind v4 default) plus warm-tinted hex overrides. Note the neutral scale is **warm** (browns/ochres), NOT cool blue-gray like Tailwind's default `slate`/`zinc`.

### Brand accent (signature orange) — `--accent-*`
| Token | Hex | Usage |
|---|---|---|
| `--accent-100` | `#EF6F2E` | Primary accent — most used (text/bg/border). Tailwind `text-accent-100`, `bg-accent-100`, `border-accent-100`. |
| `--accent-200` | `#ee6018` | Slightly deeper — `bg-accent-200` used for emphasis backgrounds. |
| `--accent-300` | `#d15010` | Darkest shade — badge backgrounds, focus-within borders. |
| (derived) | `oklch(70.5% .213 47.604)` | Tailwind `orange-500` — used for *link hover* (`hover:text-orange-500`) on footer/nav links. |

The orange is a warm vermillion/burnt-orange — not iOS-orange, not GitHub-blue. ~HSB(18°, 81%, 94%).

### Base / neutral scale (warm, dark-first)

Dark theme uses a near-black base with a *secondary* slightly-lighter panel tone; the "base-1000…900…100" scale flips between themes. Concrete hex hits found in CSS:

| Role (dark theme) | Hex / value | Note |
|---|---|---|
| Page background (`--dark-base-primary`) | `#020202` | Body bg. Near-pure black with the faintest warm tint. |
| Surface / panel (`--dark-base-secondary`) | `#0a0908` (and `#020202`) | Footer/cards, slightly lifted from bg. Very subtle. |
| Elevated / card border | `#3d3a39`, `#342f2d`, `#2e2c2b` | Mid-warm-grays for hairlines and borders. |
| Muted text (`base-700`) | `#4d4947` (and `oklch(37.1% 0 0)` neutral-700) | Secondary/muted text. |
| Body / secondary text (`base-400`) | `#a49d9a` | The `text-base-400` used for nav mono labels — warm gray. |
| Primary text (light theme) | `#fafafa` | `--light-base-primary` — also used as the *inverted* button text on dark buttons. |
| Light-theme bg | `#fafafa` | `--light-base-primary` doubles as light page bg. |

### Inversion pairs (critical to the design system)
The whole UI is built on **two opposing base tones that swap on hover/active**:
- `dark-base-primary` `#020202` ↔ `light-base-primary` `#fafafa`
- `dark-base-secondary` `#0a0908` ↔ `light-base-secondary` `#fafafa`

Buttons and hover states literally swap foreground↔background between these. This is the single most load-bearing pattern — implement it as paired tokens, not as two unrelated colors.

### Other accent hues used sparingly (status / categorization, Tailwind v4 oklch defaults)
- `emerald-400/500` — success / "available"
- `red-400/500` — error / destructive
- `yellow-400/500`, `green-400/500`, `blue-400/500` — category tags
- Several warmer secondary accents appear in the homepage visual palette: `#febc2e` (amber), `#ff6b5f` / `#ff5f57` (coral/red — also the macOS window-close red), `#a3d28b` (sage green), `#9fd0ff` (sky), `#7db5e8` (steel blue), `#f3b15d` (sand), `#e6a24d` (ochre), `#b46a35` (sienna), `#e54048` (warm red). These appear to be illustrative/product-screenshot colors, not core UI tokens.

### Signature "gradient"
- The site does **not** rely on bright CSS gradients for its identity. The only `linear-gradient` in core CSS is a semi-transparent white wash layered over the line-grid texture: `linear-gradient(#fff9,#fff9), url(/assets/bg-lines.png)` (used in light theme to dim the texture).
- There *is* a raster `/_images/gradient-11.png` preloaded on the homepage (a soft atmospheric gradient glow, likely behind the hero). Treat gradients as *ambient/atmospheric background imagery*, not as brand color.

---

## 3. Typography

### Font families
| Role | Family | Source |
|---|---|---|
| Body / sans | **Geist** (variable webfont), with fallbacks `Geist Fallback`, `ui-sans-serif, system-ui, sans-serif` | Self-hosted woff2 at `/_next/static/media/*-s.p.woff2` (Vercel Geist, MIT-licensed, available on Google Fonts and npm `geist`). |
| Headings | **Geist** (same family — *not* a display face). Distinction comes from size + tracking, not family. | Same as above. |
| Code / mono | **Geist Mono** (variable), fallbacks `Geist Mono Fallback`, `ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace` | Same self-hosted woff2 bundle. |
| Emoji fallback | `"Apple Color Emoji","Segoe UI Emoji","Segoe UI Symbol","Noto Color Emoji"` | Tailwind default. |

> **Install**: Vercel's Geist + Geist Mono are open-source. For a Vue/Vite project use `pnpm add geist` (npm package ships the fonts + a CSS import) or self-host the woff2 files. **Recommend the npm `geist` package** for the closest match.

### Type scale (concrete values from `--text-*` tokens + inline classes)
| Token | Size | Typical use |
|---|---|---|
| `--text-xs` | `12px` (.75rem) | Mono uppercase labels, eyebrow nav, metadata, badges. |
| `--text-sm` | `14px` (.875rem) | Body copy, nav links, form labels, footer links. **Most common body size.** |
| `--text-base` | `16px` (1rem) | Default body, larger paragraph text. |
| `--text-lg` | `18px` (1.125rem) | Card titles, section sub-heads. |
| `--text-xl` | `20px` (1.25rem) | Small section heads. |
| `--text-2xl` | `24px` (1.5rem) | Section heads (mobile). |
| `--text-3xl` | `30px` (1.875rem) | Section heads (desktop). |
| `--text-4xl` | `36px` (2.25rem) | Large section heads. |
| `--text-5xl` | `48px` (3rem) | — |
| Hero H1 (custom) | `40px` → `72px` (`text-[40px] lg:text-[72px]`) | Homepage hero. |
| Display (custom) | `clamp(80px, 10vw, 200px)` | Occasional oversized display number. |

### Weights
- `--font-weight-light` (300), `--font-weight-normal` (400), `--font-weight-medium` (500), `--font-weight-semibold` (600), `--font-weight-bold` (700).
- **Key choice**: Headings use `font-normal` (400), NOT bold. H1 is literally `class="… font-normal text-[72px] …"`. Bold/semibold is reserved for small emphasis (table headers, prices). This is a defining trait — do **not** bold your H1/H2 in the reskin.

### Line-height
- Default body: `line-height: 1.6` (very readable). Also seen: `1.5`, `1` (for tight UI labels), `110%`/`100%` on display headings.
- Hero H1: `leading-[100%]` (line-height 1.0 — extremely tight, oversized).
- Section H2: `leading-[110%]` → `lg:leading-[1.12]`.

### Letter-spacing — the second most defining trait
Large headings are pulled *very tight* (negative tracking); mono labels are pulled tight or set wide.

| Context | Value |
|---|---|
| Hero H1 (72px) | `lg:tracking-[-0.18rem]` |
| Hero H1 (40px) | `tracking-[-0.16rem]` |
| Section H2 (36px) | `lg:tracking-[-0.07rem]` |
| Section H2 (24px) | `tracking-[-0.045rem]` |
| Mono nav labels (12px) | `tracking-[-0.015rem]` |
| Mono card labels (18px) | `lg:tracking-[-0.0225rem]` |
| Eyebrow / uppercase labels (wide) | `.08em`, `.14em`, `.18em` (positive — spread caps) |

Rule of thumb: **big text → tighten dramatically (-0.04rem to -0.18rem); tiny uppercase mono → tighten slightly or spread wide**.

### Distinctive treatments
- **Uppercase mono labels everywhere** — section eyebrows, nav items, footer column heads ("PRODUCT", "ENTERPRISE", "RESOURCES"). Always `font-mono uppercase text-[12px] tracking-[-0.015rem]`.
- `text-pretty` is applied to virtually all text blocks (Tailwind utility for `text-wrap: pretty` — improves line breaks).
- `tabular-nums` likely in use on price figures (recommend it for any numeric data tables — status counts, quotas, $ amounts).

---

## 4. Spacing & Layout

### Container / max-width
- Page container: `max-w-[1920px]` (very wide — full-bleed feel on large monitors).
- Horizontal padding: `px-4` (16px) mobile → `lg:px-9` (36px) desktop. Implemented as `w-[calc(100%-32px)]` / `lg:w-[calc(100%-72px)]` with `mx-auto`.
- Inner content rails use `max-w-[980px]`, `max-w-[760px]`, `max-w-[600px]` for prose.
- Tailwind v4 `--container-*` tokens present: `xs=20rem`, `sm=24rem`, `xl=36rem`, `2xl=42rem`, `3xl=48rem`, `5xl=64rem`, `6xl=72rem`. (Not all used — Factory prefers explicit px max-widths.)
- `--spacing: .25rem` (4px) — Tailwind v4 default spacing multiplier.

### Section vertical rhythm
- Page sections: `py-16` to `py-24` (64–96px) between major blocks.
- Inner card/row spacing: `py-6` → `lg:py-8` (24–32px), `gap-x-4 gap-y-3` between columns.
- Footer: `my-20` (80px), internal `pt-28 pb-8`.
- Nav height ~ `h-[31px]` buttons inside a taller header bar.

### Grid
- Footer uses a 12-column CSS grid (`grid-cols-4 lg:grid-cols-12`, `gap-x-4 lg:gap-x-6`).
- Pricing rows use `flex flex-wrap items-start justify-between`.
- Homepage uses auto-rows-min grids with `grid-cols-4 lg:grid-cols-12`.

### Density
**Compact to medium** — not airy. Buttons are short (31px), text sizes lean small (14px body), gaps are tight (8–16px). Information density is higher than typical SaaS marketing sites. This transfers well to dashboards.

---

## 5. Components

### Buttons (the canonical Factory button — 2 variants + ghost)

**Primary (dark button, light text):**
```
class="group relative inline-flex cursor-pointer items-center justify-center
       border border-transparent transition-colors duration-150
       bg-base-1000 light:bg-dark-base-primary            /* dark fill */
       [&_*]:text-light-base-primary                       /* children light */
       hover:bg-light-base-secondary focus-visible:bg-light-base-secondary
       hover:[&_*]:text-dark-base-primary focus-visible:[&_*]:text-dark-base-primary
       hover:outline-light-base-secondary focus-visible:outline-light-base-secondary
       light:border-base-700
       overflow-clip rounded-sm focus-visible:outline focus-visible:outline-offset-4
       h-[31px] px-[14px] w-full"
```
**Secondary (light button, dark text):**
```
class="… bg-light-base-secondary hover:bg-dark-base-primary focus-visible:bg-dark-base-primary
        [&_*]:text-dark-base-primary hover:[&_*]:text-light-base-secondary
        hover:border-base-600 focus-visible:border-base-600
        … rounded-sm h-[31px] px-[14px] w-full"
```
**Key traits**:
- Height **31px**, horizontal padding **14px**.
- Radius **`rounded-sm` (3px)** — intentionally sharp.
- **Color inversion on hover/active** (dark fill → light fill, text flips). No elevation/lift/shadow.
- Transition: `transition-colors duration-150` (default 150ms ease-in-out).
- `focus-visible:outline focus-visible:outline-offset-4` for keyboard a11y (a 4px-offset outline ring — very visible, on-brand).
- `overflow-clip` so hover color changes don't bleed past the radius.
- `[&_*]:transition-colors` cascades the color swap to child icon/text nodes.

**Ghost / text link** (nav + footer + "Docs"):
```
class="text-pretty font-sans text-[14px] leading-[100%]
       relative flex w-fit items-center transition-colors duration-200
       hover:text-orange-500 group
       after:absolute after:-bottom-px after:left-0 after:h-px after:w-0 after:bg-current
       after:transition-all after:duration-300 after:ease-in-out hover:after:w-full
       text-base-500"
```
- 14px Geist Sans, color `text-base-500` (muted warm gray), `hover:text-orange-500`.
- **Animated underline** grows from 0→100% width on hover, 300ms `ease-in-out`, 1px tall (`after:h-px`), positioned `-bottom-px`.
- Mono variant exists (`font-mono text-[12px] uppercase tracking-[-0.015rem]`) used in the top nav row.

### Cards / panels
- Pricing plan rows are separated by **`border-b border-dashed border-base-300`** hairlines (dashed blueprint dividers — a signature motif).
- Footer is a single `rounded-xl lg:rounded-3xl` panel with `bg-dark-base-secondary` on `max-w-[1920px]`, `min-h-[430px]`+, generous internal padding (`pt-28 pb-8`).
- Cards mostly rely on **hairline borders** + background tone separation, NOT shadows. Borders use the stepped warm neutral scale (`base-600`, `base-700`, `base-800`).

### Form inputs
- Not prominent on the marketing pages, but the CSS shows `focus-within:border-accent-100` (orange border on focus) and `focus-within:border-[#a5a8ad]` (lighter warm-gray alt). Inputs would follow the same small-radius (3px), hairline-border, dark-fill convention.

### Top navigation / header
- Two rows of links: a **mono uppercase 12px row** ("Product", "Enterprise", "Pricing", "News", "Company", "Careers", "Docs") and a **sans 14px row** (footer-style ghost links).
- All links share the same animated-underline pseudo-element (`after:w-0 → hover:after:w-full`).
- "Log In" / "Contact Sales" CTAs are mono uppercase 12px text — text links, not buttons.
- Header uses Radix UI primitives (`data-radix-collection-item`).

### Badges / tags
- Mono 14px, uppercase, with an icon slot: `[&>[data-slot='badge-icon']]:border-transparent [&>[data-slot='badge-icon']]:bg-accent-300 [&_*]:text-base-500` and `[&_*]:!text-base-800`.
- So badges are small mono uppercase pills with an accent-orange-filled icon and warm-gray text — terminal-chip aesthetic.

### Code blocks
- **Shiki** is the syntax highlighter (`shiki-light`, `shiki-dark`, `shiki-dark-bg` tokens present). Both light and dark themes are bundled.
- Expect code blocks rendered as sharp-cornered panels with a Shiki dark theme (background ~ very dark, syntax colors keyed to the warm palette).

### Distinctive motifs
1. **Line-grid texture background** — `bg-[url("/assets/bg-lines.png")]` is the single most recognizable Factory element. A faint grid (graph-paper / CAD blueprint) sits behind the entire page. In light mode it's washed with `linear-gradient(rgba(255,255,255,.6), …)`. This is the brand's visual signature.
2. **Dashed hairline dividers** (`border-dashed border-base-300`) in pricing/feature tables — reinforces the blueprint metaphor.
3. **Terminal-cursor blink** — `@keyframes desktopAppCursorBlink` plus `desktopAppFadeUp` / `desktopAppPopupIn` / `desktopAppBackdropIn` suggest animated demo panels styled like a desktop IDE/terminal window.
4. **Mono-uppercase eyebrow labels** above every section.
5. **Color-inversion hover** (no lift, no scale, no shadow) — interaction feedback is purely tonal.

---

## 6. Radius / Shadow / Border

### Border-radius scale (Tailwind v4 `--radius-*`, all SMALL)
| Token | Value | Note |
|---|---|---|
| `--radius-xs` | `2px` (.125rem) | |
| `--radius-sm` | **`3px`** (.1875rem) | **Buttons, inputs, small chips**. Dominant radius. |
| `--radius-md` | `4px` (.25rem) | |
| `--radius-lg` | `6px` (.375rem) | |
| `--radius-xl` | `8px` (.5rem) | |
| `--radius-2xl` | `10px` (.625rem) | |
| `--radius-3xl` | `12px` (.75rem) | |
| `rounded-full` | pill | Used on avatars / status dots / some badges. |
| `rounded-[20px]`, `rounded-2xl`, `rounded-3xl`, `rounded-[32px]`, `rounded-[40px]` | larger radii | Used sparingly on big panels (footer `rounded-xl lg:rounded-3xl`). |

> **The system is intentionally low-radius.** Buttons at 3px feel engineered, not playful. Do NOT port this to a 12–16px radius dashboard — keep buttons/inputs at 3–6px, cards at 8–12px max.

### Shadow
- **Essentially no shadows.** The only `box-shadow` rule is the Tailwind v4 composite variable (`--tw-ring-shadow`, etc.), and it's used for focus rings (outline-based), not elevation.
- Hover/elevation is achieved via **color inversion + hairline border changes**, never drop-shadows.
- No `shadow-lg` / `shadow-2xl` marketing-card glow.

### Border treatment
- Hairlines only: 1px borders (`border`), color from the stepped warm-neutral scale (`base-600/700/800` dark side; `base-300/400` light side).
- **Dashed** dividers in tabular content (`border-dashed border-base-300`).
- Accent border on focus: `focus-within:border-accent-100` (orange).
- `border-transparent` used as the resting state on buttons so the box size doesn't shift when a real border appears.
- Top borders (`border-t`) used to delimit sections; `last:border-b-0` to clean up table rows.

---

## 7. Motion

### Default transition
- `--default-transition-duration: .15s`
- `--default-transition-timing-function: cubic-bezier(.4,0,.2,1)` (Tailwind's standard `ease-in-out`).
- `--ease-in-out: cubic-bezier(.4,0,.2,1)`, `--ease-out: cubic-bezier(0,0,.2,1)`.

### Specific transitions observed
- Buttons / color changes: `transition-colors duration-150` (150ms, ease-in-out).
- Nav/footer ghost links: `transition-colors duration-200`.
- Animated underline: `after:transition-all after:duration-300 after:ease-in-out` (**300ms** ease-in-out) — slightly slower so the underline draw feels deliberate.
- Pricing hover: `hover:border-dashed`, `hover:border-t-light-base-primary` (a dashed top border appears on hover for the active plan row).

### Hover behavior (very restrained)
- **No transforms, no scale, no translate-Y lift, no shadow change.** Hover = color swap only (dark↔light inversion, or `hover:text-orange-500`).
- `will-change-transform` declared on buttons (perf hint), but no actual transform animates.

### Keyframes (animations beyond hover)
Found in CSS: `accordion-down`, `accordion-up`, `carouselSlide`, `delayedFadeIn`, `desktopAppBackdropIn`, `desktopAppCursorBlink`, `desktopAppFadeUp`, `desktopAppPopupIn`, `enter`, `enterFromLeft`, `enterFromRight`, `exit`, `exitToLeft`, `exitToRight`, `fadeIn`.
- Implication: page sections **fade-in / fade-up on scroll** (`delayedFadeIn`, `desktopAppFadeUp`); a carousel auto-slides (`carouselSlide`); the product demo simulates a desktop IDE with a blinking cursor (`desktopAppCursorBlink`) and pop-in panels (`desktopAppPopupIn`).
- Standard utilities: `--animate-pulse` (2s ease-in-out infinite), `--animate-spin` (1s linear).

> Transferable rule: motion is **subtle and tonal** — fades (150–300ms) and a single "cursor blink" character moment. Avoid bouncy springs, large parallax, or scale-on-hover.

---

## 8. Transferability Notes (for an API-gateway SaaS admin/user console)

### Transfers well — adopt directly
1. **Warm-tinted near-black dark theme** (`#020202` bg, `#0a0908` panels, warm-gray text scale). Far more distinctive than the default `zinc-950` Tailwind dark; reads as "developer tool" not "generic dashboard". Make it the **default** theme.
2. **Single restrained orange accent** (`#EF6F2E` / `#ee6018` / `#d15010`) for primary actions, active states, links, focus rings, and key metrics. Easy to slot into a Tailwind `accent-*` token family.
3. **Geist + Geist Mono** via the `geist` npm package — clean, modern, designed for dev-tool UIs. Mono is perfect for API keys, request IDs, JSON, code snippets, log lines (your API gateway has tons of these).
4. **Color-inversion hover buttons** (31px tall, 3px radius, 150ms color swap, no lift). Compact and dense — ideal for tables and toolbars in a console. Add a slightly larger 36–40px variant for primary page CTAs.
5. **Mono uppercase eyebrow labels** at 12px (`PRODUCT`, `ACCOUNT`, `USAGE`) — great for grouping sidebar sections, table column groups, and stat-card captions.
6. **Hairline warm-gray borders + dashed dividers** for tables — your account/key/quota tables will look like an engineering instrument, not a spreadsheet.
7. **`focus-visible:outline` with 4px offset** — accessible and on-brand.
8. **Shiki for code blocks** (your docs / request examples / log viewers).
9. **Subtle fade-up on section mount** (150–300ms) — fine for dashboard panels, but keep it optional and respect `prefers-reduced-motion`.

### Adapt — don't copy verbatim
1. **Max-width 1920px full-bleed** is a marketing choice. For a console, use a constrained app shell (e.g., fixed left sidebar + `max-w-[1400px]` content rail) so dense tables don't sprawl.
2. **Large hero H1 at 72px with -0.18rem tracking** is for landing pages. In the console, cap headings at `text-3xl`/`text-2xl` with the same tight negative tracking for personality without wasting space.
3. **Line-grid blueprint background** is the brand signature — use it *subtly* behind the login page, empty states, or the marketing/landing sections of your console, NOT behind data-dense tables where it would hurt readability.

### Marketing-specific — do NOT force-fit
1. **Carousel slides, `desktopApp*` IDE demo animations, cursor blink** — these are Factory's product-demo showpieces. They have no place in an admin console and will feel gimmicky if imitated.
2. **Animated underline draw on every nav link** — fine for top-of-app nav, but noisy if applied to every sidebar item or table cell link. Reserve for primary nav and footer.
3. **Atmospheric gradient PNGs** behind hero sections — pure marketing eye-candy.
4. **`max-w-[1920px]` ultra-wide hero layouts** — console users want their data constrained and scannable, not cinematic.
5. **`font-normal` 72px hero headlines** — keep the *tight tracking* but you don't need display-scale type in a dashboard.

### Concrete Tailwind token mapping (starter for the reskin)
```css
/* tailwind theme extension — dark-first, warm neutral, vermillion accent */
--color-bg:           #020202;   /* page */
--color-surface:      #0a0908;   /* cards/panels */
--color-surface-2:    #14110f;   /* elevated (extrapolated) */
--color-border:       #342f2d;   /* hairline */
--color-border-strong:#3d3a39;   /* emphasized hairline */
--color-fg:           #fafafa;   /* primary text / light-base */
--color-fg-muted:     #a49d9a;   /* base-400, secondary text */
--color-fg-subtle:    #62666d;   /* base-600-ish, tertiary */
--color-accent-100:   #ef6f2e;   /* primary accent */
--color-accent-200:   #ee6018;   /* deeper accent */
--color-accent-300:   #d15010;   /* darkest accent / badges */
--color-accent-hover: oklch(70.5% .213 47.604); /* orange-500 link hover */

--font-sans:  'Geist', 'Geist Fallback', ui-sans-serif, system-ui, sans-serif;
--font-mono:  'Geist Mono', 'Geist Mono Fallback', ui-monospace, SFMono-Regular, Menlo, monospace;

--radius-xs:2px; --radius-sm:3px; --radius-md:4px; --radius-lg:6px; --radius-xl:8px; --radius-2xl:10px; --radius-3xl:12px;
--default-transition-duration:.15s;
--default-transition-timing-function:cubic-bezier(.4,0,.2,1);
```

---

## Caveats / Not Found

- **No rendered screenshots** — chrome-devtools MCP was unavailable in this environment. All values come from the live production HTML/CSS; they are exact for what Factory ships but cannot capture the *feel* of scroll animations, video, or hover micro-interactions. Recommend a manual visual QA pass on factory.ai before finalizing tokens.
- **Light theme** is supported but secondary; only the most-used dark-theme hexes were enumerated above. The `.light:` namespace flips `--color-background` → `#fafafa` and reuses the same accent palette. A full light-theme spec would need a second pass.
- **Geist weights/axes** — the CSS only declares `@font-face` for `Geist`, `Geist Mono`, and their `Fallback`s; the exact variable-font axis ranges (wght 100–900?) were not enumerated. The `geist` npm package documents the full range.
- **Exact `bg-lines.png`** is shipped as a 1×1 px asset that expands via `background-size`/repeat — the visible grid is likely defined by background-size on the element using `bg-[url(...)]`; the precise tile/repeat sizing wasn't captured per-element. Visually confirm the grid spacing you want before baking a custom asset.
- **Pricing CTA button color per tier** (Pro vs Enterprise highlight) was not separately enumerated; the same 2-variant button system applies.
- Did **not** capture `/product`, `/enterprise`, or `docs.factory.ai` pages — `/` and `/pricing` gave sufficient component coverage, but a `/product` page may reveal more code-block / terminal-panel styling worth a follow-up look.
