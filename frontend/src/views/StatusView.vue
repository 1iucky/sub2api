<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-gray-100">
    <PublicTopNav />

    <main class="mx-auto max-w-[1240px] px-4 py-8 lg:px-8">
      <section>
        <p class="font-mono text-[11px] uppercase tracking-[0.18em] text-primary-500">
          {{ t('channelStatus.publicEyebrow') }}
        </p>
        <h1 class="mt-3 text-3xl font-semibold tracking-[-0.03em] text-gray-950 dark:text-white sm:text-4xl">
          {{ t('channelStatus.publicTitle') }}
        </h1>
        <p class="mt-3 max-w-[64ch] text-sm leading-6 text-gray-600 dark:text-dark-300">
          {{ t('channelStatus.publicDescription') }}
        </p>
      </section>

      <section class="mt-8 flex flex-col gap-3 rounded-md border border-gray-200 bg-white p-3 dark:border-dark-800 dark:bg-dark-900 md:flex-row md:items-center md:justify-between">
        <div class="flex flex-wrap items-center gap-2">
          <label class="relative block">
            <Icon name="search" size="sm" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              v-model.trim="search"
              type="search"
              class="input h-9 min-w-[240px] pl-9"
              :placeholder="t('channelStatus.searchPlaceholder')"
            />
          </label>
        </div>

        <div class="flex flex-wrap items-center gap-2">
          <button
            class="btn-secondary btn-sm h-9 font-mono uppercase tracking-[0.12em]"
            :disabled="loading"
            @click="manualReload"
          >
            <Icon name="refresh" size="sm" :class="loading ? 'animate-spin' : ''" />
            {{ t('common.refresh') }}
          </button>
        </div>
      </section>

      <section v-if="initialLoading" class="mt-5 rounded-md border border-gray-200 bg-white dark:border-dark-800 dark:bg-dark-900">
        <div
          v-for="n in 4"
          :key="n"
          class="h-36 animate-pulse border-b border-gray-200 last:border-b-0 dark:border-dark-800"
        ></div>
      </section>

      <section v-else-if="filteredItems.length === 0" class="mt-6 rounded-md border border-dashed border-gray-300 bg-white p-10 text-center dark:border-dark-700 dark:bg-dark-900">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('channelStatus.empty.title') }}</h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">{{ t('channelStatus.empty.description') }}</p>
      </section>

      <section v-else class="relative mt-5">
        <div
          v-if="loading"
          class="absolute inset-0 z-10 grid place-items-center rounded-md bg-gray-50/70 backdrop-blur-sm dark:bg-dark-950/70"
        >
          <span class="inline-flex items-center gap-2 rounded-full border border-gray-200 bg-white px-4 py-2 text-sm text-gray-600 shadow-sm dark:border-dark-700 dark:bg-dark-900 dark:text-dark-300">
            <Icon name="refresh" size="sm" class="animate-spin" />
            {{ t('common.loading') }}
          </span>
        </div>
        <div class="overflow-hidden rounded-md border border-gray-200 bg-white dark:border-dark-800 dark:bg-dark-900">
          <div class="divide-y divide-gray-200 dark:divide-dark-800">
            <article
              v-for="item in filteredItems"
              :key="item.id"
              data-testid="status-monitor-card"
              class="px-5 py-4"
            >
              <div class="flex items-start justify-between gap-4">
                <div class="min-w-0">
                  <h3 class="truncate text-lg font-semibold tracking-[-0.01em] text-gray-950 dark:text-white">{{ item.name }}</h3>
                  <p class="mt-0.5 truncate font-mono text-xs text-gray-500 dark:text-dark-400">{{ item.primary_model }}</p>
                </div>
                <span class="shrink-0 text-sm font-medium sm:text-base" :class="statusTextClass(item.primary_status)">
                  {{ statusDisplayLabel(item.primary_status) }}
                </span>
              </div>

              <div class="mt-3 flex h-9 items-end gap-[4px]">
                <span
                  v-for="(bar, idx) in uptimeBars(item)"
                  :key="idx"
                  data-testid="status-uptime-bar"
                  class="min-w-[3px] flex-1 rounded-[1px]"
                  :class="bar.colorClass"
                  :style="{ height: bar.heightPct + '%' }"
                  :title="bar.title"
                ></span>
              </div>

              <div class="mt-3 grid grid-cols-[auto_minmax(24px,1fr)_auto_minmax(24px,1fr)_auto] items-center gap-3 text-sm font-medium text-gray-500 dark:text-dark-400">
                <span>{{ t('channelStatus.windowAgo.7d') }}</span>
                <span class="h-px bg-gray-300 dark:bg-dark-700"></span>
                <span class="whitespace-nowrap tabular-nums">{{ t('channelStatus.uptimeValue', { value: formatAvailability(item.availability_7d) }) }}</span>
                <span class="h-px bg-gray-300 dark:bg-dark-700"></span>
                <span>{{ t('channelStatus.today') }}</span>
              </div>
            </article>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore, useAuthStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import PublicTopNav from '@/components/home/PublicTopNav.vue'
