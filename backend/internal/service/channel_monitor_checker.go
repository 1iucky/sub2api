package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/servertiming"
	"github.com/Wei-Shaw/sub2api/internal/util/httputil"
	"github.com/tidwall/gjson"
)

// monitorHTTPClient 共享一个 http.Client，避免每次检测重建 transport。
// 自定义 Transport 在 dial 时强制再次校验 IP，防止 DNS rebinding 绕过 validateEndpoint。
var monitorHTTPClient = newSSRFSafeHTTPClient(monitorRequestTimeout)

// monitorPingHTTPClient 用于 endpoint origin 的 HEAD ping，超时更短。
var monitorPingHTTPClient = newSSRFSafeHTTPClient(monitorPingTimeout)

// newSSRFSafeHTTPClient 返回一个使用 safeDialContext 的 http.Client。
// 仅供监控模块对外发起请求使用——所有目标都应是公网 endpoint。
func newSSRFSafeHTTPClient(timeout time.Duration) *http.Client {
	tr := &http.Transport{
		DialContext:           safeDialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          16,
		IdleConnTimeout:       monitorIdleConnTimeout,
		TLSHandshakeTimeout:   monitorTLSHandshakeTimeout,
		ResponseHeaderTimeout: monitorResponseHeaderTimeout,
	}
	return &http.Client{Timeout: timeout, Transport: servertiming.WrapRoundTripper(tr)}
}

// CheckOptions 承载一次检测的自定义入参。
// 所有字段都是可选（零值即等价于"用默认行为"）。
type CheckOptions struct {
	// APIMode 仅对 OpenAI provider 生效；空串等同 chat_completions。
	APIMode string
	// ExtraHeaders 用户自定义 HTTP 头（merge 到 adapter 默认 headers，用户优先）。
	ExtraHeaders map[string]string
	// BodyOverrideMode: off | merge | replace
	BodyOverrideMode string
	// BodyOverride 在 merge 模式下做浅合并（key 命中黑名单时静默丢弃），
	// 在 replace 模式下直接当作完整 body。
	BodyOverride map[string]any
}

// runCheckForModel 对单个 (provider, model) 做一次完整检测。
// 不返回 error：所有失败都包装进 CheckResult.Status=error/failed。
//
// opts 承载模板 / 监控快照带来的自定义配置。nil 等同于 "off + 无 extra headers"。
func runCheckForModel(ctx context.Context, provider, endpoint, apiKey, model string, opts *CheckOptions) *CheckResult {
	res := &CheckResult{
		Model:     model,
		Status:    MonitorStatusError,
		CheckedAt: time.Now(),
	}

	challenge := generateChallenge()
	mode := bodyOverrideMode(opts)

	start := time.Now()
	respText, rawBody, statusCode, respHeaders, err := callProvider(ctx, provider, endpoint, apiKey, model, challenge.Prompt, opts)
	latency := time.Since(start)
	latencyMs := int(latency / time.Millisecond)
	res.LatencyMs = &latencyMs

	if err != nil {
		res.Status = MonitorStatusError
		res.Message = truncateMessage(sanitizeErrorMessage(err.Error()))
		return res
	}
	if statusCode < 200 || statusCode >= 300 {
		// 错误路径：用 rawBody 而非 respText（gjson textPath 抽取在错误响应里通常为空，
		// 会丢掉真正的上游错误信息，例如 `{"error":{"message":"No available accounts ..."}}`）。
		res.Status = MonitorStatusError
		res.Message = formatMonitorUpstreamHTTPError(statusCode, respHeaders, rawBody)
		return res
	}

	// Replace 模式：跳过 challenge 校验（用户 body 是静态的，challenge 没法嵌入）。
	// 改用「HTTP 2xx + 响应文本（adapter.textPath 抽取）非空」作为 operational 判定。
	// 响应文本为空则降级为 failed（视为上游回了 200 但没实际内容）。
	if mode == MonitorBodyOverrideModeReplace {
		if strings.TrimSpace(respText) == "" {
			res.Status = MonitorStatusFailed
			res.Message = truncateMessage("replace-mode: upstream returned 2xx with empty text")
			return res
		}
		return finalizeOperationalOrDegraded(res, latency, latencyMs)
	}

	if !validateChallenge(respText, challenge.Expected) {
		res.Status = MonitorStatusFailed
		res.Message = truncateMessage(sanitizeErrorMessage(fmt.Sprintf("challenge mismatch (expected %s, got %q)", challenge.Expected, respText)))
		return res
	}

	return finalizeOperationalOrDegraded(res, latency, latencyMs)
}

