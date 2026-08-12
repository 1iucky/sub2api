<template>
  <header class="sticky top-0 z-50 border-b border-gray-200/70 bg-white/85 backdrop-blur-md dark:border-dark-800/70 dark:bg-dark-950/85">
    <nav class="mx-auto flex h-14 max-w-[1920px] items-center justify-between px-4 lg:px-9" aria-label="Main">
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
        <span class="font-mono text-[13px] uppercase tracking-[0.14em] text-gray-900 transition-colors duration-200 group-hover:text-primary-500 dark:text-gray-100">
          <span class="sr-only">{{ siteName }}</span>
          <span aria-hidden="true">
            <span
              v-for="(ch, i) in siteName"
              :key="i"
              class="brand-hover__letter"
              :style="{ '--brand-i': i }"
            >{{ ch }}</span>
          </span>
        </span>
      </router-link>

      <div class="hidden items-center gap-6 md:flex">
        <div class="nav-hover-group flex items-center gap-6">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="nav-link"
          >{{ t('home.nav.docs') }}</a>
          <component
            :is="isStatusActive ? 'span' : RouterLink"
            :to="isStatusActive ? undefined : '/status'"
            class="nav-link"
            :class="isStatusActive ? 'nav-link-active cursor-default' : ''"
            :aria-current="isStatusActive ? 'page' : undefined"
            @click="isStatusActive && $event.preventDefault()"
          >{{ t('home.nav.status') }}</component>
          <component
            :is="isModelsActive ? 'span' : RouterLink"
            :to="isModelsActive ? undefined : '/models'"
            class="nav-link"
            :class="isModelsActive ? 'nav-link-active cursor-default' : ''"
            :aria-current="isModelsActive ? 'page' : undefined"
            @click="isModelsActive && $event.preventDefault()"
          >{{ t('home.nav.models') }}</component>
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
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="inline-flex h-[31px] items-center justify-center rounded-sm border border-dark-800 bg-dark-950 px-3.5 font-mono text-[12px] uppercase tracking-[0.12em] text-white transition-colors duration-150 hover:border-primary-500 hover:text-primary-500 dark:border-dark-700 dark:bg-dark-900 dark:text-gray-100"
        >
          <span
            v-if="isAuthenticated"
            class="mr-1.5 flex h-4 w-4 items-center justify-center rounded-sm bg-primary-500 text-[9px] font-semibold text-white"
          >{{ userInitial || '·' }}</span>
          <span>{{ isAuthenticated ? t('home.dashboard') : t('home.login') }}</span>
        </router-link>
      </div>

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

  <Transition name="sf-drawer">
    <div v-if="mobileOpen" class="fixed inset-0 z-[70] md:hidden" :aria-hidden="!mobileOpen">
      <div class="absolute inset-0 bg-dark-950/60 backdrop-blur-sm" @click="mobileOpen = false"></div>
      <aside class="sf-drawer-panel absolute left-0 top-0 flex h-full w-[280px] max-w-[80vw] flex-col border-r border-gray-200 bg-white pt-4 dark:border-dark-800 dark:bg-dark-950">
        <div class="flex items-center justify-between px-5 pb-4">
          <span class="font-mono text-[13px] uppercase tracking-[0.14em] text-gray-900 dark:text-gray-100">{{ siteName }}</span>
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
            class="mobile-nav-link"
          >{{ t('home.nav.docs') }}</a>
          <component
            :is="isStatusActive ? 'span' : RouterLink"
            :to="isStatusActive ? undefined : '/status'"
            class="mobile-nav-link"
            :class="isStatusActive ? 'mobile-nav-link-active cursor-default' : ''"
            :aria-current="isStatusActive ? 'page' : undefined"
            @click="handleMobileNav(isStatusActive, $event)"
          >{{ t('home.nav.status') }}</component>
          <component
            :is="isModelsActive ? 'span' : RouterLink"
            :to="isModelsActive ? undefined : '/models'"
            class="mobile-nav-link"
            :class="isModelsActive ? 'mobile-nav-link-active cursor-default' : ''"
            :aria-current="isModelsActive ? 'page' : undefined"
            @click="handleMobileNav(isModelsActive, $event)"
          >{{ t('home.nav.models') }}</component>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="mt-2 inline-flex items-center justify-center rounded-sm border border-dark-800 bg-dark-950 px-3 py-2 font-mono text-[12px] uppercase tracking-[0.12em] text-white dark:border-dark-700 dark:bg-dark-900"
            @click="mobileOpen = false"
          >{{ isAuthenticated ? t('home.dashboard') : t('home.login') }}</router-link>
        </nav>
      </aside>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import { useBrandHover } from '@/composables/useBrandHover'
import { useTheme } from '@/composables/useTheme'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()
const route = useRoute()
const authStore = useAuthStore()
const appStore = useAppStore()
const { brandClass, onEnter, onLeave } = useBrandHover()
const { isDark, toggleTheme } = useTheme()

const mobileOpen = ref(false)
const siteName = computed(() => appStore.cachedPublicSettings?.site_name?.trim() || appStore.siteName.trim() || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))
const userInitial = computed(() => authStore.user?.email?.charAt(0).toUpperCase() || '')
const isStatusActive = computed(() => route.path === '/status')
const isModelsActive = computed(() => route.path === '/models')

function handleMobileNav(isActive: boolean, event: Event) {
  if (isActive) {
    event.preventDefault()
    return
  }
  mobileOpen.value = false
}
</script>

<style scoped>
.nav-link {
  @apply font-mono text-[12px] uppercase tracking-[0.14em] text-gray-700 transition-colors duration-200 hover:text-primary-500 dark:text-gray-200;
}

.nav-link-active {
  @apply text-primary-600 dark:text-primary-400;
}

.nav-hover-group:hover .nav-link:not(:hover):not(.nav-link-active) {
  opacity: 0.42;
}

.mobile-nav-link {
  @apply rounded-sm px-3 py-2 font-mono text-[12px] uppercase tracking-[0.12em] text-gray-600 transition-colors duration-150 hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white;
}

.mobile-nav-link-active {
  @apply bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-300;
}

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
  transition: transform 0.22s ease-in-out;
}

.sf-drawer-enter-from .sf-drawer-panel,
.sf-drawer-leave-to .sf-drawer-panel {
  transform: translateX(-100%);
}
</style>
