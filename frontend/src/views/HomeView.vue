<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="hasHomeContent" class="min-h-screen">
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Compact Home Page -->
  <div
    v-else-if="compactHomeEnabled"
    data-testid="compact-home"
    class="flex min-h-screen flex-col bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white"
  >
    <header class="border-b border-gray-200 px-4 py-4 sm:px-6 dark:border-dark-800">
      <nav class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-3 sm:gap-4">
        <div class="flex min-w-0 flex-1 items-center gap-3">
          <img
            :src="siteLogo || '/logo.svg'"
            alt="Logo"
            class="h-9 w-9 shrink-0 rounded-lg object-contain"
          />
          <span class="min-w-0 truncate text-base font-semibold">{{ siteName }}</span>
        </div>
        <div class="flex max-w-full shrink-0 flex-wrap items-center justify-end gap-2">
          <LocaleSwitcher />
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>
          <button
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="inline-flex min-h-10 shrink-0 items-center justify-center rounded-lg bg-gray-900 px-4 py-2 text-sm font-medium text-white hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-200"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="flex min-w-0 flex-1 items-center justify-center px-4 py-16 sm:px-6">
      <div class="min-w-0 max-w-2xl text-center">
        <img
          :src="siteLogo || '/logo.svg'"
          alt="Logo"
          class="mx-auto mb-6 h-20 w-20 rounded-2xl object-contain"
        />
        <h1 class="[overflow-wrap:anywhere] text-3xl font-bold md:text-4xl">{{ siteName }}</h1>
        <p class="mt-4 whitespace-pre-wrap [overflow-wrap:anywhere] text-base text-gray-600 dark:text-dark-300">{{ siteSubtitle }}</p>
        <router-link
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="mt-8 inline-flex min-h-10 items-center justify-center rounded-lg bg-primary-600 px-5 py-2.5 text-sm font-medium text-white hover:bg-primary-700"
        >
          {{ isAuthenticated ? t('home.goToDashboard') : t('home.login') }}
        </router-link>
      </div>
    </main>

    <footer class="min-w-0 border-t border-gray-200 px-4 py-5 text-center text-sm text-gray-500 [overflow-wrap:anywhere] sm:px-6 dark:border-dark-800 dark:text-dark-400">
      &copy; {{ currentYear }} {{ siteName }}
    </footer>
  </div>

  <!-- Default Home Page (v2 — factory.ai motion rebuild) -->
  <div
    v-else
    class="relative flex min-h-screen flex-col overflow-x-hidden bg-gray-50 dark:bg-dark-950"
  >
    <!-- Backdrop: Factory-inspired ambient grid, confined to the first
         screen. From the marquee downward the surface is clean near-black /
         warm off-white. -->
    <div
      class="bg-blueprint-fade bg-factory-surface-grid pointer-events-none absolute inset-x-0 top-0 h-screen"
    ></div>

    <!-- ============================ NAV ============================ -->
    <header
      class="sf-nav fixed inset-x-0 top-0 z-60 border-b border-transparent transition-colors duration-200"
      :class="navScrolled
        ? 'border-gray-200/70 bg-white/80 backdrop-blur-md dark:border-dark-800/70 dark:bg-dark-950/80'
        : 'border-transparent'"
    >
      <nav
        class="mx-auto flex h-14 max-w-[1920px] items-center justify-between px-4 lg:px-9"
        aria-label="Main"
      >
        <!-- Brand wordmark -->
        <router-link
          to="/home"
          class="brand-hover group flex items-center gap-2.5"
          :class="brandClass"
          :aria-label="siteName"
          @mouseenter="onEnter"
          @mouseleave="onLeave"
          @focusin="onEnter"
          @focusout="onLeave"
        >
          <span class="brand-hover__logo block h-6 w-6 overflow-hidden rounded-sm">
            <img :src="siteLogo || '/logo.svg'" alt="" class="h-full w-full object-contain" />
          </span>
          <span
            class="font-mono text-[13px] uppercase tracking-[0.14em] text-gray-900 transition-colors duration-200 group-hover:text-primary-500 dark:text-gray-100"
          >
            <span class="sr-only">SiliconBase</span>
            <span aria-hidden="true">
              <span
                v-for="(ch, i) in 'SiliconBase'"
                :key="i"
                class="brand-hover__letter"
                :style="{ '--brand-i': i }"
                >{{ ch }}</span
              >
            </span>
          </span>
        </router-link>

        <!-- Desktop nav links + actions -->
        <div class="hidden items-center gap-6 md:flex">
          <div class="home-nav-hover-group flex items-center gap-6">
            <a
              v-if="docUrl"
              :href="docUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="home-nav-link font-mono text-[12px] uppercase tracking-[0.14em] text-gray-700 transition-colors duration-200 hover:text-primary-500 dark:text-gray-200"
            >{{ t('home.nav.docs') }}</a>
            <router-link
              to="/status"
              class="home-nav-link font-mono text-[12px] uppercase tracking-[0.14em] text-gray-800 transition-colors duration-200 hover:text-primary-500 dark:text-gray-100"
            >{{ t('home.nav.status') }}</router-link>
            <router-link
              to="/models"
              class="home-nav-link font-mono text-[12px] uppercase tracking-[0.14em] text-gray-800 transition-colors duration-200 hover:text-primary-500 dark:text-gray-100"
            >{{ t('home.nav.models') }}</router-link>
          </div>
          <LocaleSwitcher />
          <button
            @click="toggleTheme"
            class="rounded-sm p-1.5 text-gray-700 transition-colors duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-200 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            :aria-label="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>

          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="sf-btn group relative inline-flex h-[31px] items-center justify-center overflow-clip rounded-sm border border-dark-800 bg-dark-950 px-3.5 font-mono text-[12px] uppercase tracking-[0.12em] text-white transition-colors duration-150 hover:border-primary-500 hover:text-primary-500 dark:border-dark-700 dark:bg-dark-900 dark:text-gray-100"
          >
            <span
              class="mr-1.5 flex h-4 w-4 items-center justify-center rounded-sm bg-primary-500 text-[9px] font-semibold text-white"
            >{{ userInitial || '·' }}</span>
            <span>{{ t('home.dashboard') }}</span>
          </router-link>
          <router-link
            v-else
            to="/login"
            class="sf-btn group relative inline-flex h-[31px] items-center justify-center overflow-clip rounded-sm border border-dark-800 bg-dark-950 px-3.5 font-mono text-[12px] uppercase tracking-[0.12em] text-white transition-colors duration-150 hover:border-primary-500 hover:text-primary-500 dark:border-dark-700 dark:bg-dark-900 dark:text-gray-100"
          >{{ t('home.login') }}</router-link>
        </div>

        <!-- Mobile actions -->
        <div class="flex items-center gap-2 md:hidden">
          <LocaleSwitcher />
          <button
            @click="toggleTheme"
            class="rounded-sm p-1.5 text-gray-700 transition-colors duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-200 dark:hover:bg-dark-800 dark:hover:text-white"
            :aria-label="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
          <button
            @click="mobileOpen = true"
            class="rounded-sm p-1.5 text-gray-700 transition-colors duration-150 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-dark-800"
            :aria-label="t('home.nav.openMenu')"
            :aria-expanded="mobileOpen"
          >
            <Icon name="menu" size="md" />
          </button>
        </div>
      </nav>
    </header>

    <!-- Mobile off-canvas drawer -->
    <Transition name="sf-drawer">
      <div
        v-if="mobileOpen"
        class="fixed inset-0 z-70 md:hidden"
        :aria-hidden="!mobileOpen"
      >
        <div
          class="absolute inset-0 bg-dark-950/60 backdrop-blur-sm"
          @click="mobileOpen = false"
        ></div>
        <aside
          class="sf-drawer-panel absolute left-0 top-0 flex h-full w-[280px] max-w-[80vw] flex-col border-r border-gray-200 bg-white pt-4 dark:border-dark-800 dark:bg-dark-950"
        >
          <div class="flex items-center justify-between px-5 pb-4">
            <span
              class="font-mono text-[13px] uppercase tracking-[0.14em] text-gray-900 dark:text-gray-100"
            >SiliconBase</span>
            <button
              @click="mobileOpen = false"
              class="rounded-sm p-1.5 text-gray-500 hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
              :aria-label="t('home.nav.closeMenu')"
            >
              <Icon name="x" size="md" />
            </button>
          </div>
          <nav class="flex flex-col gap-1 px-3" aria-label="Mobile">
            <a
              v-if="docUrl"
              :href="docUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="rounded-sm px-3 py-2 font-mono text-[12px] uppercase tracking-[0.12em] text-gray-600 transition-colors duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white"
            >{{ t('home.nav.docs') }}</a>
            <router-link
              to="/status"
              class="rounded-sm px-3 py-2 font-mono text-[12px] uppercase tracking-[0.12em] text-gray-600 transition-colors duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white"
              @click="mobileOpen = false"
            >{{ t('home.nav.status') }}</router-link>
            <router-link
              to="/models"
              class="rounded-sm px-3 py-2 font-mono text-[12px] uppercase tracking-[0.12em] text-gray-600 transition-colors duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white"
              @click="mobileOpen = false"
            >{{ t('home.nav.models') }}</router-link>
            <router-link
              :to="isAuthenticated ? dashboardPath : '/login'"
              class="mt-2 inline-flex items-center justify-center rounded-sm border border-dark-800 bg-dark-950 px-3 py-2 font-mono text-[12px] uppercase tracking-[0.12em] text-white dark:border-dark-700 dark:bg-dark-900"
            >{{ isAuthenticated ? t('home.dashboard') : t('home.login') }}</router-link>
          </nav>
        </aside>
      </div>
    </Transition>

    <!-- ============================ HERO ============================ -->
    <section
      class="relative z-10 mx-auto w-full max-w-[1920px] px-4 pb-16 pt-28 lg:px-9 lg:pb-24 lg:pt-36"
    >
      <div class="grid grid-cols-1 items-center gap-10 lg:grid-cols-12 lg:gap-6">
        <!-- Text column (lg:col-span-5) -->
        <div class="lg:col-span-5 lg:self-center">
          <p
            data-home-hero-intro
            :style="{ '--sf-delay': '0ms' }"
            class="eyebrow mb-5"
          >{{ t('home.hero2.eyebrow') }}</p>
          <h1
            data-home-hero-intro
            :style="{ '--sf-delay': '120ms' }"
            class="mb-6 max-w-[14ch] text-[clamp(40px,7vw,72px)] font-normal leading-[100%] tracking-[-0.04em] text-gray-900 lg:tracking-[-0.06em] dark:text-white"
          >
            {{ t('home.hero2.title') }}
          </h1>
          <p
            data-home-hero-intro
            :style="{ '--sf-delay': '240ms' }"
            class="mb-9 max-w-[46ch] font-mono text-[14px] leading-[150%] tracking-[-0.01em] text-gray-600 lg:text-[16px] dark:text-dark-300"
          >
            {{ t('home.hero2.subhead') }}
          </p>

          <div
            data-home-hero-intro
            :style="{ '--sf-delay': '360ms' }"
            class="flex flex-wrap items-center gap-3"
          >
            <router-link
              :to="isAuthenticated ? dashboardPath : '/register'"
              class="sf-btn group relative inline-flex h-[40px] items-center justify-center overflow-clip rounded-sm border border-dark-800 bg-dark-950 px-5 font-mono text-[13px] uppercase tracking-[0.12em] text-white transition-colors duration-150 hover:border-primary-500 hover:text-primary-500 dark:border-dark-700 dark:bg-dark-900 dark:text-gray-100"
            >
              <span class="relative z-10 flex items-center gap-2">
                {{ isAuthenticated ? t('home.goToDashboard') : t('home.hero2.ctaPrimary') }}
                <Icon name="arrowRight" size="sm" :stroke-width="2" />
              </span>
            </router-link>
            <a
              v-if="docUrl"
              :href="docUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="sf-btn group relative inline-flex h-[40px] items-center justify-center overflow-clip rounded-sm border border-gray-300 bg-transparent px-5 font-mono text-[13px] uppercase tracking-[0.12em] text-gray-700 transition-colors duration-150 hover:border-primary-500 hover:text-primary-500 dark:border-dark-700 dark:text-gray-200"
            >
              <span class="relative z-10 flex items-center gap-2">
                {{ t('home.hero2.ctaSecondary') }}
                <Icon name="arrowRight" size="sm" :stroke-width="2" />
              </span>
            </a>
          </div>
        </div>

        <!-- Demo column (lg:col-span-7) — bespoke inline SVG gateway dashboard -->
        <div class="lg:col-span-7">
          <div
            data-hero-demo
            class="sf-dashboard-frame relative w-full overflow-hidden rounded-[16px] border border-gray-200 bg-white shadow-[0_45px_120px_rgba(0,0,0,0.18)] lg:-mr-16 dark:border-dark-800 dark:bg-dark-900 dark:shadow-[0_45px_120px_rgba(0,0,0,0.55)]"
          >
            <GatewayDashboard :t="t" />
          </div>
        </div>
      </div>
    </section>

    <!-- ============================ COMPATIBILITY MARQUEES ============================ -->
    <section class="relative z-10 overflow-hidden border-y border-dashed border-gray-200 py-6 dark:border-dark-800">
      <div class="mx-auto mb-4 flex max-w-[1920px] items-center gap-3 px-4 lg:px-9">
        <span class="eyebrow whitespace-nowrap">{{ t('home.marquee2.providersLabel') }}</span>
        <span class="h-px flex-1 bg-gray-200 dark:bg-dark-800"></span>
      </div>
      <div class="sf-marquee relative overflow-hidden">
        <div class="sf-marquee-track flex w-max items-center gap-4 pr-4">
          <template v-for="n in 3" :key="`provider-${n}`">
            <span
              v-for="item in providerMarqueeItems"
              :key="`provider-${n}-${item.name}`"
              class="sf-logo-pill group"
            >
              <span :class="['sf-brand-mark', `sf-brand-mark--${item.tone}`]">{{ item.mark }}</span>
              <span>{{ item.name }}</span>
            </span>
          </template>
        </div>
        <div class="sf-marquee-fade pointer-events-none absolute inset-y-0 left-0 w-24 bg-gradient-to-r from-gray-50 to-transparent dark:from-dark-950"></div>
        <div class="sf-marquee-fade pointer-events-none absolute inset-y-0 right-0 w-24 bg-gradient-to-l from-gray-50 to-transparent dark:from-dark-950"></div>
      </div>

      <div class="mx-auto mb-4 mt-7 flex max-w-[1920px] items-center gap-3 px-4 lg:px-9">
        <span class="eyebrow whitespace-nowrap">{{ t('home.marquee2.toolsLabel') }}</span>
        <span class="h-px flex-1 bg-gray-200 dark:bg-dark-800"></span>
      </div>
      <div class="sf-marquee relative overflow-hidden">
        <div class="sf-marquee-track sf-marquee-track--reverse flex w-max items-center gap-4 pr-4">
          <template v-for="n in 3" :key="`tool-${n}`">
            <span
              v-for="item in toolMarqueeItems"
              :key="`tool-${n}-${item.name}`"
              class="sf-logo-pill group"
            >
              <span :class="['sf-brand-mark', `sf-brand-mark--${item.tone}`]">{{ item.mark }}</span>
              <span>{{ item.name }}</span>
            </span>
          </template>
        </div>
        <div class="sf-marquee-fade pointer-events-none absolute inset-y-0 left-0 w-24 bg-gradient-to-r from-gray-50 to-transparent dark:from-dark-950"></div>
        <div class="sf-marquee-fade pointer-events-none absolute inset-y-0 right-0 w-24 bg-gradient-to-l from-gray-50 to-transparent dark:from-dark-950"></div>
      </div>
    </section>

    <!-- ============================ DEFINING ============================ -->
    <section class="relative z-10 mx-auto w-full max-w-[1920px] px-4 py-20 lg:px-9 lg:py-28">
      <div class="mb-10 lg:mb-14">
        <h2 class="max-w-[22ch] text-[clamp(30px,4.4vw,52px)] font-normal leading-[110%] tracking-[-0.035em] text-gray-900 lg:tracking-[-0.05em] dark:text-white">
          {{ t('home.defining2.title') }}
        </h2>
      </div>

      <div class="grid grid-cols-1 gap-4 lg:grid-cols-3">
        <article v-reveal :style="{ '--sf-delay': '0ms' }" class="sf-defining-card">
          <p class="font-mono text-[12px] text-primary-500">{{ t('home.defining2.cards.models.index') }}</p>
          <h3 class="mt-3 text-xl font-normal tracking-tight text-gray-900 dark:text-white">
            {{ t('home.defining2.cards.models.title') }}
          </h3>
          <div class="sf-defining-radar mt-9">
            <svg viewBox="0 0 320 280" class="h-full w-full" aria-hidden="true">
              <g transform="translate(160 132)">
                <g class="sf-model-radar-grid">
                  <polygon points="0,-92 80,-46 80,46 0,92 -80,46 -80,-46" />
                  <polygon points="0,-69 60,-34.5 60,34.5 0,69 -60,34.5 -60,-34.5" />
                  <polygon points="0,-46 40,-23 40,23 0,46 -40,23 -40,-23" />
                  <polygon points="0,-23 20,-11.5 20,11.5 0,23 -20,11.5 -20,-11.5" />
                  <line x1="0" y1="-92" x2="0" y2="92" />
                  <line x1="80" y1="-46" x2="-80" y2="46" />
                  <line x1="80" y1="46" x2="-80" y2="-46" />
                </g>
                <polygon class="sf-model-radar-model sf-model-radar-model--a" points="0,-87 58,-36 58,31 0,46 -66,38 -58,-31" />
                <polygon class="sf-model-radar-model sf-model-radar-model--b" points="0,-72 50,-26 69,39 0,56 -46,25 -68,-39" />
                <polygon class="sf-model-radar-model sf-model-radar-model--c" points="0,-58 42,-20 47,27 0,61 -53,30 -38,-22" />
                <polygon class="sf-model-radar-best-path" points="0,-88 73,-42 67,38 0,78 -72,42 -75,-43" />
                <polygon class="sf-model-radar-best-path sf-model-radar-best-path--runner" points="0,-88 73,-42 67,38 0,78 -72,42 -75,-43" />
                <g class="sf-model-radar-best-dots">
                  <circle cx="0" cy="-88" r="3" />
                  <circle cx="73" cy="-42" r="3" />
                  <circle cx="67" cy="38" r="3" />
                  <circle cx="0" cy="78" r="3" />
                  <circle cx="-72" cy="42" r="3" />
                  <circle cx="-75" cy="-43" r="3" />
                </g>
              </g>
              <g font-family="Geist Mono Variable, ui-monospace, monospace" font-size="10" class="fill-gray-500 dark:fill-dark-400">
                <text x="142" y="18">{{ t('home.defining2.cards.models.axes.coding') }}</text>
                <text x="236" y="82">{{ t('home.defining2.cards.models.axes.reasoning') }}</text>
                <text x="236" y="180">{{ t('home.defining2.cards.models.axes.speed') }}</text>
                <text x="146" y="252">{{ t('home.defining2.cards.models.axes.cost') }}</text>
                <text x="24" y="180">{{ t('home.defining2.cards.models.axes.context') }}</text>
                <text x="42" y="82">{{ t('home.defining2.cards.models.axes.tools') }}</text>
              </g>
            </svg>
          </div>
          <div class="sf-model-radar-legend mt-7">
            <span><i class="sf-model-radar-legend-line sf-model-radar-legend-line--a"></i>{{ t('home.defining2.cards.models.legend.primary') }}</span>
            <span><i class="sf-model-radar-legend-line sf-model-radar-legend-line--b"></i>{{ t('home.defining2.cards.models.legend.backup') }}</span>
            <span><i class="sf-model-radar-legend-line sf-model-radar-legend-line--c"></i>{{ t('home.defining2.cards.models.legend.third') }}</span>
            <span class="text-primary-500"><i class="sf-model-radar-legend-line sf-model-radar-legend-line--best"></i>{{ t('home.defining2.cards.models.legend.best') }}</span>
          </div>
        </article>

        <article v-reveal :style="{ '--sf-delay': '120ms' }" class="sf-defining-card">
          <p class="font-mono text-[12px] text-primary-500">{{ t('home.defining2.cards.runtime.index') }}</p>
          <h3 class="mt-3 text-xl font-normal tracking-tight text-gray-900 dark:text-white">
            {{ t('home.defining2.cards.runtime.title') }}
          </h3>
          <div class="mt-9 flex flex-1 flex-col">
            <div
              v-for="row in runtimeRows"
              :key="row.key"
              class="flex items-center justify-between border-t border-gray-200 py-4 first:border-t-0 dark:border-dark-800"
            >
              <div>
                <p class="text-[15px] text-gray-900 dark:text-white">{{ t(`home.defining2.cards.runtime.rows.${row.key}.title`) }}</p>
                <p class="mt-1 font-mono text-[12px] text-gray-500 dark:text-dark-400">{{ t(`home.defining2.cards.runtime.rows.${row.key}.desc`) }}</p>
              </div>
              <Icon name="check" size="sm" class="text-primary-500" />
            </div>
          </div>
          <p class="mt-auto pt-7 font-mono text-[12px] text-gray-500 dark:text-dark-400">
            {{ t('home.defining2.cards.runtime.note') }}
          </p>
        </article>

        <article v-reveal :style="{ '--sf-delay': '240ms' }" class="sf-defining-card">
          <p class="font-mono text-[12px] text-primary-500">{{ t('home.defining2.cards.workflow.index') }}</p>
          <h3 class="mt-3 text-xl font-normal tracking-tight text-gray-900 dark:text-white">
            {{ t('home.defining2.cards.workflow.title') }}
          </h3>
          <div class="sf-sdlc-orbit mt-8">
            <svg class="sf-sdlc-path" viewBox="0 0 340 340" aria-hidden="true">
              <circle cx="170" cy="170" r="138" />
            </svg>
            <span class="sf-sdlc-center">{{ t('home.defining2.cards.workflow.center') }}</span>
            <div class="sf-sdlc-orbit-stage">
              <span
                v-for="(node, index) in workflowNodes"
                :key="node.key"
                class="sf-sdlc-node-anchor"
                :style="{ '--sf-node-angle': `${node.angle}deg`, '--sf-node-angle-inverse': `${-node.angle}deg`, '--sf-node-delay': `${index * 160}ms` }"
              >
                <span class="sf-sdlc-node">
                  {{ t(`home.defining2.cards.workflow.nodes.${node.key}`) }}
                </span>
              </span>
            </div>
          </div>
          <p class="mt-auto pt-7 font-mono text-[12px] leading-relaxed text-gray-500 dark:text-dark-400">
            {{ t('home.defining2.cards.workflow.note') }}
          </p>
        </article>
      </div>
    </section>

    <!-- ============================ BENTO ============================ -->
    <section class="relative z-10 mx-auto w-full max-w-[1920px] px-4 py-20 lg:px-9 lg:py-28">
      <div class="mb-10 flex flex-col gap-3">
        <p class="eyebrow">{{ t('home.bento2.eyebrow') }}</p>
        <h2 class="max-w-[20ch] text-[clamp(28px,4vw,44px)] font-normal leading-[110%] tracking-[-0.03em] text-gray-900 lg:tracking-[-0.04em] dark:text-white">
          {{ t('home.bento2.title') }}
        </h2>
      </div>

      <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
        <!-- 01 Pooling -->
        <article
          v-reveal
          :style="{ '--sf-delay': '0ms' }"
          class="group relative flex flex-col gap-5 rounded-lg border border-gray-200 bg-white p-8 transition-colors duration-150 hover:border-gray-300 dark:border-dark-800 dark:bg-dark-900 dark:hover:border-dark-700"
        >
          <div class="flex items-center gap-3">
            <span class="font-mono text-[13px] text-primary-500">{{ t('home.bento2.cards.pooling.index') }}</span>
            <span class="h-px flex-1 bg-gray-200 dark:bg-dark-800"></span>
          </div>
          <PoolingChart />
          <h3 class="text-lg font-normal tracking-tight text-gray-900 dark:text-white">
            {{ t('home.bento2.cards.pooling.title') }}
          </h3>
          <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-300">
            {{ t('home.bento2.cards.pooling.desc') }}
          </p>
        </article>

        <!-- 02 Routing -->
        <article
          v-reveal
          :style="{ '--sf-delay': '120ms' }"
          class="group relative flex flex-col gap-5 rounded-lg border border-gray-200 bg-white p-8 transition-colors duration-150 hover:border-gray-300 dark:border-dark-800 dark:bg-dark-900 dark:hover:border-dark-700"
        >
          <div class="flex items-center gap-3">
            <span class="font-mono text-[13px] text-primary-500">{{ t('home.bento2.cards.routing.index') }}</span>
            <span class="h-px flex-1 bg-gray-200 dark:bg-dark-800"></span>
          </div>
          <RoutingDiagram />
          <h3 class="text-lg font-normal tracking-tight text-gray-900 dark:text-white">
            {{ t('home.bento2.cards.routing.title') }}
          </h3>
          <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-300">
            {{ t('home.bento2.cards.routing.desc') }}
          </p>
        </article>

        <!-- 03 Billing -->
        <article
          v-reveal
          :style="{ '--sf-delay': '240ms' }"
          class="group relative flex flex-col gap-5 rounded-lg border border-gray-200 bg-white p-8 transition-colors duration-150 hover:border-gray-300 dark:border-dark-800 dark:bg-dark-900 dark:hover:border-dark-700"
        >
          <div class="flex items-center gap-3">
            <span class="font-mono text-[13px] text-primary-500">{{ t('home.bento2.cards.billing.index') }}</span>
            <span class="h-px flex-1 bg-gray-200 dark:bg-dark-800"></span>
          </div>
          <BillingSparkline />
          <h3 class="text-lg font-normal tracking-tight text-gray-900 dark:text-white">
            {{ t('home.bento2.cards.billing.title') }}
          </h3>
          <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-300">
            {{ t('home.bento2.cards.billing.desc') }}
          </p>
        </article>
      </div>
    </section>

    <!-- ============================ CTA BAND ============================ -->
    <section class="relative z-10 mx-auto w-full max-w-[1920px] px-4 py-16 lg:px-9 lg:py-24">
      <div
        v-reveal
        class="flex flex-col items-start gap-5 border-t border-dashed border-gray-200 pt-12 dark:border-dark-800"
      >
        <p class="eyebrow">{{ t('home.cta2.eyebrow') }}</p>
        <h2 class="max-w-[24ch] text-[clamp(28px,4vw,48px)] font-normal leading-[110%] tracking-[-0.03em] text-gray-900 lg:tracking-[-0.04em] dark:text-white">
          {{ t('home.cta2.title') }}
        </h2>
        <p class="max-w-[60ch] font-mono text-[14px] leading-[150%] text-gray-600 dark:text-dark-300">
          {{ t('home.cta2.note') }}
        </p>
        <router-link
          :to="isAuthenticated ? dashboardPath : '/register'"
          class="sf-btn group relative mt-2 inline-flex h-[44px] items-center justify-center overflow-clip rounded-sm border border-dark-800 bg-dark-950 px-6 font-mono text-[13px] uppercase tracking-[0.12em] text-white transition-colors duration-150 hover:border-primary-500 hover:text-primary-500 dark:border-dark-700 dark:bg-dark-900 dark:text-gray-100"
        >
          <span class="relative z-10 flex items-center gap-2">
            {{ t('home.cta2.button') }}
            <Icon name="arrowRight" size="sm" :stroke-width="2" />
          </span>
        </router-link>
      </div>
    </section>

    <PublicFooter :doc-url="docUrl" :site-logo="siteLogo" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import { useBrandHover } from '@/composables/useBrandHover'
