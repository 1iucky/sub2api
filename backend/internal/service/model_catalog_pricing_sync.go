package service

import (
	"sort"
	"strings"
	"time"
)

type PricingCatalogEntry struct {
	ModelID      string
	Provider     string
	Mode         string
	Capabilities map[string]any
	Pricing      map[string]any
	Metadata     map[string]any
}

func (s *PricingService) ListCatalogEntries() []PricingCatalogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]PricingCatalogEntry, 0, len(s.pricingData))
	for modelID, pricing := range s.pricingData {
		if strings.TrimSpace(modelID) == "" || pricing == nil {
			continue
		}
		out = append(out, PricingCatalogEntry{
			ModelID:      modelID,
			Provider:     strings.ToLower(strings.TrimSpace(pricing.LiteLLMProvider)),
			Mode:         pricing.Mode,
			Capabilities: capabilitiesFromLiteLLM(pricing),
			Pricing:      pricingMapFromLiteLLM(pricing),
			Metadata: map[string]any{
				"max_input_tokens":      pricing.MaxInputTokens,
				"max_output_tokens":     pricing.MaxOutputTokens,
				"max_tokens":            pricing.MaxTokens,
				"supports_service_tier": pricing.SupportsServiceTier,
				"source_provider":       pricing.LiteLLMProvider,
			},
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ModelID < out[j].ModelID
	})
	return out
}

func capabilitiesFromLiteLLM(p *LiteLLMModelPricing) map[string]any {
	caps := map[string]any{
		"assistant_prefill": p.SupportsAssistantPrefill,
		"computer_use":      p.SupportsComputerUse,
		"function_calling":  p.SupportsFunctionCalling,
		"pdf_input":         p.SupportsPDFInput,
		"prompt_caching":    p.SupportsPromptCaching,
		"reasoning":         p.SupportsReasoning,
		"response_schema":   p.SupportsResponseSchema,
		"service_tier":      p.SupportsServiceTier,
		"tool_choice":       p.SupportsToolChoice,
		"vision":            p.SupportsVision,
		"web_search":        p.SupportsWebSearch,
	}
	if p.OutputCostPerImage > 0 || p.OutputCostPerImageToken > 0 {
		caps["image_output"] = true
	}
	if p.Mode != "" {
		caps["mode"] = p.Mode
	}
	if p.LongContextInputTokenThreshold > 0 {
		caps["long_context"] = true
		caps["long_context_input_token_threshold"] = p.LongContextInputTokenThreshold
	}
	if p.MaxInputTokens > 0 || p.MaxOutputTokens > 0 || p.MaxTokens > 0 {
		caps["context_limits"] = map[string]any{
			"max_input_tokens":  p.MaxInputTokens,
			"max_output_tokens": p.MaxOutputTokens,
			"max_tokens":        p.MaxTokens,
		}
	}
	return caps
}

func pricingMapFromLiteLLM(p *LiteLLMModelPricing) map[string]any {
	pricing := map[string]any{}
	putNonZero := func(key string, value float64) {
		if value > 0 {
			pricing[key] = value
		}
	}
	putNonZero("input_cost_per_token", p.InputCostPerToken)
	putNonZero("input_cost_per_token_priority", p.InputCostPerTokenPriority)
	putNonZero("output_cost_per_token", p.OutputCostPerToken)
	putNonZero("output_cost_per_token_priority", p.OutputCostPerTokenPriority)
	putNonZero("cache_creation_input_token_cost", p.CacheCreationInputTokenCost)
	putNonZero("cache_creation_input_token_cost_above_1hr", p.CacheCreationInputTokenCostAbove1hr)
	putNonZero("cache_read_input_token_cost", p.CacheReadInputTokenCost)
	putNonZero("cache_read_input_token_cost_priority", p.CacheReadInputTokenCostPriority)
	putNonZero("output_cost_per_image", p.OutputCostPerImage)
	putNonZero("output_cost_per_image_token", p.OutputCostPerImageToken)
	if p.LongContextInputCostMultiplier > 0 {
		pricing["long_context_input_cost_multiplier"] = p.LongContextInputCostMultiplier
	}
	if p.LongContextOutputCostMultiplier > 0 {
		pricing["long_context_output_cost_multiplier"] = p.LongContextOutputCostMultiplier
	}
	return pricing
}

