package server

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShannonPhases(t *testing.T) {
	dir := t.TempDir()
	phases, err := GetShannonPhases(dir)
	if err != nil {
		t.Fatalf("GetShannonPhases() error: %v", err)
	}
	if len(phases) != 6 {
		t.Fatalf("expected 6 Shannon phases, got %d", len(phases))
	}
	expected := []string{"recon", "vuln-discovery", "browser-testing", "idor-upload-testing", "exploitation", "reporting"}
	seen := make(map[string]bool)
	for _, p := range phases {
		seen[p.Name] = true
	}
	for _, name := range expected {
		if !seen[name] {
			t.Errorf("missing phase %q", name)
		}
	}
}

func TestShannonTools(t *testing.T) {
	dir := t.TempDir()
	tools, err := GetShannonTools(dir, "")
	if err != nil {
		t.Fatalf("GetShannonTools() error: %v", err)
	}
	if len(tools) == 0 {
		t.Fatal("expected at least one Shannon tool")
	}
	valid := map[string]bool{
		"recon": true, "vuln-discovery": true, "browser-testing": true,
		"idor-upload-testing": true, "exploitation": true, "reporting": true,
	}
	for _, tool := range tools {
		if !valid[tool.Phase] {
			t.Errorf("tool %q has invalid phase %q", tool.Name, tool.Phase)
		}
		if tool.Description == "" {
			t.Errorf("tool %q has empty description", tool.Name)
		}
		if tool.OfficialURL == "" {
			t.Errorf("tool %q has empty official_url", tool.Name)
		}
	}
}

func TestShannonToolsFilterByPhase(t *testing.T) {
	dir := t.TempDir()
	all, err := GetShannonTools(dir, "")
	if err != nil {
		t.Fatalf("GetShannonTools() error: %v", err)
	}
	recon, err := GetShannonTools(dir, "recon")
	if err != nil {
		t.Fatalf("GetShannonTools(recon) error: %v", err)
	}
	if len(recon) == 0 {
		t.Fatal("expected at least one recon tool")
	}
	for _, tool := range recon {
		if tool.Phase != "recon" {
			t.Errorf("expected phase recon, got %q for tool %q", tool.Phase, tool.Name)
		}
	}
	if len(recon) > len(all) {
		t.Errorf("filtered recon count %d exceeds total %d", len(recon), len(all))
	}
}

func TestInstallShannonPluginToolAsSkill(t *testing.T) {
	dir := t.TempDir()
	yamlPath, skill, err := InstallShannonPluginToolAsSkill(dir, "nmap")
	if err != nil {
		t.Fatalf("InstallShannonPluginToolAsSkill(nmap) error: %v", err)
	}
	if yamlPath == "" {
		t.Fatal("expected non-empty yamlPath")
	}
	if skill == nil {
		t.Fatal("expected non-empty skill")
	}
	if skill.Name != "nmap" {
		t.Errorf("expected skill name nmap, got %q", skill.Name)
	}
	if skill.Category != "recon" {
		t.Errorf("expected phase recon, got %q", skill.Category)
	}
	if _, err := os.Stat(yamlPath); err != nil {
		t.Errorf("skill.yaml not written at %s: %v", yamlPath, err)
	}
	if _, _, err := InstallShannonPluginToolAsSkill(dir, "nmap"); err == nil {
		t.Error("expected conflict error when installing existing skill")
	}
	if _, _, err := InstallShannonPluginToolAsSkill(dir, "nonexistent-tool"); err == nil {
		t.Error("expected not-found error for non-existent tool")
	}
	_ = filepath.Join
}
