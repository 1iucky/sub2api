import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const componentPath = resolve(dirname(fileURLToPath(import.meta.url)), '../LocaleSwitcher.vue')
const componentSource = readFileSync(componentPath, 'utf8')

describe('LocaleSwitcher personalized presentation', () => {
  it('uses an icon-only trigger with an accessible locale label', () => {
    expect(componentSource).toContain('name="languages"')
    expect(componentSource).toContain('<span class="sr-only">{{ currentLocale?.name }}</span>')
    expect(componentSource).not.toContain('currentLocale?.flag')
    expect(componentSource).not.toContain('currentLocale?.code.toUpperCase()')
  })

  it('keeps locale names and the current-locale checkmark in the menu', () => {
    expect(componentSource).toContain('v-for="locale in availableLocales"')
    expect(componentSource).toContain('locale.name')
    expect(componentSource).toContain('name="check"')
  })
})
