import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import AnnouncementBell from '../AnnouncementBell.vue'
import { useAnnouncementStore } from '@/stores/announcements'

vi.mock('@/api', () => ({
  announcementsAPI: {
    list: vi.fn(),
    markRead: vi.fn().mockResolvedValue({ message: 'ok' })
  }
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError: vi.fn(),
    showSuccess: vi.fn()
  })
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key
  })
}))

vi.mock('marked', () => ({
  marked: {
    setOptions: vi.fn(),
    parse: (value: string) => value
  }
}))

vi.mock('dompurify', () => ({
  default: {
    sanitize: (value: string) => value
  }
}))

vi.mock('@/utils/format', () => ({
  formatRelativeTime: (value: string) => value,
  formatRelativeWithDateTime: (value: string) => value
}))

const mountBell = () => mount(AnnouncementBell, {
  global: {
    stubs: {
      Icon: true,
      Teleport: true,
      Transition: false
    }
  }
})

describe('AnnouncementBell drawer', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    document.body.innerHTML = ''
  })

  it('opens a right-side categorized timeline drawer with per-tab counts', async () => {
    const store = useAnnouncementStore()
    store.announcements = [
      {
        id: 1,
        title: 'System notice',
        content: 'System content '.repeat(40),
        notify_mode: 'silent',
        category: 'announcement',
        created_at: '2026-07-07T17:03:29Z',
        updated_at: '2026-07-07T17:03:29Z'
      },
      {
        id: 2,
        title: 'New model',
        content: 'Model content',
        notify_mode: 'silent',
        category: 'model_update',
        created_at: '2026-07-06T14:44:45Z',
        updated_at: '2026-07-06T14:44:45Z'
      },
      {
        id: 3,
        title: 'Release notes',
        content: 'Release content',
        notify_mode: 'silent',
        category: 'changelog',
        created_at: '2026-07-02T20:39:50Z',
        updated_at: '2026-07-02T20:39:50Z',
        read_at: '2026-07-03T00:00:00Z'
      },
      {
        id: 4,
        title: 'Another model',
        content: 'More model content',
        notify_mode: 'silent',
        category: 'model_update',
        created_at: '2026-07-01T20:39:50Z',
        updated_at: '2026-07-01T20:39:50Z'
      }
    ]

    const wrapper = mountBell()
    await wrapper.get('button[aria-label="announcements.title"]').trigger('click')
    await flushPromises()

    expect(wrapper.get('[data-testid="announcement-drawer"]').classes()).toContain('right-0')
    expect(wrapper.text()).toContain('announcements.categories.announcement')
    expect(wrapper.text()).toContain('announcements.categories.modelUpdate')
    expect(wrapper.text()).toContain('announcements.categories.changelog')
    expect(wrapper.get('[data-testid="announcement-tab-announcement"]').text()).toContain('1')
    expect(wrapper.get('[data-testid="announcement-tab-model_update"]').text()).toContain('2')
    expect(wrapper.get('[data-testid="announcement-tab-changelog"]').text()).toContain('1')
    expect(wrapper.find('[data-testid="announcement-tab-all"]').exists()).toBe(false)
    expect(wrapper.get('[data-testid="announcement-timeline"]').exists()).toBe(true)
  })

  it('collapses long markdown announcements by default and expands them on demand', async () => {
    const store = useAnnouncementStore()
    store.announcements = [
      {
        id: 11,
        title: 'Long markdown notice',
        content: '# Heading\n\n' + 'Long markdown paragraph. '.repeat(80),
        notify_mode: 'silent',
        category: 'announcement',
        created_at: '2026-07-07T17:03:29Z',
        updated_at: '2026-07-07T17:03:29Z'
      }
    ]

    const wrapper = mountBell()
    await wrapper.get('button[aria-label="announcements.title"]').trigger('click')
    await flushPromises()

    expect(wrapper.get('[data-testid="announcement-content-11"]').classes()).toContain('max-h-40')
    expect(wrapper.get('[data-testid="announcement-toggle-11"]').text()).toContain('common.expand')

    await wrapper.get('[data-testid="announcement-toggle-11"]').trigger('click')

    expect(wrapper.get('[data-testid="announcement-content-11"]').classes()).not.toContain('max-h-40')
    expect(wrapper.get('[data-testid="announcement-toggle-11"]').text()).toContain('common.collapse')
  })
})
