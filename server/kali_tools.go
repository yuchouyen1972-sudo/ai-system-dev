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

// KaliCategory represents one of the 16 Kali Linux tool categories
// (MITRE ATT&CK style organization from kali.org/tools/)
type KaliCategory string

const (
	KaliCatRecon          KaliCategory = "recon"
	KaliCatResourceDev    KaliCategory = "resource-dev"
	KaliCatInitialAccess  KaliCategory = "initial-access"
	KaliCatExecution      KaliCategory = "execution"
	KaliCatPersistence    KaliCategory = "persistence"
	KaliCatPrivEsc        KaliCategory = "priv-esc"
	KaliCatDefenseEvasion KaliCategory = "defense-evasion"
	KaliCatCredentialAccess KaliCategory = "credential-access"
	KaliCatDiscovery      KaliCategory = "discovery"
	KaliCatLateralMovement KaliCategory = "lateral-movement"
	KaliCatCollection     KaliCategory = "collection"
	KaliCatC2             KaliCategory = "c2"
	KaliCatExfiltration   KaliCategory = "exfiltration"
	KaliCatImpact         KaliCategory = "impact"
	KaliCatForensics      KaliCategory = "forensics"
	KaliCatServices       KaliCategory = "services"
)

// AllKaliCategories returns the complete list of 16 Kali categories
func AllKaliCategories() []KaliCategory {
	return []KaliCategory{
		KaliCatRecon,
		KaliCatResourceDev,
		KaliCatInitialAccess,
		KaliCatExecution,
		KaliCatPersistence,
		KaliCatPrivEsc,
		KaliCatDefenseEvasion,
		KaliCatCredentialAccess,
		KaliCatDiscovery,
		KaliCatLateralMovement,
		KaliCatCollection,
		KaliCatC2,
		KaliCatExfiltration,
		KaliCatImpact,
		KaliCatForensics,
		KaliCatServices,
	}
}

// KaliTool represents a Kali Linux tool entry
type KaliTool struct {
	Name        string         `json:"name" yaml:"name"`
	Category    KaliCategory   `json:"category" yaml:"category"`
	Description string         `json:"description" yaml:"description"`
	OfficialURL string         `json:"official_url" yaml:"official_url"`
}

