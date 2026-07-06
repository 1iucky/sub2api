/**
 * formatScaled formats a per-token (or per-request) display price scaled by `scale`.
 *
 *   formatScaled(0.000003, 1_000_000)    → "$3"      // per 1M tokens
 *   formatScaled(0.5,        1)          → "$0.5"    // per request
 *   formatScaled(null,       1_000_000)  → "-"
 *   formatScaled(0.000003, 1_000_000, 2) → "$3.00"   // pad to ≥2 decimals
 *   formatScaled(1.25e-8,  1_000_000, 2) → "$0.0125" // longer decimals kept as-is
 *
 * Uses toPrecision(10) then strips trailing zeros to avoid IEEE 754 display noise.
 * `minFractionDigits` pads the result back up to a minimum number of decimals.
 */
export type DisplayCurrency = 'USD' | 'CNY'

export function normalizeDisplayCurrency(value: unknown): DisplayCurrency {
  return value === 'CNY' ? 'CNY' : 'USD'
}

export function currencySymbol(value: unknown): string {
  return normalizeDisplayCurrency(value) === 'CNY' ? '￥' : '$'
}

export function formatScaled(
  value: number | null,
  scale: number,
  minFractionDigitsOrCurrency: number | unknown = 0,
  displayCurrency: unknown = 'USD'
): string {
  if (value == null) return '-'
  const minFractionDigits = typeof minFractionDigitsOrCurrency === 'number'
    ? minFractionDigitsOrCurrency
    : 0
  const resolvedCurrency = typeof minFractionDigitsOrCurrency === 'number'
    ? displayCurrency
    : minFractionDigitsOrCurrency
  let s = (value * scale).toPrecision(10).replace(/\.?0+$/, '')
  if (minFractionDigits > 0 && !s.includes('e')) {
    const dot = s.indexOf('.')
    const digits = dot === -1 ? 0 : s.length - dot - 1
    if (digits < minFractionDigits) {
      s = (dot === -1 ? `${s}.` : s) + '0'.repeat(minFractionDigits - digits)
    }
  }
  return `${currencySymbol(resolvedCurrency)}${s}`
}