// finalizeOperationalOrDegraded 负责走到最后一步的 operational/degraded 判定。
// 拆出来是为了让 runCheckForModel 不超过 30 行。
func finalizeOperationalOrDegraded(res *CheckResult, latency time.Duration, latencyMs int) *CheckResult {
	if latency >= monitorDegradedThreshold {
		res.Status = MonitorStatusDegraded
		res.Message = truncateMessage(fmt.Sprintf("slow response: %dms", latencyMs))
		return res
	}
	res.Status = MonitorStatusOperational
	return res
}

// bodyOverrideMode 归一取 opts.BodyOverrideMode，nil opts / 空串都视为 off。
func bodyOverrideMode(opts *CheckOptions) string {
	if opts == nil || opts.BodyOverrideMode == "" {
		return MonitorBodyOverrideModeOff
	}
	return opts.BodyOverrideMode
}

// pingEndpointOrigin 对 endpoint 的 origin (scheme://host) 发起 HEAD 请求，返回耗时。
// 失败时返回 nil（不影响主状态判定）。
func pingEndpointOrigin(ctx context.Context, endpoint string) *int {
	origin, err := extractOrigin(endpoint)
	if err != nil || origin == "" {
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, origin, nil)
	if err != nil {
		return nil
	}
	start := time.Now()
	resp, err := monitorPingHTTPClient.Do(req)
	if err != nil {
		return nil
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, monitorPingDiscardMaxBytes))
	ms := int(time.Since(start) / time.Millisecond)
	return &ms
}

// providerAdapter 描述某个 provider 在 challenge 检测中需要的 4 件事：
//   - 拼出请求路径（含 model 占位）
//   - 序列化请求体
//   - 构造鉴权头
//   - 从响应 JSON 中按 path 提取文本（gjson path）
//
// 加新 provider 只需要在 providerAdapters 里增加一个条目，无需触碰 callProvider / validateProvider。
type providerAdapter struct {
	buildPath    func(model string) string
	buildBody    func(model, prompt string) ([]byte, error)
	buildHeaders func(apiKey string) map[string]string
	textPath     string // gjson 提取响应文本的 path
}

// providerAdapters 全部已支持的 provider。键值即 MonitorProvider* 字符串。
//
//nolint:gochecknoglobals // 适配器表是只读静态数据，初始化后不变更。
var providerAdapters = map[string]providerAdapter{
	MonitorProviderOpenAI: providerOpenAIChatAdapter,
	MonitorProviderGrok:   providerGrokChatAdapter,
	MonitorProviderAnthropic: {
		buildPath: func(string) string { return providerAnthropicPath },
		buildBody: func(model, prompt string) ([]byte, error) {
			return json.Marshal(map[string]any{
				"model":      model,
				"messages":   []map[string]string{{"role": "user", "content": prompt}},
				"max_tokens": monitorChallengeMaxTokens,
			})
		},
		buildHeaders: func(apiKey string) map[string]string {
			return map[string]string{
				"x-api-key":         apiKey,
				"anthropic-version": monitorAnthropicAPIVersion,
			}
		},
		textPath: "content.0.text",
	},
	MonitorProviderGemini: {
		// Gemini 把 model 名写在 URL path 上：/v1beta/models/{model}:generateContent
		buildPath: func(model string) string { return fmt.Sprintf(providerGeminiPathTemplate, model) },
		buildBody: func(_, prompt string) ([]byte, error) {
			return json.Marshal(map[string]any{
				"contents": []map[string]any{
					{"parts": []map[string]any{{"text": prompt}}},
				},
				"generationConfig": map[string]any{"maxOutputTokens": monitorChallengeMaxTokens},
			})
		},
		// 使用 x-goog-api-key header 而不是 ?key= query，避免 *url.Error 把 key 回填到错误日志。
		buildHeaders: func(apiKey string) map[string]string {
			return map[string]string{"x-goog-api-key": apiKey}
		},
		textPath: "candidates.0.content.parts.0.text",
	},
}