// KaliToolsList is the static reference of Kali Linux tools
// Sourced from kali.org/tools/ (16 categories, ~300+ tools; representative subset below)
var KaliToolsList = []KaliTool{
	// Reconnaissance
	{Name: "nmap", Category: KaliCatRecon, Description: "Network exploration tool and security / port scanner", OfficialURL: "https://www.kali.org/tools/nmap/"},
	{Name: "amass", Category: KaliCatRecon, Description: "In-depth attack surface mapping and asset discovery", OfficialURL: "https://www.kali.org/tools/amass/"},
	{Name: "theHarvester", Category: KaliCatRecon, Description: "E-mail, subdomain and names harvesting tool", OfficialURL: "https://www.kali.org/tools/theharvester/"},
	{Name: "sherlock", Category: KaliCatRecon, Description: "Hunt down social media accounts by username across social networks", OfficialURL: "https://www.kali.org/tools/sherlock/"},
	{Name: "gobuster", Category: KaliCatRecon, Description: "Directory/file & DNS busting tool written in Go", OfficialURL: "https://www.kali.org/tools/gobuster/"},
	{Name: "ffuf", Category: KaliCatRecon, Description: "Fast web fuzzer written in Go", OfficialURL: "https://www.kali.org/tools/ffuf/"},
	{Name: "spiderfoot", Category: KaliCatRecon, Description: "Automation tool for intelligence gathering", OfficialURL: "https://www.kali.org/tools/spiderfoot/"},
	{Name: "dnsrecon", Category: KaliCatRecon, Description: "DNS enumeration script", OfficialURL: "https://www.kali.org/tools/dnsrecon/"},
	{Name: "whatweb", Category: KaliCatRecon, Description: "Next generation web scanner", OfficialURL: "https://www.kali.org/tools/whatweb/"},
	{Name: "wafw00f", Category: KaliCatRecon, Description: "Identifies and fingerprints Web Application Firewall (WAF) products", OfficialURL: "https://www.kali.org/tools/wafw00f/"},

	// Resource Development
	{Name: "msfvenom", Category: KaliCatResourceDev, Description: "Metasploit payload generator", OfficialURL: "https://www.kali.org/tools/metasploit-framework/#msfvenom"},
	{Name: "donut", Category: KaliCatResourceDev, Description: "Shellcode loader that can execute .NET assemblies", OfficialURL: "https://www.kali.org/tools/donut-shellcode/"},
	{Name: "pyinstaller", Category: KaliCatResourceDev, Description: "Converts Python applications into stand-alone executables", OfficialURL: "https://www.kali.org/tools/pyinstaller/"},
	{Name: "apktool", Category: KaliCatResourceDev, Description: "Tool for reverse engineering Android APK files", OfficialURL: "https://www.kali.org/tools/apktool/"},
	{Name: "jadx-gui", Category: KaliCatResourceDev, Description: "GUI for decompiling Android DEX files to Java source", OfficialURL: "https://www.kali.org/tools/jadx/"},
	{Name: "ghidra", Category: KaliCatResourceDev, Description: "Software reverse engineering (SRE) framework", OfficialURL: "https://www.kali.org/tools/ghidra/"},
	{Name: "rizin", Category: KaliCatResourceDev, Description: "Reverse engineering framework", OfficialURL: "https://www.kali.org/tools/rizin/"},
	{Name: "searchsploit", Category: KaliCatResourceDev, Description: "Command line search tool for Exploit-DB", OfficialURL: "https://www.kali.org/tools/exploitdb/#searchsploit"},

	// Initial Access
	{Name: "sqlmap", Category: KaliCatInitialAccess, Description: "Automatic SQL injection and database takeover tool", OfficialURL: "https://www.kali.org/tools/sqlmap/"},
	{Name: "commix", Category: KaliCatInitialAccess, Description: "Command injection exploitation and detection tool", OfficialURL: "https://www.kali.org/tools/commix/"},
	{Name: "gophish", Category: KaliCatInitialAccess, Description: "Open-source phishing campaign platform", OfficialURL: "https://www.kali.org/tools/gophish/"},
	{Name: "setoolkit", Category: KaliCatInitialAccess, Description: "Social-Engineering Toolkit", OfficialURL: "https://www.kali.org/tools/set/"},

	// Execution
	{Name: "armitage", Category: KaliCatExecution, Description: "Graphical cyber attack management tool for Metasploit", OfficialURL: "https://www.kali.org/tools/armitage/"},
	{Name: "beef-xss", Category: KaliCatExecution, Description: "Browser Exploitation Framework Project", OfficialURL: "https://www.kali.org/tools/beef-xss/"},
	{Name: "xsser", Category: KaliCatExecution, Description: "Cross Site Scripter - XSS testing tool", OfficialURL: "https://www.kali.org/tools/xsser/"},

	// Persistence
	{Name: "weevely", Category: KaliCatPersistence, Description: "Stealth PHP web shell", OfficialURL: "https://www.kali.org/tools/weevely/"},
	{Name: "webshells", Category: KaliCatPersistence, Description: "Web shell collections", OfficialURL: "https://www.kali.org/tools/webshells/"},

	// Privilege Escalation
	{Name: "linpeas", Category: KaliCatPrivEsc, Description: "Linux Privilege Escalation Awesome Script", OfficialURL: "https://www.kali.org/tools/peass-ng/#linpeas"},
	{Name: "winpeas", Category: KaliCatPrivEsc, Description: "Windows Privilege Escalation Awesome Script", OfficialURL: "https://www.kali.org/tools/peass-ng/#winpeas"},
	{Name: "bloodyad", Category: KaliCatPrivEsc, Description: "Active Directory privilege escalation using ABCD attacks", OfficialURL: "https://www.kali.org/tools/bloodyad/"},

	// Defense Evasion
	{Name: "mimikatz", Category: KaliCatDefenseEvasion, Description: "Extract passwords, hashes, PIN codes and kerberos tickets", OfficialURL: "https://www.kali.org/tools/mimikatz/"},
	{Name: "evil-winrm", Category: KaliCatDefenseEvasion, Description: "The ultimate WinRM shell for hacking/pentesting", OfficialURL: "https://www.kali.org/tools/evil-winrm/"},
	{Name: "steghide", Category: KaliCatDefenseEvasion, Description: "Steganography tool to hide data in images and audio", OfficialURL: "https://www.kali.org/tools/steghide/"},
	{Name: "donut", Category: KaliCatDefenseEvasion, Description: "Shellcode loader for AV evasion", OfficialURL: "https://www.kali.org/tools/donut-shellcode/"},

	// Credential Access
	{Name: "hashcat", Category: KaliCatCredentialAccess, Description: "World's fastest and most advanced password recovery utility", OfficialURL: "https://www.kali.org/tools/hashcat/"},
	{Name: "john", Category: KaliCatCredentialAccess, Description: "John the Ripper password cracker", OfficialURL: "https://www.kali.org/tools/john/"},
	{Name: "hydra", Category: KaliCatCredentialAccess, Description: "Network logon cracker supporting many services", OfficialURL: "https://www.kali.org/tools/hydra/"},
	{Name: "aircrack-ng", Category: KaliCatCredentialAccess, Description: "WiFi security auditing tools suite", OfficialURL: "https://www.kali.org/tools/aircrack-ng/"},
	{Name: "wifite", Category: KaliCatCredentialAccess, Description: "Automated wireless auditing tool", OfficialURL: "https://www.kali.org/tools/wifite/"},
	{Name: "trufflehog", Category: KaliCatCredentialAccess, Description: "Find leaked secrets in git repositories", OfficialURL: "https://www.kali.org/tools/trufflehog/"},
	{Name: "kerberoast", Category: KaliCatCredentialAccess, Description: "Kerberoasting attack toolkit", OfficialURL: "https://www.kali.org/tools/kerberoast/"},

	// Discovery
	{Name: "masscan", Category: KaliCatDiscovery, Description: "TCP port scanner, scans the entire Internet in under 6 minutes", OfficialURL: "https://www.kali.org/tools/masscan/"},
	{Name: "sslscan", Category: KaliCatDiscovery, Description: "Fast SSL/TLS configuration scanner", OfficialURL: "https://www.kali.org/tools/sslscan/"},
	{Name: "wireshark", Category: KaliCatDiscovery, Description: "Network protocol analyzer", OfficialURL: "https://www.kali.org/tools/wireshark/"},
	{Name: "bloodhound", Category: KaliCatDiscovery, Description: "Active Directory trust relationship visualization", OfficialURL: "https://www.kali.org/tools/bloodhound/"},
	{Name: "fierce", Category: KaliCatDiscovery, Description: "DNS reconnaissance tool for locating non-contiguous IP space", OfficialURL: "https://www.kali.org/tools/fierce/"},
	{Name: "hping3", Category: KaliCatDiscovery, Description: "TCP/IP packet assembler/analyzer", OfficialURL: "https://www.kali.org/tools/hping3/"},

	// Lateral Movement
	{Name: "crackmapexec", Category: KaliCatLateralMovement, Description: "Swiss army knife for pentesting networks", OfficialURL: "https://www.kali.org/tools/crackmapexec/"},
	{Name: "impacket-psexec", Category: KaliCatLateralMovement, Description: "Executes commands on remote Windows systems via SMB", OfficialURL: "https://www.kali.org/tools/impacket-scripts/#impacket-psexec"},
	{Name: "rdesktop", Category: KaliCatLateralMovement, Description: "Remote Desktop Protocol client", OfficialURL: "https://www.kali.org/tools/rdesktop/"},
	{Name: "smbmap", Category: KaliCatLateralMovement, Description: "SMB enumeration tool", OfficialURL: "https://www.kali.org/tools/smbmap/"},
	{Name: "evil-winrm", Category: KaliCatLateralMovement, Description: "WinRM shell for lateral movement", OfficialURL: "https://www.kali.org/tools/evil-winrm/"},

	// Collection
	{Name: "ettercap", Category: KaliCatCollection, Description: "Network sniffer/interceptor/logger for Ethernet LAN", OfficialURL: "https://www.kali.org/tools/ettercap/"},
	{Name: "mitmproxy", Category: KaliCatCollection, Description: "Interactive HTTPS proxy", OfficialURL: "https://www.kali.org/tools/mitmproxy/"},
	{Name: "evilginx2", Category: KaliCatCollection, Description: "Man-in-the-middle attack framework for phishing credentials", OfficialURL: "https://www.kali.org/tools/evilginx2/"},
	{Name: "sslsplit", Category: KaliCatCollection, Description: "Transparent SSL/TLS interception", OfficialURL: "https://www.kali.org/tools/sslsplit/"},

	// Command and Control
	{Name: "metasploit-framework", Category: KaliCatC2, Description: "Penetration testing framework with C2 capabilities", OfficialURL: "https://www.kali.org/tools/metasploit-framework/"},
	{Name: "chisel", Category: KaliCatC2, Description: "Fast TCP/UDP tunnel over HTTP", OfficialURL: "https://www.kali.org/tools/chisel/"},
	{Name: "sshuttle", Category: KaliCatC2, Description: "Transparent proxy server over SSH", OfficialURL: "https://www.kali.org/tools/sshuttle/"},
	{Name: "netcat", Category: KaliCatC2, Description: "TCP/UDP swiss army knife", OfficialURL: "https://www.kali.org/tools/nmap/#ncat-1"},

	// Exfiltration
	{Name: "netcat", Category: KaliCatExfiltration, Description: "TCP/UDP swiss army knife for data exfiltration", OfficialURL: "https://www.kali.org/tools/nmap/#ncat-1"},
	{Name: "impacket-smbserver", Category: KaliCatExfiltration, Description: "SMB server for file exfiltration", OfficialURL: "https://www.kali.org/tools/impacket-scripts/#impacket-smbserver"},
	{Name: "goshs", Category: KaliCatExfiltration, Description: "Go-based HTTP server for file serving", OfficialURL: "https://www.kali.org/tools/goshs/"},

	// Impact
	{Name: "siege", Category: KaliCatImpact, Description: "HTTP load testing and benchmarking utility", OfficialURL: "https://www.kali.org/tools/siege/"},
	{Name: "slowhttptest", Category: KaliCatImpact, Description: "Simulates application layer denial of service attacks", OfficialURL: "https://www.kali.org/tools/slowhttptest/"},
	{Name: "thc-ssl-dos", Category: KaliCatImpact, Description: "SSL denial of service test tool", OfficialURL: "https://www.kali.org/tools/thc-ssl-dos/"},
	{Name: "goldeneye", Category: KaliCatImpact, Description: "HTTP DoS test tool", OfficialURL: "https://www.kali.org/tools/goldeneye/"},

	// Forensics
	{Name: "autopsy", Category: KaliCatForensics, Description: "Digital forensics platform", OfficialURL: "https://www.kali.org/tools/autopsy/"},
	{Name: "binwalk", Category: KaliCatForensics, Description: "Firmware analysis and extraction tool", OfficialURL: "https://www.kali.org/tools/binwalk/"},
	{Name: "foremost", Category: KaliCatForensics, Description: "File carving utility", OfficialURL: "https://www.kali.org/tools/foremost/"},
	{Name: "yara", Category: KaliCatForensics, Description: "Tool for identifying and classifying malware samples", OfficialURL: "https://www.kali.org/tools/yara/"},
	{Name: "photorec", Category: KaliCatForensics, Description: "Data recovery tool for lost files", OfficialURL: "https://www.kali.org/tools/testdisk/#photorec-1"},

	// Services and Other Tools
	{Name: "dradis", Category: KaliCatServices, Description: "Collaboration and reporting platform for security teams", OfficialURL: "https://www.kali.org/tools/dradis/"},
	{Name: "faraday", Category: KaliCatServices, Description: "Collaborative penetration testing and vulnerability management platform", OfficialURL: "https://www.kali.org/tools/python-faraday/"},
	{Name: "maltego", Category: KaliCatServices, Description: "Open source intelligence and forensics application", OfficialURL: "https://www.kali.org/tools/maltego/"},
	{Name: "gvm", Category: KaliCatServices, Description: "Greenbone Vulnerability Management vulnerability scanner", OfficialURL: "https://www.kali.org/tools/gvm/"},
}

