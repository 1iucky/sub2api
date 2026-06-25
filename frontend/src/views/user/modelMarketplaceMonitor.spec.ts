import { describe, expect, it } from 'vitest'
import { matchMonitorsByModelId } from './modelMarketplaceMonitor'

describe('matchMonitorsByModelId', () => {
  it('matches monitors by exact model id only', () => {
    const monitors = [
      { primary_model: 'deepseek-v4-flash', extra_models: [{ model: 'glm-4.5' }] },
      { primary_model: 'gpt-4o', extra_models: [{ model: 'deepseek-v4-flash' }] },
      { primary_model: 'deepseek-v4-flash-thinking', extra_models: [] },
    ] as any

    expect(matchMonitorsByModelId('deepseek-v4-flash', monitors).map(m => m.primary_model)).toEqual([
      'deepseek-v4-flash',
      'gpt-4o',
    ])
  })
})
