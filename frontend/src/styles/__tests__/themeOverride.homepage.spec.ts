import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const stylePath = resolve(dirname(fileURLToPath(import.meta.url)), '../theme-override.css')
const styleSource = readFileSync(stylePath, 'utf8')

describe('homepage factory visual styles', () => {
  it('defines the personalized homepage classes used by HomeView', () => {
    const requiredSelectors = [
      '.bg-factory-surface-grid',
      '.bg-blueprint-fade',
      '.sf-logo-pill',
      '.sf-brand-mark',
      '.sf-defining-card',
      '.sf-model-radar-grid',
      '.sf-model-radar-legend',
      '.sf-sdlc-orbit',
      '.sf-sdlc-node',
    ]

    for (const selector of requiredSelectors) {
      expect(styleSource).toContain(selector)
    }
  })

  it('keeps homepage marquee and defining-card motion covered by reduced-motion overrides', () => {
    expect(styleSource).toContain('.sf-marquee-track--reverse')
    expect(styleSource).toContain('sf-sdlc-orbit-stage')
    expect(styleSource).toContain('.brand-hover__letter')
  })
})