// KaliToolsByCategory groups tools by their Kali category
func KaliToolsByCategory() map[KaliCategory][]KaliTool {
	result := make(map[KaliCategory][]KaliTool)
	for _, cat := range AllKaliCategories() {
		result[cat] = []KaliTool{}
	}
	for _, tool := range KaliToolsList {
		result[tool.Category] = append(result[tool.Category], tool)
	}
	return result
}

// KaliCategoryInfo represents category metadata for API response
type KaliCategoryInfo struct {
	Name       string `json:"name" yaml:"name"`
	ToolCount  int    `json:"tool_count" yaml:"tool_count"`
	SkillCount int    `json:"skill_count" yaml:"skill_count"`
}

// KaliCategoryStats extends KaliCategoryStats with runtime call/finding counts.
type KaliCategoryStats struct {
	Name        string `json:"name" yaml:"name"`
	ToolCount   int    `json:"tool_count" yaml:"tool_count"`
	SkillCount  int    `json:"skill_count" yaml:"skill_count"`
	CallCount   int    `json:"call_count" yaml:"call_count"`
	FindingCount int   `json:"finding_count" yaml:"finding_count"`
}

// KaliToolInfo represents a tool entry for API response
type KaliToolInfo struct {
	Name          string `json:"name" yaml:"name"`
	Category      string `json:"category" yaml:"category"`
	Description   string `json:"description" yaml:"description"`
	OfficialURL   string `json:"official_url" yaml:"official_url"`
	SkillAvailable bool  `json:"skill_available" yaml:"skill_available"`
	SkillYAMLPath  string `json:"skill_yaml_path,omitempty" yaml:"skill_yaml_path,omitempty"`
}

