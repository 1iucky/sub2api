# Research: factory.ai Homepage — Structure & Motion Blueprint

- **Query**: Concrete, implementation-ready blueprint of the factory.ai **homepage** layout + animation feel (fills the motion gap left by `factory-ai-design-language.md`).
- **Scope**: external (live site fetch + CSS/JS analysis)
- **Date**: 2026-06-18
- **Method**: `curl` of live `https://factory.ai/` (HTTP 200, 479 KB SSR HTML) + both CSS bundles (193 KB + 3.7 KB) + representative JS chunks. `mcp__chrome-devtools__*` is **unavailable** in this env (confirmed) — analysis is from source, not screenshots. All keyframe/animation values below are extracted verbatim from the production CSS, not eyeballed.
- **Stack observed**: Next.js App Router, Tailwind v4 (`@layer`, `color-mix`, oklab), **GSAP** (minified as `i.os` — `timeline`, `autoAlpha`, `power3.out`, `stagger`) for the in-view bento/radar reveal, **plain CSS `@keyframes`** for the hero demo panel and all on-load/loop motion, **IntersectionObserver** (React hook) for scroll triggers. No Framer Motion anywhere on the homepage.

---

## 1. Section-by-section structure (top → bottom)

Container is uniform: every section is a 12-col grid, `max-w-[1920px]`, `mx-auto`, `px-4 lg:px-9`, `gap-x-4 lg:gap-x-6`, vertical `my-20` / `lg:mt-20 lg:mb-30`. Ultra-wide (1920) is the single most load-bearing layout decision.

| # | Section | Layout | Dominant content | Key class patterns |
|---|---|---|---|---|
| 0 | **Nav / header** | `fixed inset-x-0 top-0 z-60`, `flex justify-between`; left = wordmark, right = desktop `<nav aria-label="Main">` + CTA buttons. Mobile = off-canvas `-translate-x-full` drawer (`inert`). | Wordmark, nav links, 2 CTAs | `group/header fixed … z-60 px-4 lg:px-9`; nav uses `data-orientation="horizontal"` (Radix-style) |
| 1 | **Hero** (one `section`) | Same 12-col grid split `lg:col-span-5` (text) + `lg:col-span-7` (demo). Text `lg:self-center`, demo `lg:-mr-16` (bleeds off the right edge by 4rem). | H1 + eyebrow + subhead + 2 CTAs **+ animated SVG dashboard** | `sf-dashboard-frame-in` wrapper `rounded-[20px] shadow-[0_45px_120px_rgba(0,0,0,0.55)]` |
| 2 | **Social proof** | `col-span-full` centered eyebrow, then full-bleed **CSS marquee** (`overflow-x-clip`) of ~20 customer logo SVGs. | Single line of mono/SVG logos scrolling L→R | `.marquee-track` inline `animation-duration:125s; animation-play-state:running` |
| 3 | **"Defining your Software Factory"** bento | `col-span-full` heading, then `lg:grid-cols-3` row of 3 bordered cards `p-8`. | 3 cards: `01 Model Independence` (animated **radar chart** SVG), `02 Sovereign Deployment`, `03 Across the SDLC` | `border-base-800 bg-dark-base-secondary rounded-lg border p-8`; card number `font-mono text-[13px] text-accent-100` |
| 4 | **Surfaces demo** (wrapped `overflow-x-clip`) | A separate full-bleed section showing the desktop app + mobile app "cinematically" flying into place. | Large product screenshots | `data-surfaces-cinematic="desktop/mobile"` + `data-surfaces-ready` toggled by JS |
| 5 | **CTA band** ("Build with us") | Single centered block, H2 + sub-CTA link with arrow. | "Ready to build the software of the future?" + `Start building →` | mono uppercase eyebrow above H2 |
| 6 | **Footer** | `max-w-[1920px]` 12-col grid, but rendered as **one big rounded panel** `bg-dark-base-secondary`, `my-20 px-4 lg:px-9`. Columns: brand/logo, Resources, Company, Legal. Bottom row: social icons + `© 2026`. | Link columns + social row | `rounded-lg bg-dark-base-secondary` — the footer itself is a card, not a bare `<footer>` |

