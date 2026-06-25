import { getLocale } from '@/i18n'
import type { LoginAgreementDocument } from '@/types'

type TranslateFn = (key: string) => string
type LoginAgreementLocale = 'zh' | 'en'

export interface LocalizedLoginAgreementDocument extends LoginAgreementDocument {
  localized_title: string
  localized_content_md: string
}

const footerLegalTitleKeys: Record<string, string> = {
  terms: 'home.footer2.links.serviceTerms',
  'usage-policy': 'home.footer2.links.usagePolicy',
  'supported-countries': 'home.footer2.links.supportedCountries',
  'supported-regions': 'home.footer2.links.supportedCountries',
  'service-specific-terms': 'home.footer2.links.serviceSpecificTerms'
}

function normalizeAgreementLocale(locale: string = getLocale()): LoginAgreementLocale {
  return locale === 'zh' ? 'zh' : 'en'
}

export function loginAgreementFooterTitleKey(doc: LoginAgreementDocument): string {
  return footerLegalTitleKeys[doc.id] || ''
}

export function localizedLoginAgreementTitle(doc: LoginAgreementDocument, locale: string = getLocale(), t?: TranslateFn): string {
  const normalizedLocale = normalizeAgreementLocale(locale)
  const localized = doc.title_i18n?.[normalizedLocale]?.trim()
  const footerTitleKey = loginAgreementFooterTitleKey(doc)
  const footerTitle = footerTitleKey && t ? t(footerTitleKey).trim() : ''
  if (normalizedLocale === 'en') {
    return localized || footerTitle || doc.title?.trim() || doc.title_i18n?.zh?.trim() || ''
  }
  return localized || doc.title?.trim() || footerTitle || doc.title_i18n?.en?.trim() || ''
}

export function localizedLoginAgreementContent(doc: LoginAgreementDocument, locale: string = getLocale()): string {
  const normalizedLocale = normalizeAgreementLocale(locale)
  const localized = doc.content_i18n?.[normalizedLocale]?.trim()
  return localized || doc.content_md?.trim() || doc.content_i18n?.zh?.trim() || doc.content_i18n?.en?.trim() || ''
}

export function localizeLoginAgreementDocument(doc: LoginAgreementDocument, locale: string = getLocale(), t?: TranslateFn): LocalizedLoginAgreementDocument {
  const localizedTitle = localizedLoginAgreementTitle(doc, locale, t)
  const localizedContent = localizedLoginAgreementContent(doc, locale)
  return {
    ...doc,
    title: localizedTitle,
    content_md: localizedContent,
    localized_title: localizedTitle,
    localized_content_md: localizedContent,
  }
}

export function localizeLoginAgreementDocuments(docs: LoginAgreementDocument[], locale: string = getLocale(), t?: TranslateFn): LocalizedLoginAgreementDocument[] {
  return docs
    .map((doc) => localizeLoginAgreementDocument(doc, locale, t))
    .filter((doc) => doc.localized_title || doc.localized_content_md)
}
