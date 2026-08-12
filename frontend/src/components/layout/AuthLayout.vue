<template>
  <div class="relative flex min-h-screen items-center justify-center overflow-hidden bg-gray-50 p-4 dark:bg-dark-950">
    <!-- Background: warm near-black / warm off-white base -->
    <div class="absolute inset-0 bg-white dark:bg-dark-950"></div>

    <!-- Decorative Elements: blueprint line-grid (factory motif) -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <!-- Blueprint line-grid signature -->
      <div class="absolute inset-0 bg-blueprint opacity-60"></div>

      <!-- Hairline frame inset -->
      <div class="absolute inset-4 border border-gray-200/60 dark:border-dark-700/60"></div>
    </div>

    <div class="absolute right-5 top-5 z-20 sm:right-6 sm:top-6">
      <LocaleSwitcher />
    </div>

    <!-- Content Container -->
    <div class="relative z-10 w-full max-w-md">
      <!-- Logo/Brand -->
      <div
        class="brand-hover mb-8 text-center"
        :class="brandClass"
        :aria-label="siteName"
        @mouseenter="onEnter"
        @mouseleave="onLeave"
        @focusin="onEnter"
        @focusout="onLeave"
      >
        <!-- Custom Logo or Default Logo -->
        <template v-if="settingsLoaded">
          <div
            class="brand-hover__logo mb-4 inline-flex h-14 w-14 items-center justify-center overflow-hidden rounded-md border border-gray-200 bg-white p-1.5 dark:border-dark-700 dark:bg-dark-900"
          >
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <h1
            class="mb-2 text-2xl font-normal tracking-tight text-gray-900 dark:text-gray-100"
          >
            <span class="sr-only">{{ siteName }}</span>
            <span aria-hidden="true">
              <span
                v-for="(ch, i) in siteName"
                :key="i"
                class="brand-hover__letter"
                :style="{ '--brand-i': i }"
                >{{ ch }}</span
              >
            </span>
          </h1>
          <p class="text-sm text-gray-500 dark:text-dark-400">
            {{ siteSubtitle }}
          </p>
        </template>
      </div>

      <!-- Card Container -->
      <div class="rounded-lg border border-gray-200 bg-white p-8 dark:border-dark-700 dark:bg-dark-900">
        <slot />
      </div>

      <!-- Footer Links -->
      <div class="mt-6 text-center text-sm">
        <slot name="footer" />
      </div>

      <!-- Copyright -->
      <div class="mt-8 text-center font-mono text-xs text-gray-400 dark:text-dark-500">
        &copy; {{ currentYear }} {{ siteName }}. All rights reserved.
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAppStore } from '@/stores'
import { sanitizeUrl } from '@/utils/url'
import { useBrandHover } from '@/composables/useBrandHover'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'

const appStore = useAppStore()

// Brand micro-interaction (shared phase state — see useBrandHover).
// Destructure so each ref/handler is a top-level setup binding — Vue only
// auto-unwraps refs at the top level of the bindings object, not through
// nested property access. Binding `brand.brandClass` directly would pass the
// ComputedRef object to `:class` instead of its value (animation never fires).
const { brandClass, onEnter, onLeave } = useBrandHover()

const siteName = computed(() => appStore.siteName.trim() || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'Subscription to API Conversion Platform')
const settingsLoaded = computed(() => appStore.publicSettingsLoaded)

const currentYear = computed(() => new Date().getFullYear())

onMounted(() => {
  appStore.fetchPublicSettings()
})
</script>
