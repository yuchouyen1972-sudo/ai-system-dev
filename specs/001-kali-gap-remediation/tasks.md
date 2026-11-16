
# Tasks: Kali Gap Remediation

**Input**: Design documents from `/specs/001-kali-gap-remediation/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are OPTIONAL - not explicitly requested in the feature specification. Including basic Go test for kali_tools.go as quality gate.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Single project: `server/`, `web/src/` at repository root (`artex-global/`)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Initialize Spec Kit feature directory `specs/001-kali-gap-remediation/` (spec.md, plan.md, tasks.md)
- [ ] T002 [P] Create `specs/001-kali-gap-remediation/checklists/requirements.md` quality checklist

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**CRITICAL**: No user story work can begin until this phase is complete

- [ ] T003 Extend `server/skill_yaml.go` Category whitelist with 14 new Kali categories (resource-dev, initial-access, persistence, priv-esc, defense-evasion, credential-access, discovery, lateral-movement, collection, c2, exfiltration, impact, forensics, services)
- [ ] T004 [P] Create `server/kali_tools.go` with KaliTool data structure, static tool list, skill.yaml generation, and validation
- [ ] T005 [P] Register `GET /api/kali/categories`, `GET /api/kali/tools`, `POST /api/kali/tools/{tool}/install` in `server/routes.go`
- [ ] T006 Run `go test ./server/...` to verify backend changes compile and pass

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Browse Kali tools by category in ARTEX UI (Priority: P1) - MVP

**Goal**: Security operators can browse Kali Linux tools organized by 16 categories in the ARTEX Web UI skills page, see availability status, and identify gaps.

**Independent Test**: Navigate to skills page, open Kali tools directory tab, verify all 16 categories render with tool counts and availability.

### Implementation for User Story 1

- [ ] T007 [P] [US1] Create `web/src/app/(main)/system/kali/page.tsx` with 16-category tool directory listing
- [ ] T008 [US1] Add Kali tools directory tab to `web/src/app/(main)/system/skills/page.tsx`
- [ ] T009 [US1] Wire frontend to `GET /api/kali/categories` and `GET /api/kali/tools` API endpoints
- [ ] T010 [US1] Verify skill availability status renders correctly (skill_available boolean)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Install a Kali tool as an ARTEX skill (Priority: P2)

**Goal**: Security operators can install a Kali tool from the directory as an ARTEX skill, generating a skill.yaml and registering the tool's entrypoint.

**Independent Test**: Select a Kali tool without an existing skill, click install, verify new skill.yaml appears with valid category and entrypoint.

### Implementation for User Story 2

- [ ] T011 [P] [US2] Add install button and handler in `web/src/app/(main)/system/kali/page.tsx` calling `POST /api/kali/tools/{tool}/install`
- [ ] T012 [US2] Handle duplicate skill name conflict during installation (FR-008)
- [ ] T013 [US2] Verify generated skill.yaml passes ValidateSkillYAML with zero errors

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - View Kali category dashboard with statistics (Priority: P3)

**Goal**: Security managers can view a dashboard summarizing Kali category coverage, skill call statistics, and finding counts per category.

**Independent Test**: Navigate to Kali dashboard, verify all 16 category cards show skill count, call count, and finding count.

### Implementation for User Story 3

- [ ] T014 [P] [US3] Create `web/src/app/(main)/system/kali-dashboard/page.tsx` with 16 category cards
- [ ] T015 [US3] Wire dashboard to backend statistics (skill count, call count, finding count per category)
- [ ] T016 [US3] Verify dashboard renders non-negative integer statistics for all categories

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories; deployment verification

- [ ] T017 [P] Run `go test ./server/...` and `golangci-lint` for backend quality gate
- [ ] T018 [P] Run Next.js lint and build for frontend quality gate
- [ ] T019 Verify `docker compose up -d --build` succeeds and PostgreSQL container is running
- [ ] T020 [P] Update `docs/KALI_REMEDIATION_PLAN.md` with implementation status
- [ ] T021 Run `/speckit.converge` to verify implementation satisfies spec, plan, and tasks

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User stories can proceed in parallel (if staffed) or sequentially in priority order (P1 -> P2 -> P3)
  - Each story is independently testable
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 API endpoints (`GET /api/kali/tools`)
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Depends on US1 API endpoints

### Within Each User Story

- Models/data structures before services
- Services before endpoints/UI
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- T001 and T002 (Setup) can run in parallel
- T003, T004, T005 (Foundational): T004 and T005 can run in parallel (different files); T003 must complete before T005 (routes depend on skill_yaml types)
- T007, T008, T009 (US1): T007 and T008 can run in parallel (different files); T009 depends on both
- T011, T012 (US2): T011 and T012 can run in parallel; T013 depends on T011
- T014, T015 (US3): Can run in parallel; T016 depends on both
- T017, T018, T019, T020 (Polish): All marked [P] can run in parallel; T021 depends on all

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational -> Foundation ready
2. Add User Story 1 -> Test independently -> Deploy/Demo (MVP!)
3. Add User Story 2 -> Test independently -> Deploy/Demo
4. Add User Story 3 -> Test independently -> Deploy/Demo
5. Each story adds value without breaking previous stories

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence

---

## Phase 7: Convergence

**Purpose**: Close gaps between spec/plan/tasks and current implementation, per `/speckit.converge`.
**Source**: Convergence assessment of specs/001-kali-gap-remediation/ against server/ and web/ code.

- [ ] T022 [HIGH] [FR-009] Record skill installation in skill_usage ledger after InstallKaliToolAsSkill in server/kali_tools.go (partial)
- [ ] T023 [HIGH] [FR-010/T019] Run `docker compose up -d --build` and verify PostgreSQL container is running on Linux/Docker environment (partial — Docker daemon unavailable on current Windows host)
- [ ] T024 [LOW] [T006/T017] Add Go unit test `server/kali_tools_test.go` covering GetKaliCategories/GetKaliTools/InstallKaliToolAsSkill (partial)
- [ ] T025 [LOW] [T015] Wire Kali dashboard to backend statistics API for call_count and finding_count per category in web/src/app/(main)/system/kali-dashboard/page.tsx (partial)
- [ ] T026 [LOW] [plan.md] Update plan.md Project Structure route file reference from `routes.go` to `server.go` (minor)
