<template>
  <!--
    Bespoke inline SVG "gateway dashboard" for the homepage hero.
    Hand-authored vector (no raster). Anatomy:
      - theme-colored window surface + radial radar glow
      - window chrome: macOS traffic-light dots + titlebar
      - sidebar: Gateway / Pools / Keys / Settings nav rows
      - KPI tiles: RPM 1.2k · P95 240ms · Uptime 99.9% · Keys 348
      - request stream: /v1/messages -> pool -> upstream (4-6 rows)
      - 3 sparklines: claude / openai / gemini throughput
    Stagger reveal is driven purely by CSS classes + --sf-delay custom property
    (see styles/theme-override.css `[data-hero-demo] .sf-dashboard-segment`).
  -->
  <svg
    viewBox="0 0 980 620"
    class="block w-full"
    role="img"
    :aria-label="t('home.hero2.demo.title')"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <defs>
      <!-- Vermillion radar glow. -->
      <radialGradient id="sfRadarGlow" cx="50%" cy="50%" r="50%">
        <stop offset="0%" stop-color="#ef6f2e" stop-opacity="0.55" />
        <stop offset="60%" stop-color="#ef6f2e" stop-opacity="0.12" />
        <stop offset="100%" stop-color="#ef6f2e" stop-opacity="0" />
      </radialGradient>
      <!-- Sparkline gradient fill (vermillion). -->
      <linearGradient id="sfSparkFill" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0%" stop-color="#ef6f2e" stop-opacity="0.35" />
        <stop offset="100%" stop-color="#ef6f2e" stop-opacity="0" />
      </linearGradient>
    </defs>

    <!-- ====== Base panel ====== -->
    <rect x="0.5" y="0.5" width="979" height="619" rx="15.5" class="fill-white dark:fill-dark-900" />

    <!-- ====== Window chrome ====== -->
    <g class="sf-dashboard-segment" :style="{ '--sf-delay': '120ms' }">
      <rect x="0.5" y="0.5" width="979" height="40" rx="15.5" class="fill-gray-50 dark:fill-dark-950" />
      <rect x="0.5" y="24" width="979" height="17" class="fill-gray-50 dark:fill-dark-950" />
      <!-- macOS traffic lights -->
      <circle cx="20" cy="20" r="5" fill="#ff5f57" />
      <circle cx="36" cy="20" r="5" fill="#ffbd2e" />
      <circle cx="52" cy="20" r="5" fill="#28c840" />
      <!-- Title -->
      <text
        x="490"
        y="24"
        text-anchor="middle"
        class="fill-gray-500 dark:fill-dark-400"
        font-family="Geist Mono Variable, ui-monospace, monospace"
        font-size="11"
        letter-spacing="0.06em"
      >{{ t('home.hero2.demo.title') }}</text>
    </g>

    <!-- ====== Sidebar ====== -->
    <g class="sf-dashboard-segment" :style="{ '--sf-delay': '180ms' }">
      <rect x="16" y="48" width="170" height="556" rx="8" class="fill-gray-50 dark:fill-dark-950" />
      <line x1="186" y1="56" x2="186" y2="596" class="stroke-gray-200 dark:stroke-dark-800" stroke-width="1" />

      <!-- Sidebar nav rows -->
      <g class="font-mono" font-family="Geist Mono Variable, ui-monospace, monospace" font-size="11" letter-spacing="0.08em">
        <!-- Gateway (active) -->
        <rect x="28" y="72" width="146" height="28" rx="4" class="fill-primary-500/10" />
        <rect x="28" y="72" width="3" height="28" rx="1.5" fill="#ef6f2e" />
        <text x="44" y="89" class="fill-primary-600 dark:fill-primary-400" font-weight="500">{{ t('home.hero2.demo.sidebarGateway') }}</text>

        <text x="44" y="125" class="fill-gray-500 dark:fill-dark-400">{{ t('home.hero2.demo.sidebarPools') }}</text>
        <text x="44" y="153" class="fill-gray-500 dark:fill-dark-400">{{ t('home.hero2.demo.sidebarKeys') }}</text>
        <text x="44" y="181" class="fill-gray-500 dark:fill-dark-400">{{ t('home.hero2.demo.sidebarSettings') }}</text>
      </g>

      <!-- Sidebar section label -->
      <text
        x="28"
        y="226"
        class="fill-gray-400 dark:fill-dark-400"
        font-family="Geist Mono Variable, ui-monospace, monospace"
        font-size="9"
        letter-spacing="0.14em"
      >POOLS</text>

      <!-- Pool rows -->
      <g class="font-mono" font-family="Geist Mono Variable, ui-monospace, monospace" font-size="10">
        <circle cx="34" cy="248" r="3" fill="#ef6f2e" />
        <text x="46" y="252" class="fill-gray-600 dark:fill-gray-300">claude-prod</text>
        <text x="168" y="252" text-anchor="end" class="fill-primary-500">12</text>

        <circle cx="34" cy="272" r="3" class="fill-gray-300 dark:fill-dark-600" />
        <text x="46" y="276" class="fill-gray-600 dark:fill-gray-300">openai-team</text>
        <text x="168" y="276" text-anchor="end" class="fill-gray-500 dark:fill-dark-400">8</text>

        <circle cx="34" cy="296" r="3" class="fill-gray-300 dark:fill-dark-600" />
        <text x="46" y="300" class="fill-gray-600 dark:fill-gray-300">gemini-pool</text>
        <text x="168" y="300" text-anchor="end" class="fill-gray-500 dark:fill-dark-400">6</text>

        <circle cx="34" cy="320" r="3" class="fill-gray-300 dark:fill-dark-600" />
        <text x="46" y="324" class="fill-gray-600 dark:fill-gray-300">bedrock</text>
        <text x="168" y="324" text-anchor="end" class="fill-gray-500 dark:fill-dark-400">3</text>
      </g>
    </g>

    <!-- ====== KPI tiles row ====== -->
    <g class="sf-dashboard-segment" :style="{ '--sf-delay': '260ms' }">
      <!-- RPM -->
      <g transform="translate(206,64)">
        <rect width="180" height="68" rx="6" class="fill-white stroke-gray-200 dark:fill-dark-900 dark:stroke-dark-800" stroke-width="1" />
        <text x="14" y="22" class="fill-gray-400 dark:fill-dark-400" font-family="Geist Mono Variable, monospace" font-size="9" letter-spacing="0.14em">{{ t('home.hero2.demo.kpiRpm') }}</text>
        <text x="14" y="50" class="fill-gray-900 dark:fill-white" font-family="Geist Variable, sans-serif" font-size="26" font-weight="500">{{ t('home.hero2.demo.kpiRpmValue') }}</text>
        <circle class="sf-dashboard-live" cx="166" cy="18" r="3" fill="#10b981" />
      </g>
      <!-- P95 -->
      <g transform="translate(396,64)">
        <rect width="180" height="68" rx="6" class="fill-white stroke-gray-200 dark:fill-dark-900 dark:stroke-dark-800" stroke-width="1" />
        <text x="14" y="22" class="fill-gray-400 dark:fill-dark-400" font-family="Geist Mono Variable, monospace" font-size="9" letter-spacing="0.14em">{{ t('home.hero2.demo.kpiP95') }}</text>
        <text x="14" y="50" class="fill-gray-900 dark:fill-white" font-family="Geist Variable, sans-serif" font-size="26" font-weight="500">{{ t('home.hero2.demo.kpiP95Value') }}</text>
        <circle class="sf-dashboard-live" cx="166" cy="18" r="3" fill="#ef6f2e" />
      </g>
      <!-- Uptime (radar tile) -->
      <g transform="translate(586,64)">
        <rect width="180" height="68" rx="6" class="fill-white stroke-gray-200 dark:fill-dark-900 dark:stroke-dark-800" stroke-width="1" />
        <text x="14" y="22" class="fill-gray-400 dark:fill-dark-400" font-family="Geist Mono Variable, monospace" font-size="9" letter-spacing="0.14em">{{ t('home.hero2.demo.kpiUptime') }}</text>
        <text x="14" y="50" class="fill-gray-900 dark:fill-white" font-family="Geist Variable, sans-serif" font-size="26" font-weight="500">{{ t('home.hero2.demo.kpiUptimeValue') }}</text>
        <!-- radar health indicator -->
        <g transform="translate(150,42)">
          <circle r="14" fill="url(#sfRadarGlow)" />
          <circle r="11" fill="none" class="stroke-dark-800 dark:stroke-dark-700" stroke-width="0.75" opacity="0.5" />
          <circle r="7" fill="none" class="stroke-dark-800 dark:stroke-dark-700" stroke-width="0.75" opacity="0.5" />
          <g class="sf-radar-sweep">
            <line x1="0" y1="0" x2="11" y2="0" stroke="#ef6f2e" stroke-width="1.2" stroke-linecap="round" />
          </g>
          <circle r="1.5" fill="#ef6f2e" />
        </g>
      </g>
      <!-- Keys -->
      <g transform="translate(776,64)">
        <rect width="180" height="68" rx="6" class="fill-white stroke-gray-200 dark:fill-dark-900 dark:stroke-dark-800" stroke-width="1" />
        <text x="14" y="22" class="fill-gray-400 dark:fill-dark-400" font-family="Geist Mono Variable, monospace" font-size="9" letter-spacing="0.14em">{{ t('home.hero2.demo.kpiKeys') }}</text>
        <text x="14" y="50" class="fill-gray-900 dark:fill-white" font-family="Geist Variable, sans-serif" font-size="26" font-weight="500">{{ t('home.hero2.demo.kpiKeysValue') }}</text>
        <circle class="sf-dashboard-live" cx="166" cy="18" r="3" fill="#ef6f2e" />
      </g>
    </g>

    <!-- ====== Request stream panel ====== -->
    <g class="sf-dashboard-segment" :style="{ '--sf-delay': '340ms' }">
      <rect x="206" y="148" width="750" height="240" rx="6" class="fill-gray-50 stroke-gray-200 dark:fill-dark-950 dark:stroke-dark-800" stroke-width="1" />

      <!-- Stream header -->
      <text x="222" y="170" class="fill-gray-500 dark:fill-dark-400" font-family="Geist Mono Variable, monospace" font-size="9" letter-spacing="0.14em">{{ t('home.hero2.demo.streamTitle') }}</text>
      <line x1="206" y1="180" x2="956" y2="180" class="stroke-gray-200 dark:stroke-dark-800" stroke-width="1" stroke-dasharray="3 3" />

      <!-- Column headers -->
      <g class="font-mono" font-family="Geist Mono Variable, monospace" font-size="9" letter-spacing="0.1em">
        <text x="222" y="198" class="fill-gray-400 dark:fill-dark-400">ENDPOINT</text>
        <text x="430" y="198" class="fill-gray-400 dark:fill-dark-400">{{ t('home.hero2.demo.streamRoute') }}</text>
        <text x="640" y="198" class="fill-gray-400 dark:fill-dark-400">{{ t('home.hero2.demo.streamUpstream') }}</text>
        <text x="820" y="198" class="fill-gray-400 dark:fill-dark-400">STATUS</text>
      </g>

      <!-- Stream rows -->
      <g class="font-mono" font-family="Geist Mono Variable, monospace" font-size="11">
        <!-- row 1 -->
        <text x="222" y="226" class="fill-gray-700 dark:fill-gray-200">/v1/messages</text>
        <text x="430" y="226" class="fill-primary-500">claude-prod</text>
        <text x="640" y="226" class="fill-gray-500 dark:fill-dark-300">anthropic</text>
        <rect x="820" y="216" width="48" height="14" rx="3" class="fill-emerald-500/10" />
        <text x="844" y="226" text-anchor="middle" class="fill-emerald-600 dark:fill-emerald-400" font-size="9">{{ t('home.hero2.demo.streamStatus') }}</text>
        <line x1="206" y1="234" x2="956" y2="234" class="stroke-gray-200 dark:stroke-dark-800" stroke-width="0.75" stroke-dasharray="2 3" />

        <!-- row 2 -->
        <text x="222" y="254" class="fill-gray-700 dark:fill-gray-200">/v1/chat/completions</text>
        <text x="430" y="254" class="fill-primary-500">openai-team</text>
        <text x="640" y="254" class="fill-gray-500 dark:fill-dark-300">openai</text>
        <rect x="820" y="244" width="48" height="14" rx="3" class="fill-emerald-500/10" />
        <text x="844" y="254" text-anchor="middle" class="fill-emerald-600 dark:fill-emerald-400" font-size="9">{{ t('home.hero2.demo.streamStatus') }}</text>
        <line x1="206" y1="262" x2="956" y2="262" class="stroke-gray-200 dark:stroke-dark-800" stroke-width="0.75" stroke-dasharray="2 3" />

        <!-- row 3 -->
        <text x="222" y="282" class="fill-gray-700 dark:fill-gray-200">/v1beta/models</text>
        <text x="430" y="282" class="fill-primary-500">gemini-pool</text>
        <text x="640" y="282" class="fill-gray-500 dark:fill-dark-300">google</text>
        <rect x="820" y="272" width="48" height="14" rx="3" class="fill-emerald-500/10" />
        <text x="844" y="282" text-anchor="middle" class="fill-emerald-600 dark:fill-emerald-400" font-size="9">{{ t('home.hero2.demo.streamStatus') }}</text>
        <line x1="206" y1="290" x2="956" y2="290" class="stroke-gray-200 dark:stroke-dark-800" stroke-width="0.75" stroke-dasharray="2 3" />

        <!-- row 4 -->
        <text x="222" y="310" class="fill-gray-700 dark:fill-gray-200">/v1/messages</text>
        <text x="430" y="310" class="fill-primary-500">claude-prod</text>
        <text x="640" y="310" class="fill-gray-500 dark:fill-dark-300">bedrock</text>
        <rect x="820" y="300" width="48" height="14" rx="3" class="fill-amber-500/10" />
        <text x="844" y="310" text-anchor="middle" class="fill-amber-600 dark:fill-amber-400" font-size="9">retry</text>
        <line x1="206" y1="318" x2="956" y2="318" class="stroke-gray-200 dark:stroke-dark-800" stroke-width="0.75" stroke-dasharray="2 3" />

        <!-- row 5 -->
        <text x="222" y="338" class="fill-gray-700 dark:fill-gray-200">/antigravity/v1/messages</text>
        <text x="430" y="338" class="fill-primary-500">bedrock</text>
        <text x="640" y="338" class="fill-gray-500 dark:fill-dark-300">antigravity</text>
        <rect x="820" y="328" width="48" height="14" rx="3" class="fill-emerald-500/10" />
        <text x="844" y="338" text-anchor="middle" class="fill-emerald-600 dark:fill-emerald-400" font-size="9">{{ t('home.hero2.demo.streamStatus') }}</text>
        <line x1="206" y1="346" x2="956" y2="346" class="stroke-gray-200 dark:stroke-dark-800" stroke-width="0.75" stroke-dasharray="2 3" />

        <!-- row 6 (live, with cursor) -->
        <text x="222" y="366" class="fill-gray-700 dark:fill-gray-200">$ curl -X POST</text>
        <rect class="sf-dashboard-cursor" x="338" y="358" width="7" height="12" fill="#ef6f2e" />
        <text x="430" y="366" class="fill-gray-400 dark:fill-dark-400">claude-prod</text>
        <text x="640" y="366" class="fill-gray-400 dark:fill-dark-400">anthropic</text>
      </g>
    </g>

    <!-- ====== Sparkline cards ====== -->
    <g class="sf-dashboard-segment" :style="{ '--sf-delay': '420ms' }">
      <g
        v-for="card in miniChartCards"
        :key="card.key"
        :transform="`translate(${card.x},404)`"
      >
        <rect width="240" height="180" rx="6" class="fill-white stroke-gray-200 dark:fill-dark-900 dark:stroke-dark-800" stroke-width="1" />
        <text x="14" y="22" class="fill-gray-400 dark:fill-dark-400" font-family="Geist Mono Variable, monospace" font-size="9" letter-spacing="0.14em">
          {{ t(card.labelKey) }}
        </text>
        <text x="14" y="50" class="fill-gray-900 dark:fill-white" font-family="Geist Variable, sans-serif" font-size="24" font-weight="500">
          {{ card.current }}<tspan class="fill-gray-400 dark:fill-dark-400" font-size="11" font-weight="400">/min</tspan>
        </text>
        <text x="226" y="48" text-anchor="end" class="fill-emerald-600 dark:fill-emerald-400" font-family="Geist Mono Variable, ui-monospace, monospace" font-size="10">
          {{ card.delta }}
        </text>

        <g class="stroke-gray-200 dark:stroke-dark-800" stroke-width="0.75" opacity="0.8">
          <line x1="18" y1="82" x2="222" y2="82" stroke-dasharray="2 5" />
          <line x1="18" y1="112" x2="222" y2="112" stroke-dasharray="2 5" />
          <line x1="18" y1="142" x2="222" y2="142" stroke-dasharray="2 5" />
        </g>
        <path :d="miniChartAreaPath(card.values)" :fill="card.fill" opacity="0.48" />
        <path
          :d="miniChartLinePath(card.values)"
          fill="none"
          :stroke="card.stroke"
          stroke-width="1.8"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
        <g
          v-for="point in miniChartPoints(card.values)"
          :key="`${card.key}-${point.label}`"
          class="sf-mini-chart-point"
          tabindex="0"
          :transform="`translate(${point.x},${point.y})`"
        >
          <title>{{ `${card.name} ${point.label}: ${point.value}/min` }}</title>
          <line y1="-64" y2="28" class="sf-mini-chart-crosshair" />
          <circle class="sf-mini-chart-hit" r="8" fill="transparent" />
          <circle r="3" :fill="card.stroke" />
          <g class="sf-mini-chart-tooltip" :transform="point.x > 160 ? 'translate(-70,-42)' : 'translate(10,-42)'">
            <rect width="62" height="30" rx="4" class="fill-gray-950/95 dark:fill-white/95" />
            <text x="8" y="12" class="fill-gray-300 dark:fill-gray-500" font-family="Geist Mono Variable, ui-monospace, monospace" font-size="8">{{ point.label }}</text>
            <text x="8" y="24" class="fill-white dark:fill-gray-950" font-family="Geist Mono Variable, ui-monospace, monospace" font-size="10">{{ point.value }}/min</text>
          </g>
        </g>
        <line x1="18" y1="156" x2="222" y2="156" class="stroke-gray-200 dark:stroke-dark-800" stroke-width="1" />
        <g class="fill-gray-400 dark:fill-dark-500" font-family="Geist Mono Variable, ui-monospace, monospace" font-size="8">
          <text x="18" y="170">{{ miniChartLabels[0] }}</text>
          <text x="120" y="170" text-anchor="middle">{{ miniChartLabels[3] }}</text>
          <text x="222" y="170" text-anchor="end">{{ miniChartLabels[7] }}</text>
        </g>
      </g>
    </g>
  </svg>
