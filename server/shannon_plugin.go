package server

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
		"net/http"
		"os"
		"path/filepath"
		"strings"

		"github.com/Autumn-27/artex/db"
)

// ShannonPhase represents one of the 6 opencode-shannon-plugin pentest phases
type ShannonPhase string

const (
	ShannonPhaseRecon              ShannonPhase = "recon"
	ShannonPhaseVulnDiscovery      ShannonPhase = "vuln-discovery"
	ShannonPhaseBrowserTesting     ShannonPhase = "browser-testing"
	ShannonPhaseIdorUploadTesting  ShannonPhase = "idor-upload-testing"
	ShannonPhaseExploitation       ShannonPhase = "exploitation"
	ShannonPhaseReporting          ShannonPhase = "reporting"
)

// AllShannonPhases returns the complete list of 6 opencode-shannon-plugin phases
func AllShannonPhases() []ShannonPhase {
	return []ShannonPhase{
		ShannonPhaseRecon,
		ShannonPhaseVulnDiscovery,
		ShannonPhaseBrowserTesting,
		ShannonPhaseIdorUploadTesting,
		ShannonPhaseExploitation,
		ShannonPhaseReporting,
	}
}

// ShannonTool represents an opencode-shannon-plugin tool entry
type ShannonTool struct {
	Name        string         `json:"name" yaml:"name"`
	Phase       ShannonPhase   `json:"phase" yaml:"phase"`
	Description string         `json:"description" yaml:"description"`
	OfficialURL string         `json:"official_url" yaml:"official_url"`
}

