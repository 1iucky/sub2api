import { apiClient } from './client'
import type { DisplayCurrency } from '@/utils/pricing'

export interface ModelVendor {
  id: number
  name: string
  provider_key: string
  icon_key: string
  description: string
  sort_order: number
  created_at: string
  updated_at: string
}

export interface ModelPricingAssociation {
  channel_count: number
  channels: string[]
  entries: ModelPricingAssociationEntry[]
  groups: ModelPricingGroupSummary[]
  has_intervals: boolean
}

export interface ModelGroupAssociation {
  id: number
  name: string
  platform: string
  rate_multiplier: number
  is_exclusive: boolean
  status: string
}

export interface ModelPricingGroupSummary {
  id: number
  name: string
  platform: string
  rate_multiplier: number
}

export interface ModelPricingAssociationEntry {
  channel_id: number
  channel_name: string
  platform: string
  models: string[]
  billing_mode: string
  display_currency: DisplayCurrency
  input_price: number | null
  output_price: number | null
  cache_write_price: number | null
  cache_read_price: number | null
  image_output_price: number | null
  per_request_price: number | null
  intervals: ModelPricingAssociationInterval[]
  groups: ModelPricingGroupSummary[]
}

export interface ModelPricingAssociationInterval {
  display_currency?: DisplayCurrency
  min_tokens: number
  max_tokens: number | null
  tier_label: string
  input_price: number | null
  output_price: number | null
  cache_write_price: number | null
  cache_read_price: number | null
  per_request_price: number | null
  sort_order: number
}

export interface ModelPricingTier {
  label?: string
  threshold?: number
  input_cost_per_token?: number
  output_cost_per_token?: number
  cache_creation_input_token_cost?: number
  cache_read_input_token_cost?: number
}

export interface ModelPricing {
  input_cost_per_token?: number
  input_cost_per_token_priority?: number
  output_cost_per_token?: number
  output_cost_per_token_priority?: number
  cache_creation_input_token_cost?: number
  cache_creation_input_token_cost_above_1hr?: number
  cache_read_input_token_cost?: number
  cache_read_input_token_cost_priority?: number
  output_cost_per_image?: number
  output_cost_per_image_token?: number
  long_context_input_cost_multiplier?: number
  long_context_output_cost_multiplier?: number
  tiers?: ModelPricingTier[]
  [key: string]: unknown
}

export interface ModelCapabilities {
  assistant_prefill?: boolean
  computer_use?: boolean
  function_calling?: boolean
  pdf_input?: boolean
  prompt_caching?: boolean
  reasoning?: boolean
  response_schema?: boolean
  service_tier?: boolean
  tool_choice?: boolean
  vision?: boolean
  web_search?: boolean
  image_output?: boolean
  long_context?: boolean
  mode?: string
  context_limits?: {
    max_input_tokens?: number
    max_output_tokens?: number
    max_tokens?: number
  }
  [key: string]: unknown
}

export interface ModelMetadata {
  max_input_tokens?: number
  max_output_tokens?: number
  max_tokens?: number
  supports_service_tier?: boolean
  source_provider?: string
  associated_group_ids?: number[]
  monitor_keys?: string[]
  [key: string]: unknown
}

export interface ModelCatalog {
  id: number
  model_id: string
  display_name: string
  platform: string
  provider: string
  vendor_id: number | null
  vendor: ModelVendor | null
  mode: string
  description: string
  tags: string[]
  capabilities: ModelCapabilities
  endpoints: string[]
  pricing: ModelPricing
  metadata: ModelMetadata
  status: 'active' | 'disabled'
  visibility: 'public' | 'admin'
  source: string
  icon_key: string
  last_synced_at: string | null
  related_pricing: ModelPricingAssociation
  related_groups: ModelGroupAssociation[]
  created_at: string
  updated_at: string
}

export interface ModelCatalogListFilters {
  search?: string
  platform?: string
  provider?: string
  vendor_id?: number | null
  status?: string
  visibility?: string
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  pages: number
}

export async function listModels(
  page = 1,
  pageSize = 24,
  filters: ModelCatalogListFilters = {},
  options?: { signal?: AbortSignal }
): Promise<PaginatedResponse<ModelCatalog>> {
  const { data } = await apiClient.get<PaginatedResponse<ModelCatalog>>('/public/models', {
    params: {
      page,
      page_size: pageSize,
      ...filters
    },
    signal: options?.signal
  })
  return data
}

export async function listVendors(options?: { signal?: AbortSignal }): Promise<ModelVendor[]> {
  const { data } = await apiClient.get<ModelVendor[]>('/public/models/vendors', {
    signal: options?.signal
  })
  return data
}

const modelsAPI = { listModels, listVendors }
export default modelsAPI
