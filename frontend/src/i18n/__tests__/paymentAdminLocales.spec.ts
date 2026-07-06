import { describe, expect, it } from 'vitest'

import en from '../locales/en'
import zh from '../locales/zh'

describe('payment admin locale keys', () => {
  it('contains zh subscription plan request limit labels under payment admin scope', () => {
    expect(zh.payment.admin.dailyRequestLimit).toBe('每日请求限额')
    expect(zh.payment.admin.weeklyRequestLimit).toBe('每周请求限额')
    expect(zh.payment.admin.monthlyRequestLimit).toBe('每月请求限额')
    expect(zh.payment.admin.requestLimitPlaceholder).toBe('留空=不限制')
  })

  it('keeps en subscription plan request limit labels under payment admin scope', () => {
    expect(en.payment.admin.dailyRequestLimit).toBe('Daily Request Limit')
    expect(en.payment.admin.weeklyRequestLimit).toBe('Weekly Request Limit')
    expect(en.payment.admin.monthlyRequestLimit).toBe('Monthly Request Limit')
    expect(en.payment.admin.requestLimitPlaceholder).toBe('Empty = unlimited')
  })
})