</template>

<script setup lang="ts">
// Translation function injected from the parent (HomeView) so this SVG stays
// presentational and never imports the i18n plugin directly.
import type { ComposerTranslation } from 'vue-i18n'

defineProps<{
  t: ComposerTranslation
}>()

const miniChartLabels = ['09:00', '10:00', '11:00', '12:00', '13:00', '14:00', '15:00', '16:00']
const miniChartBounds = {
  x: 18,
  y: 74,
  width: 204,
  height: 78,
  baseline: 156
}

interface MiniChartCard {
  key: string
  name: string
  labelKey: string
  x: number
  current: number
  delta: string
  stroke: string
  fill: string
  values: number[]
}

const miniChartCards: MiniChartCard[] = [
  {
    key: 'claude',
    name: 'Claude',
    labelKey: 'home.hero2.demo.sparkClaude',
    x: 206,
    current: 428,
    delta: '+18%',
    stroke: '#ef6f2e',
    fill: 'url(#sfSparkFill)',
    values: [314, 332, 349, 371, 366, 397, 412, 428]
  },
  {
    key: 'openai',
    name: 'OpenAI',
    labelKey: 'home.hero2.demo.sparkOpenai',
    x: 460,
    current: 312,
    delta: '+11%',
    stroke: '#10b981',
    fill: 'rgba(16,185,129,0.18)',
    values: [246, 238, 264, 252, 286, 301, 294, 312]
  },
  {
    key: 'gemini',
    name: 'Gemini',
    labelKey: 'home.hero2.demo.sparkGemini',
    x: 714,
    current: 186,
    delta: '+7%',
    stroke: '#38bdf8',
    fill: 'rgba(56,189,248,0.16)',
    values: [126, 141, 133, 152, 148, 169, 158, 186]
  }
]

function miniChartPoints(values: number[]) {
  const min = Math.min(...values)
  const max = Math.max(...values)
  const range = Math.max(max - min, 1)
  return values.map((value, index) => ({
    label: miniChartLabels[index],
    value,
    x: miniChartBounds.x + (index / (values.length - 1)) * miniChartBounds.width,
    y: miniChartBounds.y + (1 - (value - min) / range) * miniChartBounds.height
  }))
}

function miniChartLinePath(values: number[]) {
  return miniChartPoints(values)
    .map((point, index) => `${index === 0 ? 'M' : 'L'}${point.x.toFixed(1)} ${point.y.toFixed(1)}`)
    .join(' ')
}

function miniChartAreaPath(values: number[]) {
  const points = miniChartPoints(values)
  const line = points
    .map((point, index) => `${index === 0 ? 'M' : 'L'}${point.x.toFixed(1)} ${point.y.toFixed(1)}`)
    .join(' ')
  const first = points[0]
  const last = points[points.length - 1]
  return `${line} L${last.x.toFixed(1)} ${miniChartBounds.baseline} L${first.x.toFixed(1)} ${miniChartBounds.baseline} Z`
}
</script>
