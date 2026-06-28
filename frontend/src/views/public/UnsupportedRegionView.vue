<template>
  <main class="flex min-h-screen items-center justify-center bg-[#f3f6fb] px-4 py-10 text-gray-900 dark:bg-dark-950 dark:text-gray-100">
    <section class="w-full max-w-[420px] rounded-lg border border-gray-200 bg-white px-8 py-9 text-center shadow-2xl shadow-blue-950/10 dark:border-dark-800 dark:bg-dark-900">
      <div class="flex items-center justify-center gap-3">
        <span class="flex h-12 w-12 items-center justify-center overflow-hidden rounded-xl bg-primary-500 text-base font-bold text-white shadow-lg shadow-primary-500/30">
          <img :src="siteLogo" alt="" class="h-full w-full object-contain" />
        </span>
        <div class="text-left">
          <p class="text-xl font-semibold tracking-tight text-gray-950 dark:text-white">{{ siteName }}</p>
        </div>
      </div>

      <div class="mx-auto mt-7 flex h-20 w-20 items-center justify-center rounded-full border-[6px] border-red-500 text-red-500">
        <Icon name="ban" size="xl" :stroke-width="2.5" />
      </div>

      <h1 class="mt-7 text-2xl font-semibold tracking-tight text-gray-950 dark:text-white">
        {{ t('unsupportedRegion.title') }}
      </h1>
      <p class="mt-4 text-sm leading-7 text-gray-600 dark:text-dark-300">
        {{ t('unsupportedRegion.description', { site: siteName }) }}
      </p>

      <p class="mt-6 break-words font-mono text-xs leading-6 text-gray-500 dark:text-dark-400">
        {{ t('unsupportedRegion.accessInfo', { host, location: locationLabel, path: originalPath }) }}
      </p>

      <button
        type="button"
        class="mt-7 inline-flex h-11 w-full items-center justify-center rounded-md bg-primary-500 px-4 text-sm font-semibold text-white shadow-lg shadow-primary-500/25 transition hover:bg-primary-600 focus:outline-none focus:ring-2 focus:ring-primary-500/30"
        @click="retryCheck"
      >
        {{ t('unsupportedRegion.retry') }}
      </button>
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores'
import { useTheme } from '@/composables/useTheme'
import Icon from '@/components/icons/Icon.vue'
import type { PublicSettings } from '@/types'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const { syncThemeFromDocument } = useTheme()
const loadedSettings = ref<PublicSettings | null>(appStore.cachedPublicSettings)

const siteName = computed(() => loadedSettings.value?.site_name || appStore.cachedPublicSettings?.site_name || appStore.siteName || 'SiliconBase')
const siteLogo = computed(() => loadedSettings.value?.site_logo || appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '/logo.svg')
const host = computed(() => window.location.host || '-')
const sourceIP = computed(() => queryText('ip') || '-')
const country = computed(() => queryText('country') || t('unsupportedRegion.unknownLocation'))
const region = computed(() => queryText('region'))
const originalPath = computed(() => queryText('path') || '/')
const locationLabel = computed(() => {
  const parts = [country.value, region.value].filter(Boolean)
  return `${parts.join(' / ')} · ${sourceIP.value}`
})
function queryText(key: string): string {
  const value = route.query[key]
  if (Array.isArray(value)) return String(value[0] || '').trim()
  return String(value || '').trim()
}

async function retryCheck() {
  try {
    const response = await fetch(`/region-check?ts=${Date.now()}`, {
      method: 'GET',
      cache: 'no-store',
    })

    if (response.status === 204) {
      void router.replace('/home')
    }
  } catch {
    // Keep the user on the unsupported-region page if the edge check fails.
  }
}

onMounted(() => {
  syncThemeFromDocument()
  if (!appStore.publicSettingsLoaded) {
    void appStore.fetchPublicSettings().then((settings) => {
      if (settings) {
        loadedSettings.value = settings
      }
    })
  }
})
</script>
