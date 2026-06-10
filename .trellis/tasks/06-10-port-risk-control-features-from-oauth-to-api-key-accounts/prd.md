# Port Anthropic Risk Control Features from OAuth to API Key Accounts

## Goal

将当前仅对 Anthropic OAuth/SetupToken 账号生效的风控特征模拟功能全面扩展到 Anthropic API Key 类型账号，使 API Key 账号也能通过伪装为合法 Claude Code CLI 客户端来规避上游风控检测。

## Requirements

1. TLS 指纹模拟：API Key 账号可启用 TLS 指纹（Node.js 24.x 握手特征）
2. Identity 指纹：API Key 账号可使用 X-Stainless-* / UA 统一指纹
3. Session ID 伪装/遮蔽：API Key 账号可启用 metadata.user_id 重写和 session 固定
4. Claude Code Header 伪装：API Key 账号可使用 mimic headers（跳过 OAuth 专用头）
5. Body 模仿：API Key 账号走 system 重写、metadata 注入等 body mimicry
6. count_tokens 路径同步所有上述改动
7. API Key passthrough 路径支持 TLS 指纹
8. 现有 OAuth/SetupToken 账号行为不受影响

## Acceptance Criteria

- [ ] API Key 账号设置 `enable_tls_fingerprint: true` 后 TLS 指纹生效
- [ ] API Key 账号的转发请求中包含 X-Stainless-* 统一指纹
- [ ] API Key 账号的 User-Agent 被替换为 claude-cli 风格
- [ ] API Key 账号不注入 `Anthropic-Dangerous-Direct-Browser-Access` 等 OAuth 专用头
- [ ] API Key 账号认证头仍使用 `x-api-key` 而非 `Bearer`
- [ ] 现有 OAuth 账号行为完全不变
- [ ] `make test-unit` 通过

## Definition of Done

- 单元测试覆盖新增逻辑
- golangci-lint 无新增问题
- OAuth 账号行为回归无影响

## Out of Scope

- OpenAI/Gemini/Antigravity 平台
- OAuth Token 刷新机制
- 账号调度层 preferOAuth 逻辑
- 前端管理界面开关（后续任务）

## Technical Approach

### 核心策略

引入 `IsAnthropicAccount()` 方法替换 `IsAnthropicOAuthOrSetupToken()` 作为风控守卫，同时在 gateway_service.go 中将 `account.IsOAuth()` 检查扩展为包含 API Key。

### 关键改动

| # | 文件 | 改动 |
|---|---|---|
| 1 | account.go | 新增 `IsAnthropicAccount()` |
| 2 | account.go | `IsTLSFingerprintEnabled()` 改用 `IsAnthropicAccount()` |
| 3 | account.go | `IsSessionIDMaskingEnabled()` 改用 `IsAnthropicAccount()` |
| 4 | gateway_service.go:4564 | `shouldMimicClaudeCode` 扩展到 API Key |
| 5 | gateway_service.go:6301 | Identity 指纹 + metadata 重写扩展 |
| 6 | gateway_service.go:6375 | Header passthrough skip 扩展 |
| 7 | gateway_service.go:6399 | applyClaudeOAuthHeaderDefaults 扩展 |
| 8 | gateway_service.go:6405 | applyClaudeCodeMimicHeaders 扩展（跳过 OAuth 专用头） |
| 9 | gateway_service.go:6438 | Debug info 捕获扩展 |
| 10 | gateway_service.go:6670 | Beta header 计算分支扩展 |
| 11 | gateway_service.go:5245 | API Key passthrough 路径 TLS 指纹 |
| 12 | gateway_service.go:9404-9902 | count_tokens 路径同步 |

### 适配细节

- `applyClaudeCodeMimicHeaders` 新增条件：API Key 账号跳过 `Anthropic-Dangerous-Direct-Browser-Access` 和 `x-app: cli`
- metadata 重写中 accountUUID 缺失时用确定性替代值（账号 ID hash）
- API Key 认证头保持 `x-api-key` 不变

## Decision (ADR-lite)

**Context**: 风控特性被 OAuth 硬编码守卫限制，API Key 账号无法使用
**Decision**: 引入 IsAnthropicAccount() 统一守卫，按需条件跳过 OAuth 专用元素
**Consequences**: API Key 账号获得完整风控对抗能力；需要确保 OAuth 专用元素不泄露到 API Key 路径

## Research References

* [`research/oauth-only-risk-control-features.md`](research/oauth-only-risk-control-features.md) — 12 个 OAuth 专属风控特性的完整分析
