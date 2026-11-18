type MonitorLike = {
  primary_model: string
  extra_models?: Array<{ model: string }>
  timeline: Array<{
    status: string
    latency_ms: number | null
    checked_at: string
  }>
  availability_7d: number
}

export type ModelLike = {
  id: number | string
  model_id: string
  related_pricing?: {
    channel_count?: number
    entries?: unknown[]
  }
  updated_at?: string
}

function normalizeModelId(value: string) {
  return (value || '').trim().toLowerCase()
}

export function matchMonitorsByModelId(modelId: string, monitors: MonitorLike[]) {
  const normalized = normalizeModelId(modelId)
  if (!normalized) return []
  return monitors.filter((monitor) => {
    if (normalizeModelId(monitor.primary_model) === normalized) return true
    return monitor.extra_models?.some(extra => normalizeModelId(extra.model) === normalized) ?? false
  })
}

export function dedupeModelsByModelId<T extends ModelLike>(models: T[]) {
  const byModelId = new Map<string, T>()
  for (const model of models) {
    const normalized = normalizeModelId(model.model_id)
    if (!normalized) continue
    const existing = byModelId.get(normalized)
    if (!existing || modelRank(model) > modelRank(existing)) {
      byModelId.set(normalized, model)
    }
  }
  return Array.from(byModelId.values())
}

function modelRank(model: ModelLike) {
  const pricing = model.related_pricing
  const channelCount = pricing?.channel_count || 0
  const entryCount = pricing?.entries?.length || 0
  const updatedAt = Date.parse(model.updated_at || '') || 0
  return channelCount * 1_000_000 + entryCount * 1_000 + updatedAt / 1_000_000_000_000
}
