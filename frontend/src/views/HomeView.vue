<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
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

  <!-- Default Home Page (v2 — factory.ai motion rebuild) -->
  <div
    v-else
    class="relative flex min-h-screen flex-col overflow-x-hidden bg-gray-50 dark:bg-dark-950"
  >
    <!-- Backdrop: blueprint line-grid + faint vermillion mesh (PR1 utility kept). -->
    <div class="pointer-events-none absolute inset-0 bg-blueprint"></div>
    <div
      class="pointer-events-none absolute inset-0 bg-[radial-gradient(at_40%_12%,rgba(239,111,46,0.10)_0px,transparent_55%),radial-gradient(at_85%_0%,rgba(209,80,16,0.06)_0px,transparent_50%)]"
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
          class="group flex items-center gap-2.5"
          :aria-label="siteName"
        >
          <span class="block h-6 w-6 overflow-hidden rounded-sm">
            <img :src="siteLogo || '/logo.svg'" alt="" class="h-full w-full object-contain" />
          </span>
          <span
            class="font-mono text-[13px] uppercase tracking-[0.14em] text-gray-900 transition-colors duration-200 group-hover:text-primary-500 dark:text-gray-100"
          >SiliconBase</span>
        </router-link>

        <!-- Desktop nav links + actions -->
        <div class="hidden items-center gap-6 md:flex">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="font-mono text-[12px] uppercase tracking-[0.14em] text-gray-500 transition-colors duration-200 hover:text-primary-500 dark:text-dark-400"
          >{{ t('home.nav.docs') }}</a>
          <LocaleSwitcher />
          <button
            @click="toggleTheme"
            class="rounded-sm p-1.5 text-gray-500 transition-colors duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
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
            class="rounded-sm p-1.5 text-gray-500 transition-colors duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
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

    <!-- ============================ LOGO MARQUEE ============================ -->
    <section class="relative z-10 overflow-hidden border-y border-dashed border-gray-200 py-6 dark:border-dark-800">
      <div class="mx-auto mb-4 flex max-w-[1920px] items-center gap-3 px-4 lg:px-9">
        <span class="eyebrow whitespace-nowrap">{{ t('home.marquee2.label') }}</span>
        <span class="h-px flex-1 bg-gray-200 dark:bg-dark-800"></span>
      </div>
      <div class="sf-marquee relative overflow-hidden">
        <div class="sf-marquee-track flex w-max items-center gap-12 pr-12">
          <!-- duplicated track for seamless loop -->
          <template v-for="n in 2" :key="n">
            <span
              v-for="(item, idx) in marqueeItems"
              :key="`${n}-${idx}`"
              class="group inline-flex items-center gap-2 whitespace-nowrap font-mono text-[14px] uppercase tracking-[-0.01em] text-gray-400 transition-colors duration-200 hover:text-primary-500 dark:text-dark-400"
            >
              <span class="inline-block h-1.5 w-1.5 rounded-full bg-primary-500/60 group-hover:bg-primary-500"></span>
              {{ item }}
            </span>
          </template>
        </div>
        <!-- edge fades -->
        <div class="sf-marquee-fade pointer-events-none absolute inset-y-0 left-0 w-24 bg-gradient-to-r from-gray-50 to-transparent dark:from-dark-950"></div>
        <div class="sf-marquee-fade pointer-events-none absolute inset-y-0 right-0 w-24 bg-gradient-to-l from-gray-50 to-transparent dark:from-dark-950"></div>
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

    <!-- ============================ FOOTER PANEL ============================ -->
    <footer class="relative z-10 mx-auto w-full max-w-[1920px] px-4 pb-10 lg:px-9">
      <div
        class="rounded-lg border border-gray-200 bg-white px-6 py-10 dark:border-dark-800 dark:bg-dark-900 lg:px-10 lg:py-12"
      >
        <div class="grid grid-cols-2 gap-8 lg:grid-cols-12 lg:gap-6">
          <!-- Brand block -->
          <div class="col-span-2 flex flex-col gap-3 lg:col-span-4">
            <div class="flex items-center gap-2.5">
              <span class="block h-6 w-6 overflow-hidden rounded-sm">
                <img :src="siteLogo || '/logo.svg'" alt="" class="h-full w-full object-contain" />
              </span>
              <span class="font-mono text-[13px] uppercase tracking-[0.14em] text-gray-900 dark:text-gray-100">SiliconBase</span>
            </div>
            <p class="max-w-[36ch] font-mono text-[12px] uppercase tracking-[0.12em] text-gray-500 dark:text-dark-400">
              {{ t('home.footer2.tagline') }}
            </p>
            <a
              :href="githubUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="font-mono text-[12px] uppercase tracking-[0.12em] text-primary-500 transition-colors duration-200 hover:text-primary-400"
            >{{ t('home.footer2.domain') }}</a>
          </div>

          <!-- Resources column -->
          <div class="flex flex-col gap-3 lg:col-span-3 lg:col-start-7">
            <p class="eyebrow">{{ t('home.footer2.columns.resources') }}</p>
            <a
              v-if="docUrl"
              :href="docUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="font-mono text-[13px] tracking-[-0.01em] text-gray-600 transition-colors duration-200 hover:text-primary-500 dark:text-dark-300"
            >{{ t('home.footer2.links.docs') }}</a>
            <a
              :href="githubUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="font-mono text-[13px] tracking-[-0.01em] text-gray-600 transition-colors duration-200 hover:text-primary-500 dark:text-dark-300"
            >{{ t('home.footer2.links.status') }}</a>
            <a
              :href="githubUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="font-mono text-[13px] tracking-[-0.01em] text-gray-600 transition-colors duration-200 hover:text-primary-500 dark:text-dark-300"
            >{{ t('home.footer2.links.changelog') }}</a>
          </div>

          <!-- Legal column -->
          <div class="flex flex-col gap-3 lg:col-span-3">
            <p class="eyebrow">{{ t('home.footer2.columns.legal') }}</p>
            <a
              :href="githubUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="font-mono text-[13px] tracking-[-0.01em] text-gray-600 transition-colors duration-200 hover:text-primary-500 dark:text-dark-300"
            >{{ t('home.footer2.links.privacy') }}</a>
            <a
              :href="githubUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="font-mono text-[13px] tracking-[-0.01em] text-gray-600 transition-colors duration-200 hover:text-primary-500 dark:text-dark-300"
            >{{ t('home.footer2.links.terms') }}</a>
          </div>
        </div>

        <!-- Bottom row -->
        <div class="mt-10 flex flex-col items-start justify-between gap-3 border-t border-dashed border-gray-200 pt-6 dark:border-dark-800 md:flex-row md:items-center">
          <p class="font-mono text-[11px] uppercase tracking-[0.12em] text-gray-500 dark:text-dark-400">
            {{ t('home.footer2.copyright') }}
          </p>
          <p class="font-mono text-[11px] uppercase tracking-[0.12em] text-gray-400 dark:text-dark-400">
            {{ t('home.footer2.domain') }}
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'
import GatewayDashboard from '@/components/home/GatewayDashboard.vue'
import PoolingChart from '@/components/home/PoolingChart.vue'
import RoutingDiagram from '@/components/home/RoutingDiagram.vue'
import BillingSparkline from '@/components/home/BillingSparkline.vue'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// Site settings - directly from appStore (already initialized from injected config)
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'SiliconBase')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Theme
const isDark = ref(document.documentElement.classList.contains('dark'))

// Brand / marketing URL
const githubUrl = 'https://siliconbase.link'

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => (isAdmin.value ? '/admin/dashboard' : '/dashboard'))
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})

// Marquee items — upstream providers row
const marqueeItems = computed(() => t('home.marquee2.items').split('·').map((s) => s.trim()))

// Nav scroll state (subtle border + backdrop on scroll)
const mobileOpen = ref(false)
const navScrolled = ref(false)
function handleScroll() {
  navScrolled.value = window.scrollY > 8
}

// Toggle theme
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

onMounted(() => {
  // Theme is already bootstrapped by main.ts initThemeClass(); sync local ref.
  isDark.value = document.documentElement.classList.contains('dark')

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
</style>
