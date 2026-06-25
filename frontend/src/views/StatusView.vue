<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-gray-100">
    <PublicTopNav />

    <main class="mx-auto max-w-[1440px] px-4 py-10 lg:px-8">
      <section class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_360px] lg:items-end">
        <div>
          <p class="font-mono text-[11px] uppercase tracking-[0.18em] text-primary-500">
            {{ t('channelStatus.publicEyebrow') }}
          </p>
          <h1 class="mt-4 text-4xl font-semibold tracking-[-0.03em] text-gray-950 dark:text-white sm:text-5xl">
            {{ t('channelStatus.publicTitle') }}
          </h1>
          <p class="mt-4 max-w-[64ch] text-sm leading-6 text-gray-600 dark:text-dark-300">
            {{ t('channelStatus.publicDescription') }}
          </p>
        </div>

        <div class="rounded-md border border-gray-200 bg-white p-4 dark:border-dark-800 dark:bg-dark-900">
          <div class="flex items-center justify-between gap-3">
            <div>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('channelStatus.overallLabel') }}</p>
              <p class="mt-1 text-lg font-semibold" :class="overallStatus === 'operational' ? 'text-emerald-600 dark:text-emerald-300' : 'text-amber-600 dark:text-amber-300'">
                {{ t(`channelStatus.overall.${overallStatus}`) }}
              </p>
            </div>
            <span class="flex h-11 w-11 items-center justify-center rounded-full" :class="overallStatus === 'operational' ? 'bg-emerald-100 text-emerald-600 dark:bg-emerald-500/15 dark:text-emerald-300' : 'bg-amber-100 text-amber-600 dark:bg-amber-500/15 dark:text-amber-300'">
              <Icon :name="overallStatus === 'operational' ? 'checkCircle' : 'exclamationTriangle'" size="lg" />
            </span>
          </div>
          <div class="mt-4 grid grid-cols-3 gap-2 text-sm">
            <div class="rounded-sm bg-gray-50 p-3 dark:bg-dark-800">
              <p class="font-mono text-lg font-semibold">{{ items.length }}</p>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('channelStatus.monitorCount') }}</p>
            </div>
            <div class="rounded-sm bg-gray-50 p-3 dark:bg-dark-800">
              <p class="font-mono text-lg font-semibold">{{ providerGroups.length }}</p>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('channelStatus.platformCount') }}</p>
            </div>
            <div class="rounded-sm bg-gray-50 p-3 dark:bg-dark-800">
              <p class="font-mono text-lg font-semibold">{{ degradedCount }}</p>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('channelStatus.degradedCount') }}</p>
            </div>
          </div>
        </div>
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
          <Select
            v-model="providerFilter"
            :options="providerSelectOptions"
            class="min-w-[180px]"
            searchable
          />
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

      <section v-if="loading && items.length === 0" class="mt-6 grid gap-4 lg:grid-cols-2">
        <div
          v-for="n in 4"
          :key="n"
          class="h-56 animate-pulse rounded-md border border-gray-200 bg-white dark:border-dark-800 dark:bg-dark-900"
        ></div>
      </section>

      <section v-else-if="filteredItems.length === 0" class="mt-6 rounded-md border border-dashed border-gray-300 bg-white p-10 text-center dark:border-dark-700 dark:bg-dark-900">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('channelStatus.empty.title') }}</h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">{{ t('channelStatus.empty.description') }}</p>
      </section>

      <section v-else class="mt-6 space-y-8">
        <div
          v-for="providerGroup in providerGroups"
          :key="providerGroup.provider"
          class="rounded-md border border-gray-200 bg-white dark:border-dark-800 dark:bg-dark-900"
        >
          <div class="flex flex-col gap-3 border-b border-gray-200 p-5 dark:border-dark-800 sm:flex-row sm:items-center sm:justify-between">
            <div class="flex items-center gap-3">
              <span class="flex h-10 w-10 items-center justify-center rounded-sm font-mono text-xs font-semibold uppercase" :class="providerBadgeClass(providerGroup.provider)">
                {{ providerInitial(providerGroup.provider) }}
              </span>
              <div>
                <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('channelStatus.platform') }}</p>
                <h2 class="text-xl font-semibold text-gray-950 dark:text-white">{{ providerLabel(providerGroup.provider) }}</h2>
              </div>
            </div>
            <p class="font-mono text-xs uppercase tracking-[0.12em] text-gray-500 dark:text-dark-400">
              {{ t('channelStatus.groupSummary', { groups: providerGroup.groups.length, monitors: providerGroup.total }) }}
            </p>
          </div>

          <div class="divide-y divide-gray-200 dark:divide-dark-800">
            <div
              v-for="group in providerGroup.groups"
              :key="`${providerGroup.provider}-${group.name}`"
              class="p-5"
            >
              <div class="mb-4 flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('channelStatus.group') }}</p>
                  <h3 class="text-base font-semibold text-gray-900 dark:text-white">{{ group.name || t('channelStatus.defaultGroup') }}</h3>
                </div>
                <span class="w-fit rounded-full bg-gray-100 px-2.5 py-1 text-xs text-gray-600 dark:bg-dark-800 dark:text-dark-300">
                  {{ group.items.length }} {{ t('channelStatus.monitorsUnit') }}
                </span>
              </div>

              <div class="grid gap-4 lg:grid-cols-2">
                <article
                  v-for="item in group.items"
                  :key="item.id"
                  class="rounded-md border border-gray-200 bg-white p-5 transition-colors hover:border-gray-300 dark:border-dark-700 dark:bg-dark-950/60 dark:hover:border-primary-500/50"
                >
                  <div class="flex items-start justify-between gap-3">
                    <div class="flex min-w-0 items-start gap-3">
                      <span
                        class="grid h-9 w-9 shrink-0 place-items-center rounded-md ring-1 ring-black/5 dark:ring-white/10"
                        :class="[providerGradient(item.provider), providerTintClass(item.provider)]"
                      >
                        <ProviderIcon :provider="item.provider" :size="20" />
                      </span>
                      <div class="min-w-0">
                        <h4 class="truncate text-base font-normal tracking-tight text-gray-950 dark:text-white">{{ item.name }}</h4>
                        <div class="mt-0.5 flex min-w-0 items-center gap-1.5">
                          <span
                            class="inline-flex shrink-0 items-center rounded-sm px-1.5 py-0.5 font-mono text-[10px] tracking-wide"
                            :class="providerBadgeClass(item.provider)"
                          >
                            {{ providerLabel(item.provider) }}
                          </span>
                          <span class="truncate font-mono text-xs text-gray-500 dark:text-gray-400">{{ item.primary_model }}</span>
                        </div>
                      </div>
                    </div>
                    <span class="shrink-0 rounded-sm px-2.5 py-1 font-mono text-xs tracking-wide" :class="statusBadgeClass(item.primary_status)">
                      {{ statusLabel(item.primary_status) }}
                    </span>
                  </div>

                  <MonitorMetricPair
                    primary-icon="bolt"
                    :primary-label="t('monitorCommon.dialogLatency')"
                    :primary-value="formatLatency(item.primary_latency_ms)"
                    primary-unit="ms"
                    secondary-icon="globe"
                    :secondary-label="t('monitorCommon.endpointPing')"
                    :secondary-value="formatLatency(item.primary_ping_latency_ms)"
                    secondary-unit="ms"
                  />

                  <div class="mt-4 border-t border-dashed border-gray-200 dark:border-dark-700"></div>

                  <MonitorAvailabilityRow
                    :window-label="`${t('monitorCommon.availabilityPrefix')} · ${t('channelStatus.windowTab.7d')}`"
                    :value="item.availability_7d"
                    :samples-label="item.extra_models.length ? t('monitorCommon.extraModelsCount', { n: item.extra_models.length }) : undefined"
                  />

                  <MonitorTimeline
                    :buckets="item.timeline"
                    :countdown-seconds="0"
                    :show-countdown="false"
                  />

                  <div v-if="item.extra_models.length" class="mt-4 flex flex-wrap gap-2">
                    <span
                      v-for="extra in item.extra_models.slice(0, 4)"
                      :key="extra.model"
                      class="inline-flex max-w-full items-center gap-1 rounded-sm bg-white px-2 py-1 text-xs text-gray-600 dark:bg-dark-900 dark:text-dark-300"
                    >
                      <span class="h-1.5 w-1.5 shrink-0 rounded-full" :class="timelineClass(extra.status)"></span>
                      <span class="truncate">{{ extra.model }}</span>
                    </span>
                    <span v-if="item.extra_models.length > 4" class="rounded-sm bg-white px-2 py-1 text-xs text-gray-500 dark:bg-dark-900 dark:text-dark-400">
                      +{{ item.extra_models.length - 4 }}
                    </span>
                  </div>
                </article>
              </div>
            </div>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore, useAuthStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import Select from '@/components/common/Select.vue'
