# API Contract: opencode-shannon-plugin

## GET /api/shannon-phases

List all 6 opencode-shannon-plugin pentest phases with tool and skill counts.

**Request:** None

**Response:**
```json
{
  "phases": [
    {
      "name": "recon",
      "tool_count": 150,
      "skill_count": 3
    },
    {
      "name": "vuln-discovery",
      "tool_count": 80,
      "skill_count": 1
    }
  ]
}
```

## GET /api/shannon-tools

List opencode-shannon-plugin tools, optionally filtered by phase.

**Query Parameters:**
- `phase` (optional): Filter by pentest phase name

**Response:**
```json
{
  "tools": [
    {
      "name": "nmap",
      "phase": "recon",
      "description": "Network exploration tool and security / port scanner",
      "official_url": "https://nmap.org/",
      "skill_available": true,
      "skill_yaml_path": "/skills/nmap/skill.yaml"
    }
  ]
}
```

## POST /api/shannon-tools/{tool}/install

Install an opencode-shannon-plugin tool as an ARTEX skill by generating skill.yaml.

**Request Body:** `{}`

**Response:**
```json
{
  "skill_path": "/skills/nmap/skill.yaml",
  "skill_yaml": {
    "name": "nmap",
    "category": "recon",
    ...
  }
}
```

**Errors:**
- `409 Conflict`: Skill already exists
- `404 Not Found`: Tool not found

## GET /api/shannon-stats

Return per-phase runtime statistics.

**Request:** None

**Response:**
```json
{
  "stats": [
    {
      "name": "recon",
      "tool_count": 150,
      "skill_count": 3,
      "call_count": 42,
      "finding_count": 0
    }
  ]
}
```
