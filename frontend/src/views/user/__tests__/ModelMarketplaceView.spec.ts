import { flushPromises, mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import ModelMarketplaceView from '../ModelMarketplaceView.vue'
import { listModels, listVendors } from '@/api/models'
import publicChannelMonitorAPI from '@/api/publicChannelMonitor'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string, params?: Record<string, unknown>) => {
      if (params) return `${key}:${JSON.stringify(params)}`
      return key
    },
  }),
}))

vi.mock('@/stores', () => ({
  useAppStore: () => ({
    publicSettingsLoaded: true,
    fetchPublicSettings: vi.fn(),
    showError: vi.fn(),
  }),
  useAuthStore: () => ({
    checkAuth: vi.fn(),
  }),
}))

vi.mock('@/composables/useTheme', () => ({
  useTheme: () => ({
    syncThemeFromDocument: vi.fn(),
  }),
}))

vi.mock('@/api/models', () => ({
  listModels: vi.fn(),
  listVendors: vi.fn(),
}))

vi.mock('@/api/publicChannelMonitor', () => ({
  default: {
    list: vi.fn(),
  },
}))

vi.mock('@/components/home/PublicTopNav.vue', () => ({
  default: {
    name: 'PublicTopNav',
    template: '<nav />',
  },
}))

const observerInstances: MockIntersectionObserver[] = []

class MockIntersectionObserver {
  callback: IntersectionObserverCallback
  observe = vi.fn()
  disconnect = vi.fn()
  unobserve = vi.fn()

  constructor(callback: IntersectionObserverCallback) {
    this.callback = callback
    observerInstances.push(this)
  }

  trigger(isIntersecting = true) {
    this.callback([{ isIntersecting } as IntersectionObserverEntry], this as unknown as IntersectionObserver)
  }
}

function makeModel(id: number) {
  return {
    id,
    model_id: `model-${id}`,
    display_name: `Model ${id}`,
    platform: 'anthropic',
    provider: 'anthropic',
    vendor_id: null,
    vendor: null,
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
  } as const
}

function mountView() {
  return mount(ModelMarketplaceView, {
    global: {
      stubs: {
        PublicTopNav: true,
        ModelIcon: true,
        Icon: true,
        Select: {
          name: 'Select',
          props: ['modelValue', 'options'],
          emits: ['change', 'update:modelValue'],
          template: '<div class="select-stub"></div>',
        },
        Teleport: true,
        Transition: false,
      },
    },
  })
}

describe('ModelMarketplaceView lazy loading', () => {
  beforeEach(() => {
    vi.stubGlobal('IntersectionObserver', MockIntersectionObserver)
    observerInstances.length = 0
    vi.mocked(listModels).mockReset()
    vi.mocked(listVendors).mockReset()
    vi.mocked(publicChannelMonitorAPI.list).mockReset()
    vi.mocked(listVendors).mockResolvedValue([])
    vi.mocked(publicChannelMonitorAPI.list).mockResolvedValue({ items: [] } as any)
    vi.mocked(listModels)
      .mockResolvedValueOnce({
        items: Array.from({ length: 20 }, (_, i) => makeModel(i + 1)),
        total: 40,
        page: 1,
        page_size: 20,
        pages: 2,
      })
      .mockResolvedValueOnce({
        items: Array.from({ length: 20 }, (_, i) => makeModel(i + 21)),
        total: 40,
        page: 2,
        page_size: 20,
        pages: 2,
      })
  })

  afterEach(() => {
    vi.unstubAllGlobals()
    vi.clearAllMocks()
  })

  it('loads the first marketplace page with 20 models', async () => {
    mountView()
    await flushPromises()

    expect(listModels).toHaveBeenCalledWith(
      1,
      20,
      {
        search: '',
        platform: '',
      },
      expect.objectContaining({ signal: expect.any(AbortSignal) })
    )
  })

  it('appends the next page when the lazy sentinel intersects', async () => {
    mountView()
    await flushPromises()
    await nextTick()
    await flushPromises()
    await nextTick()

    expect(observerInstances.length).toBeGreaterThan(0)

    observerInstances.at(-1)?.trigger(true)
    await flushPromises()

    expect(listModels).toHaveBeenCalledTimes(2)
    expect(listModels).toHaveBeenLastCalledWith(
      2,
      20,
      {
        search: '',
        platform: '',
      },
      expect.objectContaining({ signal: expect.any(AbortSignal) })
    )
  })
})
