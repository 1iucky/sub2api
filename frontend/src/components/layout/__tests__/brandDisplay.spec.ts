import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const root = resolve(dirname(fileURLToPath(import.meta.url)), '../../..')
const source = (path: string) => readFileSync(resolve(root, path), 'utf8')

describe('shared site brand display', () => {
  it.each([
    'components/layout/AuthLayout.vue',
    'views/HomeView.vue',
    'components/home/PublicTopNav.vue',
    'components/home/PublicFooter.vue',
  ])('%s uses the configured site name with a Sub2API fallback', (path) => {
    const component = source(path)

    expect(component).toContain("'Sub2API'")
    expect(component).not.toContain('>SiliconBase<')
    expect(component).not.toContain("'SiliconBase'")
  })
})