//nolint:gochecknoglobals // 适配器表是只读静态数据，初始化后不变更。
var providerOpenAIChatAdapter = newOpenAICompatibleChatAdapter(providerOpenAIPath)

//nolint:gochecknoglobals // 适配器表是只读静态数据，初始化后不变更。
var providerGrokChatAdapter = newOpenAICompatibleChatAdapter(providerGrokPath)

func newOpenAICompatibleChatAdapter(path string) providerAdapter {
	return providerAdapter{
		buildPath: func(string) string { return path },
		buildBody: func(model, prompt string) ([]byte, error) {
			return json.Marshal(map[string]any{
				"model":      model,
				"messages":   []map[string]string{{"role": "user", "content": prompt}},
				"max_tokens": monitorChallengeMaxTokens,
				"stream":     false,
			})
		},
		buildHeaders: func(apiKey string) map[string]string {
			return map[string]string{"Authorization": "Bearer " + apiKey}
		},
		textPath: "choices.0.message.content",
	}
}

//nolint:gochecknoglobals // 适配器表是只读静态数据，初始化后不变更。
var providerOpenAIResponsesAdapter = providerAdapter{
	buildPath: func(string) string { return providerOpenAIResponsesPath },
	buildBody: func(model, prompt string) ([]byte, error) {
		return json.Marshal(map[string]any{
			"model":             model,
			"instructions":      "You are a channel health-check endpoint. Answer the arithmetic challenge exactly and briefly.",
			"input":             prompt,
			"max_output_tokens": monitorChallengeMaxTokens,
			"stream":            false,
		})
	},
	buildHeaders: func(apiKey string) map[string]string {
		return map[string]string{"Authorization": "Bearer " + apiKey}
	},
	textPath: "output.0.content.0.text",
}

// providerAdapterFor 按 provider + api_mode 选择具体 adapter。
func providerAdapterFor(provider, apiMode string) (providerAdapter, string, bool) {
	if provider == MonitorProviderOpenAI && defaultAPIMode(apiMode) == MonitorAPIModeResponses {
		return providerOpenAIResponsesAdapter, MonitorAPIModeResponses, true
	}
	adapter, ok := providerAdapters[provider]
	return adapter, MonitorAPIModeChatCompletions, ok
}

// isSupportedProvider 校验 provider 字符串是否在 adapter 表中。
// 供 validate.go 的 validateProvider 复用，避免两份 switch 漂移。
func isSupportedProvider(p string) bool {
	_, ok := providerAdapters[p]
	return ok
}

// callProvider 通过 providerAdapters 分发到具体实现。
// opts 承载用户的自定义 headers / body 覆盖（可为 nil）。
//
// 返回值：
//   - extractedText: 按 textPath 抽出的成功文本，仅在 status 2xx 时有意义；非 2xx 时通常为空串
//   - rawBody: 完整响应体的字符串形式（已被 monitorResponseMaxBytes 截断），用于错误路径保留上游真实回包
//   - status: HTTP 状态码
//   - err: 网络 / 序列化错误
func callProvider(ctx context.Context, provider, endpoint, apiKey, model, prompt string, opts *CheckOptions) (extractedText, rawBody string, status int, headers http.Header, err error) {
	requestedAPIMode := checkAPIMode(opts)
	if err := validateAPIMode(provider, requestedAPIMode); err != nil {
		return "", "", 0, nil, err
	}
	adapter, apiMode, ok := providerAdapterFor(provider, requestedAPIMode)
	if !ok {
		return "", "", 0, nil, fmt.Errorf("unsupported provider %q", provider)
	}
	body, err := buildRequestBody(adapter, provider, apiMode, model, prompt, opts)
	if err != nil {
		return "", "", 0, nil, err
	}
	reqHeaders := mergeHeaders(adapter.buildHeaders(apiKey), opts)
	full := joinURL(endpoint, adapter.buildPath(model))
	respBytes, status, respHeaders, err := postRawJSON(ctx, full, body, reqHeaders)
	if err != nil {
		return "", "", status, respHeaders, err
	}
	if provider == MonitorProviderOpenAI && apiMode == MonitorAPIModeResponses {
		return extractOpenAIResponsesText(respBytes), string(respBytes), status, respHeaders, nil
	}
	if provider == MonitorProviderOpenAI {
		return extractOpenAIChatText(respBytes), string(respBytes), status, respHeaders, nil
	}
	if provider == MonitorProviderAnthropic {
		return extractAnthropicText(respBytes), string(respBytes), status, respHeaders, nil
	}
	return gjson.GetBytes(respBytes, adapter.textPath).String(), string(respBytes), status, respHeaders, nil
}

