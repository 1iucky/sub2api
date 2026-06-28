# Quality Guidelines

> Code quality standards for frontend development.

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

### Region Block Redirect Contract

When page-only region blocking is enforced at the edge, SPA navigation does not
reload HTML and therefore does not automatically pass through Nginx again.
Frontend API calls must handle region-block responses centrally in
`frontend/src/api/client.ts`.

**Contract**:
- If an API response is redirected to `/unsupported-region` and returns HTML,
  reject the request with `code: 'UNSUPPORTED_REGION'` and redirect the browser
  to the received `/unsupported-region?...` path.
- If an API error returns HTTP `403` or `451` with `code`, `reason`, `error`, or
  `message` containing `unsupported_region`, `region_blocked`, or
  `unsupported region`, reject with `code: 'UNSUPPORTED_REGION'` and redirect to
  `/unsupported-region`.
- Do not redirect when the current path is already `/unsupported-region`.

**Tests required**:
- API client test for HTML redirect responses whose `responseURL` points at
  `/unsupported-region`.
- Route test proving `/unsupported-region` remains a public route, including
  backend-mode public-route restrictions.

---

## Testing Requirements

<!-- What level of testing is expected -->

(To be filled by the team)

---

## Code Review Checklist

<!-- What reviewers should check -->

(To be filled by the team)