Order summary: **Nav → Hero(+demo) → Logo marquee → 3-card bento → Surfaces cinematic demo → CTA → Footer panel.** Notably **there is NO pricing teaser and NO testimonial grid** on the homepage — it's hero-demo-heavy.

---

## 2. The hero — exact composition

Text column (`lg:col-span-5`):
- **H1** — `Build Your Software Factory`. `font-normal` (not bold), `text-[40px]` mobile / `text-[72px]` desktop, `leading-[100%]` (tightest possible), `tracking-[-0.16rem]` mobile / `lg:tracking-[-0.18rem]` (very tight negative tracking), `lg:max-w-[12ch]`. Starts `invisible` then JS-revealed.
- **Subhead** — `A self-improving system for your SDLC. Ingest continuous signals and deploy production software.` `font-mono text-[14px] lg:text-[16px] leading-[120%] tracking-[-0.0175rem] text-base-500 max-w-[46ch]`. Mono body text is a signature choice.
- **CTAs** — `Download` (primary, bordered) + `Contact Sales` (ghost, with arrow icon). Both `<a class="group … border transition-colors duration-150">`; the visible label is `font-mono text-[14px] tracking-[-0.0175rem] uppercase`. **Hover behavior is unusual** (see §4).

Demo column (`lg:col-span-7`):
- Wrapper `div.sf-dashboard-frame-in` — `relative w-full max-w-[980px] rounded-[20px] shadow-[0_45px_120px_rgba(0,0,0,0.55)] lg:-mr-16`. The `-mr-16` makes the panel break out of the right gutter.
- **The whole "demo" is a single inline `<svg viewBox="0 0 1365 974">`** — not a video, not HTML. It is a hand-built mock of a "Software Factory" dashboard:
  - **Window chrome**: top-left sidebar `226×974`, macOS traffic-light dots (`#ff5f57 #ffbd2e #28c840`), `Back to app` chevron, `Software Factory` and `Automations` nav rows.
  - **6 pipeline stages** (Triage / Code Gen / Validate / Release / Document / Monitor), each with a `YOUR SOFTWARE FACTORY` mono header, KPI tiles (`THROUGHPUT 36/day`, `CYCLE TIME`, `PASS RATE 98.7%`), and a counter `DEPLOYS / 7DAY 0 12`.
  - **4 bottom sparkline cards** (`PR VALIDATIONS +71%`, `PRS MERGED SHIP RATE -43%`, `INCIDENTS PROCESSED -18%`, …) with MAX/MID/MIN grid lines.
  - **Defs**: `<pattern id="sfDashboardGrid" width="16" height="16">` for the faint 0.045-opacity white grid; `<radialGradient id="sfRadarGlow">` orange glow; `<filter id="sfDashboardShadow">` `feDropShadow dy=48 stdDeviation=34`.
  - **Animated sub-elements** each carry `class="sf-dashboard-segment" style="--sf-delay:180ms"` (sidebar), `--sf-delay:260ms` (header), etc. — staggered via CSS custom property. There is also `sf-radar-sweep` and `sf-dashboard-live` (pulsing KPI dots).

So the hero demo is: **one big SVG, CSS-keyframe-staggered segment reveal on load, plus perpetual radar sweep + dot pulse.** No JS required for the demo itself.

---

## 3. Motion catalogue (the priority)

All `@keyframes` are extracted verbatim from `/_next/static/css/2f8a58de57cfbf1c.css`.

### 3a. On-load (fire once, page entry)

