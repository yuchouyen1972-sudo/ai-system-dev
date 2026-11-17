# Implementation Plan: opencode-shannon-plugin Integration

**Branch**: `[003-opencode-shannon-plugin]` | **Date**: 2026-11-17 | **Spec**: `specs/003-opencode-shannon-plugin/spec.md`

**Input**: Feature specification from `specs/003-opencode-shannon-plugin/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Integrate opencode-shannon-plugin (AI-powered autonomous penetration testing plugin for OpenCode with 600+ Kali Linux tools, Playwright browser automation, and OWASP methodology) into ARTEX AI v0.9 by mapping its 6 pentest phases, 600+ tools, and Shannon agent architecture onto ARTEX's existing skill.yaml system, agent system, and finding/evidence model. Implement backend `server/shannon_plugin.go` with four API endpoints (`GET /api/shannon-phases`, `GET /api/shannon-tools`, `POST /api/shannon-tools/{tool}/install`, `GET /api/shannon-stats`), frontend opencode-shannon-plugin tools directory page and category dashboard, and Docker Compose + PostgreSQL deployment verification.

## Technical Context

**Language/Version**: Go 1.27.4 (backend), Node.js v24.20.0 (frontend)

**Primary Dependencies**: Go standard library + `gopkg.in/yaml.v3` (backend); Next.js, React, shadcn/ui, lucide-react, sonner, next-intl (frontend)

**Storage**: SQLite (dual graph stores, existing ARTEX backend); PostgreSQL (Docker Compose deployment, existing)

**Testing**: `go test` (backend), Next.js built-in testing (frontend)

**Target Platform**: Linux server (Docker container)

**Project Type**: web-service + CLI (ARTEX backend is a Go HTTP server; frontend is a Next.js web app)

**Performance Goals**: N/A (internal security tooling; no specific throughput requirements)

**Constraints**:
- Must not use web UI / browser as output method (AGENTS.md 第零級 — terminal/CLI and pure files only)
- Docker Compose with PostgreSQL must work (`docker compose up -d --build`)
- skill.yaml validation must pass (existing `ValidateSkillYAML`)
- No external network access required for tool directory (static tool list embedded in backend)
- opencode-shannon-plugin requires Docker container for tool execution (shannon-tools image)

**Scale/Scope**: Single Go project (`artex-global` repository); 6 opencode-shannon-plugin pentest phases; 600+ tools reference; 4 new API endpoints; 2 new frontend pages

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution file `.specify/memory/constitution.md` is an unfilled template — skipped gracefully per speckit.plan rules. No governance constraints to evaluate.

## Project Structure

### Documentation (this feature)

```text
specs/003-opencode-shannon-plugin/
├── spec.md              # Feature specification (/speckit-specify output)
├── plan.md              # This file (/speckit-plan output)
├── research.md          # Phase 0 output (/speckit-plan)
├── data-model.md        # Phase 1 output (/speckit-plan)
├── quickstart.md        # Phase 1 output (/speckit-plan)
├── contracts/           # Phase 1 output (/speckit-plan)
│   └── shannon-plugin-api.md  # API contract for /api/shannon/*
├── checklists/
│   └── requirements.md  # Spec quality checklist
└── tasks.md             # Phase 2 output (/speckit-tasks - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
backend/artex-global/
├── server/
│   ├── skill_yaml.go          # Extended: Category whitelist + validation
│   ├── kali_tools.go          # Existing: Kali tool directory (from prior feature)
│   ├── opencode_pentester.go  # Existing: opencode-pentester tool directory (from prior feature)
│   ├── shannon_plugin.go      # NEW: Shannon plugin tool directory, install, skill.yaml generation
│   └── server.go              # Extended: Register /api/shannon/* endpoints
├── web/
│   └── src/
│       └── app/
│           └── (main)/
│               └── system/
│                   ├── skills/
│                   │   └── page.tsx         # Extended: Shannon plugin directory tab
│                   ├── shannon-plugin/
│                   │   └── page.tsx         # NEW: Shannon plugin tools directory page
│                   └── shannon-plugin-dashboard/
│                       └── page.tsx         # NEW: Shannon plugin category dashboard
├── docker-compose.local.yml   # Existing: PostgreSQL + ARTEX service verification
├── docs/
│   ├── KALI_GAP_ANALYSIS.md   # Existing: Gap analysis report
├── docs/
│   ├── KALI_REMEDIATION_PLAN.md  # Existing: Remediation plan
├── docs/
│   ├── KALI_ARTEX_SEC_TESTING_GAP.md  # Existing: Sec testing gap report
├── docs/
│   └── OPCODE_PENTESTER_INTEGRATION.md  # Existing: opencode-pentester integration plan
└── docs/
    └── OPCODE_SHANNON_PLUGIN.md  # NEW: opencode-shannon-plugin integration plan
