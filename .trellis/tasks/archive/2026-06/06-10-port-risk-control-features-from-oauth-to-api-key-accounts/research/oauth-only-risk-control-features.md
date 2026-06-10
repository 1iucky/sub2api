# Research: OAuth-Only Risk Control Features

- **Query**: Identify ALL risk control features restricted to Anthropic/Claude OAuth accounts, understand how they work, and what would need to change to support API Key accounts.
- **Scope**: internal
- **Date**: 2026-06-10

## Findings

### Files Found

| File Path | Description |
|---|---|
| `backend/internal/service/account.go` | Account model, type checks, IsTLSFingerprintEnabled, IsSessionIDMaskingEnabled, IsAnthropicOAuthOrSetupToken |
| `backend/internal/service/gateway_service.go` | Main gateway forwarding logic, buildUpstreamRequest, Forward, applyClaudeCodeMimicHeaders, computeFinalAnthropicBeta |
| `backend/internal/service/identity_service.go` | Fingerprint creation/apply, RewriteUserID, RewriteUserIDWithMasking |
| `backend/internal/service/header_util.go` | Header wire casing map, resolveWireCasing, setHeaderRaw |
| `backend/internal/service/metadata_userid.go` | ParseMetadataUserID, FormatMetadataUserID (both formats) |
| `backend/internal/service/tls_fingerprint_profile_service.go` | ResolveTLSProfile, gated by IsTLSFingerprintEnabled |
| `backend/internal/pkg/claude/constants.go` | DefaultHeaders, CLICurrentVersion, Beta constants, APIKeyBetaHeader |
| `backend/internal/domain/constants.go` | AccountTypeOAuth, AccountTypeAPIKey, AccountTypeSetupToken, etc. |
| `backend/internal/service/setting_service.go` | GetGatewayForwardingSettings (fp, mpt, cch toggles) |

---

## Feature 1: TLS Fingerprint

**File**: `backend/internal/service/account.go:1583-1597`

**Current gate** (line 1583):
```go
func (a *Account) IsTLSFingerprintEnabled() bool {
    // 仅支持 Anthropic OAuth/SetupToken 账号
    if !a.IsAnthropicOAuthOrSetupToken() {
        return false
    }
    ...
}
```

`IsAnthropicOAuthOrSetupToken()` is defined at line 1576:
```go
func (a *Account) IsAnthropicOAuthOrSetupToken() bool {
    return a.Platform == PlatformAnthropic && (a.Type == AccountTypeOAuth || a.Type == AccountTypeSetupToken)
}
```

**What it does**: When enabled (`Extra["enable_tls_fingerprint"] == true`), the TLS handshake is modified to mimic a Node.js Claude CLI client. The `ResolveTLSProfile` method in `tls_fingerprint_profile_service.go:177` calls `account.IsTLSFingerprintEnabled()` and returns `nil` if false (no TLS mimicry). A non-nil profile is passed to `s.httpUpstream.DoWithTLS()` at every request dispatch point (lines 4715, 5245, 9491, 9518, 9615).

The profile can be bound to a specific TLS fingerprint template via `Extra["tls_fingerprint_profile_id"]`, or use `-1` for random selection from a pool.

**Where used**:
- `gateway_service.go:4689` - main Forward path
- `gateway_service.go:5245` - API Key passthrough path (currently always nil for API Key)
- `gateway_service.go:9491,9518` - count_tokens path
- `gateway_service.go:9615` - count_tokens API Key passthrough path

**Safe for API Key?**: Technically yes -- TLS fingerprinting is a transport-layer feature that makes the TCP/TLS handshake look like a real browser/CLI. It has no dependency on the authentication method. The `DoWithTLS` function accepts any account. The only blocker is the guard in `IsTLSFingerprintEnabled()` which hardcodes `IsAnthropicOAuthOrSetupToken()`.

---

## Feature 2: Claude Code Mimic Headers

**File**: `backend/internal/service/gateway_service.go:7023-7050`