// extractOpenAIChatText 兼容本项目“当前服务”作为监控 endpoint 时可能出现的几类回包：
// 1. 标准 Chat Completions JSON: choices[].message.content
// 2. 兼容层返回的 Responses JSON: output[].content[].text
// 3. text/event-stream 中的 Chat/Responses data 帧（即使请求 stream=false，上游兼容层也可能这样返回）
func extractOpenAIChatText(respBytes []byte) string {
	if text := extractOpenAIChatJSONText(gjson.ParseBytes(respBytes)); strings.TrimSpace(text) != "" {
		return text
	}
	if text := extractOpenAIResponsesText(respBytes); strings.TrimSpace(text) != "" {
		return text
	}
	return extractOpenAISSEText(respBytes)
}

func extractOpenAIChatJSONText(root gjson.Result) string {
	if !root.Exists() {
		return ""
	}
	var texts []string
	choices := root.Get("choices")
	if !choices.IsArray() {
		return root.Get("choices.0.message.content").String()
	}
	choices.ForEach(func(_, choice gjson.Result) bool {
		content := choice.Get("message.content")
		if text := openAIContentText(content); strings.TrimSpace(text) != "" {
			texts = append(texts, text)
		}
		return true
	})
	return strings.Join(texts, "")
}

func extractOpenAIChatDeltaText(root gjson.Result) string {
	var texts []string
	root.Get("choices").ForEach(func(_, choice gjson.Result) bool {
		content := choice.Get("delta.content")
		if text := openAIContentText(content); strings.TrimSpace(text) != "" {
			texts = append(texts, text)
		}
		return true
	})
	return strings.Join(texts, "")
}

func openAIContentText(content gjson.Result) string {
	switch {
	case !content.Exists():
		return ""
	case content.IsArray():
		var parts []string
		content.ForEach(func(_, part gjson.Result) bool {
			if text := part.Get("text").String(); strings.TrimSpace(text) != "" {
				parts = append(parts, text)
			}
			return true
		})
		return strings.Join(parts, "")
	default:
		return content.String()
	}
}

// extractOpenAIResponsesText 聚合 Responses API 的最终 assistant 文本。
// Responses 的 output 数组顺序由模型决定：reasoning / tool-call item 可能排在 message 前面，
// 因此不能假设文本永远在 output.0.content.0.text。
func extractOpenAIResponsesText(respBytes []byte) string {
	root := gjson.ParseBytes(respBytes)
	if text := extractOpenAIResponsesTextFromResult(root); strings.TrimSpace(text) != "" {
		return text
	}
	return extractOpenAIResponsesTextFromResult(root.Get("response"))
}

func extractOpenAIResponsesTextFromResult(root gjson.Result) string {
	if !root.Exists() {
		return ""
	}
	if text := root.Get("output_text").String(); strings.TrimSpace(text) != "" {
		return text
	}
	var texts []string
	outputs := root.Get("output")
	if outputs.IsArray() {
		outputs.ForEach(func(_, output gjson.Result) bool {
			outputType := output.Get("type").String()
			if outputType != "" && outputType != "message" {
				return true
			}

			content := output.Get("content")
			if !content.IsArray() {
				return true
			}

			content.ForEach(func(_, block gjson.Result) bool {
				blockType := block.Get("type").String()
				if blockType != "" && blockType != "output_text" {
					return true
				}
				if text := block.Get("text").String(); strings.TrimSpace(text) != "" {
					texts = append(texts, text)
				}
				return true
			})
			return true
		})
	}

	if len(texts) > 0 {
		return strings.Join(texts, "")
	}
	return root.Get(providerOpenAIResponsesAdapter.textPath).String()
}