```

**Structure Decision**: Single Go project (Option 1 default) — ARTEX is a unified Go backend + Next.js frontend in one repository. No separate backend/frontend split needed. Documentation lives under `specs/003-opencode-shannon-plugin/`. Source changes are confined to `server/` (Go) and `web/src/app/(main)/system/` (Next.js).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No constitution violations (constitution is unfilled template, skipped). Complexity tracking not required.

## Phase 0: Research & Decisions (research.md)

All NEEDS CLARIFICATION items from spec.md are resolved with informed guesses:

| Decision | Chosen | Rationale | Alternatives Considered |
|----------|--------|-----------|------------------------|
| Shannon tool data source | Static embedded list from README | No external network dependency; offline-capable | Sync from GitHub (rejected: network dependency) |
| skill.yaml Category extension | Add 6 new domains to existing whitelist | Minimal change; reuses existing validation | New separate skill type (rejected: over-engineering) |
| API endpoint prefix | `/api/shannon/*` | Consistent with existing `/api/kali/*` pattern | `/shannon/*` (rejected: inconsistent) |
| Skill install mechanism | POST /api/shannon-tools/{tool}/install generates skill.yaml | One-click install; reuses ValidateSkillYAML | Manual skill.yaml creation (rejected: poor UX) |
| Frontend page location | `web/src/app/(main)/system/shannon-plugin/` | Follows existing kali page routing pattern | Top-level `/shannon-plugin` route (rejected: inconsistent grouping) |
| Tool execution | Docker container (shannon-tools image) | Matches opencode-shannon-plugin design | Native binary execution (rejected: requires Bun + Node runtime) |

## Phase 1: Design & Contracts

### Data Model (data-model.md)

No new persistent data model required. The feature reuses existing:
- `SkillYAML` struct (server/skill_yaml.go) — extended with Category whitelist only
- `SkillItem` type (web/src/lib/types.ts) — already has `category`, `has_yaml`, `parameters`, `entrypoint`, `outputs`, etc.
- `skill_usage` ledger — existing, no changes required
- `Agent` type (web/src/lib/types.ts) — already has `role`, `builtin`, `enabled`, etc.

New in-memory structures (not persisted):
- `ShannonPhase` — name, toolCount, skillCount, callCount, findingCount
- `ShannonTool` — name, phase, description, officialUrl, skillAvailable
- `ShannonAgent` — name, role, file, supportedCategories

### Interface Contracts (contracts/shannon-plugin-api.md)

**GET /api/shannon-phases**
- Response: `{ phases: [{ name: string, tool_count: number, skill_count: number }] }`
- Returns all 6 opencode-shannon-plugin phases with counts.

**GET /api/shannon-tools**
- Query params: `phase` (optional filter)
- Response: `{ tools: [{ name: string, phase: string, description: string, official_url: string, skill_available: boolean, skill_yaml_path: string|null }] }`

**POST /api/shannon-tools/{tool}/install**
- Request body: `{}`
- Response: `{ skill_path: string, skill_yaml: {...} }`
- Generates skill.yaml from opencode-shannon-plugin tool definition, validates with ValidateSkillYAML, writes to skills directory.

**GET /api/shannon-stats**
- Response: `{ stats: [{ name: string, tool_count: number, skill_count: number, call_count: number, finding_count: number }] }`
- Returns per-phase runtime stats.

### Quickstart (quickstart.md)

1. Build and start ARTEX: `docker compose up -d --build`
2. Verify PostgreSQL is running: `docker compose ps`
3. Open ARTEX Web UI at `http://localhost:8787`
4. Navigate to Skills → Shannon Plugin Tools Directory
5. Browse 6 phases; click Install on any tool without a skill
6. Verify generated skill.yaml appears in the skills directory
7. Navigate to Shannon Plugin Dashboard to view category statistics