// ShannonToolsList is the static reference of opencode-shannon-plugin tools
// Sourced from opencode-shannon-plugin README (600+ Kali Linux tools, Playwright browser automation)
var ShannonToolsList = []ShannonTool{
	// Recon (偵察)
	{Name: "nmap", Phase: ShannonPhaseRecon, Description: "Network exploration tool and security / port scanner", OfficialURL: "https://nmap.org/"},
	{Name: "amass", Phase: ShannonPhaseRecon, Description: "In-depth attack surface mapping and asset discovery", OfficialURL: "https://github.com/owasp-amass/amass"},
	{Name: "subfinder", Phase: ShannonPhaseRecon, Description: "Subdomain discovery tool", OfficialURL: "https://github.com/projectdiscovery/subfinder"},
	{Name: "dnsrecon", Phase: ShannonPhaseRecon, Description: "DNS enumeration script", OfficialURL: "https://github.com/darkoperator/dnsrecon"},
	{Name: "theHarvester", Phase: ShannonPhaseRecon, Description: "E-mail, subdomain and names harvesting tool", OfficialURL: "https://github.com/laramies/theHarvester"},
	{Name: "spiderfoot", Phase: ShannonPhaseRecon, Description: "Automation tool for intelligence gathering", OfficialURL: "https://github.com/smicallef/spiderfoot"},
	{Name: "ffuf", Phase: ShannonPhaseRecon, Description: "Fast web fuzzer written in Go", OfficialURL: "https://github.com/ffuf/ffuf"},
	{Name: "gobuster", Phase: ShannonPhaseRecon, Description: "Directory/file & DNS busting tool written in Go", OfficialURL: "https://github.com/OJ/gobuster"},
	{Name: "dirsearch", Phase: ShannonPhaseRecon, Description: "Web path scanner", OfficialURL: "https://github.com/maurosoria/dirsearch"},
	{Name: "sherlock", Phase: ShannonPhaseRecon, Description: "Hunt down social media accounts by username across social networks", OfficialURL: "https://github.com/sherlock-project/sherlock"},
	{Name: "masscan", Phase: ShannonPhaseRecon, Description: "TCP port scanner, scans the entire Internet in under 6 minutes", OfficialURL: "https://github.com/robertdavidgraham/masscan"},
	{Name: "unicornscan", Phase: ShannonPhaseRecon, Description: "Asynchronous TCP and UDP port scanner", OfficialURL: "https://github.com/scannerman/unicornscan"},
	{Name: "zenmap", Phase: ShannonPhaseRecon, Description: "Official Nmap Security Scanner GUI", OfficialURL: "https://nmap.org/zenmap/"},
	{Name: "whatweb", Phase: ShannonPhaseRecon, Description: "Next generation web scanner", OfficialURL: "https://github.com/urbanadventurer/whatweb"},
	{Name: "wafw00f", Phase: ShannonPhaseRecon, Description: "Identifies and fingerprints Web Application Firewall (WAF) products", OfficialURL: "https://github.com/EnableSecurity/wafw00f"},

	// Vuln Discovery (漏洞發現)
	{Name: "nikto", Phase: ShannonPhaseVulnDiscovery, Description: "Web server scanner", OfficialURL: "https://cirt.net/nikto2/"},
	{Name: "sqlmap", Phase: ShannonPhaseVulnDiscovery, Description: "Automatic SQL injection and database takeover tool", OfficialURL: "http://sqlmap.org/"},
	{Name: "nuclei", Phase: ShannonPhaseVulnDiscovery, Description: "Fast and customizable vulnerability scanner", OfficialURL: "https://github.com/projectdiscovery/nuclei"},
	{Name: "wapiti", Phase: ShannonPhaseVulnDiscovery, Description: "Web application vulnerability scanner", OfficialURL: "https://wapiti.sc"},
	{Name: "skipfish", Phase: ShannonPhaseVulnDiscovery, Description: "Web application security scanner", OfficialURL: "https://github.com/google/skipfish"},
	{Name: "arachni", Phase: ShannonPhaseVulnDiscovery, Description: "Web application security scanner framework", OfficialURL: "https://github.com/Arachni/arachni"},
	{Name: "sslyze", Phase: ShannonPhaseVulnDiscovery, Description: "Fast SSL/TLS configuration scanner", OfficialURL: "https://github.com/nabla-c0d3/sslyze"},
	{Name: "testssl", Phase: ShannonPhaseVulnDiscovery, Description: "Testing TLS/SSL encryption", OfficialURL: "https://testssl.sh/"},
	{Name: "owasp-zap", Phase: ShannonPhaseVulnDiscovery, Description: "OWASP Zed Attack Proxy", OfficialURL: "https://www.zaproxy.org/"},
	{Name: "openvas", Phase: ShannonPhaseVulnDiscovery, Description: "Greenbone Vulnerability Scanner", OfficialURL: "https://www.greenbone.net/"},

	// Browser Testing (瀏覽器測試)
	{Name: "BrowserBruter", Phase: ShannonPhaseBrowserTesting, Description: "Browser-based form fuzzing — bypasses encrypted HTTP bodies, captchas, client-side validation", OfficialURL: "https://github.com/Ice3xSec/BrowserBruter"},
	{Name: "gowitness", Phase: ShannonPhaseBrowserTesting, Description: "Batch web screenshots with gallery UI, Nmap/Nessus integration", OfficialURL: "https://github.com/sensepost/gowitness"},
	{Name: "aquatone", Phase: ShannonPhaseBrowserTesting, Description: "A tool for domain takeover and visual reconnaissance", OfficialURL: "https://github.com/michenriksen/aquatone"},
	{Name: "eyewitness", Phase: ShannonPhaseBrowserTesting, Description: "Capture website screenshots and gather information", OfficialURL: "https://github.com/FortyNorthSecurity/EyeWitness"},
	{Name: "dalfox", Phase: ShannonPhaseBrowserTesting, Description: "Fast and powerful XSS scanner", OfficialURL: "https://github.com/hahwul/dalfox"},

	// IDOR/Upload Testing (IDOR/上傳測試)
	{Name: "burp", Phase: ShannonPhaseIdorUploadTesting, Description: "Web application security testing platform", OfficialURL: "https://portswigger.net/burp"},
	{Name: "owasp-zap", Phase: ShannonPhaseIdorUploadTesting, Description: "OWASP Zed Attack Proxy", OfficialURL: "https://www.zaproxy.org/"},
	{Name: "sqlmap", Phase: ShannonPhaseIdorUploadTesting, Description: "Automatic SQL injection and database takeover tool", OfficialURL: "http://sqlmap.org/"},
	{Name: "commix", Phase: ShannonPhaseIdorUploadTesting, Description: "Command injection exploitation and detection tool", OfficialURL: "https://github.com/commixproject/commix"},

	// Exploitation (利用)
	{Name: "hydra", Phase: ShannonPhaseExploitation, Description: "Network logon cracker supporting many services", OfficialURL: "https://github.com/vanhauser-thc/thc-hydra"},
	{Name: "hashcat", Phase: ShannonPhaseExploitation, Description: "World's fastest and most advanced password recovery utility", OfficialURL: "https://hashcat.net/hashcat/"},
	{Name: "john", Phase: ShannonPhaseExploitation, Description: "John the Ripper password cracker", OfficialURL: "https://www.openwall.com/john/"},
	{Name: "metasploit", Phase: ShannonPhaseExploitation, Description: "Penetration testing framework", OfficialURL: "https://www.metasploit.com/"},
	{Name: "beef", Phase: ShannonPhaseExploitation, Description: "Browser Exploitation Framework", OfficialURL: "https://beefproject.com/"},
	{Name: "aircrack-ng", Phase: ShannonPhaseExploitation, Description: "WiFi security auditing tools suite", OfficialURL: "https://www.aircrack-ng.org/"},
	{Name: "reaver", Phase: ShannonPhaseExploitation, Description: "WPS WiFi attack tool", OfficialURL: "https://github.com/t6x/reaver-wps-fork-t6x"},

	// Reporting (報告)
	{Name: "dradis", Phase: ShannonPhaseReporting, Description: "Collaboration and reporting platform for security teams", OfficialURL: "https://dradisframework.com/"},
	{Name: "faraday", Phase: ShannonPhaseReporting, Description: "Collaborative penetration testing and vulnerability management platform", OfficialURL: "https://faradaysec.com/"},
	{Name: "maltego", Phase: ShannonPhaseReporting, Description: "Open source intelligence and forensics application", OfficialURL: "https://www.maltego.com/"},
}