func extractOpenAISSEText(respBytes []byte) string {
	var deltas []string
	terminalText := ""
	flush := func(data string) {
		text, terminal := extractOpenAISSEDataText(data)
		if strings.TrimSpace(text) == "" {
			return
		}
		if terminal {
			terminalText = text
			return
		}
		deltas = append(deltas, text)
	}

	var dataLines []string
	for _, line := range strings.Split(string(respBytes), "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) == "" {
			if len(dataLines) > 0 {
				flush(strings.Join(dataLines, "\n"))
				dataLines = nil
			}
			continue
		}
		if data, ok := strings.CutPrefix(line, "data:"); ok {
			dataLines = append(dataLines, strings.TrimSpace(data))
		}
	}
	if len(dataLines) > 0 {
		flush(strings.Join(dataLines, "\n"))
	}

	if strings.TrimSpace(terminalText) != "" {
		return terminalText
	}
	return strings.Join(deltas, "")
}

func extractOpenAISSEDataText(data string) (text string, terminal bool) {
	data = strings.TrimSpace(data)
	if data == "" || data == "[DONE]" || !gjson.Valid(data) {
		return "", false
	}
	root := gjson.Parse(data)
	if text := extractOpenAIChatJSONText(root); strings.TrimSpace(text) != "" {
		return text, true
	}
	if text := extractOpenAIResponsesTextFromResult(root); strings.TrimSpace(text) != "" {
		return text, true
	}
	if text := extractOpenAIResponsesTextFromResult(root.Get("response")); strings.TrimSpace(text) != "" {
		return text, true
	}
	if text := extractOpenAIChatDeltaText(root); strings.TrimSpace(text) != "" {
		return text, false
	}
	if text := root.Get("delta").String(); root.Get("type").String() == "response.output_text.delta" && strings.TrimSpace(text) != "" {
		return text, false
	}
	return "", false
}

// extractAnthropicText 兼容 /v1/messages 可能返回的两种形态：
// 1. 标准 JSON：content[].text
// 2. text/event-stream：message_start / content_block_start / content_block_delta / message_delta
func extractAnthropicText(respBytes []byte) string {
	if text := extractAnthropicJSONText(gjson.ParseBytes(respBytes)); strings.TrimSpace(text) != "" {
		return text
	}
	if text := extractAnthropicSSEText(respBytes); strings.TrimSpace(text) != "" {
		return text
	}
	return gjson.GetBytes(respBytes, providerAdapters[MonitorProviderAnthropic].textPath).String()
}

func extractAnthropicJSONText(root gjson.Result) string {
	if !root.Exists() {
		return ""
	}
	if text := anthropicContentBlocksText(root.Get("content")); strings.TrimSpace(text) != "" {
		return text
	}
	if text := anthropicContentBlocksText(root.Get("response.content")); strings.TrimSpace(text) != "" {
		return text
	}
	if text := anthropicContentBlocksText(root.Get("message.content")); strings.TrimSpace(text) != "" {
		return text
	}
	return ""
}

func extractAnthropicSSEText(respBytes []byte) string {
	var parts []string

	var eventLines []string
	flush := func(lines []string) {
		if len(lines) == 0 {
			return
		}
		eventType, data := parseSSERecord(lines)
		if text := extractAnthropicSSERecordText(eventType, data); strings.TrimSpace(text) != "" {
			// Anthropic 流式文本通常分散在 content_block_delta 中，逐段拼接即可。
			parts = append(parts, text)
		}
	}

	for _, line := range strings.Split(string(respBytes), "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) == "" {
			flush(eventLines)
			eventLines = nil
			continue
		}
		eventLines = append(eventLines, line)
	}
	flush(eventLines)
	return strings.Join(parts, "")
}

func parseSSERecord(lines []string) (eventType, data string) {
	var dataLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "event:") {
			eventType = strings.TrimSpace(strings.TrimPrefix(trimmed, "event:"))
			continue
		}
		if strings.HasPrefix(trimmed, "data:") {
			dataLines = append(dataLines, strings.TrimSpace(strings.TrimPrefix(trimmed, "data:")))
		}
	}
	return eventType, strings.Join(dataLines, "\n")
}

func extractAnthropicSSERecordText(eventType, data string) string {
	data = strings.TrimSpace(data)
	if data == "" || data == "[DONE]" || !gjson.Valid(data) {
		return ""
	}

	root := gjson.Parse(data)
	if text := anthropicMessageStartText(root); strings.TrimSpace(text) != "" {
		return text
	}
	if text := anthropicContentBlockStartText(root); strings.TrimSpace(text) != "" {
		return text
	}
	if text := anthropicContentBlockDeltaText(root); strings.TrimSpace(text) != "" {
		return text
	}
	if text := anthropicMessageDeltaText(root); strings.TrimSpace(text) != "" {
		return text
	}
	// 某些兼容层会直接把最终消息包在 data 里，不带标准 SSE 事件名。
	if eventType == "message" {
		return anthropicContentBlocksText(root.Get("content"))
	}
	return ""
}

