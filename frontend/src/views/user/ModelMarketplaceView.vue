<template>
  <component :is="embedded ? AppLayout : 'div'">
    <div :class="pageShellClass">
      <PublicTopNav v-if="!embedded" />

      <main :class="mainClass">
      <section class="border-b border-stone-300 pb-10 text-center dark:border-stone-800">
        <h1 class="font-mono text-[clamp(32px,6vw,64px)] font-semibold uppercase leading-none tracking-[0.14em] text-stone-950 dark:text-stone-50">
          {{ t('models.marketplaceTitle') }}
        </h1>
        <p class="mx-auto mt-5 max-w-[48ch] text-sm leading-6 text-stone-600 dark:text-stone-400">
          {{ t('models.marketplaceDescription') }}
        </p>
      </section>

      <section class="grid gap-6 border-b border-stone-300 py-6 dark:border-stone-800 lg:grid-cols-[280px_minmax(0,1fr)]">
        <aside class="lg:sticky lg:top-20 lg:self-start">
          <div class="rounded-md border border-stone-300 bg-white/60 p-4 dark:border-stone-800 dark:bg-[#171410]">
            <div class="flex items-center justify-between">
              <h2 class="font-mono text-[13px] uppercase tracking-[0.14em] text-stone-900 dark:text-stone-100">
                {{ t('models.filters') }}
              </h2>
              <button
                class="font-mono text-[11px] uppercase tracking-[0.12em] text-stone-500 transition-colors hover:text-primary-600 dark:text-stone-400 dark:hover:text-primary-400"
                @click="resetFilters"
              >
                {{ t('models.resetFilters') }}
              </button>
            </div>

            <div class="mt-6 border-t border-stone-200 pt-5 dark:border-stone-800">
              <p class="mb-3 text-xs text-stone-500 dark:text-stone-400">{{ t('models.vendor') }}</p>
              <label class="mb-2 block">
                <span class="relative block">
                  <Icon name="search" size="xs" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-stone-400" />
                  <input
                    v-model.trim="vendorSearch"
                    type="search"
                    class="market-input h-8 w-full pl-8 pr-3 text-xs"
                    :placeholder="t('models.search')"
                  />
                </span>
              </label>
              <div class="max-h-80 overflow-auto pr-1">
                <button
                  type="button"
                  class="market-filter-option"
                  :class="filters.vendor_ids.length === 0 ? 'market-filter-option-active' : ''"
                  @click="clearVendorFilters"
                >
                  <span class="market-filter-check">
                    <Icon v-if="filters.vendor_ids.length === 0" name="check" size="xs" />
                  </span>
                  <span>{{ t('models.allVendors') }}</span>
                </button>
                <button
                  v-for="vendor in filteredVendors"
                  :key="vendor.id"
                  type="button"
                  class="market-filter-option"
                  :class="isVendorSelected(vendor.id) ? 'market-filter-option-active' : ''"
                  @click="toggleVendorFilter(vendor.id)"
                >
                  <span class="market-filter-check">
                    <Icon v-if="isVendorSelected(vendor.id)" name="check" size="xs" />
                  </span>
                  <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-sm bg-white dark:bg-stone-950">
                    <ModelIcon :model="vendor.provider_key || vendor.name" :icon-key="vendor.icon_key" size="14px" />
                  </span>
                  <span class="truncate">{{ vendor.name }}</span>
                </button>
              </div>
            </div>

            <div class="mt-6 border-t border-stone-200 pt-5 dark:border-stone-800">
              <p class="mb-3 text-xs text-stone-500 dark:text-stone-400">{{ t('models.capabilities') }}</p>
              <div>
                <button
                  v-for="capability in capabilityFilters"
                  :key="capability.key"
                  type="button"
                  class="market-filter-option"
                  :class="filters.capabilities.includes(capability.key) ? 'market-filter-option-active' : ''"
                  @click="toggleCapabilityFilter(capability.key)"
                >
                  <span class="market-filter-check">
                    <Icon v-if="filters.capabilities.includes(capability.key)" name="check" size="xs" />
                  </span>
                  <Icon :name="capability.icon" size="xs" class="text-stone-500 dark:text-stone-400" />
                  <span>{{ t(capability.labelKey) }}</span>
                </button>
              </div>
            </div>

            <div class="mt-6 border-t border-stone-200 pt-5 dark:border-stone-800">
              <p class="mb-3 text-xs text-stone-500 dark:text-stone-400">{{ t('models.priceRange') }}</p>
              <div class="grid grid-cols-2 gap-1.5">
                <button
                  v-for="band in priceBands"
                  :key="band.value"
                  type="button"
                  class="rounded-sm border px-2.5 py-2 text-left text-xs transition-colors"
                  :class="filters.priceBand === band.value
                    ? 'border-primary-500 bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-300'
                    : 'border-stone-300 bg-white/70 text-stone-700 hover:border-primary-400 dark:border-stone-800 dark:bg-[#100e0b] dark:text-stone-300'"
                  @click="filters.priceBand = band.value"
                >
                  {{ band.label }}
                </button>
              </div>
            </div>

            <div class="mt-6 border-t border-stone-200 pt-5 dark:border-stone-800">
              <p class="mb-3 text-xs text-stone-500 dark:text-stone-400">{{ t('models.contextRange') }}</p>
              <div>
                <button
                  v-for="band in contextBands"
                  :key="band.value"
                  type="button"
                  class="w-full rounded-sm border px-2.5 py-2 text-left text-xs transition-colors"
                  :class="filters.contextBand === band.value
                    ? 'border-primary-500 bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-300'
                    : 'border-stone-300 bg-white/70 text-stone-700 hover:border-primary-400 dark:border-stone-800 dark:bg-[#100e0b] dark:text-stone-300'"
                  @click="filters.contextBand = band.value"
                >
                  {{ band.label }}
                </button>
              </div>
            </div>

          </div>
        </aside>

        <div class="min-w-0">
          <div class="mb-5 flex flex-col gap-3 rounded-md border border-stone-300 bg-white/60 p-3 dark:border-stone-800 dark:bg-[#171410] xl:flex-row xl:items-end xl:justify-between">
            <label class="min-w-0 flex-1">
              <span class="input-label">{{ t('models.search') }}</span>
              <span class="relative block">
                <Icon name="search" size="sm" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-stone-400" />
                <input
                  v-model.trim="filters.search"
                  type="search"
                  class="market-input h-10 w-full pl-9 pr-3"
                  :placeholder="t('models.searchPlaceholder')"
                  @input="debouncedLoad"
                />
              </span>
            </label>
            <div class="xl:pb-2">
              <p class="font-mono text-xs uppercase tracking-[0.12em] text-stone-500 dark:text-stone-400">
                {{ t('models.resultCount', { total: pagination.total, shown: filteredModels.length }) }}
              </p>
            </div>
            <div class="flex flex-wrap items-end gap-2">
              <Select
                v-model="filters.platform"
                :options="platformOptions"
                class="min-w-[180px]"
                @change="reloadModels"
              />
              <button
                class="btn-secondary btn-sm h-9 font-mono uppercase tracking-[0.12em]"
                :disabled="loading"
                @click="reloadModels"
              >
                <Icon name="refresh" size="sm" :class="loading ? 'animate-spin' : ''" />
                {{ t('common.refresh') }}
              </button>
            </div>
          </div>

          <div v-if="initialLoading" class="grid gap-4 xl:grid-cols-3">
            <div
              v-for="n in 6"
              :key="n"
              class="h-64 animate-pulse rounded-md border border-stone-300 bg-white/60 dark:border-stone-800 dark:bg-[#171410]"
            ></div>
          </div>

          <div v-else-if="filteredModels.length === 0" class="rounded-md border border-dashed border-stone-300 bg-white/60 p-10 text-center dark:border-stone-800 dark:bg-[#171410]">
            <p class="text-sm text-stone-500 dark:text-stone-400">{{ t('models.empty') }}</p>
          </div>

          <div v-else class="relative">
            <div
              v-if="loading"
              class="absolute inset-0 z-10 grid place-items-center rounded-md bg-[#f6f2ea]/70 backdrop-blur-sm dark:bg-[#14120f]/70"
            >
              <span class="inline-flex items-center gap-2 rounded-full border border-stone-300 bg-white px-4 py-2 text-sm text-stone-600 shadow-sm dark:border-stone-700 dark:bg-[#171410] dark:text-stone-300">
                <Icon name="refresh" size="sm" class="animate-spin" />
                {{ t('common.loading') }}
              </span>
            </div>
            <div class="grid gap-3 xl:grid-cols-3">
              <article
              v-for="model in filteredModels"
              :key="model.id"
              class="group cursor-pointer rounded-md border bg-white/70 p-4 transition duration-200 hover:-translate-y-1 hover:border-primary-500/60 hover:shadow-lg dark:bg-[#171410]"
              :class="drawerModel?.id === model.id
                ? 'border-primary-500 shadow-[0_0_0_1px_rgba(245,158,11,0.25)] dark:bg-[#1b1712]'
                : 'border-stone-300 dark:border-stone-800'"
              @click="openDrawer(model)"
            >
              <div class="flex items-start justify-between gap-4">
                <div class="flex min-w-0 items-start gap-3">
                  <span class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-stone-100 text-stone-950 dark:bg-stone-50">
                    <ModelIcon :model="model.model_id" :icon-key="model.icon_key || model.vendor?.icon_key" size="22px" />
                  </span>
                  <div class="min-w-0">
                    <h3 class="truncate text-base font-semibold tracking-[-0.01em] text-stone-950 dark:text-stone-50">
                      {{ model.display_name || model.model_id }}
                    </h3>
                    <p class="mt-1 truncate text-sm text-stone-500 dark:text-stone-400">
                      {{ model.vendor?.name || providerLabel(model.provider || model.platform) }}
                    </p>
                  </div>
                </div>
                <span
                  v-if="minGroupRate(model) !== null"
                  class="shrink-0 rounded-full border border-primary-500/30 bg-primary-50 px-2.5 py-1 font-mono text-[11px] font-semibold text-primary-700 dark:bg-primary-900/20 dark:text-primary-300"
                >
                  {{ t('models.rateMultiplier') }}×{{ formatRate(minGroupRate(model) || 0) }}
                </span>
              </div>

              <div class="mt-4 flex flex-wrap items-center gap-2 text-stone-500 dark:text-stone-400">
                <span
                  v-for="capability in activeCapabilities(model).slice(0, 6)"
                  :key="capability.key"
                  class="inline-flex items-center gap-1 rounded-sm text-xs"
                  :title="t(capability.labelKey)"
                >
                  <Icon :name="capability.icon" size="xs" />
                  <span>{{ t(capability.labelKey) }}</span>
                </span>
                <span v-if="activeCapabilities(model).length > 6" class="text-xs">
                  +{{ activeCapabilities(model).length - 6 }}
                </span>
              </div>

              <div class="mt-4 flex gap-1">
                <span
                  v-for="(point, idx) in healthDots(model)"
                  :key="idx"
                  class="h-3 flex-1 rounded-full"
                  :class="healthDotClass(model, point.status)"
                  :title="point.title"
                ></span>
              </div>

              <div class="mt-4 grid grid-cols-2 gap-x-3 gap-y-2 border-t border-stone-200 pt-4 dark:border-stone-800 md:grid-cols-4">
                <div v-for="item in pricingItems(model)" :key="item.key">
                  <p class="text-[11px] text-stone-500 dark:text-stone-400">{{ item.label }}</p>
                  <p class="mt-1 font-mono text-xs font-semibold text-stone-950 dark:text-stone-100">
                    {{ item.value }}
                  </p>
                </div>
              </div>

              <div class="mt-4 flex flex-wrap items-center gap-2 text-xs">
                <span class="rounded-sm bg-stone-100 px-2 py-1 text-stone-700 dark:bg-stone-900 dark:text-stone-300">
                  {{ formatTokenLimit(getContextLimits(model).maxInputTokens, 'ctx') }}
                </span>
                <span class="rounded-sm bg-stone-100 px-2 py-1 text-stone-700 dark:bg-stone-900 dark:text-stone-300">
                  {{ formatTokenLimit(getContextLimits(model).maxOutputTokens, 'out') }}
                </span>
                <span class="rounded-sm bg-stone-100 px-2 py-1 text-stone-700 dark:bg-stone-900 dark:text-stone-300">
                  {{ model.mode || 'chat' }}
                </span>
                <span
                  v-if="hasIntervalPricing(model)"
                  class="ml-auto rounded-full bg-stone-900 px-2.5 py-1 font-mono text-[11px] font-semibold text-white dark:bg-stone-100 dark:text-stone-900"
                >
                  {{ t('models.intervalPricingBadge') }}
                </span>
              </div>
              </article>
            </div>
          </div>

          <div v-if="filteredModels.length > 0" ref="loadMoreSentinel" class="mt-6 flex min-h-10 items-center justify-center">
            <span v-if="loadingMore" class="inline-flex items-center gap-2 rounded-full border border-stone-300 bg-white px-4 py-2 text-sm text-stone-600 dark:border-stone-700 dark:bg-[#171410] dark:text-stone-300">
              <Icon name="refresh" size="sm" class="animate-spin" />
              {{ t('common.loading') }}
            </span>
            <span v-else-if="!hasMoreModels" class="font-mono text-xs uppercase tracking-[0.12em] text-stone-400">
              {{ t('models.resultCount', { total: pagination.total, shown: filteredModels.length }) }}
            </span>
          </div>
        </div>
      </section>

      <Transition name="model-drawer">
        <div
          v-if="drawerModel"
          class="fixed inset-0 z-50 bg-black/40"
          @click.self="closeDrawer"
        >
        <aside class="model-drawer-panel ml-auto flex h-full w-full max-w-[560px] flex-col overflow-y-auto border-l border-stone-300 bg-[#f6f2ea] shadow-2xl dark:border-stone-800 dark:bg-[#14120f]">
          <div class="sticky top-0 z-10 flex items-start justify-between gap-4 border-b border-stone-300 bg-[#f6f2ea]/95 p-5 backdrop-blur-xl dark:border-stone-800 dark:bg-[#14120f]/95">
            <div class="flex min-w-0 items-start gap-3">
              <span class="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-stone-100 text-stone-950 dark:bg-stone-50">
                <ModelIcon :model="drawerModel.model_id" :icon-key="drawerModel.icon_key || drawerModel.vendor?.icon_key" size="26px" />
              </span>
              <div class="min-w-0">
                <p class="text-sm text-stone-500 dark:text-stone-400">
                  {{ drawerModel.vendor?.name || providerLabel(drawerModel.provider || drawerModel.platform) }}
                </p>
                <h2 class="mt-1 break-words text-2xl font-semibold tracking-[-0.02em] text-stone-950 dark:text-stone-50">
                  {{ drawerModel.display_name || drawerModel.model_id }}
                </h2>
                <p class="mt-2 break-all font-mono text-xs text-stone-500 dark:text-stone-400">{{ drawerModel.model_id }}</p>
              </div>
            </div>
            <button class="rounded-sm p-2 text-stone-600 hover:bg-stone-200 dark:text-stone-300 dark:hover:bg-stone-800" @click="closeDrawer">
              <Icon name="x" size="sm" />
            </button>
          </div>

          <div class="space-y-7 p-5">
            <section>
              <div class="flex gap-1">
                <span
                  v-for="(point, idx) in healthDots(drawerModel)"
                  :key="idx"
                  class="h-4 flex-1 rounded-full"
                  :class="healthDotClass(drawerModel, point.status)"
                  :title="point.title"
                ></span>
              </div>
              <div class="mt-4 flex flex-wrap items-center gap-2">
                <span class="inline-flex items-center gap-1.5 rounded-sm px-2 py-1 text-xs" :class="availabilityBadge(drawerModel).class">
                  <span class="h-1.5 w-1.5 rounded-full" :class="availabilityBadge(drawerModel).dotClass"></span>
                  {{ availabilityBadge(drawerModel).label }}
                </span>
                <span class="rounded-sm bg-stone-100 px-2 py-1 text-xs text-stone-700 dark:bg-stone-900 dark:text-stone-300">
                  {{ providerLabel(drawerModel.platform) }}
                </span>
                <span class="rounded-sm bg-stone-100 px-2 py-1 text-xs text-stone-700 dark:bg-stone-900 dark:text-stone-300">
                  {{ drawerModel.mode || 'chat' }}
                </span>
              </div>
            </section>

            <section>
              <h3 class="drawer-heading">{{ t('models.pricing') }}</h3>
              <div class="mt-4 grid grid-cols-2 gap-3">
                <div v-for="item in pricingItems(drawerModel)" :key="item.key" class="rounded-md border border-stone-300 bg-white/60 p-3 dark:border-stone-800 dark:bg-[#171410]">
                  <p class="text-xs text-stone-500 dark:text-stone-400">{{ item.label }}</p>
                  <p class="mt-2 font-mono text-xl font-semibold text-stone-950 dark:text-stone-50">{{ item.value }}</p>
                  <p class="mt-1 text-xs text-stone-500">{{ t('models.perMillionTokens') }}</p>
                </div>
              </div>
              <div v-if="aggregatedPricingIntervals(drawerModel).length" class="mt-4 rounded-md border border-stone-300 bg-white/60 dark:border-stone-800 dark:bg-[#171410]">
                <div class="flex items-center justify-between gap-3 border-b border-stone-200 px-3 py-2 dark:border-stone-800">
                  <p class="font-medium text-stone-900 dark:text-stone-100">{{ t('models.tierPricing') }}</p>
                  <span class="rounded-sm bg-primary-50 px-2 py-1 text-xs text-primary-700 dark:bg-primary-900/20 dark:text-primary-300">
                    {{ t('models.intervalPricingBadge') }}
                  </span>
                </div>
                <div class="divide-y divide-stone-200 dark:divide-stone-800">
                  <div
                    v-for="tier in aggregatedPricingIntervals(drawerModel)"
                    :key="tier.key"
                    class="grid gap-3 px-3 py-3 text-sm lg:grid-cols-[minmax(120px,1fr)_minmax(0,2fr)]"
                  >
                    <div>
                      <p class="font-medium text-stone-900 dark:text-stone-100">{{ t('models.contextInput') }}</p>
                      <p class="mt-1 font-mono text-xs text-stone-500">{{ formatContextInputRange(tier.min_tokens, tier.max_tokens) }}</p>
                    </div>
                    <div class="grid grid-cols-2 gap-2 text-xs text-stone-600 dark:text-stone-300">
                      <span>{{ t('models.inputPrice') }} {{ formatPricePerMillion(tier.input_price, tier.display_currency) }}</span>
                      <span>{{ t('models.outputPrice') }} {{ formatPricePerMillion(tier.output_price, tier.display_currency) }}</span>
                      <span>{{ t('models.cacheWritePrice') }} {{ formatPricePerMillion(tier.cache_write_price, tier.display_currency) }}</span>
                      <span>{{ t('models.cacheReadPrice') }} {{ formatPricePerMillion(tier.cache_read_price, tier.display_currency) }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </section>

            <section>
              <h3 class="drawer-heading">{{ t('models.relatedGroups') }}</h3>
              <div class="mt-4 space-y-2">
                <div
                  v-for="group in drawerModel.related_pricing?.groups || []"
                  :key="group.id"
                  class="flex items-center justify-between gap-3 rounded-md border border-stone-300 bg-white/60 p-3 text-sm dark:border-stone-800 dark:bg-[#171410]"
                >
                  <div>
                    <p class="font-medium text-stone-900 dark:text-stone-100">{{ group.name }}</p>
                    <p class="text-xs text-stone-500">{{ providerLabel(group.platform) }}</p>
                  </div>
                  <span class="font-mono text-sm text-primary-600 dark:text-primary-400">×{{ formatRate(group.rate_multiplier) }}</span>
                </div>
                <p v-if="!(drawerModel.related_pricing?.groups || []).length" class="text-sm text-stone-500">{{ t('models.noRelatedGroups') }}</p>
              </div>
            </section>

            <section>
              <h3 class="drawer-heading">{{ t('models.capabilities') }}</h3>
              <div class="mt-4 grid grid-cols-2 gap-3">
                <div
                  v-for="capability in capabilityFilters"
                  :key="capability.key"
                  class="flex items-center gap-2 text-sm"
                  :class="hasCapability(drawerModel, capability.key)
                    ? 'text-stone-900 dark:text-stone-100'
                    : 'text-stone-400 dark:text-stone-600'"
                >
                  <Icon :name="capability.icon" size="sm" />
                  <span>{{ t(capability.labelKey) }}</span>
                </div>
              </div>
            </section>

            <section>
              <h3 class="drawer-heading">{{ t('models.contextLimits') }}</h3>
              <div class="mt-4 grid grid-cols-3 gap-3">
                <div v-for="limit in contextLimitItems(drawerModel)" :key="limit.key" class="rounded-md bg-stone-100 p-3 dark:bg-stone-900">
                  <p class="font-mono text-lg font-semibold text-stone-950 dark:text-stone-50">{{ limit.value }}</p>
                  <p class="mt-1 text-xs text-stone-500">{{ limit.label }}</p>
                </div>
              </div>
            </section>

          </div>
        </aside>
        </div>
      </Transition>
      </main>
    </div>
  </component>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import ModelIcon from '@/components/common/ModelIcon.vue'
import Select from '@/components/common/Select.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import PublicTopNav from '@/components/home/PublicTopNav.vue'
import { listModels, listVendors, type ModelCatalog, type ModelPricingAssociationEntry, type ModelVendor } from '@/api/models'
import publicChannelMonitorAPI from '@/api/publicChannelMonitor'
import type { UserMonitorView } from '@/api/channelMonitor'
import { useAppStore, useAuthStore } from '@/stores'
import { extractApiErrorMessage } from '@/utils/apiError'
import { currencySymbol } from '@/utils/pricing'
import { useTheme } from '@/composables/useTheme'
import { dedupeModelsByModelId, matchMonitorsByModelId } from './modelMarketplaceMonitor'

const { t } = useI18n()
const props = withDefaults(defineProps<{
  embedded?: boolean
}>(), {
  embedded: false,
})
const appStore = useAppStore()
const authStore = useAuthStore()
const { syncThemeFromDocument } = useTheme()
const pageShellClass = computed(() => props.embedded
  ? 'text-stone-950 dark:text-stone-100'
  : 'min-h-screen bg-[#f6f2ea] text-stone-950 dark:bg-[#14120f] dark:text-stone-100'
)
const mainClass = computed(() => props.embedded
  ? 'mx-auto max-w-[1440px] px-0 pb-8 pt-0'
  : 'mx-auto max-w-[1440px] px-4 pb-14 pt-10 lg:px-8'
)

const loading = ref(true)
const loadingMore = ref(false)
const models = ref<ModelCatalog[]>([])
const vendors = ref<ModelVendor[]>([])
const monitors = ref<UserMonitorView[]>([])
const drawerModel = ref<ModelCatalog | null>(null)
const hasLoadedModels = ref(false)
const vendorSearch = ref('')
const page = ref(0)
const pagination = reactive({ total: 0, pages: 1 })
const loadMoreSentinel = ref<HTMLElement | null>(null)
const filters = reactive({
  search: '',
  platform: '',
  vendor_ids: [] as string[],
  capabilities: [] as string[],
  priceBand: 'all',
  contextBand: 'all'
})
let searchTimer: number | undefined
let modelsAbortController: AbortController | null = null
let modelsRequestSeq = 0
let loadMoreObserver: IntersectionObserver | null = null
const MARKETPLACE_PAGE_SIZE = 20

type CapabilityIcon = 'eye' | 'link' | 'filter' | 'brain' | 'database' | 'globe' | 'document' | 'terminal' | 'sparkles'

const capabilityFilters: Array<{ key: string; labelKey: string; icon: CapabilityIcon }> = [
  { key: 'vision', labelKey: 'models.capabilityLabels.vision', icon: 'eye' },
  { key: 'function_calling', labelKey: 'models.capabilityLabels.functionCalling', icon: 'link' },
  { key: 'tool_choice', labelKey: 'models.capabilityLabels.toolChoice', icon: 'filter' },
  { key: 'reasoning', labelKey: 'models.capabilityLabels.reasoning', icon: 'brain' },
  { key: 'prompt_caching', labelKey: 'models.capabilityLabels.promptCaching', icon: 'database' },
  { key: 'web_search', labelKey: 'models.capabilityLabels.webSearch', icon: 'globe' },
  { key: 'pdf_input', labelKey: 'models.capabilityLabels.pdfInput', icon: 'document' },
  { key: 'computer_use', labelKey: 'models.capabilityLabels.computerUse', icon: 'terminal' },
  { key: 'image_output', labelKey: 'models.capabilityLabels.imageOutput', icon: 'sparkles' },
]

const platforms = computed(() => Array.from(new Set(models.value.map(m => m.platform).filter(Boolean))).sort())
const platformOptions = computed(() => [
  { value: '', label: t('models.allPlatforms') },
  ...platforms.value.map(platform => ({ value: platform, label: providerLabel(platform) })),
])

const filteredVendors = computed(() => {
  const query = vendorSearch.value.trim().toLowerCase()
  if (!query) return vendors.value
  return vendors.value.filter(vendor => [
    vendor.name,
    vendor.provider_key,
    vendor.icon_key,
  ].some(value => String(value || '').toLowerCase().includes(query)))
})

const priceBands = computed(() => [
  { value: 'all', label: t('models.priceBands.all') },
  { value: 'free', label: t('models.priceBands.free') },
  { value: 'lt1', label: t('models.priceBands.lt1') },
  { value: '1to5', label: t('models.priceBands.1to5') },
  { value: '5to15', label: t('models.priceBands.5to15') },
  { value: 'gte15', label: t('models.priceBands.gte15') },
])

const contextBands = computed(() => [
  { value: 'all', label: t('models.contextBands.all') },
  { value: 'lt32k', label: t('models.contextBands.lt32k') },
  { value: '32to128k', label: t('models.contextBands.32to128k') },
  { value: '128to256k', label: t('models.contextBands.128to256k') },
  { value: 'gte256k', label: t('models.contextBands.gte256k') },
])

const filteredModels = computed(() => {
  return models.value.filter((model) => {
    if (filters.vendor_ids.length > 0 && (!model.vendor_id || !filters.vendor_ids.includes(String(model.vendor_id)))) {
      return false
    }
    if (filters.capabilities.length > 0 && filters.capabilities.some(key => !hasCapability(model, key))) {
      return false
    }
    if (!matchesPriceBand(model)) return false
    if (!matchesContextBand(model)) return false
    return true
  })
})

function resetFilters() {
  filters.search = ''
  filters.platform = ''
  filters.vendor_ids = []
  filters.capabilities = []
  filters.priceBand = 'all'
  filters.contextBand = 'all'
  void reloadModels()
}

function debouncedLoad() {
  window.clearTimeout(searchTimer)
  searchTimer = window.setTimeout(() => {
    void reloadModels()
  }, 250)
}

const initialLoading = computed(() => loading.value && !hasLoadedModels.value && models.value.length === 0)
const hasMoreModels = computed(() => page.value < pagination.pages)

async function reloadModels() {
  page.value = 0
  models.value = []
  drawerModel.value = null
  await loadModelsPage(1, false)
}

async function loadNextModelsPage() {
  if (loading.value || loadingMore.value || !hasMoreModels.value) return
  await loadModelsPage(page.value + 1, true)
}

async function loadModelsPage(targetPage: number, append: boolean) {
  modelsAbortController?.abort()
  const ctrl = new AbortController()
  modelsAbortController = ctrl
  const seq = ++modelsRequestSeq
  if (append) {
    loadingMore.value = true
  } else {
    loading.value = true
  }
  try {
    const res = await listModels(targetPage, MARKETPLACE_PAGE_SIZE, {
      search: filters.search,
      platform: filters.platform,
    }, { signal: ctrl.signal })
    if (ctrl.signal.aborted || seq !== modelsRequestSeq) return
    models.value = dedupeModelsByModelId(append ? [...models.value, ...(res.items || [])] : (res.items || []))
    page.value = res.page || targetPage
    pagination.total = res.total
    pagination.pages = res.pages || 1
    await nextTick()
    setupLoadMoreObserver()
  } catch (err) {
    const e = err as { name?: string; code?: string }
    if (e?.name === 'AbortError' || e?.code === 'ERR_CANCELED') return
    appStore.showError(extractApiErrorMessage(err, t('models.loadFailed')))
  } finally {
    if (seq === modelsRequestSeq) {
      if (append) {
        loadingMore.value = false
      } else {
        loading.value = false
        hasLoadedModels.value = true
      }
      if (modelsAbortController === ctrl) modelsAbortController = null
    }
  }
}

function setupLoadMoreObserver() {
  loadMoreObserver?.disconnect()
  if (!loadMoreSentinel.value) return
  loadMoreObserver = new IntersectionObserver((entries) => {
    if (entries.some(entry => entry.isIntersecting)) {
      void loadNextModelsPage()
    }
  }, {
    rootMargin: '480px 0px',
    threshold: 0,
  })
  loadMoreObserver.observe(loadMoreSentinel.value)
}

function openDrawer(model: ModelCatalog) {
  drawerModel.value = model
}

function closeDrawer() {
  drawerModel.value = null
}

function providerLabel(value: string) {
  const normalized = (value || '').trim()
  if (!normalized) return '-'
  const known: Record<string, string> = {
    openai: 'OpenAI',
    anthropic: 'Anthropic',
    claude: 'Anthropic',
    gemini: 'Google AI Studio',
    google: 'Google',
    vertex_ai: 'Vertex AI',
    deepseek: 'DeepSeek',
    qwen: 'DashScope',
    dashscope: 'DashScope',
    xai: 'xAI',
    openrouter: 'OpenRouter',
  }
  return known[normalized.toLowerCase()] || normalized
}

function numberFrom(value: unknown): number | null {
  if (typeof value === 'number' && Number.isFinite(value)) return value
  if (typeof value === 'string' && value.trim() !== '') {
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : null
  }
  return null
}

function pricePerMillion(value: unknown): number | null {
  const raw = numberFrom(value)
  if (raw === null || raw < 0) return null
  return raw * 1_000_000
}

function formatPricePerMillion(value: unknown, displayCurrency: unknown = 'USD') {
  const price = pricePerMillion(value)
  if (price === null) return '-'
  const decimals = price >= 10 ? 2 : price >= 1 ? 3 : 4
  return `${currencySymbol(displayCurrency)}${price.toLocaleString(undefined, {
    minimumFractionDigits: price % 1 === 0 ? 0 : Math.min(2, decimals),
    maximumFractionDigits: decimals,
  })}`
}

function pricingItems(model: ModelCatalog) {
  const pricing = representativePricing(model)
  return [
    { key: 'input', label: t('models.inputPrice'), value: formatPricePerMillion(pricing?.input_price, pricing?.display_currency) },
    { key: 'output', label: t('models.outputPrice'), value: formatPricePerMillion(pricing?.output_price, pricing?.display_currency) },
    { key: 'cache-write', label: t('models.cacheWritePrice'), value: formatPricePerMillion(pricing?.cache_write_price, pricing?.display_currency) },
    { key: 'cache-read', label: t('models.cacheReadPrice'), value: formatPricePerMillion(pricing?.cache_read_price, pricing?.display_currency) },
  ]
}

function representativePricing(model: ModelCatalog): ModelPricingAssociationEntry | null {
  const entries = model.related_pricing?.entries || []
  return entries[0] || null
}

function pricingIntervals(model: ModelCatalog) {
  return (model.related_pricing?.entries || []).flatMap(entry =>
    (entry.intervals || []).map(interval => ({
      ...interval,
      display_currency: entry.display_currency,
    }))
  )
}

function minGroupRate(model: ModelCatalog) {
  const rates = (model.related_pricing?.groups || [])
    .map(group => numberFrom(group.rate_multiplier))
    .filter((value): value is number => value !== null)
  if (rates.length === 0) return null
  return Math.min(...rates)
}

function aggregatedPricingIntervals(model: ModelCatalog) {
  const byKey = new Map<string, ModelPricingAssociationEntry['intervals'][number] & { key: string }>()
  for (const tier of pricingIntervals(model)) {
    const key = [
      tier.min_tokens || 0,
      tier.max_tokens ?? 'max',
      tier.input_price ?? '',
      tier.output_price ?? '',
      tier.cache_write_price ?? '',
      tier.cache_read_price ?? '',
      tier.per_request_price ?? '',
      tier.display_currency ?? 'USD',
    ].join(':')
    if (!byKey.has(key)) {
      byKey.set(key, { ...tier, key })
    }
  }
  return Array.from(byKey.values()).sort((a, b) => {
    if ((a.sort_order || 0) !== (b.sort_order || 0)) return (a.sort_order || 0) - (b.sort_order || 0)
    return (a.min_tokens || 0) - (b.min_tokens || 0)
  })
}

function hasIntervalPricing(model: ModelCatalog) {
  return model.related_pricing?.has_intervals === true || pricingIntervals(model).length > 0
}

function hasCapability(model: ModelCatalog, key: string) {
  return model.capabilities?.[key] === true
}

function activeCapabilities(model: ModelCatalog) {
  return capabilityFilters.filter(capability => hasCapability(model, capability.key))
}

function getContextLimits(model: ModelCatalog) {
  const caps = model.capabilities?.context_limits
  return {
    maxInputTokens: numberFrom(model.metadata?.max_input_tokens) || numberFrom(caps?.max_input_tokens) || null,
    maxOutputTokens: numberFrom(model.metadata?.max_output_tokens) || numberFrom(caps?.max_output_tokens) || null,
    maxTokens: numberFrom(model.metadata?.max_tokens) || numberFrom(caps?.max_tokens) || null,
  }
}

function formatTokenLimit(value: number | null, suffix = '') {
  if (!value || value <= 0) return '-'
  const formatted = value >= 1_000_000
    ? `${Number((value / 1_000_000).toFixed(1))}M`
    : value >= 1000
      ? `${Number((value / 1000).toFixed(1))}K`
      : String(value)
  return suffix ? `${formatted} ${suffix}` : formatted
}

function healthDots(model: ModelCatalog) {
  const matches = monitorMatches(model)
  const latest = matches.flatMap(monitor => (monitor.timeline || []).slice(0, 24))
  if (latest.length === 0) {
    const badge = availabilityBadge(model)
    const status = badge.status === 'disabled' ? 'empty' : badge.status
    return Array.from({ length: 24 }, () => ({ status, title: badge.label }))
  }
  const real = latest
    .slice(0, 24)
    .reverse()
    .map(point => ({
      status: point.status,
      title: `${point.checked_at || ''} · ${point.status} · ${formatLatencyLabel(point.latency_ms)}`,
    }))
  const pad = Array.from({ length: Math.max(0, 24 - real.length) }, () => ({ status: 'empty', title: '' }))
  return [...pad, ...real]
}

function healthDotClass(model: ModelCatalog, status?: string) {
  if (model.status !== 'active') return 'bg-stone-300 dark:bg-stone-700'
  if (status === 'operational') return 'bg-emerald-500'
  if (status === 'degraded') return 'bg-amber-500'
  if (status === 'failed' || status === 'error') return 'bg-red-500'
  const badge = availabilityBadge(model)
  if (badge.status === 'operational') return 'bg-emerald-500'
  if (badge.status === 'degraded') return 'bg-amber-500'
  if (badge.status === 'failed') return 'bg-red-500'
  return 'bg-stone-300 dark:bg-stone-700'
}

function monitorMatches(model: ModelCatalog) {
  return matchMonitorsByModelId(model.model_id, monitors.value)
}

function availabilityBadge(model: ModelCatalog) {
  const matches = monitorMatches(model)
  if (model.status !== 'active') {
    return {
      status: 'disabled',
      label: t('models.availabilityDisabled'),
      class: 'bg-stone-100 text-stone-600 dark:bg-stone-900 dark:text-stone-300',
      dotClass: 'bg-stone-400'
    }
  }
  if (matches.length === 0) {
    return {
      status: 'operational',
      label: t('models.availabilityDefault'),
      class: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300',
      dotClass: 'bg-emerald-500'
    }
  }
  const avg = matches.reduce((sum, monitor) => sum + monitor.availability_7d, 0) / matches.length
  if (avg < 60) {
    return {
      status: 'failed',
      label: t('models.availabilityOperational', { value: formatAvailability(avg) }),
      class: 'bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-300',
      dotClass: 'bg-red-500'
    }
  }
  if (avg < 90) {
    return {
      status: 'degraded',
      label: t('models.availabilityOperational', { value: formatAvailability(avg) }),
      class: 'bg-amber-50 text-amber-700 dark:bg-amber-900/20 dark:text-amber-300',
      dotClass: 'bg-amber-500'
    }
  }
  return {
    status: 'operational',
    label: t('models.availabilityOperational', { value: formatAvailability(avg) }),
    class: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300',
    dotClass: 'bg-emerald-500'
  }
}

function contextLimitItems(model: ModelCatalog) {
  const limits = getContextLimits(model)
  return [
    { key: 'input', label: t('models.maxInputTokens'), value: formatTokenLimit(limits.maxInputTokens) },
    { key: 'output', label: t('models.maxOutputTokens'), value: formatTokenLimit(limits.maxOutputTokens) },
    { key: 'total', label: t('models.maxTokens'), value: formatTokenLimit(limits.maxTokens) },
  ]
}

function formatAvailability(value: number) {
  return `${Number(value || 0).toFixed(2)}%`
}

function formatLatencyLabel(value: number | null) {
  return value == null ? '-' : `${value}ms`
}

function formatRate(value: number) {
  return Number(value || 0).toLocaleString(undefined, { maximumFractionDigits: 4 })
}

function matchesPriceBand(model: ModelCatalog) {
  const band = filters.priceBand
  if (band === 'all') return true
  const pricing = representativePricing(model)
  const input = pricePerMillion(pricing?.input_price) ?? 0
  const output = pricePerMillion(pricing?.output_price) ?? 0
  if (band === 'free') return input === 0 && output === 0
  const price = Math.max(input, output)
  if (band === 'lt1') return price < 1
  if (band === '1to5') return price >= 1 && price < 5
  if (band === '5to15') return price >= 5 && price < 15
  if (band === 'gte15') return price >= 15
  return true
}

function matchesContextBand(model: ModelCatalog) {
  const band = filters.contextBand
  if (band === 'all') return true
  const maxContext = getContextLimits(model).maxTokens || getContextLimits(model).maxInputTokens || 0
  if (band === 'lt32k') return maxContext > 0 && maxContext < 32_000
  if (band === '32to128k') return maxContext >= 32_000 && maxContext < 128_000
  if (band === '128to256k') return maxContext >= 128_000 && maxContext < 256_000
  if (band === 'gte256k') return maxContext >= 256_000
  return true
}

function formatContextInputRange(minTokens: number | null, maxTokens: number | null) {
  const min = minTokens && minTokens > 0 ? formatTokenLimit(minTokens) : '0'
  const max = maxTokens && maxTokens > 0 ? formatTokenLimit(maxTokens) : '∞'
  return `[${min}, ${max})`
}

function isVendorSelected(vendorID: number) {
  return filters.vendor_ids.includes(String(vendorID))
}

function toggleVendorFilter(vendorID: number) {
  const value = String(vendorID)
  const index = filters.vendor_ids.indexOf(value)
  if (index >= 0) {
    filters.vendor_ids.splice(index, 1)
  } else {
    filters.vendor_ids.push(value)
  }
}

function clearVendorFilters() {
  filters.vendor_ids = []
}

function toggleCapabilityFilter(key: string) {
  const index = filters.capabilities.indexOf(key)
  if (index >= 0) {
    filters.capabilities.splice(index, 1)
  } else {
    filters.capabilities.push(key)
  }
}

async function loadFilterMetadata() {
  try {
    vendors.value = await listVendors()
  } catch {
    vendors.value = []
  }
}

async function loadMonitorMetadata() {
  try {
    const res = await publicChannelMonitorAPI.list()
    monitors.value = res.items || []
  } catch {
    monitors.value = []
  }
}

onMounted(() => {
  syncThemeFromDocument()
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
  void loadFilterMetadata()
  void loadMonitorMetadata()
  void reloadModels()
})

watch(loadMoreSentinel, () => {
  nextTick(() => setupLoadMoreObserver())
})

onBeforeUnmount(() => {
  window.clearTimeout(searchTimer)
  modelsAbortController?.abort()
  loadMoreObserver?.disconnect()
})
</script>

<style scoped>
.drawer-heading {
  @apply flex items-center gap-3 font-mono text-[12px] uppercase tracking-[0.14em] text-stone-600 dark:text-stone-400;
}

.market-input {
  @apply rounded-sm border border-stone-300 bg-white text-sm text-stone-900 outline-none transition-colors placeholder:text-stone-400 focus:border-primary-500 focus:ring-2 focus:ring-primary-500/25 dark:border-stone-800 dark:bg-[#100e0b] dark:text-stone-100;
}

.market-filter-option {
  @apply flex w-full items-center gap-2 rounded-sm border border-stone-300 bg-white/70 px-2.5 py-1.5 text-left text-sm text-stone-700 transition-colors hover:border-primary-400 hover:text-stone-950 dark:border-stone-800 dark:bg-[#100e0b] dark:text-stone-300 dark:hover:text-stone-100;
}

.market-filter-option-active {
  @apply border-primary-500 bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-300;
}

.market-filter-check {
  @apply flex h-4 w-4 shrink-0 items-center justify-center rounded-sm border border-stone-300 bg-white text-primary-600 dark:border-stone-700 dark:bg-stone-950 dark:text-primary-300;
}

.market-filter-option-active .market-filter-check {
  @apply border-primary-500 bg-primary-500 text-white dark:border-primary-500 dark:bg-primary-500;
}

.model-drawer-enter-active,
.model-drawer-leave-active {
  transition: background-color 0.22s ease, opacity 0.22s ease;
}

.model-drawer-enter-active .model-drawer-panel,
.model-drawer-leave-active .model-drawer-panel {
  transition: transform 0.24s ease, opacity 0.24s ease;
}

.model-drawer-enter-from,
.model-drawer-leave-to {
  opacity: 0;
}

.model-drawer-enter-from .model-drawer-panel,
.model-drawer-leave-to .model-drawer-panel {
  opacity: 0.6;
  transform: translateX(100%);
}
</style>
