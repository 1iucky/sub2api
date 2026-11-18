import { describe, expect, it } from 'vitest'

import enLanding from '../locales/en/landing'
import zhLanding from '../locales/zh/landing'
import enMisc from '../locales/en/misc'
import zhMisc from '../locales/zh/misc'
import enAdminAccounts from '../locales/en/admin/accounts'
import enAdminOverview from '../locales/en/admin/overview'
import enAdminPlugins from '../locales/en/admin/plugins'
import enAdminSettings from '../locales/en/admin/settings'
import zhAdminAccounts from '../locales/zh/admin/accounts'
import zhAdminOverview from '../locales/zh/admin/overview'
import zhAdminPlugins from '../locales/zh/admin/plugins'
import zhAdminSettings from '../locales/zh/admin/settings'

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

  it('does not expose the upstream project brand in generic admin UI copy', () => {
    const uiCopy = [
      ['en.plugins.runtimeNotice', enAdminPlugins.plugins.runtimeNotice],
      ['en.plugins.emptyHint', enAdminPlugins.plugins.emptyHint],
      ['en.plugins.currentVersion', enAdminPlugins.plugins.currentVersion],
      ['en.plugins.confirmDisable', enAdminPlugins.plugins.confirmDisable],
      ['en.plugins.confirmUntested', enAdminPlugins.plugins.confirmUntested],
      ['en.settings.linuxdo.description', enAdminSettings.settings.linuxdo.description],
      ['en.settings.dingtalk.description', enAdminSettings.settings.dingtalk.description],
      ['en.settings.scheduling.accountSchedulingThresholdsDescription', enAdminSettings.settings.scheduling.accountSchedulingThresholdsDescription],
      ['en.settings.upstreamBillingProbe.description', enAdminSettings.settings.upstreamBillingProbe.description],
      ['en.settings.payment.easypayCustomMethodsHint', enAdminSettings.settings.payment.easypayCustomMethodsHint],
      ['en.settings.openaiFastPolicy.userIdsHint', enAdminSettings.settings.openaiFastPolicy.userIdsHint],
      ['en.accounts.upstreamBilling.trustWarning', enAdminAccounts.accounts.upstreamBilling.trustWarning],
      ['en.accounts.grok.ttsTextPlaceholder', enAdminAccounts.accounts.grok.ttsTextPlaceholder],
      ['en.accounts.grok.ttsTextDefault', enAdminAccounts.accounts.grok.ttsTextDefault],
      ['en.overview.groups.openaiLive.hint', enAdminOverview.groups.openaiLive.hint],
      ['en.overview.groups.openaiLive.unsupportedMessage', enAdminOverview.groups.openaiLive.unsupportedMessage],
      ['zh.plugins.runtimeNotice', zhAdminPlugins.plugins.runtimeNotice],
      ['zh.plugins.emptyHint', zhAdminPlugins.plugins.emptyHint],
      ['zh.plugins.currentVersion', zhAdminPlugins.plugins.currentVersion],
      ['zh.plugins.confirmDisable', zhAdminPlugins.plugins.confirmDisable],
      ['zh.plugins.confirmUntested', zhAdminPlugins.plugins.confirmUntested],
      ['zh.settings.linuxdo.description', zhAdminSettings.settings.linuxdo.description],
      ['zh.settings.dingtalk.description', zhAdminSettings.settings.dingtalk.description],
      ['zh.settings.scheduling.accountSchedulingThresholdsDescription', zhAdminSettings.settings.scheduling.accountSchedulingThresholdsDescription],
      ['zh.settings.upstreamBillingProbe.description', zhAdminSettings.settings.upstreamBillingProbe.description],
      ['zh.settings.payment.easypayCustomMethodsHint', zhAdminSettings.settings.payment.easypayCustomMethodsHint],
      ['zh.settings.openaiFastPolicy.userIdsHint', zhAdminSettings.settings.openaiFastPolicy.userIdsHint],
      ['zh.accounts.upstreamBilling.trustWarning', zhAdminAccounts.accounts.upstreamBilling.trustWarning],
      ['zh.accounts.grok.ttsTextPlaceholder', zhAdminAccounts.accounts.grok.ttsTextPlaceholder],
      ['zh.accounts.grok.ttsTextDefault', zhAdminAccounts.accounts.grok.ttsTextDefault],
      ['zh.overview.groups.openaiLive.hint', zhAdminOverview.groups.openaiLive.hint],
      ['zh.overview.groups.openaiLive.unsupportedMessage', zhAdminOverview.groups.openaiLive.unsupportedMessage]
    ] as const

    for (const [label, copy] of uiCopy) {
      expect(typeof copy, label).toBe('string')
      expect(copy, label).not.toContain('Sub2API')
    }
  })
})