**Function signature**:
```go
func applyClaudeCodeMimicHeaders(req *http.Request, isStream bool)
```

**What it does**: Forces request headers to match Claude Code CLI fingerprint:
1. Calls `applyClaudeOAuthHeaderDefaults(req)` to fill missing defaults
2. Overwrites all `claude.DefaultHeaders` keys (User-Agent, X-Stainless-*, X-App, etc.)
3. Sets `Accept: application/json`
4. Sets `x-stainless-helper-method: stream` for streaming
5. Generates a fresh UUID for `x-client-request-id`

**Call sites and conditions**:
- **Line 6405-6406** (buildUpstreamRequest): `if tokenType == "oauth" && mimicClaudeCode { applyClaudeCodeMimicHeaders(req, reqStream) }`
- **Line 9883-9884** (buildCountTokensRequest): `if tokenType == "oauth" && mimicClaudeCode { applyClaudeCodeMimicHeaders(req, false) }`

The `mimicClaudeCode` flag is set at:
- **Line 4564**: `shouldMimicClaudeCode := account.IsOAuth() && !isClaudeCode`
- **Line 9405**: `shouldMimicClaudeCode := account.IsOAuth() && !isClaudeCodeCT`

**Guard chain**: `account.IsOAuth()` returns true only for `AccountTypeOAuth` or `AccountTypeSetupToken` (account.go:155-157). So mimic headers are **never applied for API Key accounts**, regardless of any other setting.

Additionally, the OAuth mimicry path **skips client header passthrough** at line 6375:
```go
if tokenType != "oauth" || !mimicClaudeCode {
    // ... whitelist header passthrough
}
```
This means the OAuth+Mimic path deliberately drops all client headers and uses only its own injected headers.

**Safe for API Key?**: Partially. The mimic headers inject a Claude Code User-Agent and OAuth-related headers (like `Anthropic-Dangerous-Direct-Browser-Access`). For API Key accounts using standard Anthropic API, the User-Agent and x-stainless-* headers could help avoid detection, but the `x-app: cli` header and OAuth-specific defaults may be counterproductive. A modified version would need to skip OAuth-specific headers while keeping the fingerprint benefits.

---

## Feature 3: Identity/Fingerprint System

**File**: `backend/internal/service/identity_service.go`

**The guard in buildUpstreamRequest** (gateway_service.go:6301):
```go
if account.IsOAuth() && s.identityService != nil {
    // 1. GetOrCreateFingerprint
    // 2. RewriteUserIDWithMasking (metadata.user_id rewrite)
}
```

Same guard in buildCountTokensRequest (gateway_service.go:9806):
```go
if account.IsOAuth() && s.identityService != nil {
    // fingerprint + metadata rewrite
}
```

**What it does**:
1. **Fingerprint creation** (identity_service.go:78-120): Generates or retrieves a cached `Fingerprint` struct containing ClientID, UserAgent, and X-Stainless-* headers. The ClientID is a random 64-char hex string cached per account (7-day TTL). The fingerprint is created from the first request's headers or defaults to `claude-cli/2.1.161` style values.

2. **Fingerprint application** (identity_service.go:179-208): `ApplyFingerprint(req, fp)` overwrites request headers with cached fingerprint values (User-Agent, X-Stainless-Lang/Package-Version/OS/Arch/Runtime/RuntimeVersion). Uses `setHeaderRaw` for exact wire casing.

3. **Metadata user_id rewrite** (identity_service.go:270-300): `RewriteUserIDWithMasking` rewrites the `metadata.user_id` field in the request body. It replaces the device_id, account_uuid, and session_id with deterministic values derived from the account. If session ID masking is enabled, the session is fixed for 15 minutes.

4. **Gateway forwarding settings** control fingerprint application:
   - `enableFP` (fingerprint unification): When true, fingerprint is applied to request headers
   - `enableMPT` (metadata passthrough): When true, metadata rewrite is skipped
   - `enableCCH` (CCH signing): When true, billing header CCH is signed

