# Current Replay State

## Confirmed State

- The branch is based on the latest fetched upstream baseline used for this migration.
- Thirty-two personalized commits are already present on the replay branch.
- Cherry-pick is paused on `b97bb87e1e`, a frontend Factory restyle commit.
- Four files are unmerged: `UsageFilters.vue` and the endpoint, group, and model distribution charts.
- Other files from the same source commit are already staged by Git and must not be discarded.
- `stash@{0}` contains Trellis runtime files from before the replay and must be restored only after the replay and verification complete.

## Resolution Rules

1. Inspect all conflict hunks and current upstream-only additions.
2. Preserve upstream data flow, props, events, locale keys, responsive branches, and feature fields.
3. Port personalized Factory classes, compact geometry, warm gray surfaces, dashed dividers, and tabular/monospace numeric presentation into the upstream structure.
4. Stage only the intended source-commit files while Cherry-pick is active; do not accidentally stage Trellis task/runtime files.
5. Run frontend type-check after the current commit and after coherent frontend batches.
6. For backend schema or wiring changes, regenerate Ent and Wire outputs before final verification.

## Remaining Source Order After Current Commit

`e9958a27e6`, `4c81bb8bc6`, `59cdc7a30e`, `97371a2ad7`, `d20ce793a8`, `09ae7bd1f6`, `93506e785a`, `e1cfe6d8a8`, `4c3bb9a754`, `b47c51d738`, `99f2dedf88`, `6f2d201575`, `e78c32e810`, `030b397ec4`, `1df30e6913`, `76a13a2ea6`, `402d48a15c`, `e7854c37cb`, `a4d35edc8b`, `333dc3cfc3`, `e66b40dd5a`, `8e5c16b9fc`, `267930dc46`, `c40494024e`, `9bfd6d9078`, `56eda76b00`, `f3d4c6041f`, `786293ecb7`, `79e9a65dcd`, `a52346abe0`, `ae53a93090`, `c9e0675d9b`, `14491476d7`, `56652cd03f`, `312b954f62`, `9ce2b5b8d3`, `af24bf108c`, `0f819718e2`, `e06234fe4c`, `03310b5f0a`.