func anthropicMessageStartText(root gjson.Result) string {
	if text := anthropicContentBlocksText(root.Get("message.content")); strings.TrimSpace(text) != "" {
		return text
	}
	if text := anthropicContentBlocksText(root.Get("content")); strings.TrimSpace(text) != "" {
		return text
	}
	return ""
}

func anthropicContentBlockStartText(root gjson.Result) string {
	block := root.Get("content_block")
	if !block.Exists() {
		return ""
	}
	if text := block.Get("text").String(); strings.TrimSpace(text) != "" {
		return text
	}
	return ""
}

func anthropicContentBlockDeltaText(root gjson.Result) string {
	delta := root.Get("delta")
	if !delta.Exists() {
		return ""
	}
	if text := delta.Get("text").String(); strings.TrimSpace(text) != "" {
		return text
	}
	return ""
}

func anthropicMessageDeltaText(root gjson.Result) string {
	if text := root.Get("delta.text").String(); strings.TrimSpace(text) != "" {
		return text
	}
	return ""
}

func anthropicContentBlocksText(content gjson.Result) string {
	switch {
	case !content.Exists():
		return ""
	case content.IsArray():
		var parts []string
		content.ForEach(func(_, block gjson.Result) bool {
			if blockType := block.Get("type").String(); blockType != "" && blockType != "text" {
				return true
			}
			if text := block.Get("text").String(); strings.TrimSpace(text) != "" {
				parts = append(parts, text)
			}
			return true
		})
		return strings.Join(parts, "")
	default:
		return content.String()
	}
}

// mergeHeaders 把用户自定义 headers 合并到 adapter 默认 headers 上。
// 用户值覆盖默认；命中黑名单（hop-by-hop / 由 http.Client 自管的）的 key 静默丢弃。
func mergeHeaders(base map[string]string, opts *CheckOptions) map[string]string {
	if opts == nil || len(opts.ExtraHeaders) == 0 {
		return base
	}
	out := make(map[string]string, len(base)+len(opts.ExtraHeaders))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range opts.ExtraHeaders {
		if IsForbiddenHeaderName(k) {
			continue
		}
		out[k] = v
	}
	return out
}

// buildRequestBody 根据 body_override_mode 构造请求 body。
//
//   - off:     adapter 默认 body
//   - merge:   adapter 默认 body 与 BodyOverride 浅合并；BodyOverride 中命中
//     bodyMergeKeyDenyList[provider] 的 key 会被静默丢弃，避免破坏 challenge / model 路由
//   - replace: 直接 marshal BodyOverride 作为完整 body
//
// 任何 mode 返回的 []byte 都已经是合法 JSON，可直接送入 postRawJSON。
func buildRequestBody(adapter providerAdapter, provider, apiMode, model, prompt string, opts *CheckOptions) ([]byte, error) {
	mode := bodyOverrideMode(opts)

	if mode == MonitorBodyOverrideModeReplace {
		if opts == nil || len(opts.BodyOverride) == 0 {
			return nil, fmt.Errorf("replace mode: body_override is empty")
		}
		if err := validateReplaceRequestBody(provider, apiMode, opts.BodyOverride); err != nil {
			return nil, err
		}
		body, err := json.Marshal(opts.BodyOverride)
		if err != nil {
			return nil, fmt.Errorf("marshal body_override (replace): %w", err)
		}
		return body, nil
	}

	defaultBody, err := adapter.buildBody(model, prompt)
	if err != nil {
		return nil, fmt.Errorf("marshal default body: %w", err)
	}
	if mode != MonitorBodyOverrideModeMerge || opts == nil || len(opts.BodyOverride) == 0 {
		return defaultBody, nil
	}

	var defaultMap map[string]any
	if err := json.Unmarshal(defaultBody, &defaultMap); err != nil {
		return nil, fmt.Errorf("unmarshal default body for merge: %w", err)
	}
	deny := bodyMergeKeyDenyList[bodyMergeDenyKey(provider, apiMode)]
	for k, v := range opts.BodyOverride {
		if deny[k] {
			continue
		}
		defaultMap[k] = v
	}
	merged, err := json.Marshal(defaultMap)
	if err != nil {
		return nil, fmt.Errorf("marshal merged body: %w", err)
	}
	return merged, nil
}

