package server

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// SkillYAML represents the skill.yaml schema for ARTEX Web UI
type SkillYAML struct {
	// Required fields
	Name        string   `yaml:"name" json:"name"`
	DisplayName I18nMap  `yaml:"display_name" json:"display_name"`
	Description I18nMap  `yaml:"description" json:"description"`
	Version     string   `yaml:"version" json:"version"`
	Category    string   `yaml:"category" json:"category"`
	Author      string   `yaml:"author" json:"author"`
	License     string   `yaml:"license" json:"license"`

	// Optional fields
	Tags          []string          `yaml:"tags,omitempty" json:"tags,omitempty"`
	I18n          *I18nDetail       `yaml:"i18n,omitempty" json:"i18n,omitempty"`
	Compatibility *Compatibility    `yaml:"compatibility,omitempty" json:"compatibility,omitempty"`
	Dependencies  *Dependencies     `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
	AllowedTools  []string          `yaml:"allowed_tools,omitempty" json:"allowed_tools,omitempty"`
	Parameters    []Parameter       `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Entrypoint    *Entrypoint       `yaml:"entrypoint,omitempty" json:"entrypoint,omitempty"`
	Outputs       []Output          `yaml:"outputs,omitempty" json:"outputs,omitempty"`
	Examples      []Example         `yaml:"examples,omitempty" json:"examples,omitempty"`
	Documentation *Documentation    `yaml:"documentation,omitempty" json:"documentation,omitempty"`
	Maintainers   []Maintainer      `yaml:"maintainers,omitempty" json:"maintainers,omitempty"`
	Status        string            `yaml:"status,omitempty" json:"status,omitempty"`
	DeprecatedSince string          `yaml:"deprecated_since,omitempty" json:"deprecated_since,omitempty"`
	Replacement   string            `yaml:"replacement,omitempty" json:"replacement,omitempty"`
	Test          *Test             `yaml:"test,omitempty" json:"test,omitempty"`
	Hooks         *Hooks            `yaml:"hooks,omitempty" json:"hooks,omitempty"`
	Permissions   *Permissions      `yaml:"permissions,omitempty" json:"permissions,omitempty"`
}

type I18nMap map[string]string

type I18nDetail struct {
	DisplayName I18nMap `yaml:"display_name,omitempty" json:"display_name,omitempty"`
	Description I18nMap `yaml:"description,omitempty" json:"description,omitempty"`
}

type Compatibility struct {
	MinArtexVersion string   `yaml:"min_artex_version,omitempty" json:"min_artex_version,omitempty"`
	MaxArtexVersion string   `yaml:"max_artex_version,omitempty" json:"max_artex_version,omitempty"`
	Platforms       []string `yaml:"platforms,omitempty" json:"platforms,omitempty"`
	Architectures   []string `yaml:"architectures,omitempty" json:"architectures,omitempty"`
}

type Dependencies struct {
	System []SystemDep `yaml:"system,omitempty" json:"system,omitempty"`
	Python []PythonDep `yaml:"python,omitempty" json:"python,omitempty"`
	Node   []NodeDep   `yaml:"node,omitempty" json:"node,omitempty"`
	Go     []GoDep     `yaml:"go,omitempty" json:"go,omitempty"`
}

type SystemDep struct {
	Name           string `yaml:"name" json:"name"`
	Version        string `yaml:"version,omitempty" json:"version,omitempty"`
	CheckCommand   string `yaml:"check_command,omitempty" json:"check_command,omitempty"`
}

type PythonDep struct {
	Package string `yaml:"package" json:"package"`
	Version string `yaml:"version,omitempty" json:"version,omitempty"`
}

type NodeDep struct {
	Package string `yaml:"package" json:"package"`
	Version string `yaml:"version,omitempty" json:"version,omitempty"`
}

type GoDep struct {
	Module  string `yaml:"module" json:"module"`
	Version string `yaml:"version,omitempty" json:"version,omitempty"`
}