// GetKaliCategories returns category metadata with tool and skill counts
func GetKaliCategories(skillsDir string) ([]KaliCategoryInfo, error) {
	toolsByCategory := KaliToolsByCategory()
	skills, err := LoadAllSkillsYAML(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("load skills: %w", err)
	}

	// Count skills per category
	skillCategoryCount := make(map[string]int)
	for _, skill := range skills {
		if skill.Category != "" {
			skillCategoryCount[skill.Category]++
		}
	}

	var result []KaliCategoryInfo
	for _, cat := range AllKaliCategories() {
		result = append(result, KaliCategoryInfo{
			Name:       string(cat),
			ToolCount:  len(toolsByCategory[cat]),
			SkillCount: skillCategoryCount[string(cat)],
		})
	}
	return result, nil
}

// GetKaliStats returns per-category runtime stats: tool_count, skill_count,
// call_count (from skill_usage ledger), finding_count (from findings, best-effort).
func GetKaliStats(pg *db.DB, skillsDir string) ([]KaliCategoryStats, error) {
	base, err := GetKaliCategories(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("base categories: %w", err)
	}

	// Aggregate call counts per skill from the ledger.
	skillCalls, err := pg.SkillStats()
	if err != nil {
		return nil, fmt.Errorf("list skill stats: %w", err)
	}
	callBySkill := make(map[string]int)
	for _, s := range skillCalls {
		callBySkill[s.Skill] = s.Calls
	}

	// Map skill directory names to their Kali category.
	skills, err := LoadAllSkillsYAML(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("load skills: %w", err)
	}
	skillToCategory := make(map[string]string)
	for name, skill := range skills {
		if skill.Category != "" {
			skillToCategory[name] = skill.Category
		}
	}

	callByCategory := make(map[string]int)
	for skillName, calls := range callBySkill {
		if cat, ok := skillToCategory[skillName]; ok {
			callByCategory[cat] += calls
		}
	}

	var result []KaliCategoryStats
	for _, cat := range AllKaliCategories() {
		name := string(cat)
		result = append(result, KaliCategoryStats{
			Name:         name,
			ToolCount:    baseCallCount(base, name),
			SkillCount:   baseSkillCount(base, name),
			CallCount:    callByCategory[name],
			FindingCount: 0, // best-effort: findings aggregation wired in a later pass
		})
	}
	return result, nil
}