import { useTheme } from '@/composables/useTheme'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'
import PublicFooter from '@/components/home/PublicFooter.vue'
import GatewayDashboard from '@/components/home/GatewayDashboard.vue'
import PoolingChart from '@/components/home/PoolingChart.vue'
import RoutingDiagram from '@/components/home/RoutingDiagram.vue'
import BillingSparkline from '@/components/home/BillingSparkline.vue'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// Brand wordmark micro-interaction (shared phase state with the sidebar +
// auth layout). Logo spins forward on hover and reverse-unwinds on leave;
// letters do a per-letter 3D flip. See composables/useBrandHover.ts.
// Destructure so each ref/handler is a top-level setup binding — Vue only
// auto-unwraps refs at the top level of the bindings object, not through
// nested property access. Binding `brand.brandClass` directly would pass the
// ComputedRef object to `:class` instead of its value (animation never fires).
const { brandClass, onEnter, onLeave } = useBrandHover()

// Site settings - directly from appStore (already initialized from injected config)
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'SiliconBase')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || t('home.heroSubtitle'))
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const hasHomeContent = computed(() => homeContent.value.trim().length > 0)
const compactHomeEnabled = computed(() => appStore.cachedPublicSettings?.compact_home_enabled === true)
const currentYear = new Date().getFullYear()

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Theme
const { isDark, toggleTheme, syncThemeFromDocument } = useTheme()

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => (isAdmin.value ? '/admin/dashboard' : '/dashboard'))
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})