// bodyMergeKeyDenyList 在 merge 模式下，禁止用户覆盖这些 provider-specific 的关键字段。
// 思路抄 check-cx 的 EXCLUDED_METADATA_KEYS：保护 challenge / model 路由不被用户误伤。
// 用户想动这些字段就用 replace 模式（已知会跳 challenge 校验）。
//
//nolint:gochecknoglobals // 静态查表，初始化后不变。
var bodyMergeKeyDenyList = map[string]map[string]bool{
	MonitorProviderOpenAI + ":" + MonitorAPIModeChatCompletions: {"model": true, "messages": true, "stream": true},
	MonitorProviderOpenAI + ":" + MonitorAPIModeResponses:       {"model": true, "instructions": true, "input": true, "stream": true},
	MonitorProviderGrok:      {"model": true, "messages": true, "stream": true},
	MonitorProviderAnthropic: {"model": true, "messages": true},
	MonitorProviderGemini:    {"contents": true},
}

func checkAPIMode(opts *CheckOptions) string {
	if opts == nil {
		return MonitorAPIModeChatCompletions
	}
	return defaultAPIMode(opts.APIMode)
}

func bodyMergeDenyKey(provider, apiMode string) string {
	if provider == MonitorProviderOpenAI {
		return provider + ":" + defaultAPIMode(apiMode)
	}
	return provider
}

func validateReplaceRequestBody(provider, apiMode string, body map[string]any) error {
	if provider != MonitorProviderOpenAI && provider != MonitorProviderGrok {
		return nil
	}
	switch defaultAPIMode(apiMode) {
	case MonitorAPIModeResponses:
		if strings.TrimSpace(stringFromAny(body["instructions"])) == "" || !hasNonEmptyBodyValue(body["input"]) {
			return fmt.Errorf("replace mode responses body: instructions and input are required")
		}
	case MonitorAPIModeChatCompletions:
		if !hasNonEmptyBodyValue(body["messages"]) {
			return fmt.Errorf("replace mode chat_completions body: messages are required")
		}
	}
	return nil
}

func stringFromAny(v any) string {
	s, _ := v.(string)
	return s
}

func hasNonEmptyBodyValue(v any) bool {
	switch val := v.(type) {
	case nil:
		return false
	case string:
		return strings.TrimSpace(val) != ""
	case []any:
		return len(val) > 0
	case []map[string]any:
		return len(val) > 0
	case []map[string]string:
		return len(val) > 0
	default:
		return true
	}
}

// postRawJSON 发送 POST + 已序列化好的 JSON 字节，限制响应体大小，返回响应字节、HTTP status、错误。
// adapter 自行 marshal 是为了精确控制字段顺序与类型，所以这里直接收 []byte 而不是 any。
func postRawJSON(ctx context.Context, fullURL string, payload []byte, headers map[string]string) ([]byte, int, http.Header, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(payload))
	if err != nil {
		return nil, 0, nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := monitorHTTPClient.Do(req)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, monitorResponseMaxBytes))
	if err != nil {
		return nil, resp.StatusCode, resp.Header.Clone(), fmt.Errorf("read body: %w", err)
	}
	return respBody, resp.StatusCode, resp.Header.Clone(), nil
}

// joinURL 把 base origin 与 path 拼成完整 URL。
// 容忍 base 末尾有/无斜杠，path 必带前导斜杠。
func joinURL(base, path string) string {
	base = strings.TrimRight(base, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return base + path
}

// extractOrigin 从一个 endpoint URL 中提取 scheme://host[:port] 部分。
func extractOrigin(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" || u.Host == "" {
		return "", errors.New("endpoint missing scheme or host")
	}
	return u.Scheme + "://" + u.Host, nil
}

// monitorSensitiveQueryParamRegex 匹配 URL query 中可能泄露凭证的参数：
// key / api_key / api-key / access_token / token / authorization / x-api-key。
// 大小写不敏感，匹配 `?name=value` 或 `&name=value` 形式（value 截到 & 或字符串末尾）。
var monitorSensitiveQueryParamRegex = regexp.MustCompile(`(?i)([?&](?:key|api[_-]?key|access[_-]?token|token|authorization|x-api-key)=)[^&\s"']+`)

