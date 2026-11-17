# Feature Specification: opencode-shannon-plugin Integration

**Feature Branch**: `[003-opencode-shannon-plugin]`

**Created**: 2026-11-17

**Status**: Implemented

**Input**: User description: "Integrate opencode-shannon-plugin (AI-powered autonomous penetration testing plugin for OpenCode with 600+ Kali Linux tools, Playwright browser automation, and OWASP methodology) into ARTEX AI as skill.yaml-based capabilities, mapping its 6 pentest phases, 600+ tools, and Shannon agent architecture onto ARTEX's skill system, agent system, and finding/evidence model."

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.

  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - Browse opencode-shannon-plugin attack phases and tools in ARTEX UI (Priority: P1)

Security operators can browse opencode-shannon-plugin's 600+ Kali Linux tools organized by 6 pentest phases directly in the ARTEX Web UI skills page, see which tools are already available as ARTEX skills, and install missing tools as ARTEX skills with one click.

**Why this priority**: This is the primary user-facing value — visibility into opencode-shannon-plugin's tool coverage and one-click skill installation is the prerequisite for all subsequent remediation work.

**Independent Test**: Can be fully tested by navigating to the skills page, opening the opencode-shannon-plugin directory tab, and verifying all 6 pentest phases render with tool counts and availability status.

**Acceptance Scenarios**:

1. **Given** the ARTEX UI is running, **When** the user opens the skills page and selects the opencode-shannon-plugin directory tab, **Then** all 6 pentest phases are listed with tool names and availability status.
2. **Given** an opencode-shannon-plugin tool exists in the directory without an ARTEX skill, **When** the user clicks install, **Then** a skill.yaml is generated in the skills directory with valid category and entrypoint.

---

### User Story 2 - Execute opencode-shannon-plugin agents through ARTEX task workflow (Priority: P2)

Security operators can create an ARTEX task that leverages opencode-shannon-plugin's agents to execute a penetration testing engagement, with results written back into ARTEX's dual-graph (asset graph + exploration graph) as assets, intents, facts, and findings.

**Why this priority**: Directly enables end-to-end penetration testing using opencode-shannon-plugin's agent orchestration within ARTEX's task and graph model.

**Independent Test**: Can be fully tested by creating a task with an authorized target, selecting opencode-shannon-plugin agents, running the task, and verifying assets and findings appear in the asset graph and exploration graph.

**Acceptance Scenarios**:

1. **Given** an authorized target is in scope, **When** the user creates a task with opencode-shannon-plugin agents, **Then** the planner dispatches intents and workers execute them, writing assets/facts/findings back to both graphs.
2. **Given** a task completes with findings, **When** the user views the task detail, **Then** each finding carries evidence (PoC, execution output) bound to the assets it was tested against.

---

### User Story 3 - View opencode-shannon-plugin phase dashboard with statistics (Priority: P3)

Security managers can view a dashboard summarizing opencode-shannon-plugin phase coverage, skill call statistics, and finding counts per pentest phase.

**Why this priority**: Provides management-level visibility into security testing coverage across all opencode-shannon-plugin pentest phases.

**Independent Test**: Can be fully tested by navigating to the opencode-shannon-plugin dashboard and verifying all 6 phase cards show skill count, call count, and finding count.

**Acceptance Scenarios**:

1. **Given** the opencode-shannon-plugin dashboard page loads, **When** the user views phase cards, **Then** each of the 6 phases shows skill count, call count, and finding count.
2. **Given** skills have been called, **When** the user refreshes the dashboard, **Then** the statistics reflect the latest usage data.

---

### Edge Cases

- What happens when an opencode-shannon-plugin tool name conflicts with an existing ARTEX skill name?
- How does the system handle opencode-shannon-plugin tools that require Docker container runtime not available in the ARTEX container?
- What happens when opencode-shannon-plugin's tool list changes between versions?
- How are skill.yaml validation errors handled when auto-generating from opencode-shannon-plugin tool definitions?

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: System MUST list all 6 opencode-shannon-plugin pentest phases in the ARTEX Web UI with tool names and availability status.
- **FR-002**: System MUST allow users to install an opencode-shannon-plugin tool as an ARTEX skill, generating a valid skill.yaml.
- **FR-003**: System MUST extend the skill.yaml Category whitelist to include 6 new opencode-shannon-plugin pentest phases (recon, vuln-discovery, browser-testing, idor-upload-testing, exploitation, reporting).
- **FR-004**: System MUST provide API endpoints GET /api/shannon-phases, GET /api/shannon-tools, POST /api/shannon-tools/{tool}/install, GET /api/shannon-stats.
- **FR-005**: System MUST validate generated skill.yaml using the existing ValidateSkillYAML function before writing.
- **FR-006**: System MUST display an opencode-shannon-plugin dashboard with per-phase skill count, call count, and finding count.
- **FR-007**: System MUST link each skill entry to its corresponding opencode-shannon-plugin tool documentation when a mapping exists.
- **FR-008**: System MUST handle duplicate skill names gracefully during opencode-shannon-plugin tool installation.
- **FR-009**: System MUST record skill installation and usage in the existing skill_usage ledger.
- **FR-010**: System MUST work with Docker Compose + PostgreSQL deployment.

### Key Entities *(include if feature involves data)*

- **ShannonPhase**: One of 6 opencode-shannon-plugin pentest phases (recon, vuln-discovery, browser-testing, idor-upload-testing, exploitation, reporting).
- **ShannonTool**: A tool entry with name, phase, description, official URL, and ARTEX skill availability status.
- **SkillYAML**: Existing ARTEX skill definition extended with phase and source mapping fields.

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: All 6 opencode-shannon-plugin pentest phases are visible in the ARTEX Web UI with correct tool listings.
- **SC-002**: Installing any opencode-shannon-plugin tool as a skill produces a skill.yaml that passes ValidateSkillYAML with zero errors.
- **SC-003**: The four new API endpoints return HTTP 200 with correct JSON payloads for valid requests.
- **SC-004**: The opencode-shannon-plugin dashboard renders all 6 phase cards with non-negative integer statistics.
- **SC-005**: `docker compose up -d --build` succeeds and PostgreSQL container is running after deployment.

## Assumptions

- Users have access to the ARTEX Web UI with appropriate permissions to view and install skills.
- The existing skill.yaml validation infrastructure (ValidateSkillYAML) is sufficient for opencode-shannon-plugin-generated skills.
- opencode-shannon-plugin tool data is sourced from its README and bundled security tools list and maintained as a static reference list.
- Docker Compose and PostgreSQL are already configured in the ARTEX deployment (per prior work state).
- No external network access is required for the tool directory to function (static tool list embedded in backend).