| Animation | Keyframe definition | Applied as (duration / easing / fill) | What it does |
|---|---|---|---|
| **`home-hero-intro`** | `0%{opacity:0;transform:translateY(14px)} to{opacity:1;transform:translate(0)}` | hero text column elements, staggered | Fade-up the H1 + subhead + CTA group on mount |
| **`home-hero-terminal-in`** | `0%{opacity:0;transform:translateY(80px)} to{opacity:1;transform:translate(0)}` | demo panel base | Big 80px rise of the whole dashboard SVG |
| **`home-hero-terminal-fade-in`** | `0%{opacity:0;transform:translateY(10px)scale(.985)} 60%{opacity:.7} to{opacity:1;transform:translate(0)scale(1)}` | demo overlay layer | Subtle scale+fade layer on top of the rise |
| **`sfDashboardFrameIn`** (scoped `<style>` in hero markup) | `0%{opacity:0;transform:scale(0.18)} 55%{opacity:1} 100%{opacity:1;transform:scale(1)}` | `.sf-dashboard-frame-in` → `animation: sfDashboardFrameIn 900ms cubic-bezier(0.16,1,0.3,1) 120ms forwards` | Panel **zooms in from 18% scale** with overshoot-ish ease (the signature hero "boom-in") |
| **`sfDashboardSegmentIn`** (scoped) | `from{opacity:0} to{opacity:1}` | `.sf-dashboard-segment` → `animation: sfDashboardSegmentIn 520ms cubic-bezier(0.22,1,0.36,1) forwards; animation-delay: var(--sf-delay,0ms)` | Each dashboard segment (sidebar `180ms`, header `260ms`, …) fades in in sequence — a **stagger driven by inline `--sf-delay`**, not a JS loop |
| **`surface-chrome-stagger-reveal`** | `0%{opacity:0;transform:translateY(-6px)} to{opacity:1;transform:translate(0)}` | `[data-surface-chrome-reveal=active] [data-chrome-stagger]` → `.52s cubic-bezier(.22,1,.36,1) both` | Window-chrome elements of the surfaces demo reveal in sequence once `active` attr is set |

### 3b. On-scroll-into-view (IntersectionObserver, JS-toggled)

Mechanism: a React `useInView`-style hook creates `new IntersectionObserver(cb, {rootMargin})`, sets a `data-ready=true` / `data-surfaces-ready=true` / `data-surface-chrome-reveal="active"` flag on first intersect, then `disconnect()`. **Reveal happens via attribute + CSS animation, not via JS tweens** (except the radar — see 3d).

| Animation | Keyframe | Duration/easing | Notes |
|---|---|---|---|
| **`surfacesDesktopCinematicIn`** | `0%{opacity:0;transform:translate3d(var(--surfaces-cinematic-shift,0),22vh,0)scale(1.32)} 18%{opacity:1} 45%{opacity:1;transform:…translate3d(…,0,0)scale(1.32)} 72%{…scale(1)} to{opacity:1;transform:translate(0)scale(1)}` | `2.2s cubic-bezier(.4,0,.2,1) both`; `--surfaces-cinematic-shift:166px` on desktop (LG+) | The desktop screenshot **rises 22vh from below at 1.32× scale, holds oversized, then scales down to 1×** — a deliberate cinematic two-beat (rise → settle). |
| **`surfacesMobileCinematicIn`** | `0%,68%{opacity:0;transform:translate(-220px)scale(.94)} 78%{opacity:1;transform:translate(-180px)scale(.96)} to{opacity:1;transform:translate(0)scale(1)}` | `2.2s cubic-bezier(.4,0,.2,1) both` | Mobile phone screenshot slides in from the left (`-220px`) **after** the desktop one (68% blank = synced delay), then eases to position. |
| Bento card reveal | (CSS `invisible`→visible via class swap by the same IO hook) | staggered fade-up | Each bento card has `invisible` and gets shown when scrolled into view |

### 3c. Continuous / loop (infinite)

