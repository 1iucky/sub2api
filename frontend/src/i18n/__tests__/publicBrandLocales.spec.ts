import { describe, expect, it } from 'vitest'

import enLanding from '../locales/en/landing'
import zhLanding from '../locales/zh/landing'
import enMisc from '../locales/en/misc'
import zhMisc from '../locales/zh/misc'

describe('public navigation and legal locales', () => {
  it('defines matching Chinese and English navigation labels', () => {
    expect(enLanding.home.nav.status).toBe('Status')
    expect(enLanding.home.nav.models).toBe('Model Marketplace')
    expect(zhLanding.home.nav.status).toBe('状态')
    expect(zhLanding.home.nav.models).toBe('模型集市')
  })

  it('defines matching Chinese and English legal labels', () => {
    expect(enLanding.home.footer2.links.serviceTerms).toBe('Service Terms')
    expect(enLanding.home.footer2.links.usagePolicy).toBe('Usage Policy')
    expect(enLanding.home.footer2.links.supportedCountries).toBe('Supported Countries and Regions')
    expect(enLanding.home.footer2.links.serviceSpecificTerms).toBe('Service-Specific Terms')
    expect(zhLanding.home.footer2.links.serviceTerms).toBe('服务条款')
    expect(zhLanding.home.footer2.links.usagePolicy).toBe('使用政策')
    expect(zhLanding.home.footer2.links.supportedCountries).toBe('支持的国家和地区')
    expect(zhLanding.home.footer2.links.serviceSpecificTerms).toBe('服务特定条款')
  })

  it('defines the Factory homepage locale groups in both languages', () => {
    expect(enLanding.home.hero2.title).toBe('One gateway for every AI subscription.')
    expect(enLanding.home.marquee2.label).toBe('Compatible with')
    expect(enLanding.home.bento2.cards.pooling.title).toBe('Multi-upstream pooling')
    expect(enLanding.home.cta2.button).toBe('Start free')
    expect(enLanding.home.footer2.columns.resources).toBe('Resources')
    expect(enLanding.home.footer2.links.models).toBe('Model Marketplace')

    expect(zhLanding.home.hero2.title).toBe('一个网关，接入所有 AI 订阅。')
    expect(zhLanding.home.marquee2.label).toBe('兼容')
    expect(zhLanding.home.bento2.cards.pooling.title).toBe('多上游账号池')
    expect(zhLanding.home.cta2.button).toBe('免费开始')
    expect(zhLanding.home.footer2.columns.resources).toBe('资源')
    expect(zhLanding.home.footer2.links.models).toBe('模型集市')
  })

  it('uses the personalized brand on user-facing setup and onboarding copy', () => {
    expect(enLanding.setup.title).toBe('SiliconBase Setup')
    expect(enLanding.setup.description).toBe('Configure your SiliconBase instance')
    expect(zhLanding.setup.title).toBe('SiliconBase 安装向导')
    expect(zhLanding.setup.description).toBe('配置您的 SiliconBase 实例')

    expect(enMisc.onboarding.admin.welcome.title).toBe('👋 Welcome to SiliconBase')
    expect(enMisc.onboarding.admin.welcome.description).not.toContain('Sub2API')
    expect(enMisc.onboarding.user.welcome.title).toBe('👋 Welcome to SiliconBase')
    expect(enMisc.onboarding.user.welcome.description).not.toContain('Sub2API')
    expect(zhMisc.onboarding.admin.welcome.title).toBe('👋 欢迎使用 SiliconBase')
    expect(zhMisc.onboarding.admin.welcome.description).not.toContain('Sub2API')
    expect(zhMisc.onboarding.user.welcome.title).toBe('👋 欢迎使用 SiliconBase')
    expect(zhMisc.onboarding.user.welcome.description).not.toContain('Sub2API')
  })
})