func baseCallCount(base []KaliCategoryInfo, name string) int {
	for _, b := range base {
		if b.Name == name {
			return b.ToolCount
		}
	}
	return 0
}

func baseSkillCount(base []KaliCategoryInfo, name string) int {
	for _, b := range base {
		if b.Name == name {
			return b.SkillCount
		}
	}
	return 0
}

// GetKaliTools returns tool entries, optionally filtered by category
func GetKaliTools(skillsDir, categoryFilter string) ([]KaliToolInfo, error) {
	toolsByCategory := KaliToolsByCategory()
	skills, err := LoadAllSkillsYAML(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("load skills: %w", err)
	}

	// Build set of skill names that have skill.yaml
	skillNames := make(map[string]bool)
	for name := range skills {
		skillNames[name] = true
	}

	var result []KaliToolInfo
	for _, cat := range AllKaliCategories() {
		if categoryFilter != "" && string(cat) != categoryFilter {
			continue
		}
		for _, tool := range toolsByCategory[cat] {
			info := KaliToolInfo{
				Name:          tool.Name,
				Category:      string(tool.Category),
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

// InstallKaliToolAsSkill installs a Kali tool as an ARTEX skill by generating skill.yaml
func InstallKaliToolAsSkill(skillsDir, toolName string) (string, *SkillYAML, error) {
	// Find the tool
	var tool *KaliTool
	for i := range KaliToolsList {
		if KaliToolsList[i].Name == toolName {
			tool = &KaliToolsList[i]
			break
		}
	}
	if tool == nil {
		return "", nil, fmt.Errorf("kali tool not found: %s", toolName)
	}

	skillDir := filepath.Join(skillsDir, toolName)

	// Check if skill already exists
	if _, err := os.Stat(skillDir); err == nil {
		return "", nil, fmt.Errorf("skill already exists: %s", toolName)
	}

	// Generate skill.yaml from Kali tool definition
	skill := &SkillYAML{
		Name:        tool.Name,
		DisplayName: I18nMap{"en": tool.Name, "zh": tool.Name},
		Description: I18nMap{"en": tool.Description, "zh": tool.Description},
		Version:     "1.0.0",
		Category:    string(tool.Category),
		Author:      "Kali Linux",
		License:     "See kali.org/tools/",
		Tags:        []string{string(tool.Category), "kali"},
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
			Type:     "binary",
			Command:  tool.Name,
			Args:     []string{"{{.target}}"},
			WorkingDir: ".",
		},
		Outputs: []Output{
			{Name: "result", Type: "text", Description: "Tool output"},
		},
		Examples: []Example{
			{Name: "Basic", Description: fmt.Sprintf("Run %s against a target", tool.Name), Parameters: map[string]interface{}{"target": "example.com"}},
		},
		Status: "stable",
	}

	// Validate before writing
	if errors := ValidateSkillYAML(skill); len(errors) > 0 {
		return "", nil, fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}

	// Create skill directory
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		return "", nil, fmt.Errorf("create skill dir: %w", err)
	}

	// Write skill.yaml
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

// kaliStats handles GET /api/kali/stats
func (s *Server) kaliStats(w http.ResponseWriter, r *http.Request) {
	pg := s.pg(w)
	if pg == nil {
		http.Error(w, "database unavailable", http.StatusInternalServerError)
		return
	}
	stats, err := GetKaliStats(pg, s.skillDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"stats": stats}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// kaliCategories handles GET /api/kali/categories
func (s *Server) kaliCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := GetKaliCategories(s.skillDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"categories": categories}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// kaliTools handles GET /api/kali/tools
func (s *Server) kaliTools(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	tools, err := GetKaliTools(s.skillDir, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"tools": tools}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// kaliToolInstall handles POST /api/kali/tools/{name}/install
func (s *Server) kaliToolInstall(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "missing tool name", http.StatusBadRequest)
		return
	}
	yamlPath, skill, err := InstallKaliToolAsSkill(s.skillDir, name)
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

	// FR-009: record skill installation in the existing skill_usage ledger.
	// Best-effort: a ledger failure must not fail the installation response.
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
