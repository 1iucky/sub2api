import { flushPromises, mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import MonitorFormDialog from '../MonitorFormDialog.vue'
import adminModelsAPI from '@/api/admin/models'
import { adminAPI } from '@/api/admin'
import type { ChannelMonitorTemplate } from '@/api/admin/channelMonitorTemplate'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (_key: string, _paramsOrFallback?: unknown, fallback?: string) => {
      if (typeof _paramsOrFallback === 'string') return _paramsOrFallback
      return fallback || _key
    },
  }),
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    cachedPublicSettings: { channel_monitor_default_interval_seconds: 60 },
    showError: vi.fn(),
    showSuccess: vi.fn(),
  }),
}))

vi.mock('@/api/admin/models', () => ({
  default: {
    list: vi.fn(),
  },
}))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    channelMonitorTemplate: {
      list: vi.fn(),
    },
    channelMonitor: {
      create: vi.fn(),
      update: vi.fn(),
    },
  },
}))

vi.mock('@/api/keys', () => ({
  keysAPI: {
    list: vi.fn(),
  },
}))

vi.mock('@/api/groups', () => ({
  userGroupsAPI: {
    getUserGroupRates: vi.fn(),
  },
}))

const selectStub = {
  name: 'Select',
  props: [
    'modelValue',
    'options',
    'placeholder',
    'searchPlaceholder',
    'emptyText',
    'loading',
    'loadingText',
    'searchable',
    'clearable',
    'creatable',
    'creatablePrefix',
    'filterOptions',
  ],
  emits: ['search', 'change', 'update:modelValue'],
  template: '<div class="select-stub"></div>',
}

const emptyTemplateResponse: { items: ChannelMonitorTemplate[] } = { items: [] }

function mountDialog() {
  return mount(MonitorFormDialog, {
    props: {
      show: true,
      monitor: null,
    },
    global: {
      stubs: {
        BaseDialog: {
          name: 'BaseDialog',
          template: '<div><slot /><slot name="footer" /></div>',
        },
        Select: selectStub,
        Toggle: true,
        ModelTagInput: true,
        MonitorKeyPickerDialog: true,
        MonitorAdvancedRequestConfig: true,
        ProviderIcon: true,
      },
    },
  })
}

describe('MonitorFormDialog model catalog search', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    vi.mocked(adminModelsAPI.list).mockResolvedValue({
      items: [
        {
          id: 1,
          model_id: 'claude-sonnet-4-5',
          display_name: 'Claude Sonnet 4.5',
          platform: 'anthropic',
          provider: 'anthropic',
          vendor_id: null,
          vendor: null,
          mode: '',
          description: '',
          tags: [],
          capabilities: {},
          endpoints: [],
          pricing: {},
          metadata: {},
          status: 'active',
          visibility: 'public',
          source: 'manual',
          icon_key: '',
          last_synced_at: null,
          related_pricing: {
            channel_count: 0,
            channels: [],
            entries: [],
            groups: [],
            has_intervals: false,
          },
          related_groups: [],
          created_at: '',
          updated_at: '',
        },
      ],
      total: 1,
      page: 1,
      page_size: 30,
      pages: 1,
    })
    vi.mocked(adminAPI.channelMonitorTemplate.list).mockResolvedValue(emptyTemplateResponse)
  })

  afterEach(() => {
    vi.useRealTimers()
    vi.clearAllMocks()
  })

  it('does not preload model catalog entries when opened', async () => {
    mountDialog()
    await flushPromises()

    expect(adminModelsAPI.list).not.toHaveBeenCalled()
  })

  it('remote-searches primary model options after keyword input', async () => {
    const wrapper = mountDialog()
    await flushPromises()

    wrapper.findAllComponents(selectStub)[0].vm.$emit('search', ' sonnet ')
    await vi.advanceTimersByTimeAsync(300)
    await flushPromises()

    expect(adminModelsAPI.list).toHaveBeenCalledTimes(1)
    expect(adminModelsAPI.list).toHaveBeenCalledWith(
      1,
      30,
      {
        search: 'sonnet',
        status: 'active',
        visibility: 'public',
        sort_by: 'model_id',
        sort_order: 'asc',
      },
      expect.objectContaining({ signal: expect.any(AbortSignal) })
    )
  })
})
