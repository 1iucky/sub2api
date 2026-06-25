<template>
  <AppLayout>
    <div class="space-y-5">
      <section class="min-w-0 space-y-4">
          <div class="card p-4">
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div class="min-w-0">
                <h3 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('admin.models.vendorManagement') }}</h3>
                <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('admin.models.vendorManagementDesc') }}</p>
              </div>
              <a
                href="https://icons.lobehub.com/components/lobe-hub"
                target="_blank"
                rel="noopener noreferrer"
                class="btn-secondary btn-sm"
              >
                {{ t('admin.models.iconReference') }}
                <Icon name="externalLink" size="xs" />
              </a>
            </div>

            <div class="mt-4 flex items-center gap-2">
              <span
                class="shrink-0"
                @mouseenter="showVendorOverflow('left', $event)"
                @mouseleave="scheduleHideVendorOverflow"
              >
                <button class="vendor-scroll-button" type="button" :disabled="hiddenLeftVendors.length === 0" @click="scrollVendorTabs('left')">
                  <Icon name="chevronLeft" size="sm" />
                </button>
              </span>
              <div
                ref="vendorTabsRef"
                class="vendor-tabs-scroll"
                @scroll.passive="handleVendorTabsScroll"
              >
                <div class="vendor-tabs">
                  <button
                    type="button"
                    class="vendor-tab"
                    :class="!filters.vendor_id ? 'vendor-tab-active' : ''"
                    @click="selectVendorTab('')"
                  >
                    <span>{{ t('models.allVendors') }}</span>
                  </button>
                  <div
                    v-for="vendor in vendors"
                    :key="vendor.id"
                    class="vendor-tab"
                    :class="filters.vendor_id === String(vendor.id) ? 'vendor-tab-active' : ''"
                  >
                    <button
                      type="button"
                      class="flex min-w-0 flex-1 items-center gap-2"
                      @click="selectVendorTab(String(vendor.id))"
                    >
                      <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-sm bg-white dark:bg-dark-950">
                        <ModelIcon :model="vendor.name" :icon-key="vendor.icon_key || vendor.provider_key" size="15px" />
                      </span>
                      <span class="max-w-[160px] truncate">{{ vendor.name }}</span>
                    </button>
                    <span class="shrink-0">
                      <button
                        type="button"
                        class="vendor-tab-action"
                        :aria-label="t('common.actions')"
                        @click.stop="toggleVendorMenu(vendor, $event)"
                      >
                        <Icon name="more" size="xs" />
                      </button>
                    </span>
                  </div>
                </div>
              </div>
              <span
                class="shrink-0"
                @mouseenter="showVendorOverflow('right', $event)"
                @mouseleave="scheduleHideVendorOverflow"
              >
                <button class="vendor-scroll-button" type="button" :disabled="hiddenRightVendors.length === 0" @click="scrollVendorTabs('right')">
                  <Icon name="chevronRight" size="sm" />
                </button>
              </span>
              <button class="btn-primary h-10 shrink-0" type="button" @click="openVendorDialog()">
                <Icon name="plus" size="xs" />
                {{ t('admin.models.addVendor') }}
              </button>
            </div>
          </div>

          <div class="card p-4">
            <div class="grid gap-3 xl:grid-cols-[minmax(220px,1fr)_180px_150px_auto]">
              <label>
                <span class="input-label">{{ t('models.search') }}</span>
                <span class="relative block">
                  <Icon name="search" size="sm" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                  <input
                    v-model.trim="filters.search"
                    class="input pl-9"
                    :placeholder="t('models.searchPlaceholder')"
                    @input="debouncedLoad"
                  />
                </span>
              </label>
              <label>
                <span class="input-label">{{ t('models.platform') }}</span>
                <Select
                  v-model="filters.platform"
                  :options="[{ value: '', label: t('models.allPlatforms') }, ...platformOptions]"
                  @change="reloadFromFirstPage"
                />
              </label>
              <label>
                <span class="input-label">{{ t('models.status') }}</span>
                <Select
                  v-model="filters.status"
                  :options="statusFilterOptions"
                  @change="reloadFromFirstPage"
                />
              </label>
              <div class="flex flex-wrap items-end justify-end gap-2">
                <button class="btn-secondary h-10 min-w-[128px]" :disabled="syncing" @click="handleSync">
                  <Icon name="sync" size="sm" :class="syncing ? 'animate-spin' : ''" />
                  {{ syncing ? t('admin.models.syncing') : t('admin.models.syncPricing') }}
                </button>
                <button class="btn-primary h-10 min-w-[112px]" @click="openCreate">
                  <Icon name="plus" size="sm" />
                  {{ t('admin.models.addModel') }}
                </button>
              </div>
            </div>
          </div>

          <div class="overflow-hidden rounded-md border border-gray-200 bg-white dark:border-dark-800 dark:bg-dark-900">
            <table class="min-w-full divide-y divide-gray-200 dark:divide-dark-800">
              <thead class="bg-gray-50 dark:bg-dark-900/70">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">{{ t('models.model') }}</th>
                  <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">{{ t('models.vendor') }}</th>
                  <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">{{ t('models.contextLimits') }}</th>
                  <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">{{ t('models.relatedChannels') }}</th>
                  <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">{{ t('models.status') }}</th>
                  <th class="w-44 px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500">{{ t('common.actions') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
                <tr v-if="loading">
                  <td colspan="6" class="px-4 py-10 text-center text-sm text-gray-500">{{ t('common.loading') }}</td>
                </tr>
                <tr v-else-if="models.length === 0">
                  <td colspan="6" class="px-4 py-10 text-center text-sm text-gray-500">{{ t('models.empty') }}</td>
                </tr>
                <tr v-for="model in models" v-else :key="model.id" class="hover:bg-gray-50 dark:hover:bg-dark-800/50">
                  <td class="px-4 py-3">
                    <div class="flex items-center gap-3">
                      <span class="flex h-9 w-9 shrink-0 items-center justify-center rounded-sm border border-gray-200 bg-gray-50 dark:border-dark-700 dark:bg-dark-800">
                        <ModelIcon :model="model.vendor?.name || model.provider || model.model_id" :icon-key="model.vendor?.icon_key || model.vendor?.provider_key" size="22px" />
                      </span>
                      <div class="min-w-0">
                        <p class="truncate text-sm font-medium text-gray-900 dark:text-white">{{ model.display_name || model.model_id }}</p>
                        <p class="truncate font-mono text-xs text-gray-500">{{ model.model_id }}</p>
                        <p class="mt-0.5 text-xs text-gray-400">{{ providerLabel(model.platform) }} · {{ model.mode || 'chat' }}</p>
                      </div>
                    </div>
                  </td>
                  <td class="px-4 py-3 text-sm text-gray-600 dark:text-dark-300">
                    {{ model.vendor?.name || model.provider || '-' }}
                  </td>
                  <td class="px-4 py-3">
                    <div class="min-w-[160px] text-xs text-gray-600 dark:text-dark-300">
                      <span>{{ formatTokenLimit(getContextLimits(model).maxTokens) }}</span>
                      <span class="mx-1 text-gray-300">/</span>
                      <span>{{ formatTokenLimit(getContextLimits(model).maxInputTokens) }} in</span>
                    </div>
                  </td>
                  <td class="px-4 py-3">
                    <span class="text-sm text-gray-700 dark:text-dark-300">
                      {{ model.related_pricing?.channel_count || 0 }} {{ t('models.channelsUnit') }}
                    </span>
                  </td>
                  <td class="px-4 py-3">
                    <span class="rounded-sm px-2 py-1 text-xs" :class="statusClass(model.status)">
                      {{ model.status === 'active' ? t('models.active') : t('models.disabled') }}
                    </span>
                  </td>
                  <td class="px-4 py-3 text-right">
                    <div class="flex justify-end gap-1 whitespace-nowrap">
                    <button class="btn-ghost btn-sm" @click="openEdit(model)">
                      <Icon name="edit" size="xs" />
                      {{ t('common.edit') }}
                    </button>
                    <button class="btn-ghost btn-sm text-red-600 hover:text-red-700" @click="handleDelete(model)">
                      <Icon name="trash" size="xs" />
                      {{ t('common.delete') }}
                    </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <Pagination
            v-if="pagination.total > 0"
            :page="pagination.page"
            :total="pagination.total"
            :page-size="pagination.page_size"
            @update:page="handlePageChange"
            @update:pageSize="handlePageSizeChange"
          />
      </section>
    </div>

    <Teleport to="body">
    <div v-if="dialogOpen" class="fixed inset-0 z-[100000010] flex items-start justify-center overflow-y-auto bg-black/50 p-4">
      <div class="my-6 w-full max-w-5xl rounded-md border border-gray-200 bg-white shadow-xl dark:border-dark-800 dark:bg-dark-900">
        <div class="sticky top-0 z-10 flex items-center justify-between border-b border-gray-200 bg-white px-5 py-4 dark:border-dark-800 dark:bg-dark-900">
          <div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ editing ? t('admin.models.editModel') : t('admin.models.addModel') }}</h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('admin.models.formHint') }}</p>
          </div>
          <button class="btn-ghost btn-icon" @click="dialogOpen = false" :aria-label="t('common.close')">
            <Icon name="x" size="sm" />
          </button>
        </div>

        <form class="grid gap-5 p-5" @submit.prevent="saveModel">
          <section class="grid gap-4 md:grid-cols-2">
            <div class="md:col-span-2 flex items-center gap-3">
              <span class="flex h-11 w-11 shrink-0 items-center justify-center rounded-sm border border-gray-200 bg-gray-50 dark:border-dark-700 dark:bg-dark-800">
                <ModelIcon :model="selectedVendor?.name || form.provider || form.model_id" :icon-key="selectedVendor?.icon_key || selectedVendor?.provider_key" size="24px" />
              </span>
              <div>
                <h4 class="font-medium text-gray-900 dark:text-white">{{ t('admin.models.basicInfo') }}</h4>
                <p class="text-sm text-gray-500">{{ t('admin.models.iconFollowsVendor') }}</p>
              </div>
            </div>

            <label>
              <span class="input-label">{{ t('models.modelId') }}</span>
              <input v-model.trim="form.model_id" class="input" required />
            </label>
            <label>
              <span class="input-label">{{ t('models.displayName') }}</span>
              <input v-model.trim="form.display_name" class="input" />
            </label>
            <div ref="platformDropdownRef" class="relative">
              <span class="input-label">{{ t('models.platform') }}</span>
              <div class="relative">
                <button
                  type="button"
                  class="platform-multiselect-trigger"
                  :class="platformDropdownOpen ? 'platform-multiselect-trigger-open' : ''"
                  @click="platformDropdownOpen = !platformDropdownOpen"
                >
                  <span class="min-w-0 flex-1 truncate text-left">{{ platformSelectionLabel }}</span>
                  <Icon
                    name="chevronDown"
                    size="sm"
                    class="shrink-0 text-gray-400 transition-transform duration-200"
                    :class="platformDropdownOpen ? 'rotate-180' : ''"
                  />
                </button>
                <div v-if="platformDropdownOpen" class="platform-multiselect-menu">
                  <button
                    v-for="option in platformOptions"
                    :key="option.value"
                    type="button"
                    class="platform-multiselect-option"
                    :class="isModelPlatformSelected(option.value) ? 'platform-multiselect-option-active' : ''"
                    @click="toggleModelPlatform(option.value)"
                  >
                    <span class="platform-multiselect-check">
                      <Icon v-if="isModelPlatformSelected(option.value)" name="check" size="xs" />
                    </span>
                    <span class="min-w-0 flex-1 truncate">{{ option.label }}</span>
                  </button>
                </div>
              </div>
              <span class="input-hint">{{ t('admin.models.platformMultiSelectHint') }}</span>
            </div>
            <label>
              <span class="input-label">{{ t('models.vendor') }}</span>
              <Select
                v-model="vendorSelection"
                :options="vendorSelectOptions"
                :placeholder="t('admin.models.selectVendor')"
                searchable
                clearable
                @change="applySelectedVendor"
              />
            </label>
            <label>
              <span class="input-label">{{ t('models.provider') }}</span>
              <input v-model.trim="form.provider" class="input" :placeholder="t('admin.models.providerPlaceholder')" />
            </label>
            <label>
              <span class="input-label">{{ t('admin.models.mode') }}</span>
              <Select v-model="form.mode" :options="modeOptions" />
            </label>
            <label>
              <span class="input-label">{{ t('models.status') }}</span>
              <Select v-model="form.status" :options="modelStatusOptions" />
            </label>
            <label>
              <span class="input-label">{{ t('admin.models.visibility') }}</span>
              <Select v-model="form.visibility" :options="visibilityOptions" />
            </label>
            <label class="md:col-span-2">
              <span class="input-label">{{ t('models.description') }}</span>
              <textarea v-model.trim="form.description" class="input min-h-24"></textarea>
            </label>
            <label class="md:col-span-2">
              <span class="input-label">{{ t('admin.models.tags') }}</span>
              <input v-model.trim="tagInput" class="input" :placeholder="t('admin.models.tagsPlaceholder')" />
              <span class="input-hint">{{ t('admin.models.commaSeparated') }}</span>
            </label>
          </section>

          <section class="border-t border-gray-200 pt-5 dark:border-dark-800">
            <h4 class="font-medium text-gray-900 dark:text-white">{{ t('models.capabilities') }}</h4>
            <div class="mt-3 grid gap-2 sm:grid-cols-2 lg:grid-cols-4">
              <label
                v-for="capability in capabilityOptions"
                :key="capability.key"
                class="flex cursor-pointer items-center gap-2 rounded-md border px-3 py-2 text-sm transition-colors"
                :class="capabilityForm[capability.key]
                  ? 'border-primary-500 bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-300'
                  : 'border-gray-200 text-gray-600 hover:border-gray-300 dark:border-dark-700 dark:text-dark-300'"
              >
                <input v-model="capabilityForm[capability.key]" type="checkbox" class="h-4 w-4 rounded border-gray-300 text-primary-500 focus:ring-primary-500" />
                <Icon :name="capability.icon" size="xs" />
                <span>{{ t(capability.labelKey) }}</span>
              </label>
            </div>
          </section>

          <section class="grid gap-4 border-t border-gray-200 pt-5 dark:border-dark-800 lg:grid-cols-3">
            <div class="lg:col-span-3">
              <h4 class="font-medium text-gray-900 dark:text-white">{{ t('models.contextLimits') }}</h4>
            </div>
            <label>
              <span class="input-label">{{ t('models.maxInputTokens') }}</span>
              <input v-model.number="contextForm.max_input_tokens" type="number" min="0" class="input" />
            </label>
            <label>
              <span class="input-label">{{ t('models.maxOutputTokens') }}</span>
              <input v-model.number="contextForm.max_output_tokens" type="number" min="0" class="input" />
            </label>
            <label>
              <span class="input-label">{{ t('models.maxTokens') }}</span>
              <input v-model.number="contextForm.max_tokens" type="number" min="0" class="input" />
            </label>
          </section>

          <div class="sticky bottom-0 flex justify-end gap-2 border-t border-gray-200 bg-white py-4 dark:border-dark-800 dark:bg-dark-900">
            <button type="button" class="btn-secondary" @click="dialogOpen = false">{{ t('common.cancel') }}</button>
            <button type="submit" class="btn-primary" :disabled="saving">
              <Icon name="check" size="sm" />
              {{ saving ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </form>
      </div>
    </div>
    </Teleport>

    <Teleport to="body">
    <div v-if="vendorDialogOpen" class="fixed inset-0 z-[100000020] flex items-start justify-center overflow-y-auto bg-black/50 p-4">
      <div class="my-6 w-full max-w-lg rounded-md border border-gray-200 bg-white p-5 shadow-xl dark:border-dark-800 dark:bg-dark-900">
        <div class="mb-4 flex items-center justify-between">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.models.vendorFormTitle') }}</h3>
          <button class="btn-ghost btn-icon" @click="vendorDialogOpen = false">
            <Icon name="x" size="sm" />
          </button>
        </div>
        <form class="grid gap-3" @submit.prevent="saveVendor">
          <div class="flex items-center gap-3 rounded-md bg-gray-50 p-3 dark:bg-dark-800">
            <span class="flex h-10 w-10 items-center justify-center rounded-sm border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900">
              <ModelIcon :model="vendorForm.name" :icon-key="vendorForm.icon_key || vendorForm.name" size="24px" />
            </span>
            <p class="text-sm text-gray-500">{{ t('admin.models.vendorIconHint') }}</p>
          </div>
          <label>
            <span class="input-label">{{ t('admin.models.vendorName') }}</span>
            <input v-model.trim="vendorForm.name" class="input" required />
          </label>
          <label>
            <span class="input-label">{{ t('models.iconKey') }}</span>
            <input v-model.trim="vendorForm.icon_key" class="input" placeholder="openai / claude / gemini" />
            <span class="input-hint">
              {{ t('admin.models.vendorIconExtra') }}
              <a
                href="https://icons.lobehub.com/components/lobe-hub"
                target="_blank"
                rel="noopener noreferrer"
                class="inline-flex items-center gap-1 text-primary-600 hover:text-primary-700 dark:text-primary-400"
              >
                {{ t('admin.models.iconReference') }}
                <Icon name="externalLink" size="xs" />
              </a>
            </span>
            <div class="mt-2 grid grid-cols-3 gap-2 sm:grid-cols-4">
              <button
                v-for="icon in vendorIconPresets"
                :key="icon"
                type="button"
                class="vendor-icon-preset"
                :class="normalizeIconKey(vendorForm.icon_key) === normalizeIconKey(icon) ? 'vendor-icon-preset-active' : ''"
                @click="vendorForm.icon_key = icon"
              >
                <ModelIcon :model="icon" :icon-key="icon" size="18px" />
                <span class="truncate">{{ icon }}</span>
              </button>
            </div>
          </label>
          <label>
            <span class="input-label">{{ t('admin.models.sortOrder') }}</span>
            <input v-model.number="vendorForm.sort_order" type="number" class="input" />
          </label>
          <label>
            <span class="input-label">{{ t('models.description') }}</span>
            <textarea v-model.trim="vendorForm.description" class="input min-h-20"></textarea>
          </label>
          <div class="mt-2 flex justify-end gap-2">
            <button type="button" class="btn-secondary" @click="vendorDialogOpen = false">{{ t('common.cancel') }}</button>
            <button type="submit" class="btn-primary" :disabled="vendorSaving">
              {{ vendorSaving ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </form>
      </div>
    </div>
    </Teleport>

    <Teleport to="body">
      <div
        v-if="vendorOverflowPanel && vendorOverflowVendors.length"
        class="vendor-hover-panel"
        :style="{ top: `${vendorOverflowPanel.top}px`, left: `${vendorOverflowPanel.left}px` }"
        @mouseenter="cancelHideVendorOverflow"
        @mouseleave="scheduleHideVendorOverflow"
      >
        <VendorHoverList
          :vendors="vendorOverflowVendors"
          @select="selectVendorTab(String($event.id))"
          @edit="openVendorDialogFromMenu"
          @delete="deleteVendorFromMenu"
        />
      </div>
    </Teleport>

    <Teleport to="body">
      <div
        v-if="vendorActionMenu"
        class="vendor-action-menu"
        :style="{ top: `${vendorActionMenu.top}px`, left: `${vendorActionMenu.left}px` }"
        @click.stop
      >
        <button type="button" class="vendor-action-item" @click="openVendorDialogFromMenu(vendorActionMenu.vendor)">
          <Icon name="edit" size="xs" />
          {{ t('common.edit') }}
        </button>
        <button type="button" class="vendor-action-item text-red-600 hover:text-red-700 dark:text-red-400" @click="deleteVendorFromMenu(vendorActionMenu.vendor)">
          <Icon name="trash" size="xs" />
          {{ t('common.delete') }}
        </button>
      </div>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch, type PropType } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import ModelIcon from '@/components/common/ModelIcon.vue'
import Icon from '@/components/icons/Icon.vue'
import Select from '@/components/common/Select.vue'
import Pagination from '@/components/common/Pagination.vue'
import adminModelsAPI, { type ModelCatalogRequest, type ModelVendorRequest } from '@/api/admin/models'
import type { ModelCatalog, ModelVendor } from '@/api/models'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import { getPersistedPageSize } from '@/composables/usePersistedPageSize'

const { t } = useI18n()
const appStore = useAppStore()

type CapabilityIcon = 'eye' | 'link' | 'filter' | 'brain' | 'database' | 'globe' | 'document' | 'terminal' | 'sparkles' | 'checkCircle' | 'cloud'

const loading = ref(false)
const syncing = ref(false)
const saving = ref(false)
const vendorSaving = ref(false)
const models = ref<ModelCatalog[]>([])
const vendors = ref<ModelVendor[]>([])
const filters = reactive({ search: '', platform: '', vendor_id: '', status: '' })
const pagination = reactive({ page: 1, page_size: getPersistedPageSize(), total: 0, pages: 1 })
const dialogOpen = ref(false)
const vendorDialogOpen = ref(false)
const editing = ref<ModelCatalog | null>(null)
const vendorSelection = ref('')
const platformSelection = ref<string[]>([])
const platformDropdownOpen = ref(false)
const platformDropdownRef = ref<HTMLElement | null>(null)
const vendorTabsRef = ref<HTMLElement | null>(null)
const firstVisibleVendorIndex = ref(0)
const lastVisibleVendorIndex = ref(-1)
const vendorOverflowPanel = ref<{ direction: 'left' | 'right'; top: number; left: number } | null>(null)
const vendorActionMenu = ref<{ vendor: ModelVendor; top: number; left: number } | null>(null)
const tagInput = ref('')
let searchTimer: number | undefined
let vendorScrollFrame = 0
let vendorOverflowHideTimer: number | undefined

const form = reactive<ModelCatalogRequest>({
  model_id: '',
  display_name: '',
  platform: '',
  provider: '',
  vendor_id: null,
  mode: 'chat',
  description: '',
  tags: [],
  capabilities: {},
  endpoints: [],
  pricing: {},
  metadata: {},
  status: 'active',
  visibility: 'public',
  source: 'manual',
  icon_key: ''
})

const vendorForm = reactive<ModelVendorRequest>({
  name: '',
  provider_key: '',
  icon_key: '',
  description: '',
  sort_order: 0
})

const capabilityForm = reactive<Record<string, boolean>>({})
const contextForm = reactive({
  max_input_tokens: null as number | null,
  max_output_tokens: null as number | null,
  max_tokens: null as number | null
})

const platformOptions = [
  { value: 'anthropic', label: 'Anthropic' },
  { value: 'openai', label: 'OpenAI Completions' },
  { value: 'gemini', label: 'Google Gemini' },
  { value: 'antigravity', label: 'Antigravity' },
]
const allowedModelPlatforms = new Set(platformOptions.map(option => option.value))

const vendorIconPresets = [
  'openai',
  'claude',
  'gemini',
  'deepseek',
	  'qwen',
	  'baichuan',
	  'openrouter',
  'mistral',
  'meta',
  'cohere',
  'moonshot',
  'xai',
  'doubao',
  'minimax',
  'perplexity',
  'ollama',
	  'cloudflare',
	  'XiaomiMiMo',
	]

const capabilityOptions: Array<{ key: string; labelKey: string; icon: CapabilityIcon }> = [
  { key: 'vision', labelKey: 'models.capabilityLabels.vision', icon: 'eye' },
  { key: 'function_calling', labelKey: 'models.capabilityLabels.functionCalling', icon: 'link' },
  { key: 'tool_choice', labelKey: 'models.capabilityLabels.toolChoice', icon: 'filter' },
  { key: 'reasoning', labelKey: 'models.capabilityLabels.reasoning', icon: 'brain' },
  { key: 'prompt_caching', labelKey: 'models.capabilityLabels.promptCaching', icon: 'database' },
  { key: 'web_search', labelKey: 'models.capabilityLabels.webSearch', icon: 'globe' },
  { key: 'pdf_input', labelKey: 'models.capabilityLabels.pdfInput', icon: 'document' },
  { key: 'computer_use', labelKey: 'models.capabilityLabels.computerUse', icon: 'terminal' },
  { key: 'image_output', labelKey: 'models.capabilityLabels.imageOutput', icon: 'sparkles' },
  { key: 'response_schema', labelKey: 'models.capabilityLabels.responseSchema', icon: 'checkCircle' },
  { key: 'service_tier', labelKey: 'models.capabilityLabels.serviceTier', icon: 'cloud' },
]

const selectedVendor = computed(() => vendors.value.find(v => String(v.id) === vendorSelection.value) || null)
const vendorSelectOptions = computed(() => vendors.value.map(vendor => ({ value: String(vendor.id), label: vendor.name })))
const platformSelectionLabel = computed(() => {
  if (platformSelection.value.length === 0) return t('admin.models.selectPlatforms')
  return platformSelection.value
    .map(value => platformOptions.find(option => option.value === value)?.label || value)
    .join(', ')
})
const statusFilterOptions = computed(() => [
  { value: '', label: t('admin.models.allStatuses') },
  { value: 'active', label: t('models.active') },
  { value: 'disabled', label: t('models.disabled') },
])
const modelStatusOptions = computed(() => statusFilterOptions.value.filter(option => option.value !== ''))
const visibilityOptions = computed(() => [
  { value: 'public', label: t('admin.models.publicVisibility') },
  { value: 'admin', label: t('admin.models.adminVisibility') },
])
const modeOptions = [
  { value: 'chat', label: 'chat' },
  { value: 'completion', label: 'completion' },
  { value: 'embedding', label: 'embedding' },
  { value: 'image', label: 'image' },
  { value: 'audio', label: 'audio' },
]

const VendorHoverList = defineComponent({
  name: 'VendorHoverList',
  props: {
    vendors: {
      type: Array as PropType<ModelVendor[]>,
      required: true,
    },
  },
  emits: ['select', 'edit', 'delete'],
  setup(props, { emit }) {
    return () => h('div', { class: 'max-h-80 w-72 overflow-auto py-1' }, props.vendors.map((vendor) => h('div', {
      key: vendor.id,
      class: 'flex items-center gap-2 px-2 py-2 hover:bg-gray-50 dark:hover:bg-dark-800',
    }, [
      h('button', {
        type: 'button',
        class: 'flex min-w-0 flex-1 items-center gap-2 text-left',
        onClick: () => emit('select', vendor),
      }, [
          h('span', { class: 'flex h-7 w-7 shrink-0 items-center justify-center rounded-sm border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900' }, [
          h(ModelIcon, { model: vendor.name, iconKey: vendor.icon_key || vendor.provider_key, size: '18px' }),
        ]),
        h('span', { class: 'min-w-0' }, [
          h('span', { class: 'block truncate text-sm font-medium text-gray-900 dark:text-white' }, vendor.name),
          h('span', { class: 'block truncate font-mono text-xs text-gray-500' }, vendor.icon_key || '-'),
        ]),
      ]),
      h('button', {
        type: 'button',
        class: 'btn-ghost btn-icon !p-1.5',
        'aria-label': t('common.edit'),
        onClick: () => emit('edit', vendor),
      }, [h(Icon, { name: 'edit', size: 'xs' })]),
      h('button', {
        type: 'button',
        class: 'btn-ghost btn-icon !p-1.5 text-red-600 hover:text-red-700',
        'aria-label': t('common.delete'),
        onClick: () => emit('delete', vendor),
      }, [h(Icon, { name: 'trash', size: 'xs' })]),
    ])))
  },
})

const hiddenLeftVendors = computed(() => vendors.value.slice(0, Math.max(0, firstVisibleVendorIndex.value)))
const hiddenRightVendors = computed(() => vendors.value.slice(Math.max(0, lastVisibleVendorIndex.value + 1)))
const vendorOverflowVendors = computed(() => {
  if (!vendorOverflowPanel.value) return []
  return vendorOverflowPanel.value.direction === 'left' ? hiddenLeftVendors.value : hiddenRightVendors.value
})

watch(platformSelection, (platforms) => {
  form.platform = platforms[0] || ''
  form.endpoints = [...platforms]
  if (!form.provider && form.platform) {
    form.provider = form.platform
  }
})

watch(vendors, () => {
  void nextTick(updateVendorScrollState)
})

function resetForm(model?: ModelCatalog) {
  editing.value = model || null
  form.model_id = model?.model_id || ''
  form.display_name = model?.display_name || ''
  form.platform = model?.platform || ''
  form.provider = model?.provider || model?.platform || ''
  form.vendor_id = model?.vendor_id ?? null
  form.mode = model?.mode || 'chat'
  form.description = model?.description || ''
  form.status = model?.status || 'active'
  form.visibility = model?.visibility || 'public'
  form.source = model?.source || 'manual'
  form.icon_key = ''
  platformSelection.value = uniqueStrings([...(model?.endpoints || []), normalizeModelPlatformForForm(model?.platform || '')])
    .map(normalizeModelPlatformForForm)
    .filter(platform => allowedModelPlatforms.has(platform))
  form.platform = platformSelection.value[0] || ''
  form.endpoints = [...platformSelection.value]
  tagInput.value = (model?.tags || []).join(', ')
  vendorSelection.value = model?.vendor_id ? String(model.vendor_id) : ''

  for (const capability of capabilityOptions) {
    capabilityForm[capability.key] = model?.capabilities?.[capability.key] === true
  }
  contextForm.max_input_tokens = numberFrom(model?.metadata?.max_input_tokens) || numberFrom(model?.capabilities?.context_limits?.max_input_tokens)
  contextForm.max_output_tokens = numberFrom(model?.metadata?.max_output_tokens) || numberFrom(model?.capabilities?.context_limits?.max_output_tokens)
  contextForm.max_tokens = numberFrom(model?.metadata?.max_tokens) || numberFrom(model?.capabilities?.context_limits?.max_tokens)
}

function openCreate() {
  resetForm()
  dialogOpen.value = true
}

function openEdit(model: ModelCatalog) {
  resetForm(model)
  dialogOpen.value = true
}

function openVendorDialog(vendor?: ModelVendor) {
  vendorActionMenu.value = null
  vendorOverflowPanel.value = null
  vendorForm.name = vendor?.name || ''
  vendorForm.provider_key = vendor?.provider_key || ''
  vendorForm.icon_key = vendor?.icon_key || ''
  vendorForm.description = vendor?.description || ''
  vendorForm.sort_order = vendor?.sort_order || 0
  vendorDialogOpen.value = true
}

function applySelectedVendor() {
  const vendor = selectedVendor.value
  form.vendor_id = vendor ? vendor.id : null
  if (vendor) {
    form.provider = vendor.provider_key || form.provider
  }
}

function isModelPlatformSelected(value: string) {
  return platformSelection.value.includes(value)
}

function toggleModelPlatform(value: string) {
  const index = platformSelection.value.indexOf(value)
  if (index >= 0) {
    platformSelection.value = platformSelection.value.filter(item => item !== value)
    return
  }
  platformSelection.value = [...platformSelection.value, value]
}

function selectVendorTab(vendorID: string) {
  filters.vendor_id = vendorID
  pagination.page = 1
  vendorActionMenu.value = null
  vendorOverflowPanel.value = null
  void loadModels()
}

function scrollVendorTabs(direction: 'left' | 'right') {
  if (vendors.value.length === 0) return
  const nextIndex = direction === 'left'
    ? Math.max(0, firstVisibleVendorIndex.value - 1)
    : Math.min(vendors.value.length - 1, firstVisibleVendorIndex.value + 1)
  const container = vendorTabsRef.value
  if (!container) return
  const tabStrip = container.firstElementChild as HTMLElement | null
  const target = tabStrip?.children.item(nextIndex + 1) as HTMLElement | null
  target?.scrollIntoView({ behavior: 'smooth', inline: 'start', block: 'nearest' })
  window.setTimeout(updateVendorScrollState, 260)
}

function reloadFromFirstPage() {
  pagination.page = 1
  void loadModels()
}

function toggleVendorMenu(vendor: ModelVendor, event: MouseEvent) {
  if (vendorActionMenu.value?.vendor.id === vendor.id) {
    vendorActionMenu.value = null
    return
  }
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
  vendorActionMenu.value = {
    vendor,
    top: rect.bottom + 8,
    left: clampFloatingLeft(rect.right - 128, 128)
  }
}

function openVendorDialogFromMenu(vendor: ModelVendor) {
  vendorActionMenu.value = null
  vendorOverflowPanel.value = null
  openVendorDialog(vendor)
}

function deleteVendorFromMenu(vendor: ModelVendor) {
  vendorActionMenu.value = null
  vendorOverflowPanel.value = null
  void handleDeleteVendor(vendor)
}

function showVendorOverflow(direction: 'left' | 'right', event: MouseEvent) {
  cancelHideVendorOverflow()
  const vendorsForDirection = direction === 'left' ? hiddenLeftVendors.value : hiddenRightVendors.value
  if (vendorsForDirection.length === 0) {
    vendorOverflowPanel.value = null
    return
  }
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
  vendorOverflowPanel.value = {
    direction,
    top: rect.bottom + 8,
    left: clampFloatingLeft(direction === 'left' ? rect.left : rect.right - 288, 288)
  }
}

function scheduleHideVendorOverflow() {
  window.clearTimeout(vendorOverflowHideTimer)
  vendorOverflowHideTimer = window.setTimeout(() => {
    vendorOverflowPanel.value = null
  }, 140)
}

function cancelHideVendorOverflow() {
  window.clearTimeout(vendorOverflowHideTimer)
}

function handleVendorTabsScroll() {
  vendorActionMenu.value = null
  if (vendorScrollFrame) window.cancelAnimationFrame(vendorScrollFrame)
  vendorScrollFrame = window.requestAnimationFrame(updateVendorScrollState)
}

function updateVendorScrollState() {
  vendorScrollFrame = 0
  const container = vendorTabsRef.value
  const tabStrip = container?.firstElementChild as HTMLElement | null
  if (!container || !tabStrip || vendors.value.length === 0) {
    firstVisibleVendorIndex.value = 0
    lastVisibleVendorIndex.value = -1
    return
  }

  const containerRect = container.getBoundingClientRect()
  const visibleIndexes: number[] = []
  for (let index = 0; index < vendors.value.length; index += 1) {
    const child = tabStrip.children.item(index + 1) as HTMLElement | null
    if (!child) continue
    const rect = child.getBoundingClientRect()
    if (rect.right > containerRect.left + 2 && rect.left < containerRect.right - 2) {
      visibleIndexes.push(index)
    }
  }

  if (visibleIndexes.length === 0) {
    firstVisibleVendorIndex.value = 0
    lastVisibleVendorIndex.value = -1
    return
  }
  firstVisibleVendorIndex.value = visibleIndexes[0]
  lastVisibleVendorIndex.value = visibleIndexes[visibleIndexes.length - 1]
}

function clampFloatingLeft(left: number, width: number) {
  if (typeof window === 'undefined') return left
  return Math.max(8, Math.min(left, window.innerWidth - width - 8))
}

function closeVendorFloatingMenus() {
  vendorActionMenu.value = null
}

function handleDocumentClick(event: MouseEvent) {
  closeVendorFloatingMenus()
  const target = event.target as Node
  if (platformDropdownOpen.value && platformDropdownRef.value && !platformDropdownRef.value.contains(target)) {
    platformDropdownOpen.value = false
  }
}

function debouncedLoad() {
  window.clearTimeout(searchTimer)
  searchTimer = window.setTimeout(() => {
    pagination.page = 1
    void loadModels()
  }, 250)
}

async function loadModels() {
  loading.value = true
  try {
    const res = await adminModelsAPI.list(pagination.page, pagination.page_size, {
      search: filters.search,
      platform: filters.platform,
      vendor_id: filters.vendor_id ? Number(filters.vendor_id) : undefined,
      status: filters.status,
      sort_by: 'updated_at',
      sort_order: 'desc'
    })
    models.value = res.items || []
    pagination.total = res.total
    pagination.pages = res.pages || 1
    pagination.page = res.page || pagination.page
    pagination.page_size = res.page_size || pagination.page_size
  } catch (err) {
    appStore.showError(extractApiErrorMessage(err, t('models.loadFailed')))
  } finally {
    loading.value = false
  }
}

async function loadReferenceData() {
  vendors.value = await adminModelsAPI.listVendors()
}

async function handleSync() {
  syncing.value = true
  try {
    const result = await adminModelsAPI.syncPricing()
    appStore.showSuccess(t('admin.models.syncResult', { ...result }))
    await Promise.all([loadReferenceData(), loadModels()])
  } catch (err) {
    appStore.showError(extractApiErrorMessage(err, t('admin.models.syncFailed')))
  } finally {
    syncing.value = false
  }
}

async function saveModel() {
  saving.value = true
  try {
    applySelectedVendor()
    const request = buildModelRequest()
    if (editing.value) {
      await adminModelsAPI.update(editing.value.id, request)
    } else {
      await adminModelsAPI.create(request)
    }
    dialogOpen.value = false
    await loadModels()
  } catch (err) {
    appStore.showError(extractApiErrorMessage(err, t('admin.models.saveFailed')))
  } finally {
    saving.value = false
  }
}

async function saveVendor() {
  vendorSaving.value = true
  try {
    const request: ModelVendorRequest = {
      ...vendorForm,
      provider_key: vendorForm.provider_key || normalizeVendorKey(vendorForm.name),
    }
    const vendor = await adminModelsAPI.upsertVendor(request)
    vendorDialogOpen.value = false
    await loadReferenceData()
    if (!vendorSelection.value && form.provider === vendor.provider_key) {
      vendorSelection.value = String(vendor.id)
      applySelectedVendor()
    }
  } catch (err) {
    appStore.showError(extractApiErrorMessage(err, t('admin.models.vendorSaveFailed')))
  } finally {
    vendorSaving.value = false
  }
}

async function handleDeleteVendor(vendor: ModelVendor) {
  if (!window.confirm(t('admin.models.confirmDeleteVendor', { name: vendor.name }))) return
  try {
    await adminModelsAPI.deleteVendor(vendor.id)
    if (vendorSelection.value === String(vendor.id)) {
      vendorSelection.value = ''
      applySelectedVendor()
    }
    await Promise.all([loadReferenceData(), loadModels()])
  } catch (err) {
    appStore.showError(extractApiErrorMessage(err, t('admin.models.vendorDeleteFailed')))
  }
}

async function handleDelete(model: ModelCatalog) {
  if (!window.confirm(t('admin.models.confirmDelete', { name: model.display_name || model.model_id }))) return
  try {
    await adminModelsAPI.remove(model.id)
    await loadModels()
  } catch (err) {
    appStore.showError(extractApiErrorMessage(err, t('admin.models.deleteFailed')))
  }
}

function handlePageChange(page: number) {
  pagination.page = page
  void loadModels()
}

function handlePageSizeChange(pageSize: number) {
  pagination.page_size = pageSize
  pagination.page = 1
  void loadModels()
}

function buildModelRequest(): ModelCatalogRequest {
  const capabilities: Record<string, unknown> = {}
  for (const capability of capabilityOptions) {
    capabilities[capability.key] = capabilityForm[capability.key] === true
  }
  capabilities.context_limits = compactRecord({
    max_input_tokens: contextForm.max_input_tokens,
    max_output_tokens: contextForm.max_output_tokens,
    max_tokens: contextForm.max_tokens,
  })

  const metadata = compactRecord({
    max_input_tokens: contextForm.max_input_tokens,
    max_output_tokens: contextForm.max_output_tokens,
    max_tokens: contextForm.max_tokens,
  })

  return {
    ...form,
    platform: platformSelection.value[0] || form.platform,
    vendor_id: form.vendor_id || null,
    tags: splitCSV(tagInput.value),
    capabilities,
    pricing: {},
    metadata,
    endpoints: [...platformSelection.value],
    icon_key: '',
  }
}

function compactRecord(record: Record<string, unknown>) {
  const out: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(record)) {
    if (value === '' || value === null || value === undefined) continue
    if (Array.isArray(value) && value.length === 0) continue
    out[key] = value
  }
  return out
}

function splitCSV(value: string) {
  return value.split(',').map(item => item.trim()).filter(Boolean)
}

function uniqueStrings(values: string[]) {
  return Array.from(new Set(values.map(value => value.trim()).filter(Boolean)))
}

function numberFrom(value: unknown): number | null {
  if (typeof value === 'number' && Number.isFinite(value)) return value
  if (typeof value === 'string' && value.trim() !== '') {
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : null
  }
  return null
}

function getContextLimits(model: ModelCatalog) {
  const caps = model.capabilities?.context_limits
  return {
    maxInputTokens: numberFrom(model.metadata?.max_input_tokens) || numberFrom(caps?.max_input_tokens) || null,
    maxOutputTokens: numberFrom(model.metadata?.max_output_tokens) || numberFrom(caps?.max_output_tokens) || null,
    maxTokens: numberFrom(model.metadata?.max_tokens) || numberFrom(caps?.max_tokens) || null,
  }
}

function formatTokenLimit(value: number | null) {
  if (!value || value <= 0) return '-'
  if (value >= 1_000_000) return `${Number((value / 1_000_000).toFixed(1))}M`
  if (value >= 1000) return `${Number((value / 1000).toFixed(1))}K`
  return String(value)
}

function providerLabel(value: string) {
  const option = platformOptions.find(item => item.value === value)
  return option?.label || value || '-'
}

function normalizeModelPlatformForForm(value: string) {
  const normalized = value.trim().toLowerCase()
  if (normalized === 'anthropic' || normalized === 'claude') return 'anthropic'
  if (['gemini', 'google', 'vertex_ai', 'vertex-ai', 'vertex'].includes(normalized)) return 'gemini'
  if (normalized === 'antigravity') return 'antigravity'
  if (!normalized) return ''
  return 'openai'
}

function normalizeIconKey(value?: string) {
  return (value || '').trim().toLowerCase()
}

function normalizeVendorKey(value: string) {
  return value.trim().toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '')
}

function statusClass(status: string) {
  return status === 'active'
    ? 'bg-green-50 text-green-700 dark:bg-green-900/20 dark:text-green-300'
    : 'bg-gray-100 text-gray-600 dark:bg-dark-800 dark:text-dark-300'
}

onMounted(async () => {
  document.addEventListener('click', handleDocumentClick)
  window.addEventListener('resize', updateVendorScrollState)
  await Promise.all([loadReferenceData(), loadModels()])
  await nextTick()
  updateVendorScrollState()
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleDocumentClick)
  window.removeEventListener('resize', updateVendorScrollState)
  window.clearTimeout(searchTimer)
  window.clearTimeout(vendorOverflowHideTimer)
  if (vendorScrollFrame) window.cancelAnimationFrame(vendorScrollFrame)
})
</script>