type LogoTone = 'warm' | 'green' | 'blue' | 'violet' | 'slate'
interface LogoMarqueeItem {
  name: string
  mark: string
  tone: LogoTone
}

const providerMarqueeItems: LogoMarqueeItem[] = [
  { name: 'Claude', mark: '✦', tone: 'warm' },
  { name: 'OpenAI', mark: 'OA', tone: 'green' },
  { name: 'Gemini', mark: 'G', tone: 'blue' },
  { name: 'Bedrock', mark: 'AWS', tone: 'slate' },
  { name: 'Antigravity', mark: 'AG', tone: 'violet' },
  { name: 'DeepSeek', mark: 'DS', tone: 'blue' },
  { name: 'Moonshot', mark: 'K', tone: 'slate' },
  { name: 'Qwen', mark: 'Q', tone: 'violet' },
  { name: 'OpenRouter', mark: 'OR', tone: 'green' }
]

const toolMarqueeItems: LogoMarqueeItem[] = [
  { name: 'Claude Code', mark: 'CC', tone: 'warm' },
  { name: 'Codex', mark: 'CX', tone: 'green' },
  { name: 'OpenCode', mark: 'OC', tone: 'slate' },
  { name: 'Trae', mark: 'T', tone: 'violet' },
  { name: 'WorkBuddy', mark: 'WB', tone: 'blue' },
  { name: 'Cursor', mark: '↗', tone: 'slate' },
  { name: 'Windsurf', mark: 'W', tone: 'blue' },
  { name: 'Cline', mark: 'CL', tone: 'green' },
  { name: 'Roo Code', mark: 'RC', tone: 'warm' },
  { name: 'Aider', mark: 'AI', tone: 'violet' },
  { name: 'GitHub Copilot', mark: 'GH', tone: 'slate' },
  { name: 'Gemini CLI', mark: 'G', tone: 'blue' },
  { name: 'Continue', mark: 'CT', tone: 'green' },
  { name: 'Zed', mark: 'Z', tone: 'slate' }
]