| Animation | Keyframe | Application | Purpose |
|---|---|---|---|
| **`sfDashboardPulse`** (scoped) | `0%,100%{opacity:.35} 50%{opacity:1}` | `.sf-dashboard-live` → `1200ms ease-in-out infinite` | KPI "live" dots breathe |
| **`sfRadarSweep`** (scoped) | `to{transform:rotate(360deg)}` | `.sf-radar-sweep` (transform-box:fill-box, origin:center) → `3500ms linear infinite` | Radar arm sweeps the hero dashboard's monitor tile |
| **`slidePattern`** | `0%{background-position:0 0} to{background-position:28.28px -28.28px}` | button hover hatch overlay → `animate-[slidePattern_2000ms_linear_infinite]`, **`paused` by default, `group-hover:running`** | 45° diagonal hatch slides across the button only while hovered |
| **`carouselSlide`** | `to{translate:var(--destination-x) var(--destination-y)}` | logo marquee (`.marquee-track`) inline `animation-duration:125s` | Endless horizontal logo scroll |
| **`desktopAppCursorBlink`** | `0%,45%{opacity:1} 55%,to{opacity:.15}` | `.animate-[desktopAppCursorBlink_900ms_ease-in-out_infinite]` | Terminal cursor blink |
| **`spin`** / **`pulse`** | Tailwind defaults | misc icons | standard |

### 3d. GSAP-driven radar/bento stagger (the one JS-tweened piece)

The `01 Model Independence` radar chart (`chunk_6785.js` + shared GSAP) uses a **GSAP timeline** triggered by `IntersectionObserver({threshold:0.3, rootMargin:"0px 0px -80px 0px"})`:

```
let f = gsap.timeline({paused:true, defaults:{ease:"power3.out"}});
f.to(rings,   {scale:1, autoAlpha:1, duration:0.55, stagger:0.06})
 .to(spokes,  {scale:1, autoAlpha:1, duration:0.40, stagger:0.05}, "-=0.2")
 .to(series,  {scale:1, autoAlpha:1, duration:0.50, stagger:0.12}, "+=0.05")
 .addLabel("lines","<")
 .to(labels,  {scale:1, autoAlpha:1, ...});
```

All radar elements start at `gsap.set(...,{scale:0, autoAlpha:0, svgOrigin:"170 150"})` (explode out from the radar center). On `isIntersecting`, `f.play()`.

