import { apiClient } from './client'
import type { UserMonitorDetail, UserMonitorListResponse } from './channelMonitor'

export async function list(options?: { signal?: AbortSignal }): Promise<UserMonitorListResponse> {
  const { data } = await apiClient.get<UserMonitorListResponse>('/public/channel-monitors', {
    signal: options?.signal,
  })
  return data
}

export async function status(id: number): Promise<UserMonitorDetail> {
  const { data } = await apiClient.get<UserMonitorDetail>(`/public/channel-monitors/${id}/status`)
  return data
}

export const publicChannelMonitorAPI = { list, status }
export default publicChannelMonitorAPI