const runtimeRows = [
  { key: 'saas' },
  { key: 'hybrid' },
  { key: 'selfHosted' },
  { key: 'isolated' }
] as const

const workflowNodes = [
  { key: 'plan', angle: -90 },
  { key: 'execute', angle: -45 },
  { key: 'validate', angle: 0 },
  { key: 'ship', angle: 45 },
  { key: 'monitor', angle: 90 },
  { key: 'automate', angle: 135 },
  { key: 'signal', angle: 180 },
  { key: 'triage', angle: -135 }
] as const

// Nav scroll state (subtle border + backdrop on scroll)
const mobileOpen = ref(false)
const navScrolled = ref(false)
function handleScroll() {
  navScrolled.value = window.scrollY > 8
}

onMounted(() => {
  // Theme is already bootstrapped by main.ts initThemeClass(); sync local ref.
  syncThemeFromDocument()

  // Auth + settings (kept from original logic)
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }

  // Nav scroll listener
  handleScroll()
  window.addEventListener('scroll', handleScroll, { passive: true })
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
/* Mobile drawer transition */
.sf-drawer-enter-active,
.sf-drawer-leave-active {
  transition: opacity 0.2s ease-in-out;
}
.sf-drawer-enter-from,
.sf-drawer-leave-to {
  opacity: 0;
}
.sf-drawer-enter-active .sf-drawer-panel,
.sf-drawer-leave-active .sf-drawer-panel {
  transition: transform 0.25s cubic-bezier(0.22, 1, 0.36, 1);
}
.sf-drawer-enter-from .sf-drawer-panel,
.sf-drawer-leave-to .sf-drawer-panel {
  transform: translateX(-100%);
}

/* Respect reduced motion for the drawer + nav transition */
@media (prefers-reduced-motion: reduce) {
  .sf-drawer-enter-active,
  .sf-drawer-leave-active,
  .sf-drawer-enter-active .sf-drawer-panel,
  .sf-drawer-leave-active .sf-drawer-panel {
    transition: none;
  }
}

.home-nav-hover-group:hover .home-nav-link:not(:hover) {
  opacity: 0.42;
}
</style>
