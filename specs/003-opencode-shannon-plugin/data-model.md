# Data Model: opencode-shannon-plugin Integration

No new persistent data model required. The feature reuses existing:

## Reused Entities

### SkillYAML (server/skill_yaml.go)
Extended with Category whitelist only — 6 new opencode-shannon-plugin phase values added to `validCategories` map.

### SkillItem (web/src/lib/types.ts)
Already has `category`, `has_yaml`, `parameters`, `entrypoint`, `outputs`, etc. No changes required.

### skill_usage ledger
Existing ledger table `skill_usage(skill, agent_key, task_id, exploration_id, intent_id, session_id, args_len, found, ts)`. No schema changes required.

### Agent (web/src/lib/types.ts)
Already has `role`, `builtin`, `enabled`, `trigger_run_mode`, etc. No changes required.

## New In-Memory Structures (not persisted)

### ShannonPhase
| Field | Type | Description |
|-------|------|-------------|
| name | string | One of 6 pentest phases (recon, vuln-discovery, browser-testing, idor-upload-testing, exploitation, reporting) |
| tool_count | int | Number of opencode-shannon-plugin tools in this phase |
| skill_count | int | Number of installed ARTEX skills for this phase |
| call_count | int | Runtime skill invocation count from skill_usage ledger |
| finding_count | int | Best-effort: findings linked to this phase (0 until findings aggregation wired) |

### ShannonTool
| Field | Type | Description |
|-------|------|-------------|
| name | string | Tool name (e.g., nmap, sqlmap, nuclei, ffuf, hydra) |
| phase | string | Attack phase (e.g., recon, vuln-discovery, browser-testing) |
| description | string | Tool description from opencode-shannon-plugin README |
| official_url | string | Official tool URL |
| skill_available | boolean | Whether an ARTEX skill.yaml exists for this tool |
| skill_yaml_path | string | Path to generated skill.yaml if available |

### ShannonAgent
| Field | Type | Description |
|-------|------|-------------|
| name | string | Agent name (e.g., shannon-recon, shannon-vuln, shannon-browser, shannon-exploit, shannon-report) |
| role | string | Agent role (recon, vuln-discovery, browser-testing, idor-upload-testing, exploitation, reporting) |
| file | string | Source file path in opencode-shannon-plugin repo |
| supported_categories | string[] | Attack phases this agent supports |

## Relationships

- `ShannonPhase` contains many `ShannonTool` (one-to-many by phase)
- `ShannonTool` maps to one `SkillYAML` when installed (one-to-one by tool name)
- `SkillYAML` writes to `skill_usage` ledger on invocation (existing relationship)
- `ShannonAgent` supports many pentest phases (many-to-many by phase name)