// ShannonToolsByPhase groups tools by their phase
func ShannonToolsByPhase() map[ShannonPhase][]ShannonTool {
	result := make(map[ShannonPhase][]ShannonTool)
	for _, d := range AllShannonPhases() {
		result[d] = []ShannonTool{}
	}
	for _, tool := range ShannonToolsList {
		result[tool.Phase] = append(result[tool.Phase], tool)
	}
	return result
}

// ShannonPhaseInfo represents phase metadata for API response
type ShannonPhaseInfo struct {
	Name       string `json:"name" yaml:"name"`
	ToolCount  int    `json:"tool_count" yaml:"tool_count"`
	SkillCount int    `json:"skill_count" yaml:"skill_count"`
}

// ShannonPhaseStats extends ShannonPhaseInfo with runtime call/finding counts.
type ShannonPhaseStats struct {
	Name         string `json:"name" yaml:"name"`
	ToolCount    int    `json:"tool_count" yaml:"tool_count"`
	SkillCount   int    `json:"skill_count" yaml:"skill_count"`
	CallCount    int    `json:"call_count" yaml:"call_count"`
	FindingCount int    `json:"finding_count" yaml:"finding_count"`
}

// ShannonToolInfo represents a tool entry for API response
type ShannonToolInfo struct {
	Name          string `json:"name" yaml:"name"`
	Phase         string `json:"phase" yaml:"phase"`
	Description   string `json:"description" yaml:"description"`
	OfficialURL   string `json:"official_url" yaml:"official_url"`
	SkillAvailable bool  `json:"skill_available" yaml:"skill_available"`
	SkillYAMLPath  string `json:"skill_yaml_path,omitempty" yaml:"skill_yaml_path,omitempty"`
}

// GetShannonPhases returns phase metadata with tool and skill counts
func GetShannonPhases(skillsDir string) ([]ShannonPhaseInfo, error) {
	toolsByPhase := ShannonToolsByPhase()
	skills, err := LoadAllSkillsYAML(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("load skills: %w", err)
	}

	skillPhaseCount := make(map[string]int)
	for _, skill := range skills {
		if skill.Category != "" {
			skillPhaseCount[skill.Category]++
		}
	}

	var result []ShannonPhaseInfo
	for _, d := range AllShannonPhases() {
		result = append(result, ShannonPhaseInfo{
			Name:       string(d),
			ToolCount:  len(toolsByPhase[d]),
			SkillCount: skillPhaseCount[string(d)],
		})
	}
	return result, nil
}

// GetShannonTools returns tool entries, optionally filtered by phase
func GetShannonTools(skillsDir string, phaseFilter string) ([]ShannonToolInfo, error) {
	toolsByPhase := ShannonToolsByPhase()
	skills, err := LoadAllSkillsYAML(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("load skills: %w", err)
	}

	skillNames := make(map[string]bool)
	for name := range skills {
		skillNames[name] = true
	}

	var result []ShannonToolInfo
	for _, d := range AllShannonPhases() {
		if phaseFilter != "" && string(d) != phaseFilter {
			continue
		}
		for _, tool := range toolsByPhase[d] {
			info := ShannonToolInfo{
				Name:          tool.Name,
				Phase:         string(tool.Phase),
				Description:   tool.Description,
				OfficialURL:   tool.OfficialURL,
				SkillAvailable: skillNames[tool.Name],
			}
			if skillNames[tool.Name] {
				info.SkillYAMLPath = filepath.Join(skillsDir, tool.Name, "skill.yaml")
			}
			result = append(result, info)
		}
	}
	return result, nil
}

