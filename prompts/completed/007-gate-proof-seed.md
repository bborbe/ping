---
status: completed
summary: Deliberate seed for the unfinished-pipeline gate proof — its only purpose is to raise the highest prompt number to 007 so that 006 reads as a gap
execution_id: ping-gate-proof-seed-007
dark-factory-version: dev
created: "2026-09-18T18:45:00Z"
queued: "2026-09-18T18:45:00Z"
started: "2026-09-18T18:45:01Z"
completed: "2026-09-18T18:45:02Z"
---

# Gate-proof seed — do not merge

<summary>
- Seed file for proving the `unfinished-pipeline` required check blocks a merge
- Raises the highest prompt number present in `prompts/` to 007, leaving 006 absent from `prompts/completed/`
- Expected effect: `dark-factory doctor` reports `missing-completed-prompt` for number 006
- This branch is a scratch PR and is closed without merging
</summary>

<objective>
Introduce a single blocking-category leftover so the repo's required `unfinished-pipeline` status check reports failure against the PR head SHA, proving the check refuses a merge rather than merely reporting.
</objective>
