import { describe, expect, it } from 'vitest'

import { localizedLoginAgreementTitle } from '@/utils/loginAgreement'
import enLanding from '../locales/en/landing'
import zhLanding from '../locales/zh/landing'

describe('public navigation and legal locales', () => {
  it('defines matching Chinese and English labels', () => {
    expect(enLanding.home.nav.status).toBe('Status')
    expect(enLanding.home.nav.models).toBe('Model Marketplace')
    expect(zhLanding.home.nav.status).toBe('状态')
    expect(zhLanding.home.nav.models).toBe('模型集市')
    expect(enLanding.home.footer2.links.serviceTerms).toBe('Service Terms')
    expect(enLanding.home.footer2.links.usagePolicy).toBe('Usage Policy')
    expect(enLanding.home.footer2.links.supportedCountries).toBe('Supported Countries and Regions')
    expect(enLanding.home.footer2.links.serviceSpecificTerms).toBe('Service-Specific Terms')
    expect(zhLanding.home.footer2.links.serviceTerms).toBe('服务条款')
    expect(zhLanding.home.footer2.links.usagePolicy).toBe('使用政策')
    expect(zhLanding.home.footer2.links.supportedCountries).toBe('支持的国家和地区')
    expect(zhLanding.home.footer2.links.serviceSpecificTerms).toBe('服务特定条款')
  })

  it('localizes login agreement document titles in English', () => {
    const translate = (key: string) => {
      const labels: Record<string, string> = {
        'home.footer2.links.serviceTerms': enLanding.home.footer2.links.serviceTerms,
      }
      return labels[key] || key
    }

    expect(localizedLoginAgreementTitle(
      { id: 'terms', title: '服务条款', content_md: '' },
      'en',
      translate,
    )).toBe('Service Terms')
  })
})