import PublicTopNav from '@/components/home/PublicTopNav.vue'
import ProviderIcon from '@/components/user/monitor/ProviderIcon.vue'
import MonitorMetricPair from '@/components/user/monitor/MonitorMetricPair.vue'
import MonitorAvailabilityRow from '@/components/user/monitor/MonitorAvailabilityRow.vue'
import MonitorTimeline from '@/components/user/monitor/MonitorTimeline.vue'
import { useTheme } from '@/composables/useTheme'
import { useChannelMonitorFormat } from '@/composables/useChannelMonitorFormat'
import { providerGradient } from '@/composables/useChannelMonitorFormat'
import { STATUS_OPERATIONAL } from '@/constants/channelMonitor'
import { list as listPublicChannelMonitorViews } from '@/api/publicChannelMonitor'
import type { UserMonitorView } from '@/api/channelMonitor'
import { extractApiErrorMessage } from '@/utils/apiError'

type OverallStatus = 'operational' | 'degraded'
type Provider = UserMonitorView['provider'] | string

interface MonitorGroup {
  name: string
  items: UserMonitorView[]
}

interface ProviderGroup {
  provider: Provider
  total: number
  groups: MonitorGroup[]
}

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const { syncThemeFromDocument } = useTheme()
const {
  statusLabel,
  statusBadgeClass,
  providerLabel,
  providerBadgeClass,
  formatLatency,
} = useChannelMonitorFormat()