**Safe for API Key?**: The fingerprint system itself is account-type agnostic. It caches X-Stainless-* headers per account and applies them consistently. The metadata user_id rewrite requires an `account_uuid` which exists in Extra for OAuth accounts (from the Claude session). For API Key accounts, this field may not exist, but the rewrite function gracefully handles missing values. The fingerprint application (header-only) would work for any account type.

---

## Feature 4: Session ID

**How session_id is generated**: There are two contexts:

### 4a. Sticky Session Hash (routing)
**File**: `gateway_service.go:724-741`

`GenerateSessionHash` tries in order:
1. Extract `session_id` from `metadata.user_id` (if the client is Claude Code)
2. Hash cacheable content with `cache_control: ephemeral`
3. Fallback: hash session context + system + all messages

This works for all account types -- no OAuth gate.

### 4b. Mimicry Session ID (in metadata.user_id)
**File**: `gateway_service.go:1271-1306` (buildOAuthMetadataUserID)

This is only called in the `shouldMimicClaudeCode` path:
```go
shouldMimicClaudeCode := account.IsOAuth() && !isClaudeCode
```

Inside `buildOAuthMetadataUserID`:
- Uses `buildStableSessionSeed` (line 1437) with accountID + client discriminator + first user text
- Derives a deterministic UUID via SHA-256
- Formats as `user_{deviceid}_account_{uuid}_session_{sessionid}` (legacy) or JSON (new format >=2.1.78)

### 4c. Session ID Masking
**File**: `account.go:1648-1661`

```go
func (a *Account) IsSessionIDMaskingEnabled() bool {
    if !a.IsAnthropicOAuthOrSetupToken() {
        return false
    }
    ...
}
```

Gated to OAuth/SetupToken only. When enabled, the session ID in `metadata.user_id` is replaced with a fixed value for 15 minutes (identity_service.go:278).

**Safe for API Key?**: The stable session seed generation is purely deterministic (accountID + client context + first user text). It doesn't depend on OAuth. The session ID masking cache mechanism (GetMaskedSessionID/SetMaskedSessionID) is also account-type agnostic. The only blockers are the guards in `IsSessionIDMaskingEnabled` and the `account.IsOAuth()` check in buildUpstreamRequest.

---

## Feature 5: OAuth Token Handling

**File**: `gateway_service.go:3874-3922`

**GetAccessToken** dispatches by account type:
- `AccountTypeOAuth, AccountTypeSetupToken` -> `getOAuthToken()` which returns `(token, "oauth", nil)`
- `AccountTypeAPIKey` -> returns `(apiKey, "apikey", nil)`
- `AccountTypeBedrock` -> returns `("", "bedrock", nil)`
- `AccountTypeServiceAccount` -> uses `claudeTokenProvider.GetAccessToken()`

The token type string ("oauth" vs "apikey") is then used as the primary discriminator throughout the codebase:

- **Line 6364**: Auth header format: `Bearer` for oauth, `x-api-key` for apikey
- **Line 6375**: Header passthrough skipped for `tokenType == "oauth" && mimicClaudeCode`
- **Line 6399**: `applyClaudeOAuthHeaderDefaults(req)` only when `tokenType == "oauth"`
- **Line 6405**: `applyClaudeCodeMimicHeaders` only when `tokenType == "oauth" && mimicClaudeCode`
- **Line 6670**: `computeFinalAnthropicBeta` branches on `tokenType == "oauth"` vs API-key path
- **Line 6438**: Debug info only captured for `tokenType == "oauth"`

The OAuth token refresh is handled externally by a `TokenRefreshService` that runs in the background. It's not inline in the request path.

**For API Key accounts**: There is no token refresh -- the API key is static. The `tokenType == "apikey"` branch just reads the key from credentials and passes it through as `x-api-key` header.

---

## Feature 6: Header Wire Casing

**File**: `backend/internal/service/header_util.go`

