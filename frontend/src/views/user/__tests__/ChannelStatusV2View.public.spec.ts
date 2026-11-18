import { flushPromises, shallowMount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import ChannelStatusV2View from '../ChannelStatusV2View.vue'

const { getDimensions, getSnapshot, getMatrix, getModels, getErrors, getUsers } = vi.hoisted(() => ({
  getDimensions: vi.fn(),
  getSnapshot: vi.fn(),
  getMatrix: vi.fn(),
  getModels: vi.fn(),
  getErrors: vi.fn(),
  getUsers: vi.fn(),
}))

vi.mock('@/api/channelMonitorV2', async () => ({
  ...await vi.importActual<typeof import('@/api/channelMonitorV2')>('@/api/channelMonitorV2'),
  getDimensions,
  getSnapshot,
  getMatrix,
  getModels,
  getErrors,
  getUsers,
}))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => ({ isAdmin: false }),
}))
vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ cachedPublicSettings: {}, showError: vi.fn() }),
}))
vi.mock('vue-router', () => ({
  useRoute: () => ({ query: {} }),
  useRouter: () => ({ replace: vi.fn() }),
}))
vi.mock('vue-i18n', async () => ({
  ...await vi.importActual<typeof import('vue-i18n')>('vue-i18n'),
  useI18n: () => ({ t: (key: string) => key, te: () => false, locale: { value: 'en' } }),
}))

const emptyLatency = { sample_count: 0, p50_ms: null, p90_ms: null, p95_ms: null, avg_ms: null }
const snapshotFixture = {
  config: { refresh_interval_seconds: 300 },
  coverage: { data_through: '2026-09-29T08:00:00Z', coverage_complete: true },
  metrics: {
    success_requests: 0,
    error_requests: 0,
    request_count: 0,
    token_count: 0,
    rpm: 0,
    tpm: 0,
    error_rate: 0,
    cache_rate: 0,
    cache_rate_numerator: 0,
    cache_rate_denominator: 0,
    ttft: emptyLatency,
    duration: emptyLatency,
  },
  health: { overall: 'unknown', error_rate: 'unknown', ttft: 'unknown', minimum_sample: 0 },
  trend: [],
}
const matrixFixture = { coverage: {}, group_by: 'platform', items: [] }

const mountView = (isPublic: boolean) =>
  shallowMount(ChannelStatusV2View, {
    props: { isPublic },
    global: {
      stubs: {
        AppLayout: { template: '<div data-testid="app-layout"><slot /></div>' },
      },
    },
  })

describe('ChannelStatusV2View public mode', () => {
  beforeEach(() => {
    getDimensions.mockReset().mockResolvedValue({ platforms: [], groups: [], models: [] })
    getSnapshot.mockReset().mockResolvedValue(snapshotFixture)
    getMatrix.mockReset().mockResolvedValue(matrixFixture)
    getModels.mockReset().mockResolvedValue({ items: [] })
    getErrors.mockReset().mockResolvedValue({ items: [] })
    getUsers.mockReset().mockResolvedValue({ items: [] })
  })
  afterEach(() => vi.clearAllMocks())

  it('loads through the user scope by default and keeps AppLayout', async () => {
    const wrapper = mountView(false)
    await flushPromises()

    expect(getSnapshot).toHaveBeenCalledWith(expect.anything(), 'user', expect.anything())
    expect(getMatrix).toHaveBeenCalledWith(expect.anything(), expect.anything(), 'user', expect.anything())
    expect(wrapper.find('[data-testid="app-layout"]').exists()).toBe(true)
    wrapper.unmount()
  })

  it('loads through the public scope and skips AppLayout when isPublic', async () => {
    const wrapper = mountView(true)
    await flushPromises()

    expect(getDimensions).toHaveBeenCalledWith(expect.anything(), 'public', expect.anything())
    expect(getSnapshot).toHaveBeenCalledWith(expect.anything(), 'public', expect.anything())
    expect(getMatrix).toHaveBeenCalledWith(expect.anything(), expect.anything(), 'public', expect.anything())
    expect(getModels).toHaveBeenCalledWith(expect.anything(), 'public', expect.anything())
    expect(getUsers).not.toHaveBeenCalled()
    expect(wrapper.find('[data-testid="app-layout"]').exists()).toBe(false)
    wrapper.unmount()
  })
})
