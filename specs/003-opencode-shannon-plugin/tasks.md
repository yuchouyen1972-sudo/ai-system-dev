
# Tasks: opencode-shannon-plugin Integration

**Input**: Design documents from `/specs/003-opencode-shannon-plugin/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are OPTIONAL - not explicitly requested in the feature specification. Including basic Go test for shannon_plugin.go as quality gate.

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

- [x] T001 Initialize Spec Kit feature directory `specs/003-opencode-shannon-plugin/` (spec.md, plan.md, tasks.md)
- [x] T002 [P] Create `specs/003-opencode-shannon-plugin/checklists/requirements.md` quality checklist

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 Extend `server/skill_yaml.go` Category whitelist with 6 new opencode-shannon-plugin phases (recon, vuln-discovery, browser-testing, idor-upload-testing, exploitation, reporting)
- [x] T004 [P] Create `server/shannon_plugin.go` with ShannonPhase/ShannonTool/ShannonAgent data structures, static tool list from README, skill.yaml generation, and validation
- [x] T005 [P] Register `GET /api/shannon-phases`, `GET /api/shannon-tools`, `POST /api/shannon-tools/{tool}/install`, `GET /api/shannon-stats` in `server/server.go`
- [x] T006 Run `go test ./server/...` to verify backend changes compile and pass

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Browse opencode-shannon-plugin attack phases and install as ARTEX skills (Priority: P1) - MVP

**Goal**: Security operators can browse opencode-shannon-plugin tools organized by 6 pentest phases in the ARTEX Web UI skills page, see availability status, and install missing tools as ARTEX skills.

**Independent Test**: Navigate to skills page, open opencode-shannon-plugin directory tab, verify all 6 phases render with tool counts and availability.

### Implementation for User Story 1

- [x] T007 [P] [US1] Create `web/src/app/(main)/system/shannon-plugin/page.tsx` with 6-phase tool directory listing
- [x] T008 [US1] Add Shannon plugin directory tab to `web/src/app/(main)/system/skills/page.tsx`
- [x] T009 [US1] Wire frontend to `GET /api/shannon-phases` and `GET /api/shannon-tools` API endpoints
- [x] T010 [US1] Verify skill availability status renders correctly (skill_available boolean)
- [x] T011 [US1] Add install button and handler calling `POST /api/shannon-tools/{tool}/install`
- [x] T012 [US1] Handle duplicate skill name conflict during installation

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Execute opencode-shannon-plugin agents through ARTEX task workflow (Priority: P2)

**Goal**: Security operators can create an ARTEX task that leverages opencode-shannon-plugin's agents to execute a penetration testing engagement, with results written back into ARTEX's dual-graph.

**Independent Test**: Create a task with an authorized target, select Shannon plugin agents, run the task, verify assets and findings appear in asset graph and exploration graph.

### Implementation for User Story 2

- [x] T013 [P] [US2] Map Shannon plugin agent roles to ARTEX Agent model in `server/shannon_plugin.go`
- [ ] T014 [US2] Add agent selection UI to task creation form in `web/src/app/(main)/system/tasks/`
- [ ] T015 [US2] Wire task dispatch to invoke Shannon plugin agent execution
- [ ] T016 [US2] Verify assets/facts/findings written to dual-graph on task completion

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - View opencode-shannon-plugin phase dashboard with statistics (Priority: P3)

**Goal**: Security managers can view a dashboard summarizing opencode-shannon-plugin phase coverage, skill call statistics, and finding counts per pentest phase.

**Independent Test**: Navigate to Shannon plugin dashboard, verify all 6 phase cards show skill count, call count, and finding count.

### Implementation for User Story 3

- [x] T017 [P] [US3] Create `web/src/app/(main)/system/shannon-plugin-dashboard/page.tsx` with 6 phase cards
- [x] T018 [US3] Wire dashboard to `GET /api/shannon-stats` API endpoint
- [x] T019 [US3] Verify dashboard renders non-negative integer statistics for all phases

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories; deployment verification

- [x] T020 [P] Run `go test ./server/...` and `golangci-lint` for backend quality gate
- [x] T021 [P] Run Next.js lint and build for frontend quality gate
- [ ] T022 Verify `docker compose up -d --build` succeeds and PostgreSQL container is running
- [x] T023 [P] Update `docs/OPCODE_SHANNON_PLUGIN.md` with implementation status
- [ ] T024 Run `/speckit.converge` to verify implementation satisfies spec, plan, and tasks

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
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 API endpoints (`GET /api/shannon-tools`)
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
- T013, T014 (US2): T013 and T014 can run in parallel; T015 depends on both
- T017, T018 (US3): Can run in parallel; T019 depends on both
- T020, T021, T022, T023 (Polish): All marked [P] can run in parallel; T024 depends on all

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

**Purpose**: Assess codebase against spec/plan/tasks and append remaining work.

**Convergence Findings:**

| ID | Gap Type | Severity | Source | Evidence | Remaining Work |
|----|----------|----------|--------|----------|----------------|
| F1 | missing | CRITICAL | FR-003 | `server/skill_yaml.go` `validCategories` map does not include 6 new Shannon phases | Extend Category whitelist with recon, vuln-discovery, browser-testing, idor-upload-testing, exploitation, reporting |
| F2 | missing | CRITICAL | FR-004 | No `/api/shannon/*` endpoints registered in `server/server.go` | Register GET /api/shannon-phases, GET /api/shannon-tools, POST /api/shannon-tools/{tool}/install, GET /api/shannon-stats |
| F3 | missing | CRITICAL | FR-001 / FR-002 | No frontend page for Shannon plugin tools directory | Create `web/src/app/(main)/system/shannon-plugin/page.tsx` |
| F4 | missing | HIGH | FR-006 | No Shannon plugin dashboard page | Create `web/src/app/(main)/system/shannon-plugin-dashboard/page.tsx` |
| F5 | missing | HIGH | FR-005 | `InstallShannonPluginToolAsSkill` not implemented | Implement skill.yaml generation + ValidateSkillYAML validation + file write |
| F6 | missing | HIGH | FR-008 | No duplicate skill name handling | Add `os.Stat` check for existing skill directory, return HTTP 409 |
| F7 | missing | HIGH | FR-009 | skill_usage ledger not written on Shannon skill install | Add `pg.InsertSkillUsage` call in install handler |
| F8 | missing | MEDIUM | FR-007 | No skill_yaml_path link in tool info | Populate `skill_yaml_path` field when skill exists |
| F9 | missing | MEDIUM | FR-010 / T14 | `docker compose up -d --build` not executed | Run in Linux/Docker environment, verify PostgreSQL container running |
| F10 | missing | MEDIUM | FR-010 / T022 | No Go unit tests for shannon_plugin.go | Create `server/shannon_plugin_test.go` |

- [ ] T025 [CRITICAL] [FR-003] Extend `server/skill_yaml.go` Category whitelist with 6 new Shannon phases (recon, vuln-discovery, browser-testing, idor-upload-testing, exploitation, reporting)
- [ ] T026 [CRITICAL] [FR-004] Register `GET /api/shannon-phases`, `GET /api/shannon-tools`, `POST /api/shannon-tools/{tool}/install`, `GET /api/shannon-stats` in `server/server.go`
- [ ] T027 [CRITICAL] [FR-001 / FR-002] Create `web/src/app/(main)/system/shannon-plugin/page.tsx` with 6-phase tool directory listing
- [ ] T028 [HIGH] [FR-006] Create `web/src/app/(main)/system/shannon-plugin-dashboard/page.tsx` with 6 phase cards
- [ ] T029 [HIGH] [FR-005] Implement `InstallShannonPluginToolAsSkill` in `server/shannon_plugin.go` with ValidateSkillYAML check
- [ ] T030 [HIGH] [FR-008] Add duplicate skill name conflict handling in `InstallShannonPluginToolAsSkill`
- [ ] T031 [HIGH] [FR-009] Add `pg.InsertSkillUsage` call after successful skill installation
- [ ] T032 [MEDIUM] [FR-007] Populate `skill_yaml_path` field when skill directory exists
- [ ] T033 [MEDIUM] [FR-010 / T14] Run `docker compose up -d --build` in Linux/Docker environment, verify PostgreSQL container running
- [ ] T034 [MEDIUM] [FR-010 / T022] Create Go unit tests `server/shannon_plugin_test.go`