### Easing summary
- **Primary ease**: `cubic-bezier(0.22, 1, 0.36, 1)` — aggressive ease-out (Tailwind's `ease-out` cousin), used on every segment/popup/fade-up. This is *the* factory motion signature.
- **Hero zoom-in**: `cubic-bezier(0.16, 1, 0.3, 1)` — even more dramatic ease-out.
- **Cinematic surfaces**: `cubic-bezier(0.4, 0, 0.2, 1)` — Material "standard" ease (slower in-out) for the long 2.2s fly-in.
- **GSAP radar**: `power3.out`.

### Stagger convention
- Hero dashboard segments: **`var(--sf-delay)`** = 180ms, 260ms, … (hand-tuned per segment).
- Bento/Radar GSAP: `stagger: 0.05–0.12` between sibling tweens, with relative offsets (`-=0.2`, `+=0.05`).

---

## 4. Hover micro-interactions

- **Buttons** (primary + ghost, both `class="group relative inline-flex … border transition-colors duration-150"`):
  - Label stays; an **absolutely-positioned hatch overlay** (`background-image:repeating-linear-gradient(45deg, …)`) sits `opacity-0` and **animates in on hover**: `group-hover:animate-[delayedFadeIn_100ms_ease-out_forwards]` for the fade, plus the hatch **slides**: `paused … animate-[slidePattern_2000ms_linear_infinite]` flipped to `group-hover:running`. So hovering a button = a diagonal hatch **fades in and starts scrolling**. `delayedFadeIn` keyframe: `0%,80%{opacity:0} to{opacity:1}` (the 80% blank gives a deliberate ~80ms hesitation before the hatch appears).
  - Focus-visible gets the same: `group-focus-visible:animate-[delayedFadeIn_…​]`.
- **Logo strip items**: `text-base-700 group-hover:text-light-base-primary transition-colors duration-200` — each logo brightens individually on hover (per-logo `group`).
- **Marquee**: `animation-play-state:running` (stays running; not paused on hover of the track — only individual logos change color).
- **Cards**: bordered `border-base-800`, hover not heavily animated on the homepage bento (rely on the in-view reveal for the motion).
- **Nav mega-menu**: `megamenuSlideDown / megamenuSlideUp / navContentFadeIn / navContentFadeOut` keyframes (not on homepage scroll, only on hover-open of nav dropdowns).

---

## 5. Distinguishing layout patterns (what makes it NOT generic SaaS)

These 7 are the port-over targets:

1. **Ultra-wide 1920px container** with a 12-col grid and `lg:px-9` — every section shares it; nothing is the SaaS-default `max-w-7xl mx-auto`.
2. **Oversized tight-tracking headline**: `text-[72px] tracking-[-0.18rem] leading-[100%] font-normal` (not bold). Negative tracking + 100% line-height + normal weight = the "engineering poster" look.
3. **Mono, uppercase, letter-spaced eyebrows** above every section (e.g. `YOUR SOFTWARE FACTORY`, `Trusted by leading engineering teams`, card numbers `01 02 03`) in `font-mono text-[13px]`.
4. **Blueprint line-grid background**: `bg-[url("/assets/bg-lines.png")]` (a static PNG asset, faint white lines) behind the whole page. The hero dashboard also draws its own 16px `<pattern>` grid.
5. **The hero "demo" is a bespoke inline SVG**, not a screenshot or video — hand-built window chrome + KPI tiles + radar + sparklines, all CSS-animated. This is the biggest non-generic move.
6. **Dashed hairline dividers + bordered cards** (`border-base-800`) and a **footer that is itself a big rounded panel** (`rounded-lg bg-dark-base-secondary`) rather than a bare strip.
7. **Demo panel bleeds off-grid** (`lg:-mr-16`) and uses a heavy drop-shadow (`shadow-[0_45px_120px_rgba(0,0,0,0.55)]`) to feel like a floating object, not a contained image.
8. **Two-beat cinematic scroll-reveal** (`surfacesDesktopCinematicIn`: rise-at-scale → hold → settle) instead of a plain fade-up — gives scroll a "movie" cadence.

---

## 6. Reduced-motion handling

Robust and **per-feature**, all in `@media (prefers-reduced-motion: reduce)`:

```css
/* 1. kill the page transition */
.page-transition-in { animation: none; }

/* 2. kill ALL desktopApp* (terminal/cursor/popup) animations */
[class*=animate-\[desktopApp] { animation: none !important; }

/* 3. force surfaces visible, no cinematic */
[data-surfaces-cinematic] { opacity:1 !important; filter:none !important;
                            animation:none !important; transform:none !important; }

/* 4. hero intro + terminal shells just shown */
[data-home-hero-intro], [data-terminal-shell] { opacity:1 !important;
                            animation:none !important; transform:none !important; }

/* 5. bento chrome stagger shown */
[data-surface-chrome-reveal=active] [data-chrome-stagger] { animation:none !important; }
```

Plus the hero's own scoped `<style>` wraps its keyframes in `@media (prefers-reduced-motion: no-preference)`, so under reduce the dashboard is shown statically at full opacity/scale. The GSAP radar path does an explicit early-return: `if (matchMedia("(prefers-reduced-motion: reduce)").matches) return gsap.set([...elements], {autoAlpha:1})`. **Pattern: CSS keyframes gated behind `no-preference`; JS paths short-circuit via `matchMedia`.**

---

## 7. Transferable blueprint for OUR homepage (AI API gateway)

Our product: pools/distributes upstream AI subscriptions (Claude/OpenAI/Gemini) with load balancing, billing, per-key quotas. Stack: **Vue 3 + Tailwind v3.4, no GSAP, no Framer** → use **CSS keyframes + a tiny `useInView` composable (IntersectionObserver)**. Keep all easing on `cubic-bezier(0.22,1,0.36,1)`.

### Proposed section order

| # | Section | Content (mapped to us) | Animation |
|---|---|---|---|
| 0 | **Nav** (fixed, `z-60`, `lg:px-9`) | Logo "SiliconBase" + links (Gateway, Pools, Pricing, Docs) + `Sign in` / `Get API key` | Static; mobile off-canvas drawer with `-translate-x-full` |
| 1 | **Hero** (12-col, `lg:col-span-5` text + `lg:col-span-7` demo, demo `lg:-mr-16`) | H1 "One gateway for every AI subscription" (`text-[72px] tracking-[-0.18rem] leading-[100%] font-normal`), mono subhead, 2 CTAs (`Get API key` bordered / `Read docs` ghost-arrow) | `home-hero-intro` fade-up stagger on text (H1 → subhead → CTAs, `--delay` 0/120/240ms); demo panel uses `sfDashboardFrameIn`-style **18%→100% scale zoom-in**, `900ms cubic-bezier(0.16,1,0.3,1) 120ms` |
| 1b | **Hero demo** (inline SVG, NOT a screenshot) | A mock "gateway dashboard" SVG: window chrome (macOS dots + `Gateway`/`Pools`/`Keys` nav), a **6-row request stream** (`/v1/messages` → pool → upstream), live KPI tiles (`RPM 1.2k`, `P95 latency 240ms`, `Uptime 99.9%`, `Keys 348`), and 3 sparklines (Claude/OpenAI/Gemini throughput). `<pattern>` 16px grid bg. | Segment stagger via `--sf-delay` (sidebar 180ms, header 260ms, KPI tiles +120ms each); perpetual `sfDashboardPulse` on "live" dots (`1200ms ease-in-out infinite`); optional `sfRadarSweep` on a "health" tile |
| 2 | **Social proof / logo marquee** (full-bleed, `overflow-hidden`) | "Powering AI traffic for" + customer/upstream logo SVGs (or "Compatible with Claude · OpenAI · Gemini · Bedrock · Antigravity") | CSS `marquee` `animation-duration:120s linear infinite`; per-logo `text-base-700 hover:text-foreground transition-colors duration-200` |
| 3 | **Bento "Defining your gateway"** (`lg:grid-cols-3`, bordered `p-8` cards) | `01 Multi-upstream pooling` (radar or bar chart of model coverage), `02 Sticky + load-balanced` (animated routing diagram), `03 Quotas & billing` (key/usage tiles). Card number `font-mono text-[13px] text-accent-100`. | On-scroll reveal: cards `invisible`→`visible` via `useInView`; staggered `desktopAppFadeUp 320ms cubic-bezier(0.22,1,0.36,1) both` with `+120ms` per card. Radar (if kept) can be pure-CSS `stroke-dashoffset` animate OR a lightweight IO-driven stagger — **skip GSAP, do CSS-only** |
| 4 | **Live request demo** (cinematic) — *optional, our analog of "Surfaces"* | A wide terminal/log panel showing a real `/v1/messages` call being routed (typed text + upstream selection). | Two-beat reveal mimicking `surfacesDesktopCinematicIn`: rise `translateY(22vh) scale(1.32)` → settle, `2.2s cubic-bezier(0.4,0,0.2,1)`, gated behind `prefers-reduced-motion: no-preference` |
| 5 | **CTA band** | "Ship AI features without managing quotas" + `Start free →` | Eyebrow + H2 fade-up on scroll |
| 6 | **Footer panel** (`rounded-lg bg-dark-base-secondary my-20`) | Brand, Resources (Docs, Status, Changelog), Company, Legal (Privacy, Terms), social row, `© 2026` | Static |

### Drop (marketing-only showpieces not worth the build cost for a gateway)
- The 2.2s cinematic "surfaces" multi-screenshot fly-in — high effort, low info for a B2B API product; replace with a simpler in-view fade-up.
- The GSAP radar explode-out — replace with a static radar or a CSS `stroke-dashoffset` draw.
- The full customer-logo marquee if we don't have real logos — replace with a "Compatible with" upstream-model row (still marquee, lower stakes).

### Tailwind v3.4 implementation notes (we are NOT on Tailwind v4)
- factory uses Tailwind **v4** features we must downgrade: `color-mix(in oklab, …)` → precompute the mixed color as a hex/rgb token; `@layer` oklab gradients → standard `bg-gradient-to-*`. Our existing SiliconBase tokens already cover colors.
- Arbitrary-value classes they lean on (`text-[72px]`, `tracking-[-0.18rem]`, `shadow-[0_45px_120px_rgba(0,0,0,0.55)]`, `lg:-mr-16`, `max-w-[1920px]`) **all work as-is in Tailwind v3.4**.
- The `group-hover:running` / `paused` animation-play-state trick works in v3 via `animation-play-state` utilities — verify our config exposes them (add `[animation-play-state:paused]` arbitrary props if needed).
- **Blueprint bg**: ship a `/assets/bg-lines.png` (or generate via `background-image: repeating-linear-gradient(0deg,…), repeating-linear-gradient(90deg,…)` at ~`rgba(255,255,255,0.04)`) — CSS version avoids a binary asset and matches factory's look closely.
- **The hero demo SVG**: build once as a Vue component `<GatewayHeroDashboard>` with `<defs><pattern>` grid + `<g class="segment" :style="{'--sf-delay':'180ms'}">` groups; keyframes live in a global `assets/css/animations.css` so they're reusable.

### Composable to add (one file)
`composables/useInView.ts` — `ref` + `new IntersectionObserver(cb, {rootMargin:'0px 0px -80px 0px', threshold:0.3})`, sets `inView=true` once and `disconnect()`. Mirror factory's `rootMargin`. Pair with a `v-reveal` directive that sets `data-ready="true"` (matches factory's attribute-toggle pattern) so CSS keyframes fire purely on attribute change.

