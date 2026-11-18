export const DEFAULT_SITE_NAME = 'SiliconBase'
export const UPSTREAM_DEFAULT_SITE_NAME = 'Sub2API'

export function normalizeSiteName(siteName: unknown): string {
  if (typeof siteName !== 'string') return DEFAULT_SITE_NAME
  const normalized = siteName.trim()
  if (normalized === UPSTREAM_DEFAULT_SITE_NAME) return DEFAULT_SITE_NAME
  return normalized || DEFAULT_SITE_NAME
}
