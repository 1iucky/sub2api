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
