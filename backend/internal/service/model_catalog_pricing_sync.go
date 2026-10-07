package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// catalogLiteLLMPricing is the catalog-owned view of a LiteLLM pricing entry.
// It retains the full capability facet set (supports_*, max_*_tokens) that the
// upstream LiteLLMModelPricing struct no longer parses, so the catalog sync has
// zero runtime coupling to upstream pricing internals.
type catalogLiteLLMPricing struct {
	InputCostPerToken                   float64 `json:"input_cost_per_token"`
	InputCostPerTokenPriority           float64 `json:"input_cost_per_token_priority"`
	OutputCostPerToken                  float64 `json:"output_cost_per_token"`
	OutputCostPerTokenPriority          float64 `json:"output_cost_per_token_priority"`
	CacheCreationInputTokenCost         float64 `json:"cache_creation_input_token_cost"`
	CacheCreationInputTokenCostPriority float64 `json:"cache_creation_input_token_cost_priority"`
	CacheCreationInputTokenCostAbove1hr float64 `json:"cache_creation_input_token_cost_above_1hr"`
	CacheReadInputTokenCost             float64 `json:"cache_read_input_token_cost"`
	CacheReadInputTokenCostPriority     float64 `json:"cache_read_input_token_cost_priority"`
	LongContextInputTokenThreshold      int     `json:"long_context_input_token_threshold"`
	LongContextInputCostMultiplier      float64 `json:"long_context_input_cost_multiplier"`
	LongContextOutputCostMultiplier     float64 `json:"long_context_output_cost_multiplier"`
	MaxInputTokens                      int     `json:"max_input_tokens"`
	MaxOutputTokens                     int     `json:"max_output_tokens"`
	MaxTokens                           int     `json:"max_tokens"`
	SupportsAssistantPrefill            bool    `json:"supports_assistant_prefill"`
	SupportsComputerUse                 bool    `json:"supports_computer_use"`
	SupportsFunctionCalling             bool    `json:"supports_function_calling"`
	SupportsPDFInput                    bool    `json:"supports_pdf_input"`
	SupportsPromptCaching               bool    `json:"supports_prompt_caching"`
	SupportsReasoning                   bool    `json:"supports_reasoning"`
	SupportsResponseSchema              bool    `json:"supports_response_schema"`
	SupportsServiceTier                 bool    `json:"supports_service_tier"`
	SupportsToolChoice                  bool    `json:"supports_tool_choice"`
	SupportsVision                      bool    `json:"supports_vision"`
	SupportsWebSearch                   bool    `json:"supports_web_search"`
	OutputCostPerImage                  float64 `json:"output_cost_per_image"`
	OutputCostPerImageToken             float64 `json:"output_cost_per_image_token"`
	LiteLLMProvider                     string  `json:"litellm_provider"`
	Mode                                string  `json:"mode"`
}

type PricingCatalogEntry struct {
	ModelID      string
	Provider     string
	Mode         string
	Capabilities map[string]any
	Pricing      map[string]any
	Metadata     map[string]any
}

// listCatalogEntries fetches and parses the LiteLLM pricing dataset into
// catalog entries. Remote fetch failures fall back to the configured local file.
func (s *ModelCatalogService) listCatalogEntries(ctx context.Context) ([]PricingCatalogEntry, error) {
	body, err := s.fetchCatalogPricingData(ctx)
	if err != nil {
		return nil, err
	}
	parsed := map[string]catalogLiteLLMPricing{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("parse catalog pricing data: %w", err)
	}

	out := make([]PricingCatalogEntry, 0, len(parsed))
	for modelID, pricing := range parsed {
		if strings.TrimSpace(modelID) == "" {
			continue
		}
		p := pricing
		out = append(out, PricingCatalogEntry{
			ModelID:      modelID,
			Provider:     strings.ToLower(strings.TrimSpace(p.LiteLLMProvider)),
			Mode:         p.Mode,
			Capabilities: capabilitiesFromCatalogLiteLLM(&p),
			Pricing:      pricingMapFromCatalogLiteLLM(&p),
			Metadata: map[string]any{
				"max_input_tokens":      p.MaxInputTokens,
				"max_output_tokens":     p.MaxOutputTokens,
				"max_tokens":            p.MaxTokens,
				"supports_service_tier": p.SupportsServiceTier,
				"source_provider":       p.LiteLLMProvider,
			},
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ModelID < out[j].ModelID
	})
	return out, nil
}

// fetchCatalogPricingData downloads the remote pricing JSON; on any failure it
// falls back to reading the configured fallback file.
func (s *ModelCatalogService) fetchCatalogPricingData(ctx context.Context) ([]byte, error) {
	remoteURL := ""
	fallbackFile := ""
	if s.cfg != nil {
		remoteURL = strings.TrimSpace(s.cfg.Pricing.RemoteURL)
		fallbackFile = strings.TrimSpace(s.cfg.Pricing.FallbackFile)
	}

	if remoteURL != "" {
		body, err := fetchRemotePricingJSON(ctx, remoteURL)
		if err == nil {
			return body, nil
		}
		if fallbackFile == "" {
			return nil, fmt.Errorf("download catalog pricing data: %w", err)
		}
	}

	if fallbackFile == "" {
		return nil, fmt.Errorf("no catalog pricing source configured")
	}
	body, err := os.ReadFile(fallbackFile)
	if err != nil {
		return nil, fmt.Errorf("read catalog pricing fallback file: %w", err)
	}
	return body, nil
}

var catalogPricingHTTPClient = &http.Client{Timeout: 30 * time.Second}

func fetchRemotePricingJSON(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := catalogPricingHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return io.ReadAll(io.LimitReader(resp.Body, 64<<20))
}

func capabilitiesFromCatalogLiteLLM(p *catalogLiteLLMPricing) map[string]any {
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

func pricingMapFromCatalogLiteLLM(p *catalogLiteLLMPricing) map[string]any {
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
