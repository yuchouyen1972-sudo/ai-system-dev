package server

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKaliCategories(t *testing.T) {
	dir := t.TempDir()
	cats, err := GetKaliCategories(dir)
	if err != nil {
		t.Fatalf("GetKaliCategories() error: %v", err)
	}
	if len(cats) != 16 {
		t.Fatalf("expected 16 Kali categories, got %d", len(cats))
	}
	// Verify expected category names are present.
	expected := []string{"recon", "resource-dev", "initial-access", "execution",
		"persistence", "priv-esc", "defense-evasion", "credential-access",
		"discovery", "lateral-movement", "collection", "c2", "exfiltration",
		"impact", "forensics", "services"}
	seen := make(map[string]bool)
	for _, c := range cats {
		seen[c.Name] = true
	}
	for _, name := range expected {
		if !seen[name] {
			t.Errorf("missing category %q", name)
		}
	}
}

func TestKaliToolsByCategory(t *testing.T) {
	dir := t.TempDir()
	tools, err := GetKaliTools(dir, "")
	if err != nil {
		t.Fatalf("GetKaliTools() error: %v", err)
	}
	if len(tools) == 0 {
		t.Fatal("expected at least one Kali tool")
	}
	// Verify each tool has a valid category.
	valid := map[string]bool{
		"recon": true, "resource-dev": true, "initial-access": true, "execution": true,
		"persistence": true, "priv-esc": true, "defense-evasion": true, "credential-access": true,
		"discovery": true, "lateral-movement": true, "collection": true, "c2": true,
		"exfiltration": true, "impact": true, "forensics": true, "services": true,
	}
	for _, tool := range tools {
		if !valid[tool.Category] {
			t.Errorf("tool %q has invalid category %q", tool.Name, tool.Category)
		}
		if tool.Description == "" {
			t.Errorf("tool %q has empty description", tool.Name)
		}
		if tool.OfficialURL == "" {
			t.Errorf("tool %q has empty official_url", tool.Name)
		}
	}
}

func TestKaliToolsFilterByCategory(t *testing.T) {
	dir := t.TempDir()
	all, err := GetKaliTools(dir, "")
	if err != nil {
		t.Fatalf("GetKaliTools() error: %v", err)
	}
	recon, err := GetKaliTools(dir, "recon")
	if err != nil {
		t.Fatalf("GetKaliTools(recon) error: %v", err)
	}
	if len(recon) == 0 {
		t.Fatal("expected at least one recon tool")
	}
	for _, tool := range recon {
		if tool.Category != "recon" {
			t.Errorf("expected category recon, got %q for tool %q", tool.Category, tool.Name)
		}
	}
	if len(recon) > len(all) {
		t.Errorf("filtered recon count %d exceeds total %d", len(recon), len(all))
	}
}

func TestInstallKaliToolAsSkill(t *testing.T) {
	dir := t.TempDir()
	yamlPath, skill, err := InstallKaliToolAsSkill(dir, "nmap")
	if err != nil {
		t.Fatalf("InstallKaliToolAsSkill(nmap) error: %v", err)
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
		t.Errorf("expected category recon, got %q", skill.Category)
	}
	// Verify file was written.
	if _, err := os.Stat(yamlPath); err != nil {
		t.Errorf("skill.yaml not written at %s: %v", yamlPath, err)
	}
	// Verify installing again returns conflict (skill already exists).
	if _, _, err := InstallKaliToolAsSkill(dir, "nmap"); err == nil {
		t.Error("expected conflict error when installing existing skill")
	}
	// Verify non-existent tool returns error.
	if _, _, err := InstallKaliToolAsSkill(dir, "nonexistent-tool"); err == nil {
		t.Error("expected not-found error for non-existent tool")
	}
	_ = filepath.Join
}
