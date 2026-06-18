<template>
  <div class="relative flex min-h-screen items-center justify-center overflow-hidden bg-gray-50 p-4 dark:bg-dark-950">
    <!-- Background: warm near-black / warm off-white base -->
    <div class="absolute inset-0 bg-white dark:bg-dark-950"></div>

    <!-- Decorative Elements: warm mesh + blueprint line-grid (factory motif) -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <!-- Warm ambient mesh (vermillion, low opacity) -->
      <div class="absolute inset-0 bg-mesh-gradient opacity-80"></div>

      <!-- Blueprint line-grid signature -->
      <div class="absolute inset-0 bg-blueprint opacity-60"></div>

      <!-- Hairline frame inset -->
      <div class="absolute inset-4 border border-gray-200/60 dark:border-dark-700/60"></div>
    </div>

    <!-- Content Container -->
    <div class="relative z-10 w-full max-w-md">
      <!-- Logo/Brand -->
      <div class="mb-8 text-center">
        <!-- Custom Logo or Default Logo -->
        <template v-if="settingsLoaded">
          <div
            class="mb-4 inline-flex h-14 w-14 items-center justify-center overflow-hidden rounded-md border border-gray-200 bg-white p-1.5 dark:border-dark-700 dark:bg-dark-900"
          >
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <div class="eyebrow mb-2">SiliconBase</div>
          <h1 class="mb-2 text-2xl font-normal tracking-tight text-gray-900 dark:text-gray-100">
            {{ siteName }}
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

const appStore = useAppStore()

const siteName = computed(() => appStore.siteName || 'SiliconBase')
const siteLogo = computed(() => sanitizeUrl(appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'Subscription to API Conversion Platform')
const settingsLoaded = computed(() => appStore.publicSettingsLoaded)

const currentYear = computed(() => new Date().getFullYear())

onMounted(() => {
  appStore.fetchPublicSettings()
})
</script>
