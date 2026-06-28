import { beforeEach, describe, expect, it, vi } from 'vitest'

const authStore = vi.hoisted(() => ({
  checkAuth: vi.fn(),
  isAuthenticated: false,
  isAdmin: false,
  isSimpleMode: false,
}))

const appStore = vi.hoisted(() => ({
  siteName: 'Sub2API',
  backendModeEnabled: false,
  cachedPublicSettings: null as null | Record<string, unknown>,
}))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => authStore,
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => appStore,
}))

vi.mock('@/stores/adminSettings', () => ({
  useAdminSettingsStore: () => ({
    customMenuItems: [],
  }),
}))

vi.mock('@/stores/adminCompliance', () => ({
  useAdminComplianceStore: () => ({
    fetchStatus: vi.fn(),
  }),
}))

vi.mock('@/composables/useNavigationLoading', () => ({
  useNavigationLoadingState: () => ({
    startNavigation: vi.fn(),
    endNavigation: vi.fn(),
    isLoading: { value: false },
  }),
}))

vi.mock('@/composables/useRoutePrefetch', () => ({
  useRoutePrefetch: () => ({
    triggerPrefetch: vi.fn(),
    cancelPendingPrefetch: vi.fn(),
    resetPrefetchState: vi.fn(),
  }),
}))

vi.mock('@/api/setup', () => ({
  getSetupStatus: vi.fn(),
}))

vi.mock('@/router/title', () => ({
  resolveDocumentTitle: vi.fn((title: unknown) => String(title || 'Sub2API')),
}))

describe('unsupported region route', () => {
  beforeEach(() => {
    appStore.backendModeEnabled = false
  })

  it('registers the unsupported-region page as a public route', async () => {
    const { default: router } = await import('@/router')
    const route = router.getRoutes().find((record) => record.name === 'UnsupportedRegion')

    expect(route?.path).toBe('/unsupported-region')
    expect(route?.meta.requiresAuth).toBe(false)
    expect(route?.meta.title).toBe('Unsupported Region')
  })

  it('keeps the unsupported-region page reachable when backend mode blocks public pages', async () => {
    appStore.backendModeEnabled = true

    const { default: router } = await import('@/router')
    await router.push('/unsupported-region?ip=203.0.113.10&country=China&region=Guangdong')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/unsupported-region')
  })
})