// InstallShannonPluginToolAsSkill installs an opencode-shannon-plugin tool as an ARTEX skill
func InstallShannonPluginToolAsSkill(skillsDir string, toolName string) (string, *SkillYAML, error) {
	var tool *ShannonTool
	for i := range ShannonToolsList {
		if ShannonToolsList[i].Name == toolName {
			tool = &ShannonToolsList[i]
			break
		}
	}
	if tool == nil {
		return "", nil, fmt.Errorf("shannon plugin tool not found: %s", toolName)
	}

	skillDir := filepath.Join(skillsDir, tool.Name)

	if _, err := os.Stat(skillDir); err == nil {
		return "", nil, fmt.Errorf("skill already exists: %s", toolName)
	}

	skill := &SkillYAML{
		Name:        tool.Name,
		DisplayName: I18nMap{"en": tool.Name, "zh": tool.Name},
		Description: I18nMap{"en": tool.Description, "zh": tool.Description},
		Version:     "1.0.0",
		Category:    string(tool.Phase),
		Author:      "opencode-shannon-plugin",
		License:     "AGPL-3.0",
		Tags:        []string{string(tool.Phase), "shannon-plugin"},
		Compatibility: &Compatibility{
			MinArtexVersion: "0.3.0",
			Platforms:       []string{"linux"},
			Architectures:   []string{"amd64"},
		},
		AllowedTools: []string{fmt.Sprintf("Bash(*%s*)", tool.Name)},
		Parameters: []Parameter{
			{
				Name:        "target",
				Type:        "string",
				Label:       I18nMap{"en": "Target", "zh": "目標"},
				Description: I18nMap{"en": "Target for the tool", "zh": "工具目標"},
				Required:    true,
			},
		},
		Entrypoint: &Entrypoint{
			Type:         "binary",
			Command:      tool.Name,
			Args:         []string{"{{.target}}"},
			WorkingDir:   ".",
		},
		Outputs: []Output{
			{Name: "result", Type: "text", Description: "Tool output"},
		},
		Examples: []Example{
			{Name: "Basic", Description: fmt.Sprintf("Run %s against a target", tool.Name), Parameters: map[string]interface{}{"target": "example.com"}},
		},
		Status: "stable",
	}

	if errors := ValidateSkillYAML(skill); len(errors) > 0 {
		return "", nil, fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}

	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		return "", nil, fmt.Errorf("create skill dir: %w", err)
	}

	yamlPath := filepath.Join(skillDir, "skill.yaml")
	data, err := yaml.Marshal(skill)
	if err != nil {
		return "", nil, fmt.Errorf("marshal skill.yaml: %w", err)
	}
	if err := os.WriteFile(yamlPath, data, 0o644); err != nil {
		return "", nil, fmt.Errorf("write skill.yaml: %w", err)
	}

	return yamlPath, skill, nil
}

// ShannonAgent represents a Shannon plugin agent role
type ShannonAgent struct {
	Name                string   `json:"name" yaml:"name"`
	Role                string   `json:"role" yaml:"role"`
	File                string   `json:"file" yaml:"file"`
	SupportedPhases     []string `json:"supported_phases" yaml:"supported_phases"`
}

