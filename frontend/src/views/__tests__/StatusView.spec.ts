import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import StatusView from '../StatusView.vue'

const { isChannelMonitorV1Mode } = vi.hoisted(() => ({
  isChannelMonitorV1Mode: vi.fn(),
}))

vi.mock('@/stores', () => ({
  useAppStore: () => ({
    publicSettingsLoaded: true,
    fetchPublicSettings: vi.fn(),
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

vi.mock('@/utils/featureFlags', () => ({
  isChannelMonitorV1Mode,
}))

vi.mock('@/components/home/PublicTopNav.vue', () => ({
  default: { name: 'PublicTopNav', template: '<nav data-testid="public-top-nav" />' },
}))

vi.mock('@/views/user/ChannelStatusV1View.vue', () => ({
  default: {
    name: 'ChannelStatusV1View',
    props: { isPublic: Boolean },
    template: '<div data-testid="channel-status-v1" :data-is-public="String(isPublic)" />',
  },
}))

vi.mock('@/views/user/ChannelStatusV2View.vue', () => ({
  default: {
    name: 'ChannelStatusV2View',
    props: { isPublic: Boolean },
    template: '<div data-testid="channel-status-v2" :data-is-public="String(isPublic)" />',
  },
}))

describe('StatusView', () => {
  beforeEach(() => {
    isChannelMonitorV1Mode.mockReset()
  })

  it('renders the console channel-status v1 view in public mode when monitor mode is v1', () => {
    isChannelMonitorV1Mode.mockReturnValue(true)
    const wrapper = mount(StatusView)

    expect(wrapper.find('[data-testid="public-top-nav"]').exists()).toBe(true)
    const inner = wrapper.find('[data-testid="channel-status-v1"]')
    expect(inner.exists()).toBe(true)
    expect(inner.attributes('data-is-public')).toBe('true')
    expect(wrapper.find('[data-testid="channel-status-v2"]').exists()).toBe(false)
  })

  it('renders the console channel-status v2 view in public mode when monitor mode is v2', () => {
    isChannelMonitorV1Mode.mockReturnValue(false)
    const wrapper = mount(StatusView)

    expect(wrapper.find('[data-testid="public-top-nav"]').exists()).toBe(true)
    const inner = wrapper.find('[data-testid="channel-status-v2"]')
    expect(inner.exists()).toBe(true)
    expect(inner.attributes('data-is-public')).toBe('true')
    expect(wrapper.find('[data-testid="channel-status-v1"]').exists()).toBe(false)
  })
})
