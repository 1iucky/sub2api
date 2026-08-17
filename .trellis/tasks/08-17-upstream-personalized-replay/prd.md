# Replay Personalized Changes onto Latest Upstream

## Goal

Complete the ordered replay of local personalized commits onto the latest `upstream/main` baseline while preserving upstream business correctness and the local product's complete frontend branding and Factory visual system.

## Requirements

- Continue from the existing `codex/upstream-main-personalized-20260814` branch and the in-progress Cherry-pick without resetting or replacing prior replay work.
- Treat upstream backend core business logic and contracts as authoritative; integrate local backend capabilities only as compatible, orthogonal extensions.
- Preserve upstream frontend features, fields, events, responsive behavior, locale structure, and newly added pages.
- Apply the local personalized Factory style and branding consistently to both existing and newly introduced upstream frontend surfaces.
- Keep user-facing brand text dynamic through site settings or SiliconBase branding; do not rename technical compatibility identifiers such as `SUB2API_API_KEY`.
- Replay remaining personalized commits in source order, resolving each conflict by understanding both sides rather than selecting an entire old file.
- Stop and request a product decision only when backend or product behavior has a severe, non-derivable semantic conflict.
- Preserve the backup branch and existing stash entries throughout the replay.

## Acceptance Criteria

- [ ] All intended personalized non-merge commits are replayed in order or explicitly documented as obsolete because upstream already contains an equivalent change.
- [ ] No unresolved conflict markers or unmerged paths remain.
- [ ] Upstream-added frontend behavior remains present in every conflict-resolved component.
- [ ] Personalized Factory styling and branding remain consistent across affected existing and upstream-added pages.
- [ ] Generated Ent and Wire outputs are refreshed when required by replayed backend commits.
- [ ] Frontend type-check and relevant tests pass.
- [ ] Backend targeted and broad tests pass, excluding only clearly unrelated pre-existing failures.
- [ ] The pre-replay Trellis runtime stash is restored without dropping any other stash or backup branch.

## Definition of Done

- The Cherry-pick sequence is complete and the branch is clean except for intentional Trellis runtime bookkeeping.
- `git diff --check`, conflict-marker scans, frontend verification, and backend verification are green.
- Final branch/commit counts and any skipped or manually adapted commits are reported.
- No remote push, branch deletion, backup deletion, or stash deletion is performed.

## Technical Approach

Use ordered `git cherry-pick -x` replay from `codex/custom-theme-upstream-main-20260717`. For each conflict, compare the upstream-stage version, personalized source commit, current resolved tree, and surrounding consumers. Retain upstream structure and behavior, then port the local Factory presentation into that structure. Verify after coherent frontend/backend batches instead of waiting until the entire replay ends.

## Decision (ADR-lite)

**Context**: Directly merging a long-lived personalized branch into current upstream produced a broad, high-risk conflict surface.

**Decision**: Use fresh-upstream ordered replay with an untouched backup anchor. Resolve backend conflicts upstream-first and frontend visual conflicts personalization-first while preserving upstream functionality.

**Consequences**: The replay requires more manual conflict analysis but gives commit-level rollback, clearer verification boundaries, and substantially lower risk of silently losing either upstream features or personalized behavior.

## Out of Scope

- Pushing the integration branch or opening a pull request.
- Deleting source, backup, or historical integration branches.
- Redesigning business behavior not touched by the replayed commits.
- Restoring local features whose upstream backend contract was intentionally removed, unless the user explicitly chooses to restore that contract.

## Technical Notes

- Baseline: `upstream/main` at `fbfdcef81` when this replay branch was created.
- Current replay branch: `codex/upstream-main-personalized-20260814`.
- Backup anchor: `codex/backup-upstream-batched-20260814` at `ac54093ae`.
- Personalized source: `codex/custom-theme-upstream-main-20260717`.
- Current paused source commit: `b97bb87e1e`.
- The current conflict set is limited to admin usage filters and three shared distribution charts.
- Split locale files under `frontend/src/i18n/locales/en/` and `zh/` must remain the canonical structure.
