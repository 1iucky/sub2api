import { mount, flushPromises, type VueWrapper } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import AppHeader from '../AppHeader.vue'
import { i18n, loadLocaleMessages } from '@/i18n'
import router from '@/router'

vi.mock('@/router', async () => {
  const { createRouter, createMemoryHistory } = await import('vue-router')
  return {
    default: createRouter({
      history: createMemoryHistory(),
      routes: [{ path: '/', component: { template: '<div />' } }]
    })
  }
})

describe('AppHeader personalized control interactions', () => {
  let wrapper: VueWrapper | undefined

  beforeEach(async () => {
    localStorage.clear()
    document.documentElement.classList.remove('dark')
    document.documentElement.lang = 'en'
    i18n.global.locale.value = 'en'
    await loadLocaleMessages('en')
    await loadLocaleMessages('zh')
    // Vitest uses the runtime-only i18n build, which needs compiled messages.
    i18n.global.setLocaleMessage('en', {
      common: { toggleMenu: () => 'Toggle Menu' },
      nav: { darkMode: () => 'Dark Mode', lightMode: () => 'Light Mode' }
    })
    i18n.global.setLocaleMessage('zh', {
      common: { toggleMenu: () => '\u5207\u6362\u83dc\u5355' },
      nav: { darkMode: () => '\u6df1\u8272\u6a21\u5f0f', lightMode: () => '\u6d45\u8272\u6a21\u5f0f' }
    })
    const pinia = createPinia()
    setActivePinia(pinia)
    await router.push('/')
    await router.isReady()
    wrapper = mount(AppHeader, {
      global: { plugins: [pinia, router, i18n] }
    })
  })

  afterEach(() => {
    wrapper?.unmount()
    localStorage.clear()
    document.documentElement.classList.remove('dark')
    document.documentElement.lang = 'en'
    i18n.global.locale.value = 'en'
  })

  it('toggles the global theme from the header and persists both choices', async () => {
    await wrapper!.get('button[aria-label="Dark Mode"]').trigger('click')
    expect(document.documentElement.classList.contains('dark')).toBe(true)
    expect(localStorage.getItem('theme')).toBe('dark')

    await wrapper!.get('button[aria-label="Light Mode"]').trigger('click')
    expect(document.documentElement.classList.contains('dark')).toBe(false)
    expect(localStorage.getItem('theme')).toBe('light')
  })

  it('switches locale from the icon menu and updates header labels and persistence', async () => {
    await wrapper!.get('button[title="English"]').trigger('click')
    const chineseOption = wrapper!.findAll('button').find(button => button.text() === '\u4e2d\u6587')
    expect(chineseOption).toBeDefined()
    await chineseOption!.trigger('click')
    await flushPromises()

    expect(document.documentElement.lang).toBe('zh')
    expect(localStorage.getItem('sub2api_locale')).toBe('zh')
    expect(wrapper!.get('button[title="\u4e2d\u6587"]').find('.sr-only').text()).toBe('\u4e2d\u6587')
    expect(wrapper!.get('button[aria-label="\u6df1\u8272\u6a21\u5f0f"]').exists()).toBe(true)
    await vi.waitFor(() => {
      expect(wrapper!.findAll('button').filter(button => button.text() === '\u4e2d\u6587')).toHaveLength(1)
    })
  })
})
