![VaporTrace Logo](../../assets/images/VaporTrace_Logo.png)

# 27 - Neuro F6 + MCP Upgrade Plan

## Purpose
Define a concrete implementation plan to make F6 (Neuro tab) operationally useful, integrate an MCP server safely, and fix tactical plan persistence expectations around `list-plan`.

---

## Current State (Why F6 Feels Weak)
- F6 mostly shows AI/log output, but not actionable decision context.
- There is no full traceability from AI suggestion to executed action and outcome.
- No explicit tool-routing visibility (what data was used, from where, and why).
- `list-plan` uses in-memory `ActionBuffer`, so plan state is session-local.

---

## Why `list-plan` Is Not Persistent Today
Root cause:
- `list-plan` reads `engine.ActionBuffer` (memory), not a persisted table.
- `ActionBuffer` is rebuilt by `analyze`, modified by `edit/drop/commit`, and cleared/reset in runtime flows.
- Restarting process or `reset_db` destroys runtime state and therefore plan visibility.

Technical references:
- `pkg/engine/core.go` (`ActionBuffer` global + `list-plan` command)
- `pkg/ui/dashboard.go` (F5 renders directly from `ActionBuffer`)

---

## Target Outcome
F6 becomes a "Decision Cockpit":
- Shows recommendation, evidence, confidence, and reasoning chain.
- Shows model/provider/tool usage and latency/cost.
- Links each recommendation to execution result and finding IDs.
- Allows operator actions: approve, reject, mutate, rerun.

And plan management becomes durable:
- `list-plan` survives process restarts.
- Each action has lifecycle history (`PENDING`, `RUNNING`, `SUCCESS`, `FAILED`, `DROPPED`, timestamps).

---

## F6 Upgrade Requirements

### 1) Data Model
Add persistent entities (SQLite):
- `neuro_sessions` (session id, target, provider, model, started_at, ended_at)
- `neuro_events` (session_id, phase, prompt_hash, response_summary, token_usage, latency_ms, created_at)
- `tactical_actions` (action_id, session_id, type, target, payload, confidence, reasoning, status, created_at, updated_at)
- `tactical_action_events` (action_id, event_type, details, finding_ref, created_at)

### 2) UX in F6
Mandatory panes:
- "Decision Feed": latest recommendations sorted by confidence and freshness.
- "Evidence Panel": matched endpoints/findings/traffic excerpts used for decision.
- "Execution Link": mapping from recommendation -> committed action -> outcome.
- "Provider Health": active provider, fallback status, recent errors, average latency.

### 3) Explainability
For each recommendation:
- Why this action (reasoning summary).
- Inputs used (discovery, findings, loot, traffic, external intel).
- Confidence decomposition (heuristic + model contribution).

### 4) Operator Controls
- Approve selected AI action to F5.
- Reject and store rationale.
- "Regenerate with constraints" (target, risk budget, allowed modules).

---

## MCP Server Integration Plan

## Do You Need MCP First?
No. Build internal persistence + explainability first.  
Then add MCP to extend external tools cleanly.

### MCP Phase Scope
Phase A (read-only tools):
- CVE lookup
- Service fingerprint enrichment
- Passive intel adapters

Phase B (controlled action tools):
- Run bounded scanners with strict scope guard
- Enforce approval before any active exploit action

### MCP Safety Controls
- Hard scope policy gate before tool execution.
- Allowlist of tools per mode (lab/safe/aggressive).
- Audit log per MCP call: arguments, result hash, latency, status.
- Timeout and retry budget per tool class.

### MCP Contract
Each tool call must emit:
- `tool_name`
- `input_hash`
- `output_summary`
- `confidence_impact`
- `linked_action_id`

---

## `list-plan` Persistence Upgrade (Required)

### Minimal Path (Fast)
Persist only current active plan:
- Save `ActionBuffer` into `tactical_actions` after `analyze`.
- Update row status on `edit/drop/commit`.
- Load latest active plan at startup and before `list-plan`.

Pros:
- Quick implementation.
- Immediate persistence benefit.

Cons:
- Limited historical analytics.

### Full Path (Recommended)
Event-sourced action lifecycle:
- `tactical_actions` holds current state.
- `tactical_action_events` tracks every change.
- `list-plan` queries current state.
- New command `plan-history` shows lifecycle timeline.

Pros:
- Full auditability.
- Better debugging and explainability for F6.

Cons:
- More schema and code work.

---

## Implementation Phases

### Phase 1 - Persistence Foundation
- Add new DB tables.
- Repository layer for `tactical_actions`.
- `analyze` writes persisted actions.
- `list-plan` reads persisted state.

### Phase 2 - Commit Telemetry Linkage
- Store execution start/end timestamps and outcome.
- Link findings generated during execution to `action_id`.
- Keep F5/F6 views synchronized from DB-backed state.

### Phase 3 - F6 Decision Cockpit
- Render decision feed + evidence + provider health.
- Add approve/reject/regenerate controls.
- Show reasoning + confidence decomposition.

### Phase 4 - MCP Integration
- Add MCP client wrapper and policy gate.
- Register read-only tools first.
- Add audit trail and timeout budgets.

### Phase 5 - Hardening
- Concurrency safety tests.
- Recovery tests after crash/restart.
- Scope bypass regression tests.

---

## Acceptance Criteria
- `list-plan` survives restart and shows same pending/running/completed state.
- Every committed action has outcome and evidence linkage.
- F6 displays why/how each recommendation was made.
- MCP calls are policy-gated, auditable, and tied to action IDs.

---

## Immediate Next Steps
1. Implement `tactical_actions` + `tactical_action_events` schema.
2. Replace in-memory-only read path in `list-plan` with DB-first load.
3. Add F6 "Decision Feed" and "Execution Link" panels.
4. Add MCP read-only adapters behind scope gate.