<style scoped>
.vendor-tab {
  @apply inline-flex h-10 min-w-fit items-center justify-center gap-2 rounded-sm border border-gray-200 bg-white px-3 text-sm font-medium text-gray-700 transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-300 dark:hover:text-primary-300;
}

.vendor-tab-active {
  @apply border-primary-500 bg-primary-50 text-primary-700 shadow-sm dark:bg-primary-900/20 dark:text-primary-300;
}

.vendor-icon-preset {
  @apply inline-flex min-w-0 items-center gap-2 rounded-sm border border-gray-200 bg-white px-2 py-2 text-left text-xs text-gray-600 transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-300 dark:hover:text-primary-300;
}

.vendor-icon-preset-active {
  @apply border-primary-500 bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-300;
}

.vendor-tabs {
  @apply flex h-10 min-w-max flex-nowrap items-center gap-2;
}

.vendor-tabs-scroll {
  @apply h-10 min-w-0 flex-1 overflow-x-auto overflow-y-hidden;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.vendor-tabs-scroll::-webkit-scrollbar {
  display: none;
}

.vendor-tab-action {
  @apply -mr-1 flex h-6 w-6 items-center justify-center rounded-sm text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white;
}

.vendor-action-menu {
  @apply fixed z-[100000050] w-32 overflow-hidden rounded-md border border-gray-200 bg-white py-1 shadow-xl dark:border-dark-700 dark:bg-dark-900;
}

.vendor-action-item {
  @apply flex w-full items-center gap-2 px-3 py-2 text-left text-xs text-gray-700 transition-colors hover:bg-gray-50 dark:text-dark-200 dark:hover:bg-dark-800;
}

.vendor-scroll-button {
  @apply flex h-10 w-10 items-center justify-center rounded-sm border border-gray-200 bg-white text-gray-600 transition-colors hover:border-primary-300 hover:text-primary-700 disabled:cursor-not-allowed disabled:opacity-40 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-300 dark:hover:text-primary-300;
}

.vendor-hover-panel {
  @apply fixed z-[100000050] rounded-md border border-gray-200 bg-white shadow-2xl dark:border-dark-700 dark:bg-dark-900;
}

.platform-multiselect-trigger {
  @apply flex h-10 w-full items-center justify-between gap-2 rounded-sm border border-gray-200 bg-white px-3 text-sm text-gray-900 transition-colors hover:border-gray-300 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/30 dark:border-dark-600 dark:bg-dark-800 dark:text-gray-100 dark:hover:border-dark-500;
}

.platform-multiselect-trigger-open {
  @apply border-primary-500 ring-2 ring-primary-500/30;
}

.platform-multiselect-menu {
  @apply absolute left-0 right-0 top-[calc(100%+0.25rem)] z-30 max-h-64 overflow-y-auto rounded-md border border-gray-200 bg-white py-1 shadow-lg dark:border-dark-700 dark:bg-dark-800;
}

.platform-multiselect-option {
  @apply flex w-full items-center gap-2 px-3 py-2 text-left text-sm text-gray-700 transition-colors hover:bg-gray-50 dark:text-dark-200 dark:hover:bg-dark-700;
}

.platform-multiselect-option-active {
  @apply bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-300;
}

.platform-multiselect-check {
  @apply flex h-4 w-4 shrink-0 items-center justify-center rounded-sm border border-gray-300 bg-white text-primary-600 dark:border-dark-600 dark:bg-dark-950 dark:text-primary-300;
}

.platform-multiselect-option-active .platform-multiselect-check {
  @apply border-primary-500 bg-primary-500 text-white;
}
</style>
