import { describe, expect, it } from 'vitest'
import { currencySymbol, formatScaled } from '@/utils/pricing'
import { formatTokenPricePerMillion } from '@/utils/usagePricing'

describe('pricing currency display', () => {
  it('maps display currency to symbols and defaults to USD', () => {
    expect(currencySymbol('USD')).toBe('$')
    expect(currencySymbol('CNY')).toBe('￥')
    expect(currencySymbol(null)).toBe('$')
  })

  it('formats configured model prices with the requested currency symbol', () => {
    expect(formatScaled(0.000003, 1_000_000, 'CNY')).toBe('￥3')
    expect(formatScaled(0.5, 1, 'USD')).toBe('$0.5')
  })

  it('formats usage unit prices with the usage log snapshot currency', () => {
    expect(formatTokenPricePerMillion(0.000003, 1, { displayCurrency: 'CNY' })).toBe('￥3.0000')
    expect(formatTokenPricePerMillion(0.000003, 1, { displayCurrency: 'USD' })).toBe('$3.0000')
  })
})