const PROVIDER_TINT: Record<string, string> = {
  openai: 'text-emerald-600 dark:text-emerald-300',
  anthropic: 'text-orange-600 dark:text-orange-300',
  gemini: 'text-sky-600 dark:text-sky-300',
}

const items = ref<UserMonitorView[]>([])
const loading = ref(false)
const search = ref('')
const providerFilter = ref('')
let abortController: AbortController | null = null

const providerOptions = computed<string[]>(() => {
  return Array.from(new Set(items.value.map(item => String(item.provider || '')).filter(Boolean))).sort()
})
const providerSelectOptions = computed(() => [
  { value: '', label: t('channelStatus.allProviders') },
  ...providerOptions.value.map(provider => ({ value: provider, label: providerLabel(provider) })),
])

const filteredItems = computed(() => {
  const q = search.value.trim().toLowerCase()
  return items.value.filter((item) => {
    if (providerFilter.value && item.provider !== providerFilter.value) return false
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

const providerGroups = computed<ProviderGroup[]>(() => {
  const providerMap = new Map<string, Map<string, UserMonitorView[]>>()
  for (const item of filteredItems.value) {
    const provider = item.provider || 'unknown'
    const groupName = item.group_name || t('channelStatus.defaultGroup')
    if (!providerMap.has(provider)) providerMap.set(provider, new Map())
    const groupMap = providerMap.get(provider)!
    if (!groupMap.has(groupName)) groupMap.set(groupName, [])
    groupMap.get(groupName)!.push(item)
  }
  return Array.from(providerMap.entries())
    .sort(([a], [b]) => providerLabel(a).localeCompare(providerLabel(b)))
    .map(([provider, groupMap]) => ({
      provider,
      total: Array.from(groupMap.values()).reduce((sum, groupItems) => sum + groupItems.length, 0),
      groups: Array.from(groupMap.entries())
        .sort(([a], [b]) => a.localeCompare(b))
        .map(([name, groupItems]) => ({
          name,
          items: groupItems.slice().sort((a, b) => a.name.localeCompare(b.name)),
        })),
    }))
})

const degradedCount = computed(() => items.value.filter(item => item.primary_status !== STATUS_OPERATIONAL).length)

const overallStatus = computed<OverallStatus>(() => {
  if (items.value.length === 0) return 'operational'
  return degradedCount.value > 0 ? 'degraded' : 'operational'
})

watch(providerOptions, (options) => {
  if (providerFilter.value && !options.includes(providerFilter.value)) {
    providerFilter.value = ''
  }
})

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
      abortController = null
    }
  }
}

function manualReload() {
  void reload(false)
}

function providerInitial(provider: Provider) {
  const label = providerLabel(provider)
  if (!label || label === '-') return 'AI'
  return label
    .split(/\s+/)
    .map(part => part.charAt(0))
    .join('')
    .slice(0, 2)
    .toUpperCase()
}

function providerTintClass(provider: string) {
  return PROVIDER_TINT[provider] ?? 'text-gray-500 dark:text-gray-300'
}

function timelineClass(status: string) {
  switch (status) {
    case 'operational':
      return 'bg-emerald-500'
    case 'degraded':
      return 'bg-amber-500'
    case 'failed':
    case 'error':
      return 'bg-red-500'
    default:
      return 'bg-gray-300 dark:bg-dark-700'
  }
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
