import { flushPromises, mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import PricingEntryCard from '../PricingEntryCard.vue'
import adminModelsAPI from '@/api/admin/models'
import channelsAPI from '@/api/admin/channels'
import type { PricingFormEntry } from '../types'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (_key: string, _paramsOrFallback?: unknown, fallback?: string) => {
      if (typeof _paramsOrFallback === 'string') return _paramsOrFallback
      return fallback || _key
    },
  }),
}))

vi.mock('@/api/admin/models', () => ({
  default: {
    list: vi.fn(),
  },
}))

vi.mock('@/api/admin/channels', () => ({
  default: {
    getModelDefaultPricing: vi.fn(),
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
    'filterOptions',
  ],
  emits: ['search', 'change'],
  template: '<div data-testid="model-select"></div>',
}

function makeEntry(overrides: Partial<PricingFormEntry> = {}): PricingFormEntry {
  return {
    models: [],
    billing_mode: 'token',
    input_price: null,
    output_price: null,
    cache_write_price: null,
    cache_read_price: null,
    image_output_price: null,
    per_request_price: null,
    intervals: [],
    ...overrides,
  }
}

function mountCard(entry = makeEntry()) {
  return mount(PricingEntryCard, {
    props: {
      entry,
      platform: 'anthropic',
    },
    global: {
      stubs: {
        Select: selectStub,
        Icon: true,
        IntervalRow: true,
        ModelTagInput: true,
        Teleport: true,
      },
    },
  })
}

describe('PricingEntryCard model catalog search', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    vi.mocked(adminModelsAPI.list).mockResolvedValue({
      items: [
        {
          id: 1,
          model_id: 'deepseek-v4-flash',
          display_name: 'DeepSeek V4 Flash',
          platform: 'deepseek',
          provider: 'deepseek',
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
    vi.mocked(channelsAPI.getModelDefaultPricing).mockResolvedValue({ found: false })
  })

  afterEach(() => {
    vi.useRealTimers()
    vi.clearAllMocks()
  })

  it('does not preload models until the user types a keyword', () => {
    mountCard()

    expect(adminModelsAPI.list).not.toHaveBeenCalled()
  })

  it('remote-searches models by keyword with a bounded page size', async () => {
    const wrapper = mountCard()

    wrapper.findComponent(selectStub).vm.$emit('search', ' deepseek ')
    await vi.advanceTimersByTimeAsync(300)
    await flushPromises()

    expect(adminModelsAPI.list).toHaveBeenCalledTimes(1)
    expect(adminModelsAPI.list).toHaveBeenCalledWith(
      1,
      30,
      {
        search: 'deepseek',
        status: 'active',
        visibility: 'public',
        sort_by: 'model_id',
        sort_order: 'asc',
      },
      expect.objectContaining({ signal: expect.any(AbortSignal) })
    )
  })
})
