/**
 * formatScaled formats a per-token (or per-request) display price scaled by `scale`.
 *
 *   formatScaled(0.000003, 1_000_000) → "$3"        // per 1M tokens
 *   formatScaled(0.5,        1)        → "$0.5"      // per request
 *   formatScaled(null,       1_000_000) → "-"
 *
 * Uses toPrecision(10) then strips trailing zeros to avoid IEEE 754 display noise.
 */
export type DisplayCurrency = 'USD' | 'CNY'

export function normalizeDisplayCurrency(value: unknown): DisplayCurrency {
  return value === 'CNY' ? 'CNY' : 'USD'
}

export function currencySymbol(value: unknown): string {
  return normalizeDisplayCurrency(value) === 'CNY' ? '￥' : '$'
}

export function formatScaled(value: number | null, scale: number, displayCurrency: unknown = 'USD'): string {
  if (value == null) return '-'
  return `${currencySymbol(displayCurrency)}${(value * scale).toPrecision(10).replace(/\.?0+$/, '')}`
}