// AllShannonAgents returns the complete list of opencode-shannon-plugin agents
func AllShannonAgents() []ShannonAgent {
	return []ShannonAgent{
		{Name: "shannon-recon", Role: "reconnaissance", File: "plugins/shannon/agents/recon-agent.md", SupportedPhases: []string{"recon", "vuln-discovery", "scanning-analysis"}},
		{Name: "shannon-vuln", Role: "vulnerability discovery", File: "plugins/shannon/agents/vuln-agent.md", SupportedPhases: []string{"vuln-discovery", "scanning-analysis", "server-side", "web-applications"}},
		{Name: "shannon-browser", Role: "browser testing", File: "plugins/shannon/agents/browser-agent.md", SupportedPhases: []string{"browser-testing", "web-applications", "client-side", "server-side"}},
		{Name: "shannon-idor", Role: "IDOR/upload testing", File: "plugins/shannon/agents/idor-agent.md", SupportedPhases: []string{"idor-upload-testing", "server-side", "web-applications"}},
		{Name: "shannon-exploit", Role: "exploitation", File: "plugins/shannon/agents/exploit-agent.md", SupportedPhases: []string{"exploitation", "injection", "server-side", "web-applications", "client-side"}},
		{Name: "shannon-report", Role: "reporting", File: "plugins/shannon/agents/report-agent.md", SupportedPhases: []string{"reporting", "recon", "discovery", "scanning-analysis", "injection", "web-applications", "client-side", "server-side", "cloud-containers", "system-post-exploitation", "cve-exploit-research", "ai-llm-security", "secrets-sca-secrets", "secrets-auth-middleware", "secrets-access-control"}},
	}
}

// GetShannonAgents returns all agent definitions
func GetShannonAgents() []ShannonAgent {
	return AllShannonAgents()
}

// shannonPhases handles GET /api/shannon-phases
func (s *Server) shannonPhases(w http.ResponseWriter, r *http.Request) {
	phases, err := GetShannonPhases(s.skillDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"phases": phases}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// shannonTools handles GET /api/shannon-tools
func (s *Server) shannonTools(w http.ResponseWriter, r *http.Request) {
	phase := r.URL.Query().Get("phase")
	tools, err := GetShannonTools(s.skillDir, phase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"tools": tools}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// shannonStats handles GET /api/shannon-stats
func (s *Server) shannonStats(w http.ResponseWriter, r *http.Request) {
	pg := s.pg(w)
	if pg == nil {
		http.Error(w, "database unavailable", http.StatusInternalServerError)
		return
	}
	stats, err := GetShannonStats(pg, s.skillDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"stats": stats}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// shannonToolInstall handles POST /api/shannon-tools/{name}/install
func (s *Server) shannonToolInstall(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "missing tool name", http.StatusBadRequest)
		return
	}
	yamlPath, skill, err := InstallShannonPluginToolAsSkill(s.skillDir, name)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"skill_path": yamlPath,
		"skill_yaml": skill,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Record skill installation in the existing skill_usage ledger.
	if pg := s.pg(w); pg != nil {
		if err := pg.InsertSkillUsage(&db.SkillUsage{
			Skill:         skill.Name,
			AgentKey:      "api",
			TaskID:        0,
			ExplorationID: 0,
			IntentID:      0,
			SessionID:     "",
			ArgsLen:       0,
			Found:         true,
		}); err != nil {
			log.Printf("[skillusage] install %s: %v", skill.Name, err)
		}
	}
}

// GetShannonStats returns per-phase runtime stats: tool_count, skill_count,
// call_count (from skill_usage ledger), finding_count (best-effort).
func GetShannonStats(pg *db.DB, skillsDir string) ([]ShannonPhaseStats, error) {
	base, err := GetShannonPhases(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("base phases: %w", err)
	}

	skillCalls, err := pg.SkillStats()
	if err != nil {
		return nil, fmt.Errorf("list skill stats: %w", err)
	}
	callBySkill := make(map[string]int)
	for _, s := range skillCalls {
		callBySkill[s.Skill] = s.Calls
	}

	skills, err := LoadAllSkillsYAML(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("load skills: %w", err)
	}
	skillToPhase := make(map[string]string)
	for name, skill := range skills {
		if skill.Category != "" {
			skillToPhase[name] = skill.Category
		}
	}

	callByPhase := make(map[string]int)
	for skillName, calls := range callBySkill {
		if phase, ok := skillToPhase[skillName]; ok {
			callByPhase[phase] += calls
		}
	}

	var result []ShannonPhaseStats
	for _, d := range AllShannonPhases() {
		name := string(d)
		result = append(result, ShannonPhaseStats{
			Name:         name,
			ToolCount:    baseShannonToolCount(base, name),
			SkillCount:   baseShannonSkillCount(base, name),
			CallCount:    callByPhase[name],
			FindingCount: 0,
		})
	}
	return result, nil
}

func baseShannonToolCount(base []ShannonPhaseInfo, name string) int {
	for _, b := range base {
		if b.Name == name {
			return b.ToolCount
		}
	}
	return 0
}

func baseShannonSkillCount(base []ShannonPhaseInfo, name string) int {
	for _, b := range base {
		if b.Name == name {
			return b.SkillCount
		}
	}
	return 0
}