**What it does**: Maintains a map (`headerWireCasing`) of canonical-lowercase header names to their exact wire-format casing as observed in real Claude CLI traffic captures. For example:
- `"x-stainless-os"` -> `"X-Stainless-OS"` (not Go's canonical `"X-Stainless-Os"`)
- `"anthropic-beta"` -> `"anthropic-beta"` (lowercase, not `"Anthropic-Beta"`)
- `"x-app"` -> `"x-app"` (lowercase)

Also maintains `headerWireOrder` for debug logging (matching real capture order).

**Where used**:
- `resolveWireCasing(key)` is called in whitelist passthrough (lines 6379, 6468, 6482, 9743, 9859)
- `applyClaudeCodeMimicHeaders` uses `resolveWireCasing(key)` for all injected headers (line 7038)
- `applyClaudeOAuthHeaderDefaults` uses it (line 6590)

**Safe for API Key?**: Yes. This is a pure utility function with no account-type dependency. It would apply equally to API Key accounts to ensure headers match real CLI traffic patterns.

---

## Feature 7: All Account Type Checks in gateway_service.go

Complete list of account type checks that gate features:

| Line | Check | Purpose | Feature Gated |
|---|---|---|---|
| 2945 | `accounts[idx].account.Type == AccountTypeOAuth` | Scheduling: preferOAuth filter | Account selection |
| 2972 | `a.Type == AccountTypeOAuth` | Sort priority tiebreaker | Account selection |
| 3035 | `acc.Type == AccountTypeOAuth` | Same-priority OAuth preference | Account selection |
| 3100 | `a.Type == AccountTypeOAuth` | Sort within priority groups | Account selection |
| 3245 | `acc.Type == AccountTypeOAuth` | preferOAuth selection in chat completions path | Account selection |
| 3359 | `acc.Type == AccountTypeOAuth` | preferOAuth selection in responses path | Account selection |
| 3505 | `acc.Type == AccountTypeOAuth` | Gemini-specific OAuth preference | Account selection |
| 3620 | `acc.Type == AccountTypeOAuth` | Gemini-specific OAuth preference | Account selection |
| 3863 | `account.Type != AccountTypeAPIKey` | Model normalization (short->long ID) | Model mapping |
| 3864 | `account.Type == AccountTypeServiceAccount` | Vertex model normalization | Model mapping |
| 3877 | `AccountTypeOAuth, AccountTypeSetupToken` | GetAccessToken dispatch | Token acquisition |
| 3880 | `AccountTypeAPIKey` | GetAccessToken dispatch | Token acquisition |
| 3907 | `account.Type == AccountTypeOAuth` | ClaudeTokenProvider for Anthropic OAuth | Token refresh |
| 4633 | `account.Type == AccountTypeAPIKey` | Use account-level model mapping | Model mapping |
| 4639 | `account.Type == AccountTypeServiceAccount` | Vertex model mapping | Model mapping |
| 4651 | `account.Type != AccountTypeAPIKey` | Standard Anthropic model normalization | Model mapping |
| 6262 | `account.Type == AccountTypeServiceAccount` | Vertex request builder | Request building |
| 6269 | `account.Type == AccountTypeAPIKey` | Custom base URL for API Key | Request building |
| 6301 | `account.IsOAuth()` | **Fingerprint + metadata rewrite** | Identity system |
| 6364 | `tokenType == "oauth"` | Auth header format (Bearer vs x-api-key) | Authentication |
| 6375 | `tokenType != "oauth" \|\| !mimicClaudeCode` | Skip header passthrough for OAuth mimicry | Header passthrough |
| 6399 | `tokenType == "oauth"` | `applyClaudeOAuthHeaderDefaults` | Header defaults |
| 6405 | `tokenType == "oauth" && mimicClaudeCode` | `applyClaudeCodeMimicHeaders` | Mimicry |
| 6438 | `tokenType == "oauth"` | Debug info capture | Debugging |
| 6670 | `tokenType == "oauth"` | Beta header computation branch | Beta headers |
| 6806 | `account.IsOAuth()` | Beta policy scope matching | Policy evaluation |
| 6871 | `betaPolicyScopeMatches(scope, isOAuth, isBedrock)` | Beta policy scope (OAuth/APIKey/Bedrock/All) | Policy evaluation |
| 6988 | `account.IsOAuth()` | Beta policy block checking | Policy evaluation |
| 7069 | `account.Type == AccountTypeAPIKey` | Signature rectifier settings for API Key | Error rectification |
| 7091 | `account.Type == AccountTypeAPIKey` | Signature error pattern matching | Error rectification |
| 9404 | `shouldMimicClaudeCode := account.IsOAuth() && !isClaudeCodeCT` | Mimicry in count_tokens | Body normalization |
| 9442 | `account.Type == AccountTypeAPIKey` | Model mapping in count_tokens | Model mapping |
| 9448 | `account.Type != AccountTypeAPIKey` | Standard Anthropic normalization in count_tokens | Model mapping |
| 9773 | `account.Type == AccountTypeAPIKey` | Custom base URL for count_tokens | Request building |
| 9806 | `account.IsOAuth()` | Fingerprint in count_tokens | Identity system |
| 9878 | `tokenType == "oauth"` | applyClaudeOAuthHeaderDefaults in count_tokens | Header defaults |
| 9883 | `tokenType == "oauth" && mimicClaudeCode` | applyClaudeCodeMimicHeaders in count_tokens | Mimicry |
| 9902 | `tokenType == "oauth"` | Debug info in count_tokens | Debugging |

Additionally in `account.go`:
| Line | Check | Purpose |
|---|---|---|
| 155-157 | `IsOAuth()` | Returns true for OAuth or SetupToken |
| 1576-1577 | `IsAnthropicOAuthOrSetupToken()` | Platform == Anthropic AND (OAuth or SetupToken) |
| 1583-1585 | `IsTLSFingerprintEnabled()` | Gated by `IsAnthropicOAuthOrSetupToken()` |
| 1648-1650 | `IsSessionIDMaskingEnabled()` | Gated by `IsAnthropicOAuthOrSetupToken()` |
| 1665-1667 | `IsCustomBaseURLEnabled()` | Gated by `IsAnthropicOAuthOrSetupToken()` |
| 1688-1690 | `IsCacheTTLOverrideEnabled()` | Gated by `IsAnthropicOAuthOrSetupToken()` |
| 734 | `GetBaseURL()` | Only returns non-empty for `AccountTypeAPIKey` |

---

## Feature 8: Main Request Forwarding Flow

### Entry Point: `Forward()` (line 4485)

The flow branches early based on account type:

```
Forward()
  |
  +-- IsAnthropicAPIKeyPassthroughEnabled? -> forwardAnthropicAPIKeyPassthroughWithInput()
  |     (line 4496-4513)
  |     - Only for API Key accounts with anthropic_passthrough enabled
  |     - buildUpstreamRequestAnthropicAPIKeyPassthrough (line 5429)
  |     - Pure header passthrough, no fingerprint/mimicry
  |     - No metadata rewrite
  |     - TLS fingerprint passed but IsTLSFingerprintEnabled returns false for API Key
  |
  +-- IsBedrock? -> forwardBedrock()
  |     (line 4516-4518)
  |
  +-- Standard path (OAuth / API Key non-passthrough):
        |
        +-- Beta policy evaluation (line 4520-4532)
        |
        +-- shouldMimicClaudeCode = account.IsOAuth() && !isClaudeCode
        |     (line 4564)
        |
        +-- IF shouldMimicClaudeCode:
        |     +-- rewriteSystemForNonClaudeCode (system prompt rewrite)
        |     +-- GetOrCreateFingerprint (identity_service)
        |     +-- buildOAuthMetadataUserID
        |     +-- normalizeClaudeOAuthRequestBody (billing attribution block)
        |     +-- rewriteMessageCacheControlIfEnabled
        |     +-- applyToolNameRewriteToBody / applyToolsLastCacheBreakpoint
        |
        +-- enforceCacheControlLimit
        |
        +-- Model mapping (different paths for API Key vs OAuth)
        |     API Key: account.GetMappedModel (explicit mapping)
        |     OAuth/SetupToken: claude.NormalizeModelID (standard Anthropic mapping)
        |
        +-- GetAccessToken -> token, tokenType
        |
        +-- buildUpstreamRequest (line 6261)
              |
              +-- Fingerprint + metadata rewrite: only if account.IsOAuth()
              |     (line 6301)
              |
              +-- computeFinalAnthropicBeta: branches on tokenType
              |     "oauth" -> FullClaudeCodeMimicryBetas or getBetaHeader
              |     "apikey" -> strip client beta or inject APIKeyBetaHeader
              |
              +-- Header passthrough: skipped for oauth+mimic
              |
              +-- ApplyFingerprint: only if fingerprint != nil (set by OAuth path)
              |
              +-- applyClaudeOAuthHeaderDefaults: only if tokenType == "oauth"
              |
              +-- applyClaudeCodeMimicHeaders: only if tokenType == "oauth" && mimicClaudeCode
              |
              +-- X-Claude-Code-Session-Id sync from body
```

### API Key Passthrough Path (forwardAnthropicAPIKeyPassthroughWithInput, line 5194)

This is a separate, simpler path:
- No fingerprint application
- No metadata user_id rewrite
- No system prompt rewrite
- No Claude Code mimicry
- Direct header passthrough with `x-api-key` auth
- `buildUpstreamRequestAnthropicAPIKeyPassthrough` (line 5429) is a stripped-down request builder

---

## Summary: OAuth-Gated Risk Control Features

### Hard-gated by `IsAnthropicOAuthOrSetupToken()` (account.go:1576):

1. **TLS Fingerprint** (`IsTLSFingerprintEnabled`) - line 1583
2. **Session ID Masking** (`IsSessionIDMaskingEnabled`) - line 1648
3. **Custom Base URL** (`IsCustomBaseURLEnabled`) - line 1665
4. **Cache TTL Override** (`IsCacheTTLOverrideEnabled`) - line 1688

### Hard-gated by `account.IsOAuth()` in gateway_service.go:

5. **Identity Fingerprint** (GetOrCreateFingerprint + ApplyFingerprint) - line 6301, 9806
6. **Metadata user_id Rewrite** (RewriteUserIDWithMasking) - line 6318, 9813
7. **OAuth Header Defaults** (applyClaudeOAuthHeaderDefaults) - line 6399, 9878
8. **Claude Code Mimic Headers** (applyClaudeCodeMimicHeaders) - line 6405, 9883
9. **Claude Code Body Mimicry** (system rewrite, metadata injection, billing block) - line 4566-4621
10. **OAuth Beta Header Computation** (FullClaudeCodeMimicryBetas) - line 6670-6678

### Hard-gated by `tokenType == "oauth"`:

11. **Debug info capture** (buildClaudeMimicDebugLine) - line 6438, 9902
12. **Header passthrough skip** (for mimic path) - line 6375

---

## Caveats / Not Found

- The `claude.DefaultHeaders` map includes `Anthropic-Dangerous-Direct-Browser-Access: true` which is specific to OAuth/browser-based auth. For API Key accounts, this header may be irrelevant or incorrect.
- The billing attribution block (`buildBillingAttributionBlockJSON`) is only injected in the mimic path, which requires OAuth. API Key accounts currently do not get billing attribution blocks.
- The `account_uuid` field in Extra is populated during OAuth login flows. For API Key accounts, this field may not exist, which would cause metadata user_id rewrite to silently skip.
- The `claude.ModelIDOverrides` (short-to-long ID mapping) is only applied for non-APIKey accounts (line 4651). API Key accounts use explicit account-level model mapping or pass through the raw model name.
- The scheduling layer has a `preferOAuth` preference that prioritizes OAuth accounts in selection -- this is orthogonal to risk control but worth noting.
