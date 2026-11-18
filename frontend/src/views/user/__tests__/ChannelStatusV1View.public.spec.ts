import { flushPromises, shallowMount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import ChannelStatusV1View from '../ChannelStatusV1View.vue'

const { listAuthed, detailAuthed, listPublic, detailPublic } = vi.hoisted(() => ({
  listAuthed: vi.fn(),
  detailAuthed: vi.fn(),
  listPublic: vi.fn(),
  detailPublic: vi.fn(),
}))
vi.mock('@/api/channelMonitor', () => ({ list: listAuthed, status: detailAuthed }))
vi.mock('@/api/publicChannelMonitor', () => ({ list: listPublic, status: detailPublic }))
vi.mock('@/stores/app', () => ({ useAppStore: () => ({ cachedPublicSettings: { channel_monitor_enabled: true }, showError: vi.fn() }) }))
vi.mock('vue-i18n', async () => ({
  ...await vi.importActual<typeof import('vue-i18n')>('vue-i18n'),
  useI18n: () => ({ t: (key: string) => key }),
}))

const mountView = (isPublic: boolean) => shallowMount(ChannelStatusV1View, {
  props: { isPublic },
  global: { stubs: {
    AppLayout: { template: '<div data-testid="app-layout"><slot /></div>' },
    MonitorHero: true,
    MonitorCardGrid: true,
    MonitorDetailDialog: true,
  } },
})

describe('ChannelStatusV1View public mode', () => {
  beforeEach(() => {
    localStorage.clear()
    listAuthed.mockReset().mockResolvedValue({ items: [] })
    listPublic.mockReset().mockResolvedValue({ items: [] })
  })
  afterEach(() => { vi.clearAllMocks(); localStorage.clear() })

  it('loads monitors through the authenticated API by default and keeps AppLayout', async () => {
    const wrapper = mountView(false)
    await flushPromises()

    expect(listAuthed).toHaveBeenCalledTimes(1)
    expect(listPublic).not.toHaveBeenCalled()
    expect(wrapper.find('[data-testid="app-layout"]').exists()).toBe(true)
    wrapper.unmount()
  })

  it('loads monitors through the public API and skips AppLayout when isPublic', async () => {
    const wrapper = mountView(true)
    await flushPromises()

    expect(listPublic).toHaveBeenCalledTimes(1)
    expect(listAuthed).not.toHaveBeenCalled()
    expect(wrapper.find('[data-testid="app-layout"]').exists()).toBe(false)
    wrapper.unmount()
  })
})
