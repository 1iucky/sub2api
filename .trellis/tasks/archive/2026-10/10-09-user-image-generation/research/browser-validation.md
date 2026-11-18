# Browser verification evidence

Environment: native Codex in-app browser; Vite worktree frontend at `127.0.0.1:13000`, isolated Python API fixture at `127.0.0.1:18080`. All credentials, costs, records and generated pixels are test fixtures. No production account or real image provider was contacted. This verifies page interactions and rendered layout, not server authentication or billing.

Observed using DOM-based browser controls:

- User can open `/image-studio`; existing navigation/layout shows a separate Image Studio entry.
- Enter/reload and switching from GPT to Grok do not create Keys. Fixture state after the first generated image reported `key_creates:0`, `generations:1`.
- Grok selection changes controls to resolution/aspect ratio; GPT selection exposes size/quality/output format. The Grok POST contained `resolution:1k`, `aspect_ratio:16:9`, `n:1`, `response_format:b64_json`, model and prompt, without GPT-only fields.
- Quick-create opens the real shared API Key form. After group data loads, Grok Images and Other provider are preselected; quota/custom Key/IP/rate/expiry fields remain available. Cancel sends no create request.
- Both GPT and Grok simulated submissions succeeded and appeared in server-owned fixture history; returned actual PNG bytes load through JWT Blob preview.
- Download button fetched authenticated PNG content; fixture download counter increased. Native in-app browser's `waitForEvent('download')` timed out for the Blob/anchor download, so a saved browser-download file was not confirmed. MIME/filename/Blob delivery is separately covered by frontend tests and backend content tests; do not claim browser filesystem completion.
- Opening an expired historical image shows Image cleaned up, retains prompt/model/time/confirmed cost, and has no Download control. Generation parameters remain in the expandable details.
- Refresh retains the fixture's prior generated history, independent of the selected Key.
- Light/dark styles and 390x844 viewport inspected visually. Controls/history stack vertically without horizontal clipping. Desktop1280x900 verification uses existing layout and two columns.

Initial unrelated dashboard/announcement console errors were fixture-shape errors, not studio product errors. Fixture announcement shape was corrected. HMR during worker edits resets component state; final screenshot is captured after implementation settles.

Screenshot proof will be saved outside Git to avoid adding test artifacts to the product branch.
