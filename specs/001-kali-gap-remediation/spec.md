# Feature Specification: Kali Gap Remediation

**Feature Branch**: `[001-kali-gap-remediation]`

**Created**: 2026-11-16

**Status**: Implemented

**Input**: User description: "Supplement Kali Linux functionality to ARTEX AI v0.9 — map all 16 Kali tool categories against ARTEX AI existing features, identify gaps, and implement missing capabilities (skill.yaml category extension, backend kali_tools API, frontend Kali tools directory page and dashboard, Docker/PostgreSQL verification)."

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

### User Story 1 - Browse Kali tools by category in ARTEX UI (Priority: P1)

Security operators can browse Kali Linux tools organized by the 16 MITRE ATT&CK style categories directly in the ARTEX Web UI skills page, see which tools are already available as ARTEX skills, and identify which are missing.

**Why this priority**: This is the primary user-facing value of the feature — visibility into Kali coverage gaps is the prerequisite for all subsequent remediation work.

**Independent Test**: Can be fully tested by navigating to the skills page, opening the Kali tools directory tab, and verifying all 16 categories render with correct tool counts.

**Acceptance Scenarios**:

1. **Given** the ARTEX UI is running, **When** the user opens the skills page and selects the Kali tools directory tab, **Then** all 16 Kali categories are listed with tool names and availability status.
2. **Given** a Kali tool exists in the directory, **When** the user views its entry, **Then** the entry shows whether an ARTEX skill.yaml exists for it and a link to kali.org/tools/.

---

### User Story 2 - Install a Kali tool as an ARTEX skill (Priority: P2)

Security operators can install a Kali tool from the directory as an ARTEX skill, which generates a skill.yaml and registers the tool's entrypoint in ARTEX.

**Why this priority**: Directly addresses the largest gap category (9 of 16 categories have no ARTEX mapping) by enabling one-click skill creation from Kali tools.

**Independent Test**: Can be fully tested by selecting a Kali tool without an existing skill, clicking install, and verifying a new skill.yaml appears in the skills directory with valid category and entrypoint.

**Acceptance Scenarios**:

1. **Given** a Kali tool entry without an ARTEX skill, **When** the user clicks install, **Then** a POST /api/kali/tools/{tool}/install call is made and a skill.yaml is generated in the skills directory.
2. **Given** a newly installed skill, **When** the user opens the skill detail page, **Then** the skill.yaml displays with correct category, entrypoint type, and parameters.

---

### User Story 3 - View Kali category dashboard with statistics (Priority: P3)

Security managers can view a dashboard summarizing Kali category coverage, skill call statistics, and finding counts per category.

**Why this priority**: Provides management-level visibility into security testing coverage across all Kali domains.

**Independent Test**: Can be fully tested by navigating to the Kali dashboard and verifying category statistics, call counts, and finding counts render correctly.

**Acceptance Scenarios**:

1. **Given** the Kali dashboard page loads, **When** the user views category cards, **Then** each of the 16 categories shows skill count, call count, and finding count.
2. **Given** skills have been called, **When** the user refreshes the dashboard, **Then** the statistics reflect the latest usage data.

---

### Edge Cases

- What happens when a Kali tool name conflicts with an existing ARTEX skill name?
- How does the system handle Kali tools that require system dependencies not available in the ARTEX container?
- What happens when the kali.org/tools/ source is unreachable during tool directory sync?
- How are skill.yaml validation errors handled when auto-generating from a Kali tool?

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: System MUST list all 16 Kali Linux tool categories in the ARTEX Web UI with tool names and availability status.
- **FR-002**: System MUST allow users to install a Kali tool as an ARTEX skill, generating a valid skill.yaml.
- **FR-003**: System MUST extend the skill.yaml Category whitelist to include all 14 new Kali categories (resource-dev, initial-access, persistence, priv-esc, defense-evasion, credential-access, discovery, lateral-movement, collection, c2, exfiltration, impact, forensics, services).
- **FR-004**: System MUST provide API endpoints GET /api/kali/tools, POST /api/kali/tools/{tool}/install, and GET /api/kali/categories.
- **FR-005**: System MUST validate generated skill.yaml using the existing ValidateSkillYAML function before writing.
- **FR-006**: System MUST display a Kali category dashboard with per-category skill count, call count, and finding count.
- **FR-007**: System MUST link each skill entry to its corresponding kali.org/tools/ page when a Kali mapping exists.
- **FR-008**: System MUST handle duplicate skill names gracefully during Kali tool installation.
- **FR-009**: System MUST record skill installation and usage in the existing skill_usage ledger.
- **FR-010**: System MUST work with Docker Compose + PostgreSQL deployment.

### Key Entities *(include if feature involves data)*

- **KaliCategory**: One of 16 MITRE ATT&CK style categories (recon, resource-dev, initial-access, execution, persistence, priv-esc, defense-evasion, credential-access, discovery, lateral-movement, collection, c2, exfiltration, impact, forensics, services).
- **KaliTool**: A tool entry with name, category, description, official URL, and ARTEX skill availability status.
- **SkillYAML**: Existing ARTEX skill definition extended with Kali category and source mapping fields.

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: All 16 Kali categories are visible in the ARTEX Web UI with correct tool listings.
- **SC-002**: Installing any Kali tool as a skill produces a skill.yaml that passes ValidateSkillYAML with zero errors.
- **SC-003**: The three new API endpoints return HTTP 200 with correct JSON payloads for valid requests.
- **SC-004**: The Kali category dashboard renders all 16 category cards with non-negative integer statistics.
- **SC-005**: `docker compose up -d --build` succeeds and PostgreSQL container is running after deployment.

## Assumptions

- Users have access to the ARTEX Web UI with appropriate permissions to view and install skills.
- The existing skill.yaml validation infrastructure (ValidateSkillYAML) is sufficient for Kali-generated skills.
- Kali tool data is sourced from kali.org/tools/ and maintained as a static reference list.
- Docker Compose and PostgreSQL are already configured in the ARTEX deployment (per prior work state).
- No external network access is required for the skill directory to function (static tool list embedded in backend).
