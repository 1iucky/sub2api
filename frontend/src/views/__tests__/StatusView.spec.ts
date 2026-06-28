import { flushPromises, mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import StatusView from '../StatusView.vue'
import { list as listPublicChannelMonitorViews } from '@/api/publicChannelMonitor'
import type { UserMonitorListResponse, UserMonitorView } from '@/api/channelMonitor'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string, params?: Record<string, unknown>) => {
      if (key === 'channelStatus.historyWindow') return 'Uptime over the past 90 days.'
      if (key === 'channelStatus.viewHistoricalUptime') return 'View historical uptime.'
      if (key === 'channelStatus.windowAgo.7d') return '7 days ago'
      if (key === 'channelStatus.daysAgo' && params) return `${params.days} days ago`
      if (key === 'channelStatus.today') return 'Today'
      if (key === 'channelStatus.uptimeValue' && params) return `${params.value} % uptime`
      if (key === 'channelStatus.status.operational') return 'Operational'
      if (key === 'channelStatus.status.degraded') return 'Degraded'
      if (key === 'channelStatus.status.failed') return 'Unavailable'
      if (key === 'channelStatus.status.error') return 'Error'
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

vi.mock('@/api/publicChannelMonitor', () => ({
  list: vi.fn(),
}))

vi.mock('@/components/home/PublicTopNav.vue', () => ({
  default: {
    name: 'PublicTopNav',
    template: '<nav />',
  },
}))

function createDeferred<T>() {
  let resolve!: (value: T) => void
  let reject!: (reason?: unknown) => void
  const promise = new Promise<T>((res, rej) => {
    resolve = res
    reject = rej
  })
  return { promise, resolve, reject }
}

function makeMonitor(overrides: Partial<UserMonitorView> = {}): UserMonitorView {
  return {
    id: 1,
    name: 'Claude Console',
    provider: 'anthropic',
    group_name: 'Claude',
    primary_model: 'claude-sonnet-4',
    primary_status: 'operational',
    primary_latency_ms: 180,
    primary_ping_latency_ms: 80,
    availability_7d: 99.57,
    extra_models: [],
    timeline: [
      {
        status: 'operational',
        latency_ms: 180,
        ping_latency_ms: 80,
        checked_at: '2026-06-26T00:00:00Z',
      },
      {
        status: 'degraded',
        latency_ms: 420,
        ping_latency_ms: 120,
        checked_at: '2026-06-25T00:00:00Z',
      },
    ],
    ...overrides,
  }
}

function mountView() {
  return mount(StatusView, {
    global: {
      stubs: {
        PublicTopNav: true,
        ProviderIcon: true,
        Icon: true,
      },
    },
  })
}

describe('StatusView', () => {
  beforeEach(() => {
    vi.mocked(listPublicChannelMonitorViews).mockReset()
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  it('shows initial loading skeleton while the status list is pending', async () => {
    const deferred = createDeferred<UserMonitorListResponse>()
    vi.mocked(listPublicChannelMonitorViews).mockReturnValue(deferred.promise)

    const wrapper = mountView()
    await nextTick()

    expect(wrapper.text()).not.toContain('channelStatus.empty.title')
    expect(wrapper.findAll('.animate-pulse')).toHaveLength(4)

    deferred.resolve({ items: [] })
    await flushPromises()
  })

  it('renders compact uptime rows with the 7-day availability window label', async () => {
    vi.mocked(listPublicChannelMonitorViews).mockResolvedValue({
      items: [
        makeMonitor(),
        makeMonitor({ id: 2, name: 'Claude Code', group_name: 'Claude', availability_7d: 99.93 }),
      ],
    })

    const wrapper = mountView()
    await flushPromises()

    const text = wrapper.text()
    expect(text).toContain('Claude Console')
    expect(text).toContain('Operational')
    expect(text).toContain('7 days ago')
    expect(text).toContain('99.57 % uptime')
    expect(text).not.toContain('checks ago')
    expect(text).toContain('Today')
    expect(wrapper.text()).not.toContain('Uptime over the past 90 days.')
    expect(wrapper.findAll('[data-testid="status-uptime-bar"]')).toHaveLength(4)
    expect(wrapper.findAll('[data-testid="status-monitor-card"]')).toHaveLength(2)
    expect(wrapper.findAll('[data-testid="status-monitor-group"]')).toHaveLength(0)
  })
})
