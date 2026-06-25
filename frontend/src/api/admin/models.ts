import { apiClient } from '../client'
import type {
  ModelCatalog,
  ModelCatalogListFilters,
  ModelVendor,
  PaginatedResponse,
} from '../models'

export interface ModelCatalogRequest {
  model_id: string
  display_name?: string
  platform?: string
  provider?: string
  vendor_id?: number | null
  mode?: string
  description?: string
  tags?: string[]
  capabilities?: Record<string, unknown>
  endpoints?: string[]
  pricing?: Record<string, unknown>
  metadata?: Record<string, unknown>
  status?: 'active' | 'disabled'
  visibility?: 'public' | 'admin'
  source?: string
  icon_key?: string
}

export interface ModelVendorRequest {
  name: string
  provider_key?: string
  icon_key?: string
  description?: string
  sort_order?: number
}

export interface SyncPricingResult {
  total: number
  created: number
  updated: number
}

export async function list(
  page = 1,
  pageSize = 20,
  filters: ModelCatalogListFilters = {},
  options?: { signal?: AbortSignal }
): Promise<PaginatedResponse<ModelCatalog>> {
  const { data } = await apiClient.get<PaginatedResponse<ModelCatalog>>('/admin/models', {
    params: {
      page,
      page_size: pageSize,
      ...filters
    },
    signal: options?.signal
  })
  return data
}

export async function getById(id: number): Promise<ModelCatalog> {
  const { data } = await apiClient.get<ModelCatalog>(`/admin/models/${id}`)
  return data
}

export async function create(req: ModelCatalogRequest): Promise<ModelCatalog> {
  const { data } = await apiClient.post<ModelCatalog>('/admin/models', req)
  return data
}

export async function update(id: number, req: ModelCatalogRequest): Promise<ModelCatalog> {
  const { data } = await apiClient.put<ModelCatalog>(`/admin/models/${id}`, req)
  return data
}

export async function remove(id: number): Promise<void> {
  await apiClient.delete(`/admin/models/${id}`)
}

export async function syncPricing(): Promise<SyncPricingResult> {
  const { data } = await apiClient.post<SyncPricingResult>('/admin/models/sync-pricing')
  return data
}

export async function listVendors(options?: { signal?: AbortSignal }): Promise<ModelVendor[]> {
  const { data } = await apiClient.get<ModelVendor[]>('/admin/models/vendors', {
    signal: options?.signal
  })
  return data
}

export async function upsertVendor(req: ModelVendorRequest): Promise<ModelVendor> {
  const { data } = await apiClient.post<ModelVendor>('/admin/models/vendors', req)
  return data
}

export async function deleteVendor(id: number): Promise<void> {
  await apiClient.delete(`/admin/models/vendors/${id}`)
}

const adminModelsAPI = { list, getById, create, update, remove, syncPricing, listVendors, upsertVendor, deleteVendor }
export default adminModelsAPI
