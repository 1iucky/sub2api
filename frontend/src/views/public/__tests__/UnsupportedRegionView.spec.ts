import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import UnsupportedRegionView from '../UnsupportedRegionView.vue'

const replace = vi.fn()
const fetchMock = vi.fn()
const fetchPublicSettings = vi.fn()
let publicSettings = {
  site_name: 'AI Wanwu API',
  site_logo: '/logo.png',
}
let publicSettingsLoaded = true
let storeSiteLogo = ''

vi.mock('vue-router', () => ({
  useRoute: () => ({
    query: {
      ip: '203.0.113.10',
      country: 'China',
      region: 'Guangdong',
      path: '/usage',
    },
  }),
  useRouter: () => ({
    replace,
  }),
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string, params?: Record<string, unknown>) => {
      if (key === 'unsupportedRegion.accessInfo') {
        return `Access: ${params?.host} · ${params?.location} · ${params?.path}`
      }
      const messages: Record<string, string> = {
        'unsupportedRegion.subtitle': 'Overseas API service platform',
        'unsupportedRegion.title': 'This region is not supported',
        'unsupportedRegion.description': 'The service is unavailable in this region.',
        'unsupportedRegion.retryHint': 'If you think this is incorrect, retry later.',
        'unsupportedRegion.retry': 'Retry check',
      }
      return messages[key] ?? key
    },
  }),
}))

vi.mock('@/stores', () => ({
  useAppStore: () => ({
    cachedPublicSettings: publicSettings,
    siteName: 'Fallback',
    siteLogo: storeSiteLogo,
    publicSettingsLoaded,
    fetchPublicSettings,
  }),
}))

vi.mock('@/composables/useTheme', () => ({
  useTheme: () => ({
    syncThemeFromDocument: vi.fn(),
  }),
}))

describe('UnsupportedRegionView', () => {
  beforeEach(() => {
    replace.mockReset()
    fetchMock.mockReset()
    fetchPublicSettings.mockReset()
    fetchMock.mockResolvedValue({ status: 403 })
    vi.stubGlobal('fetch', fetchMock)
    publicSettings = {
      site_name: 'AI Wanwu API',
      site_logo: '/logo.png',
    }
    publicSettingsLoaded = true
    storeSiteLogo = ''
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: {
        host: 'example.test',
      },
    })
  })

  it('renders site branding and source IP attribution from query params', () => {
    const wrapper = mount(UnsupportedRegionView)

    expect(wrapper.text()).toContain('AI Wanwu API')
    expect(wrapper.text()).toContain('This region is not supported')
    expect(wrapper.text()).toContain('203.0.113.10')
    expect(wrapper.text()).toContain('China / Guangdong')
    expect(wrapper.text()).toContain('/usage')
    expect(wrapper.find('img').attributes('src')).toBe('/logo.png')
  })

  it('uses the same default logo as the public homepage when no custom logo is configured', () => {
    publicSettings = {
      site_name: 'AI Wanwu API',
      site_logo: '',
    }

    const wrapper = mount(UnsupportedRegionView)

    expect(wrapper.find('img').attributes('src')).toBe('/logo.svg')
    expect(wrapper.text()).not.toContain('AI API')
    expect(wrapper.text()).not.toContain('If you think this is incorrect')
  })

  it('updates the logo after public settings are fetched asynchronously', async () => {
    publicSettingsLoaded = false
    publicSettings = {
      site_name: 'AI Wanwu API',
      site_logo: '',
    }
    fetchPublicSettings.mockImplementation(async () => {
      publicSettings = {
        site_name: 'AI Wanwu API',
        site_logo: '/uploaded-logo.png',
      }
      return publicSettings
    })

    const wrapper = mount(UnsupportedRegionView)
    await flushPromises()

    expect(fetchPublicSettings).toHaveBeenCalledTimes(1)
    expect(wrapper.find('img').attributes('src')).toBe('/uploaded-logo.png')
  })

  it('checks the server-side region gate before returning to home', async () => {
    fetchMock.mockResolvedValueOnce({ status: 204 })
    const wrapper = mount(UnsupportedRegionView)

    await wrapper.find('button').trigger('click')

    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(fetchMock.mock.calls[0][0]).toMatch(/^\/region-check\?ts=\d+$/)
    expect(fetchMock.mock.calls[0][1]).toMatchObject({
      method: 'GET',
      cache: 'no-store',
    })
    expect(replace).toHaveBeenCalledWith('/home')
  })

  it('stays on the unsupported-region page when the server-side region gate still rejects', async () => {
    const wrapper = mount(UnsupportedRegionView)

    await wrapper.find('button').trigger('click')

    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(replace).not.toHaveBeenCalled()
  })
})