type Parameter struct {
	Name        string      `yaml:"name" json:"name"`
	Type        string      `yaml:"type" json:"type"`
	Label       I18nMap     `yaml:"label,omitempty" json:"label,omitempty"`
	Description I18nMap     `yaml:"description,omitempty" json:"description,omitempty"`
	Required    bool        `yaml:"required,omitempty" json:"required,omitempty"`
	Default     interface{} `yaml:"default,omitempty" json:"default,omitempty"`
	Enum        []EnumItem  `yaml:"enum,omitempty" json:"enum,omitempty"`
	Validation  *Validation `yaml:"validation,omitempty" json:"validation,omitempty"`
	Secret      bool        `yaml:"secret,omitempty" json:"secret,omitempty"`
}

type EnumItem struct {
	Value string `yaml:"value" json:"value"`
	Label I18nMap `yaml:"label,omitempty" json:"label,omitempty"`
}

type Validation struct {
	Pattern     string `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	Min         *float64 `yaml:"min,omitempty" json:"min,omitempty"`
	Max         *float64 `yaml:"max,omitempty" json:"max,omitempty"`
	MinLength   *int     `yaml:"min_length,omitempty" json:"min_length,omitempty"`
	MaxLength   *int     `yaml:"max_length,omitempty" json:"max_length,omitempty"`
}

type Entrypoint struct {
	Type        string            `yaml:"type" json:"type"`
	Command     string            `yaml:"command" json:"command"`
	Args        []string          `yaml:"args,omitempty" json:"args,omitempty"`
	Env         map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
	WorkingDir  string            `yaml:"working_dir,omitempty" json:"working_dir,omitempty"`
}

type Output struct {
	Name        string      `yaml:"name" json:"name"`
	Type        string      `yaml:"type" json:"type"`
	Description string      `yaml:"description,omitempty" json:"description,omitempty"`
	Schema      interface{} `yaml:"schema,omitempty" json:"schema,omitempty"`
}

type Example struct {
	Name        string                 `yaml:"name" json:"name"`
	Description string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Parameters  map[string]interface{} `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

type Documentation struct {
	Readme    string `yaml:"readme,omitempty" json:"readme,omitempty"`
	Wiki      string `yaml:"wiki,omitempty" json:"wiki,omitempty"`
	Changelog string `yaml:"changelog,omitempty" json:"changelog,omitempty"`
}

type Maintainer struct {
	Name  string `yaml:"name" json:"name"`
	Email string `yaml:"email,omitempty" json:"email,omitempty"`
	GitHub string `yaml:"github,omitempty" json:"github,omitempty"`
}

type Test struct {
	Command  string `yaml:"command,omitempty" json:"command,omitempty"`
	Timeout  int    `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Fixtures []Fixture `yaml:"fixtures,omitempty" json:"fixtures,omitempty"`
}

type Fixture struct {
	Name           string                 `yaml:"name" json:"name"`
	Input          map[string]interface{} `yaml:"input" json:"input"`
	ExpectedOutput map[string]interface{} `yaml:"expected_output" json:"expected_output"`
}

type Hooks struct {
	PreInstall string `yaml:"pre_install,omitempty" json:"pre_install,omitempty"`
	PostInstall string `yaml:"post_install,omitempty" json:"post_install,omitempty"`
	PreRun     string `yaml:"pre_run,omitempty" json:"pre_run,omitempty"`
	PostRun    string `yaml:"post_run,omitempty" json:"post_run,omitempty"`
	OnError    string `yaml:"on_error,omitempty" json:"on_error,omitempty"`
}

type Permissions struct {
	Network    bool     `yaml:"network,omitempty" json:"network,omitempty"`
	Filesystem *FSPerm  `yaml:"filesystem,omitempty" json:"filesystem,omitempty"`
	Process    bool     `yaml:"process,omitempty" json:"process,omitempty"`
	Privileged bool     `yaml:"privileged,omitempty" json:"privileged,omitempty"`
}

type FSPerm struct {
	Read  []string `yaml:"read,omitempty" json:"read,omitempty"`
	Write []string `yaml:"write,omitempty" json:"write,omitempty"`
}

// LoadSkillYAML loads and parses skill.yaml from a skill directory
func LoadSkillYAML(skillDir string) (*SkillYAML, error) {
	yamlPath := filepath.Join(skillDir, "skill.yaml")
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // skill.yaml not found, not an error
		}
		return nil, fmt.Errorf("read skill.yaml: %w", err)
	}

	var skill SkillYAML
	if err := yaml.Unmarshal(data, &skill); err != nil {
		return nil, fmt.Errorf("parse skill.yaml: %w", err)
	}

	// Set defaults
	if skill.Status == "" {
		skill.Status = "stable"
	}
	if skill.Test != nil && skill.Test.Timeout == 0 {
		skill.Test.Timeout = 300
	}

	return &skill, nil
}

// LoadAllSkillsYAML loads skill.yaml for all skills in the skills directory
func LoadAllSkillsYAML(skillsDir string) (map[string]*SkillYAML, error) {
	result := make(map[string]*SkillYAML)

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("read skills dir: %w", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		skillDir := filepath.Join(skillsDir, e.Name())
		skill, err := LoadSkillYAML(skillDir)
		if err != nil {
			// Log error but continue with other skills
			fmt.Printf("[skill_yaml] failed to load %s: %v\n", e.Name(), err)
			continue
		}
		if skill != nil {
			result[e.Name()] = skill
		}
	}

	return result, nil
}

// ValidateSkillYAML validates a SkillYAML struct
func ValidateSkillYAML(skill *SkillYAML) []string {
	var errors []string

	// Required fields
	if skill.Name == "" {
		errors = append(errors, "name is required")
	}
	if skill.DisplayName == nil || len(skill.DisplayName) == 0 {
		errors = append(errors, "display_name is required")
	}
	if skill.Description == nil || len(skill.Description) == 0 {
		errors = append(errors, "description is required")
	}
	if skill.Version == "" {
		errors = append(errors, "version is required")
	}
	if skill.Category == "" {
		errors = append(errors, "category is required")
	}
	if skill.Author == "" {
		errors = append(errors, "author is required")
	}
	if skill.License == "" {
		errors = append(errors, "license is required")
	}

	// Validate category
	validCategories := map[string]bool{
		"recon": true, "exploit": true, "post-exploit": true,
		"utility": true, "custom": true,
		// Kali Linux 16 categories (MITRE ATT&CK style)
		"resource-dev":      true,
		"initial-access":    true,
		"persistence":       true,
		"priv-esc":          true,
		"defense-evasion":   true,
		"credential-access": true,
		"discovery":         true,
		"lateral-movement":  true,
		"collection":        true,
		"c2":                true,
		"exfiltration":      true,
		"impact":            true,
		"forensics":         true,
		"services":          true,
		// opencode-pentester 16 attack domains
		"scanning-analysis":   true,
		"injection":           true,
		"authentication-tokens": true,
		"api-security":        true,
		"web-applications":    true,
		"client-side":         true,
		"server-side":         true,
		"cloud-containers":    true,
		"system-post-exploitation": true,
		"cve-exploit-research":  true,
		"ai-llm-security":     true,
		"secrets-sca-secrets": true,
		"secrets-auth-middleware": true,
		"secrets-access-control": true,
		// opencode-shannon-plugin 6 pentest phases (recon already exists)
		"vuln-discovery":     true,
		"browser-testing":    true,
		"idor-upload-testing": true,
		"exploitation":       true,
		"reporting":          true,
	}
	if skill.Category != "" && !validCategories[skill.Category] {
		errors = append(errors, fmt.Sprintf("invalid category: %s (must be one of: recon, exploit, post-exploit, utility, custom, resource-dev, initial-access, persistence, priv-esc, defense-evasion, credential-access, discovery, lateral-movement, collection, c2, exfiltration, impact, forensics, services, scanning-analysis, injection, authentication-tokens, api-security, web-applications, client-side, server-side, cloud-containers, system-post-exploitation, cve-exploit-research, ai-llm-security, secrets-sca-secrets, secrets-auth-middleware, secrets-access-control, vuln-discovery, browser-testing, idor-upload-testing, exploitation, reporting)", skill.Category))
	}

	// Validate status
	validStatuses := map[string]bool{
		"stable": true, "beta": true, "alpha": true,
		"deprecated": true, "archived": true,
	}
	if skill.Status != "" && !validStatuses[skill.Status] {
		errors = append(errors, fmt.Sprintf("invalid status: %s (must be one of: stable, beta, alpha, deprecated, archived)", skill.Status))
	}

	// Validate entrypoint
	if skill.Entrypoint != nil {
		validTypes := map[string]bool{
			"script": true, "binary": true, "docker": true, "mcp": true,
		}
		if skill.Entrypoint.Type != "" && !validTypes[skill.Entrypoint.Type] {
			errors = append(errors, fmt.Sprintf("invalid entrypoint type: %s (must be one of: script, binary, docker, mcp)", skill.Entrypoint.Type))
		}
	}

	// Validate parameters
	for i, param := range skill.Parameters {
		if param.Name == "" {
			errors = append(errors, fmt.Sprintf("parameter[%d]: name is required", i))
		}
		validTypes := map[string]bool{
			"string": true, "number": true, "boolean": true,
			"array": true, "object": true, "file": true, "secret": true,
		}
		if param.Type != "" && !validTypes[param.Type] {
			errors = append(errors, fmt.Sprintf("parameter[%s]: invalid type %s (must be one of: string, number, boolean, array, object, file, secret)", param.Name, param.Type))
		}
	}

	// Validate outputs
	for i, out := range skill.Outputs {
		if out.Name == "" {
			errors = append(errors, fmt.Sprintf("output[%d]: name is required", i))
		}
		validTypes := map[string]bool{
			"json": true, "text": true, "file": true, "asset": true, "finding": true,
		}
		if out.Type != "" && !validTypes[out.Type] {
			errors = append(errors, fmt.Sprintf("output[%s]: invalid type %s (must be one of: json, text, file, asset, finding)", out.Name, out.Type))
		}
	}

	return errors
}

// WriteSkillYAML writes a SkillYAML to skill.yaml file
func WriteSkillYAML(skillDir string, skill *SkillYAML) error {
	// Validate before writing
	if errors := ValidateSkillYAML(skill); len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}

	yamlPath := filepath.Join(skillDir, "skill.yaml")
	data, err := yaml.Marshal(skill)
	if err != nil {
		return fmt.Errorf("marshal skill.yaml: %w", err)
	}

	if err := os.WriteFile(yamlPath, data, 0o644); err != nil {
		return fmt.Errorf("write skill.yaml: %w", err)
	}

	return nil
}

// GetSkillYAMLSchema returns the JSON Schema for skill.yaml
func GetSkillYAMLSchema() (map[string]interface{}, error) {
	// Create a sample skill with all fields to generate schema
	sample := &SkillYAML{
		Name:        "example-skill",
		DisplayName: I18nMap{"en": "Example Skill", "zh": "示例技能"},
		Description: I18nMap{"en": "An example skill", "zh": "一個示例技能"},
		Version:     "1.0.0",
		Category:    "recon",
		Author:      "Author Name",
		License:     "MIT",
		Tags:        []string{"recon", "api"},
		I18n: &I18nDetail{
			DisplayName: I18nMap{"en": "Example Skill", "zh": "示例技能"},
			Description: I18nMap{"en": "An example skill", "zh": "一個示例技能"},
		},
		Compatibility: &Compatibility{
			MinArtexVersion: "0.3.0",
			Platforms:       []string{"linux", "darwin", "windows"},
			Architectures:   []string{"amd64", "arm64"},
		},
		Dependencies: &Dependencies{
			System: []SystemDep{{Name: "python3", Version: "3.9+", CheckCommand: "python3 --version"}},
			Python: []PythonDep{{Package: "requests"}},
			Node:   []NodeDep{{Package: "puppeteer", Version: "21+"}},
		},
		AllowedTools: []string{"Bash(python3:*)", "Bash(node:*)"},
		Parameters: []Parameter{
			{
				Name:        "target_url",
				Type:        "string",
				Label:       I18nMap{"en": "Target URL", "zh": "目標 URL"},
				Description: I18nMap{"en": "The target URL", "zh": "目標 URL"},
				Required:    true,
				Validation:  &Validation{Pattern: "^https?://.+"},
			},
		},
		Entrypoint: &Entrypoint{
			Type:       "script",
			Command:    "python3",
			Args:       []string{"scripts/main.py", "{{.target_url}}"},
			WorkingDir: ".",
		},
		Outputs: []Output{
			{Name: "result", Type: "json", Description: "Scan result"},
		},
		Examples: []Example{
			{Name: "Basic", Description: "Basic usage", Parameters: map[string]interface{}{"target_url": "https://example.com"}},
		},
		Documentation: &Documentation{Readme: "README.md"},
		Maintainers: []Maintainer{{Name: "Author", Email: "author@example.com"}},
		Status:      "stable",
		Test:        &Test{Command: "python3 -m pytest", Timeout: 300},
		Hooks:       &Hooks{PreRun: "scripts/pre_run.sh"},
		Permissions: &Permissions{Network: true, Process: true},
	}

	data, err := yaml.Marshal(sample)
	if err != nil {
		return nil, err
	}

	var schema map[string]interface{}
	if err := yaml.Unmarshal(data, &schema); err != nil {
		return nil, err
	}

	return schema, nil
}

// ScanSkillDir scans a skill directory for skill.yaml and SKILL.md
func ScanSkillDir(skillDir string) (hasYAML bool, hasMD bool, yamlPath string, mdPath string) {
	yamlPath = filepath.Join(skillDir, "skill.yaml")
	mdPath = filepath.Join(skillDir, "SKILL.md")

	_, err := os.Stat(yamlPath)
	hasYAML = err == nil

	_, err = os.Stat(mdPath)
	hasMD = err == nil

	return
}

// SkillYAMLToLegacy converts SkillYAML to legacy SKILL.md frontmatter format
func SkillYAMLToLegacy(skill *SkillYAML) (frontmatter map[string]interface{}, body string) {
	frontmatter = map[string]interface{}{
		"name":        skill.Name,
		"description": getFirstI18n(skill.Description),
	}
	if skill.License != "" {
		frontmatter["license"] = skill.License
	}
	if skill.Compatibility != nil && skill.Compatibility.MinArtexVersion != "" {
		frontmatter["compatibility"] = skill.Compatibility.MinArtexVersion
	}
	if len(skill.AllowedTools) > 0 {
		// Extract MCPs from allowed_tools (simplified)
		var mcps []string
		for _, tool := range skill.AllowedTools {
			if strings.HasPrefix(tool, "MCP(") {
				mcps = append(mcps, strings.TrimPrefix(strings.TrimSuffix(tool, ")"), "MCP("))
			}
		}
		if len(mcps) > 0 {
			frontmatter["mcps"] = strings.Join(mcps, ", ")
		}
	}

	// Generate markdown body from parameters, entrypoint, etc.
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## %s\n\n", getFirstI18n(skill.DisplayName)))
	sb.WriteString(getFirstI18n(skill.Description) + "\n\n")

	if len(skill.Parameters) > 0 {
		sb.WriteString("### Parameters\n\n")
		for _, p := range skill.Parameters {
			req := ""
			if p.Required {
				req = " (required)"
			}
			sb.WriteString(fmt.Sprintf("- **%s**%s: %s\n", p.Name, req, getFirstI18n(p.Description)))
		}
		sb.WriteString("\n")
	}

	if skill.Entrypoint != nil {
		sb.WriteString("### Entrypoint\n\n")
		sb.WriteString(fmt.Sprintf("- Type: %s\n", skill.Entrypoint.Type))
		sb.WriteString(fmt.Sprintf("- Command: %s\n", skill.Entrypoint.Command))
		if len(skill.Entrypoint.Args) > 0 {
			sb.WriteString(fmt.Sprintf("- Args: %s\n", strings.Join(skill.Entrypoint.Args, " ")))
		}
		sb.WriteString("\n")
	}

	if len(skill.Examples) > 0 {
		sb.WriteString("### Examples\n\n")
		for _, ex := range skill.Examples {
			sb.WriteString(fmt.Sprintf("#### %s\n%s\n\n", ex.Name, ex.Description))
			if len(ex.Parameters) > 0 {
				sb.WriteString("Parameters:\n")
				for k, v := range ex.Parameters {
					sb.WriteString(fmt.Sprintf("- %s: %v\n", k, v))
				}
				sb.WriteString("\n")
			}
		}
	}

	return frontmatter, sb.String()
}

func getFirstI18n(m I18nMap) string {
	for _, v := range m {
		return v
	}
	return ""
}

// LegacyToSkillYAML converts legacy SKILL.md to SkillYAML (best effort)
func LegacyToSkillYAML(frontmatter map[string]interface{}, body string) *SkillYAML {
	skill := &SkillYAML{
		Status: "stable",
	}

	if v, ok := frontmatter["name"].(string); ok {
		skill.Name = v
	}
	if v, ok := frontmatter["description"].(string); ok {
		skill.Description = I18nMap{"en": v, "zh": v}
		skill.DisplayName = I18nMap{"en": v, "zh": v}
	}
	if v, ok := frontmatter["license"].(string); ok {
		skill.License = v
	}
	if v, ok := frontmatter["compatibility"].(string); ok {
		skill.Compatibility = &Compatibility{MinArtexVersion: v}
	}
	if v, ok := frontmatter["mcps"].(string); ok {
		mcps := strings.Split(v, ",")
		for i := range mcps {
			mcps[i] = strings.TrimSpace(mcps[i])
		}
		skill.AllowedTools = append(skill.AllowedTools, mcps...)
	}

	return skill
}

// SkillDirInfo represents information about a skill directory
type SkillDirInfo struct {
	Name        string   `json:"name"`
	HasYAML     bool     `json:"has_yaml"`
	HasMD       bool     `json:"has_md"`
	YAMLPath    string   `json:"yaml_path,omitempty"`
	MDPath      string   `json:"md_path,omitempty"`
	Files       []string `json:"files"`
	YAMLContent *SkillYAML `json:"yaml_content,omitempty"`
	MDContent   string   `json:"md_content,omitempty"`
}

// GetSkillDirInfo gets complete information about a skill directory
func GetSkillDirInfo(skillsDir, skillName string) (*SkillDirInfo, error) {
	skillDir := filepath.Join(skillsDir, skillName)
	info := &SkillDirInfo{Name: skillName}

	hasYAML, hasMD, yamlPath, mdPath := ScanSkillDir(skillDir)
	info.HasYAML = hasYAML
	info.HasMD = hasMD
	info.YAMLPath = yamlPath
	info.MDPath = mdPath

	// List all files
	filepath.WalkDir(skillDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(skillDir, path)
		info.Files = append(info.Files, rel)
		return nil
	})

	// Load skill.yaml if exists
	if hasYAML {
		if skill, err := LoadSkillYAML(skillDir); err == nil {
			info.YAMLContent = skill
		}
	}

	// Load SKILL.md if exists
	if hasMD {
		if data, err := os.ReadFile(mdPath); err == nil {
			info.MDContent = string(data)
		}
	}

	return info, nil
}