import { describe, expect, it } from 'vitest'
import { dedupeModelsByModelId, matchMonitorsByModelId } from './modelMarketplaceMonitor'

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

  it('deduplicates marketplace models by normalized model id and keeps the better priced row', () => {
    const models = [
      {
        id: 1,
        model_id: 'GLM-5.2',
        related_pricing: { channel_count: 1, entries: [{}] },
        updated_at: '2026-06-01T00:00:00Z',
      },
      {
        id: 2,
        model_id: 'glm-5.2',
        related_pricing: { channel_count: 2, entries: [{}, {}] },
        updated_at: '2026-06-02T00:00:00Z',
      },
      {
        id: 3,
        model_id: 'deepseek-v4-flash',
        related_pricing: { channel_count: 1, entries: [{}] },
      },
    ] as any

    expect(dedupeModelsByModelId(models).map(model => model.id)).toEqual([2, 3])
  })
})