### Reduced-motion (port verbatim)
Wrap every decorative keyframe set in `@media (prefers-reduced-motion: no-preference)`; add a global reduce block:
```css
@media (prefers-reduced-motion: reduce) {
  [data-reveal], [data-cinematic], [data-hero-demo] { opacity:1 !important; transform:none !important; animation:none !important; }
}
```
and in `useInView`, early-return `inView=true` when `matchMedia('(prefers-reduced-motion: reduce)').matches` is true.

---

## Caveats / Not Found

- `mcp__chrome-devtools__*` is unavailable in this environment (confirmed by absence of tool + tool-name search); no visual screenshots were taken. All values are from source and are therefore *exact* for CSS but the **visual feel** (exact scroll-sync timing of the cinematic surfaces, the precise moment the dashboard zoom "lands") can only be confirmed by loading the live page in a browser. Recommend the implementer open factory.ai once to eyeball the hero zoom-in cadence before locking timing.
- The two "cinematic surfaces" keyframes (`surfacesDesktopCinematicIn`/`Mobile`) reference product screenshots that load in a separate route section; the homepage SSR markup in this fetch contains the hero + bento + footer fully, but the surfaces section's image payloads are lazy-loaded (`rootMargin:300px` lazy hook confirmed in JS). Structural analysis of that section is from CSS + JS, not from rendered images.
- Exact customer-logo list in the marquee was not enumerated (there are ~20 SVG logos); only the mechanism + first two (a "P"-style wordmark, "Podium") were inspected. Not material to the blueprint.
- GSAP is bundled but minified as `i.os`; confirmed by API usage (`timeline`, `autoAlpha`, `power3.out`, `stagger`, `svgOrigin`) — these tokens are unambiguous GSAP signatures. We are deliberately **not** adopting GSAP for our stack (§7 uses CSS-only equivalents).
