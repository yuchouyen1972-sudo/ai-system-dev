# Implementation Plan: Kali Gap Remediation

**Branch**: `[001-kali-gap-remediation]` | **Date**: 2026-11-16 | **Spec**: `specs/001-kali-gap-remediation/spec.md`

**Input**: Feature specification from `specs/001-kali-gap-remediation/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

Supplement Kali Linux functionality to ARTEX AI v0.9 by mapping all 16 Kali tool categories against ARTEX AI existing features, identifying 9 completely missing categories and 6 partially covered categories, and implementing missing capabilities: skill.yaml Category whitelist extension (14 new categories), backend `server/kali_tools.go` with three API endpoints (`GET /api/kali/tools`, `POST /api/kali/tools/{tool}/install`, `GET /api/kali/categories`), frontend Kali tools directory page and category dashboard, and Docker Compose + PostgreSQL deployment verification.

## Technical Context

**Language/Version**: Go 1.21+ (backend), Node.js v24.20.0 (frontend)

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
- No external network access required for skill directory (static tool list embedded in backend)

**Scale/Scope**: Single Go project (`artex-global` repository); 16 Kali categories; ~300+ tools reference; 3 new API endpoints; 2 new frontend pages

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution file `.specify/memory/constitution.md` is an unfilled template — skipped gracefully per speckit.plan rules. No governance constraints to evaluate.

## Project Structure

### Documentation (this feature)

```text
specs/001-kali-gap-remediation/
├── spec.md              # Feature specification (/speckit-specify output)
├── plan.md              # This file (/speckit-plan output)
├── research.md          # Phase 0 output (/speckit-plan)
├── data-model.md        # Phase 1 output (/speckit-plan)
├── quickstart.md        # Phase 1 output (/speckit-plan)
├── contracts/           # Phase 1 output (/speckit-plan)
│   └── kali-api.md      # API contract for /api/kali/*
├── checklists/
│   └── requirements.md  # Spec quality checklist
└── tasks.md             # Phase 2 output (/speckit-tasks - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
backend/artex-global/
├── server/
│   ├── skill_yaml.go          # Extended: Category whitelist + validation
│   ├── kali_tools.go          # NEW: Kali tool directory, install, skill.yaml generation
│   └── server.go              # Extended: Register /api/kali/* endpoints in Server.Handler()
├── web/
│   └── src/
│       └── app/
│           └── (main)/
│               └── system/
│                   ├── skills/
│                   │   └── page.tsx         # Extended: Kali tools directory tab
│                   ├── kali/
│                   │   └── page.tsx         # NEW: Kali tools directory page
│                   └── kali-dashboard/
│                       └── page.tsx         # NEW: Kali category dashboard
├── docker-compose.local.yml   # Extended: PostgreSQL + ARTEX service verification
├── docs/
│   └── KALI_GAP_ANALYSIS.md   # Existing: Gap analysis report
└── docs/
    └── KALI_REMEDIATION_PLAN.md  # Existing: Remediation plan
```

**Structure Decision**: Single Go project (Option 1 default) — ARTEX is a unified Go backend + Next.js frontend in one repository. No separate backend/frontend split needed. Documentation lives under `specs/001-kali-gap-remediation/`. Source changes are confined to `server/` (Go) and `web/src/app/(main)/system/` (Next.js).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No constitution violations (constitution is unfilled template, skipped). Complexity tracking not required.

## Phase 0: Research & Decisions (research.md)

All NEEDS CLARIFICATION items from spec.md are resolved with informed guesses:

| Decision | Chosen | Rationale | Alternatives Considered |
|----------|--------|-----------|------------------------|
| Kali tool data source | Static embedded list | No external network dependency; offline-capable | Sync from kali.org/tools/ (rejected: network dependency) |
| skill.yaml Category extension | Add 14 new categories to existing whitelist | Minimal change; reuses existing validation | New separate skill type (rejected: over-engineering) |
| API endpoint prefix | `/api/kali/*` | Consistent with existing `/api/*` pattern | `/kali/*` (rejected: inconsistent) |
| Skill install mechanism | POST /api/kali/tools/{tool}/install generates skill.yaml | One-click install; reuses ValidateSkillYAML | Manual skill.yaml creation (rejected: poor UX) |
| Frontend page location | `web/src/app/(main)/system/kali/` | Follows existing skills page routing pattern | Top-level `/kali` route (rejected: inconsistent grouping) |

## Phase 1: Design & Contracts

### Data Model (data-model.md)

No new persistent data model required. The feature reuses existing:
- `SkillYAML` struct (server/skill_yaml.go) — extended with Category whitelist only
- `SkillItem` type (web/src/lib/types.ts) — already has `category`, `has_yaml`, `parameters`, `entrypoint`, etc.
- `skill_usage` ledger — existing, no changes needed

New in-memory structures (not persisted):
- `KaliTool` — name, category, description, officialUrl, skillAvailable (bool)
- `KaliCategory` — name, toolCount, skillCount, callCount, findingCount

### Interface Contracts (contracts/kali-api.md)

**GET /api/kali/categories**
- Response: `{ categories: [{ name: string, tool_count: number, skill_count: number }] }`
- Returns all 16 Kali categories with counts.

**GET /api/kali/tools**
- Query params: `category` (optional filter)
- Response: `{ tools: [{ name: string, category: string, description: string, official_url: string, skill_available: boolean, skill_yaml_path: string|null }] }`

**POST /api/kali/tools/{tool}/install**
- Request body: `{}`
- Response: `{ skill_path: string, skill_yaml: {...} }`
- Generates skill.yaml from Kali tool definition, validates with ValidateSkillYAML, writes to skills directory.

### Quickstart (quickstart.md)

1. Build and start ARTEX: `docker compose up -d --build`
2. Verify PostgreSQL is running: `docker compose ps`
3. Open ARTEX Web UI at `http://localhost:8787`
4. Navigate to Skills → Kali Tools Directory
5. Browse 16 categories; click Install on any tool without a skill
6. Verify generated skill.yaml appears in the skills directory
7. Navigate to Kali Dashboard to view category statistics