// monitorAPIKeyPatterns 匹配常见 provider 的 API key 字面量。
// 顺序敏感：sk-ant- 必须放在 sk- 之前，否则会被通用 sk- 模式先消费。
var monitorAPIKeyPatterns = []struct {
	pattern *regexp.Regexp
	replace string
}{
	// Anthropic（带前缀，必须先匹配）：sk-ant-xxxxxxx
	{regexp.MustCompile(`sk-ant-[A-Za-z0-9_-]{20,}`), "sk-ant-***REDACTED***"},
	// OpenAI / Anthropic 通用 sk-: sk-xxxxxxx
	{regexp.MustCompile(`sk-[A-Za-z0-9-]{20,}`), "sk-***REDACTED***"},
	// xAI API Key：xai-xxxxxxx
	{regexp.MustCompile(`xai-[A-Za-z0-9_-]{6,}`), "xai-***REDACTED***"},
	// Gemini / Google API Key：固定前缀 + 35 位
	{regexp.MustCompile(`AIza[A-Za-z0-9_-]{35}`), "AIza***REDACTED***"},
	// JWT 三段式（Bearer 后常出现）：eyJxxx.eyJxxx.signature
	{regexp.MustCompile(`eyJ[A-Za-z0-9_-]{8,}\.eyJ[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}`), "eyJ***REDACTED.JWT***"},
}

// sanitizeErrorMessage 擦除错误/响应文本中可能泄露的 API key。
// 处理两类来源：
//  1. URL query 中的 ?key= / ?api_key= 等（Go *url.Error 会回填完整 URL）
//  2. 上游 HTTP body 文本里直接出现的 sk-* / xai-* / AIza* / JWT 等密钥碎片
//
// 注意：与 gemini_messages_compat_service.go 的 sanitizeUpstreamErrorMessage 关注点类似但参数集更广，
// 监控模块独立维护，避免互相耦合。
func sanitizeErrorMessage(msg string) string {
	if msg == "" {
		return msg
	}
	msg = monitorSensitiveQueryParamRegex.ReplaceAllString(msg, `${1}REDACTED`)
	for _, p := range monitorAPIKeyPatterns {
		msg = p.pattern.ReplaceAllString(msg, p.replace)
	}
	return msg
}

func formatMonitorUpstreamHTTPError(statusCode int, headers http.Header, rawBody string) string {
	body := []byte(rawBody)
	if httputil.IsCloudflareChallengeResponse(statusCode, headers, body) {
		msg := httputil.FormatCloudflareChallengeMessage(
			fmt.Sprintf("upstream HTTP %d: Cloudflare challenge / access restricted by upstream", statusCode),
			headers,
			body,
		)
		return truncateMessage(sanitizeErrorMessage(msg))
	}

	bodySnippet := truncateForErrorBody(rawBody)
	return truncateMessage(sanitizeErrorMessage(fmt.Sprintf("upstream HTTP %d: %s", statusCode, bodySnippet)))
}

// truncateMessage 把消息按 monitorMessageMaxBytes 截断，避免 DB 列溢出与日志过长。
func truncateMessage(msg string) string {
	if len(msg) <= monitorMessageMaxBytes {
		return msg
	}
	const ellipsis = "...(truncated)"
	cutoff := monitorMessageMaxBytes - len(ellipsis)
	if cutoff < 0 {
		cutoff = 0
	}
	return msg[:cutoff] + ellipsis
}

// truncateForErrorBody 把上游错误响应 body 压到 monitorErrorBodySnippetMaxBytes 以内，
// 并顺手把连续空白折成一个空格：上游 HTML 错误页常含大量缩进/换行，保留会浪费预算。
// 被 truncateMessage 做最终总截断兜底，所以这里只负责 body 自身的精简。
func truncateForErrorBody(body string) string {
	body = strings.Join(strings.Fields(body), " ")
	if len(body) <= monitorErrorBodySnippetMaxBytes {
		return body
	}
	const ellipsis = "...(body truncated)"
	cutoff := monitorErrorBodySnippetMaxBytes - len(ellipsis)
	if cutoff < 0 {
		cutoff = 0
	}
	return body[:cutoff] + ellipsis
}