func modelCatalogInputFromPricing(entry PricingCatalogEntry, vendorID *int64, syncedAt time.Time) ModelCatalogUpsert {
	platform := providerToPlatform(entry.Provider)
	iconKey := iconKeyForProvider(entry.Provider)
	return ModelCatalogUpsert{
		ModelID:      entry.ModelID,
		DisplayName:  entry.ModelID,
		Platform:     platform,
		Provider:     entry.Provider,
		VendorID:     vendorID,
		Mode:         entry.Mode,
		Tags:         modelTagsFromPricing(entry),
		Capabilities: entry.Capabilities,
		Endpoints:    endpointsForMode(entry.Mode),
		Pricing:      entry.Pricing,
		Metadata:     entry.Metadata,
		Status:       ModelCatalogStatusActive,
		Visibility:   ModelCatalogVisibilityPublic,
		Source:       ModelCatalogSourceLiteLLM,
		IconKey:      iconKey,
		LastSyncedAt: &syncedAt,
	}
}

func modelTagsFromPricing(entry PricingCatalogEntry) []string {
	tags := []string{}
	if entry.Mode != "" {
		tags = append(tags, entry.Mode)
	}
	if entry.Provider != "" {
		tags = append(tags, entry.Provider)
	}
	if v, ok := entry.Capabilities["prompt_caching"].(bool); ok && v {
		tags = append(tags, "cache")
	}
	if v, ok := entry.Capabilities["image_output"].(bool); ok && v {
		tags = append(tags, "image")
	}
	if v, ok := entry.Capabilities["long_context"].(bool); ok && v {
		tags = append(tags, "long-context")
	}
	return tags
}

func endpointsForMode(mode string) []string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "embedding":
		return []string{"embeddings"}
	case "image_generation":
		return []string{"images"}
	case "audio_transcription", "audio_speech":
		return []string{"audio"}
	default:
		return []string{"chat", "responses"}
	}
}

func defaultVendorForProvider(provider string) ModelVendorUpsert {
	key := strings.ToLower(strings.TrimSpace(provider))
	if key == "" {
		key = "unknown"
	}
	name := vendorNameForProvider(key)
	return ModelVendorUpsert{
		Name:        name,
		ProviderKey: key,
		IconKey:     iconKeyForProvider(key),
		SortOrder:   vendorSortOrder(key),
	}
}

func vendorNameForProvider(provider string) string {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "anthropic":
		return "Anthropic"
	case "openai":
		return "OpenAI"
	case "google", "gemini", "vertex_ai":
		return "Google"
	case "azure":
		return "Azure OpenAI"
	case "bedrock", "aws":
		return "AWS Bedrock"
	case "xai":
		return "xAI"
	case "deepseek":
		return "DeepSeek"
	case "mistral":
		return "Mistral AI"
	case "cohere":
		return "Cohere"
	case "openrouter":
		return "OpenRouter"
	default:
		if provider == "" {
			return "Unknown"
		}
		return strings.ToUpper(provider[:1]) + provider[1:]
	}
}

func iconKeyForProvider(provider string) string {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "anthropic":
		return "claude"
	case "openai", "azure":
		return "openai"
	case "google", "gemini", "vertex_ai":
		return "gemini"
	case "deepseek":
		return "deepseek"
	case "mistral":
		return "mistral"
	case "cohere":
		return "cohere"
	case "xai":
		return "xai"
	case "openrouter":
		return "openrouter"
	case "bedrock", "aws":
		return "aws"
	default:
		return provider
	}
}

func providerToPlatform(provider string) string {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "anthropic", "claude":
		return PlatformAnthropic
	case "openai", "azure":
		return PlatformOpenAI
	case "google", "gemini", "vertex_ai":
		return PlatformGemini
	case "antigravity":
		return PlatformAntigravity
	default:
		return PlatformOpenAI
	}
}

func NormalizeModelCatalogPlatform(platform string) string {
	switch strings.ToLower(strings.TrimSpace(platform)) {
	case PlatformAnthropic, "claude":
		return PlatformAnthropic
	case PlatformGemini, "google", "vertex_ai", "vertex-ai", "vertex":
		return PlatformGemini
	case PlatformAntigravity:
		return PlatformAntigravity
	case PlatformOpenAI, "":
		return PlatformOpenAI
	default:
		return PlatformOpenAI
	}
}

func vendorSortOrder(provider string) int {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "openai":
		return 10
	case "anthropic":
		return 20
	case "google", "gemini", "vertex_ai":
		return 30
	default:
		return 100
	}
}
