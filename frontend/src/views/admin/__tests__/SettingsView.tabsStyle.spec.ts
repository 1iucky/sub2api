import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const sourcePath = resolve(dirname(fileURLToPath(import.meta.url)), '../SettingsView.vue')
const source = readFileSync(sourcePath, 'utf8')

function extractRule(selector: string) {
  const escapedSelector = selector.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const match = source.match(new RegExp(`${escapedSelector}\\s*\\{[\\s\\S]*?\\n\\}`))
  return match?.[0] ?? ''
}

describe('SettingsView personalized tab styles', () => {
  it('uses warm personalized dark surfaces for the settings tabs shell', () => {
    const darkShellRule = extractRule('.dark .settings-tabs-shell')
    const darkHoverRule = extractRule('.dark .settings-tab::before')
    const darkTabStyles = `${darkShellRule}\n${darkHoverRule}`

    expect(darkShellRule).toContain('rgb(10 9 8 / 0.86)')
    expect(darkTabStyles).not.toContain('rgb(15 23 42')
    expect(darkTabStyles).not.toContain('rgb(30 41 59')
    expect(darkTabStyles).not.toContain('rgb(51 65 85')
  })

  it('uses the personalized primary accent instead of the old cyan-blue underline', () => {
    const activeIndicatorRule = extractRule('.settings-tab-active::after')

    expect(activeIndicatorRule).toContain('#ef6f2e')
    expect(activeIndicatorRule).not.toContain('#14b8a6')
    expect(activeIndicatorRule).not.toContain('#0ea5e9')
  })
})