import { useTheme } from '@/composables/useTheme'
import { useChannelMonitorFormat } from '@/composables/useChannelMonitorFormat'
import { list as listPublicChannelMonitorViews } from '@/api/publicChannelMonitor'
import type { MonitorTimelinePoint, UserMonitorView } from '@/api/channelMonitor'
import { extractApiErrorMessage } from '@/utils/apiError'

interface UptimeBar {
  colorClass: string
  heightPct: number
  title: string
}

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const { syncThemeFromDocument } = useTheme()
const {
  statusLabel,
  formatLatency,
  formatRelativeTime,
} = useChannelMonitorFormat()

const STATUS_HEIGHT: Record<string, number> = {
  operational: 100,
  degraded: 72,
  failed: 48,
  error: 48,
  empty: 100,
}

const STATUS_COLOR: Record<string, string> = {
  operational: 'bg-[#67aa2b]',
  degraded: 'bg-[#dda51f]',
  failed: 'bg-[#ef4b3f]',
  error: 'bg-[#ef4b3f]',
  empty: 'bg-gray-300 dark:bg-dark-700',
}

const items = ref<UserMonitorView[]>([])
const loading = ref(true)
const hasLoaded = ref(false)
const search = ref('')
let abortController: AbortController | null = null

const filteredItems = computed(() => {
  const q = search.value.trim().toLowerCase()
  return items.value.filter((item) => {
    if (!q) return true
    return [
      item.name,
      item.provider,
      item.group_name,
      item.primary_model,
      ...item.extra_models.map(extra => extra.model),
    ].some(value => String(value || '').toLowerCase().includes(q))
  })
})

const initialLoading = computed(() => loading.value && !hasLoaded.value && items.value.length === 0)

async function reload(silent = false) {
  if (abortController) abortController.abort()
  const ctrl = new AbortController()
  abortController = ctrl
  if (!silent) loading.value = true
  try {
    const res = await listPublicChannelMonitorViews({ signal: ctrl.signal })
    if (ctrl.signal.aborted || abortController !== ctrl) return
    items.value = res.items || []
  } catch (err: unknown) {
    const e = err as { name?: string; code?: string }
    if (e?.name === 'AbortError' || e?.code === 'ERR_CANCELED') return
    appStore.showError(extractApiErrorMessage(err, t('channelStatus.loadError')))
  } finally {
    if (abortController === ctrl) {
      if (!silent) loading.value = false
      hasLoaded.value = true
      abortController = null
    }
  }
}

function manualReload() {
  void reload(false)
}

function statusTextClass(status: string) {
  switch (status) {
    case 'operational':
      return 'text-[#5fac24] dark:text-emerald-300'
    case 'degraded':
      return 'text-amber-600 dark:text-amber-300'
    case 'failed':
    case 'error':
      return 'text-red-600 dark:text-red-300'
    default:
      return 'text-gray-500 dark:text-dark-400'
  }
}

function statusDisplayLabel(status: string) {
  if (status === 'operational') return t('channelStatus.status.operational')
  if (status === 'degraded') return t('channelStatus.status.degraded')
  if (status === 'failed') return t('channelStatus.status.failed')
  if (status === 'error') return t('channelStatus.status.error')
  return statusLabel(status as MonitorTimelinePoint['status'])
}

function uptimeBars(item: UserMonitorView): UptimeBar[] {
  const real = [...(item.timeline || [])]
    .reverse()
    .map(point => uptimeBarFromPoint(point))

  if (real.length > 0) return real
  return [{
    colorClass: STATUS_COLOR.empty,
    heightPct: STATUS_HEIGHT.empty,
    title: '',
  }]
}

function uptimeBarFromPoint(point: MonitorTimelinePoint): UptimeBar {
  const status = point.status || 'empty'
  const relative = point.checked_at ? formatRelativeTime(point.checked_at) : ''
  const label = statusDisplayLabel(point.status)
  const latency = formatLatency(point.latency_ms)
  return {
    colorClass: STATUS_COLOR[status] ?? STATUS_COLOR.empty,
    heightPct: STATUS_HEIGHT[status] ?? STATUS_HEIGHT.empty,
    title: [relative, label, latency ? `${latency}ms` : ''].filter(Boolean).join(' · '),
  }
}

function formatAvailability(value: number) {
  return Number(value || 0).toFixed(2)
}

onMounted(() => {
  syncThemeFromDocument()
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
  void reload(false)
})

onBeforeUnmount(() => {
  if (abortController) abortController.abort()
})
</script>
