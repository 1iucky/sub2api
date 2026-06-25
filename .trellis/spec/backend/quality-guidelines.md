# Quality Guidelines

> Code quality standards for backend development.

---

## Overview

<!--
Document your project's quality standards here.

Questions to answer:
- What patterns are forbidden?
- What linting rules do you enforce?
- What are your testing requirements?
- What code review standards apply?
-->

(To be filled by the team)

---

## Forbidden Patterns

<!-- Patterns that should never be used and why -->

(To be filled by the team)

---

## Required Patterns

<!-- Patterns that must always be used -->

(To be filled by the team)

---

## Testing Requirements

<!-- What level of testing is expected -->

### Scenario: Channel Monitor Provider Response Parsing

#### 1. Scope / Trigger
- Trigger: Changes to channel monitor health checks, provider adapters, or gateway-backed monitor endpoints.
- The monitor request may be sent as OpenAI Chat, OpenAI Responses, Anthropic Messages, or Gemini, but the response body can still be JSON or SSE depending on the upstream gateway and compatibility layer.

#### 2. Signatures
- `runCheckForModel(ctx, provider, endpoint, apiKey, model, opts)` must validate the challenge against extracted assistant text, not against a single hard-coded JSON path.
- `callProvider(...)` should return both extracted text and raw body/status/header data so success parsing and error formatting stay separate.

#### 3. Contracts
- OpenAI Chat monitor responses may be:
  - Chat JSON: `choices[].message.content`
  - Responses JSON: `output[].content[].text`
  - SSE data frames containing either Chat deltas or Responses completion events
- Anthropic Messages monitor responses may be:
  - Messages JSON: `content[].text`
  - SSE frames containing `message_start`, `content_block_start`, `content_block_delta`, `message_delta`, and `message_stop`
- Error paths must use the raw response body and headers, not extracted assistant text, because provider error bodies usually do not match the success text paths.

#### 4. Validation & Error Matrix
- `2xx + extracted text matches challenge` -> operational or degraded based on latency.
- `2xx + extracted text empty/mismatched` -> failed with `challenge mismatch`.
- `replace body override + 2xx + non-empty extracted text` -> operational/degraded without challenge validation.
- `replace body override + 2xx + empty extracted text` -> failed.
- `non-2xx + Cloudflare challenge headers/body` -> error with sanitized Cloudflare-friendly message.
- `non-2xx + normal upstream body` -> error with sanitized raw body snippet.

#### 5. Good/Base/Bad Cases
- Good: Add provider parsing with JSON and SSE regression tests for the same challenge flow.
- Base: Preserve the provider adapter's default request body and headers when parsing response variants.
- Bad: Parsing only `adapter.textPath` for a provider that can be routed through the local gateway compatibility layer.

#### 6. Tests Required
- Add unit tests around `runCheckForModel` for each newly supported response shape.
- Assert both monitor status and request path/body/header behavior where the provider adapter is part of the change.
- Include at least one error-formatting regression when changing raw-body handling.

#### 7. Wrong vs Correct

Wrong:
```go
return gjson.GetBytes(respBytes, adapter.textPath).String()
```

Correct:
```go
if provider == MonitorProviderAnthropic {
    return extractAnthropicText(respBytes)
}
```

### Scenario: Channel Monitor Retry and Model Catalog Linkage

#### 1. Scope / Trigger
- Trigger: Changes to channel monitor create/update/list contracts, monitor execution retry behavior, or model catalog marketplace availability.
- This flow crosses database schema, service validation, admin API payloads, public monitor views, and frontend marketplace matching.

#### 2. Signatures
- `ChannelMonitor.RetryCount int` persists as `channel_monitors.retry_count` and is serialized as `retry_count`.
- `ChannelMonitorCreateParams.RetryCount` and `ChannelMonitorUpdateParams.RetryCount` must be validated before persistence.
- Model marketplace monitor matching must compare catalog `model_id` against monitor `primary_model` and every `extra_models[].model`.

#### 3. Contracts
- `retry_count` means additional attempts after the first failed monitor check; `0` performs only one attempt.
- Valid retry range is `0..5`; default is `0`.
- Marketplace availability must not infer monitor links from provider/platform, pricing platform, display name, or channel name.
- A single monitor applies to all catalog models whose `model_id` equals its primary model or any attached extra model.

#### 4. Validation & Error Matrix
- `retry_count < 0` or `retry_count > 5` -> `CHANNEL_MONITOR_INVALID_RETRY_COUNT`.
- First attempt success -> no retry.
- Attempt fails and retry budget remains -> execute the same model check again.
- All attempts fail -> return the final failed check result.
- No monitor matches a marketplace model by `model_id` -> frontend displays the default fully available status.

#### 5. Good/Base/Bad Cases
- Good: Store `retry_count` in the ent schema, migration, repository create/update, service params, and admin handler DTOs together.
- Base: Keep retry execution local to monitor checks so scheduler/concurrency behavior remains unchanged.
- Bad: Matching marketplace monitor data by provider/platform, because models such as DeepSeek-compatible Anthropic endpoints may intentionally use a different protocol endpoint.

#### 6. Tests Required
- Unit test retry helper behavior: succeeds after retry, stops on first success, and exhausts retry budget.
- Regression test marketplace monitor matching by exact `model_id`, including extra models.
- Type-check frontend after changing monitor DTO fields.

#### 7. Wrong vs Correct

Wrong:
```ts
monitor.provider === model.platform && monitor.primary_model === model.model_id
```

Correct:
```ts
monitor.primary_model === model.model_id ||
monitor.extra_models.some(extra => extra.model === model.model_id)
```

---

## Code Review Checklist

<!-- What reviewers should check -->

(To be filled by the team)
