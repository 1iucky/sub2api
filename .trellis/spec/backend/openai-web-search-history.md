# OAuth Web Search History Compatibility

## 1. Scope / Trigger

Upstream PR #7939 fixes `response protection is unavailable` when OpenAI OAuth
requests replay `web_search_call` history without declaring a search tool.
Local summary compaction and Responses Lite compaction can trigger this even
when other conversations on the same upstream account work.

## 2. Signatures

- `ensureOpenAIOAuthWebSearchToolForHistory(reqBody map[string]any, responsesLite bool) bool`
- `ensureOpenAIOAuthWebSearchToolForHistoryBody(body []byte, responsesLite bool) ([]byte, bool, error)`

## 3. Contracts

Keep search history intact. For affected OAuth requests, declare
`{"type":"web_search","external_web_access":false}` only if neither top-level
`tools` nor input `additional_tools` already declares a `web_search*` tool.
Standard Responses uses top-level tools; Lite uses `additional_tools`, before
any trailing `compaction_trigger`. With no caller tools and an absent, null,
empty, auto or none choice, set `tool_choice` to `none`.
API-key accounts and the dedicated `/responses/compact` shape remain unchanged.

## 4. Validation & Error Matrix

- Search history without a declaration: upstream may terminate with the error above.
- Lite with a top-level hosted tool: upstream rejects the tool placement.
- Existing declaration or no search history: compatibility helper is a no-op.

## 5. Good / Base / Bad Cases

- Good: search history plus a correctly placed cached-only declaration.
- Base: ordinary history passes through unchanged.
- Bad: retrying identical malformed history across accounts without normalization.

## 6. Tests Required

Run the unit-tagged `openai_oauth_web_search_history*` tests. Verify map/raw-body
parity, standard/Lite placement, trigger-last ordering, tool-choice preservation,
API-key/compact exclusions, and actual Forward wiring for transform/passthrough.

## 7. Wrong vs Correct

Wrong: infer that this message always proves an unavailable safety service or a
bad account. Correct: inspect replayed history and tool declarations; the error
also represents the request compatibility condition reproduced in PR #7939.
