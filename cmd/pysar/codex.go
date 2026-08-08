package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// codexSkillInvokeRe matches /ps and /ps-* skill invocations in the shared
// corpus so Codex packaging can rewrite them to $ps / $ps-* (haft CL1
// dialect; dec-20260808-codex-host-v4-ac3eae46). Source corpus stays
// host-neutral — this is install-time packaging only.
var codexSkillInvokeRe = regexp.MustCompile(`/ps(-[a-z0-9-]+)?`)

func transformCodexSkill(content []byte) []byte {
	s := string(content)
	s = strings.NewReplacer(
		"Slash commands", "Explicit skill invocations",
		"slash commands", "explicit skill invocations",
		"Slash command", "Explicit skill",
		"slash command", "explicit skill",
	).Replace(s)
	s = codexSkillInvokeRe.ReplaceAllString(s, `$$ps$1`)
	return []byte(s)
}

func codexAllowImplicit(skillName string) bool {
	// Only the orchestrator may be implicit; stage/onboarding skills stay
	// explicit-only (dec-20260808-codex-host-v4-ac3eae46).
	return skillName == "ps"
}

func codexPolicyYAML(skillName string) []byte {
	return []byte(fmt.Sprintf(
		"policy:\n  allow_implicit_invocation: %t\n",
		codexAllowImplicit(skillName),
	))
}

// writeCodexSkills installs the shared skill corpus into ~/.agents/skills
// with Codex packaging: $ps- invocation rewrite + per-skill agents/openai.yaml
// (dec-20260808-codex-host-v4-ac3eae46 / haft CL1).
func writeCodexSkills(home string, force bool) (syncResult, error) {
	skillsDir := filepath.Join(home, ".agents", "skills")
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		return syncResult{}, err
	}

	manifestPath := filepath.Join(skillsDir, skillManifestFile)
	manifest, err := loadSkillManifest(manifestPath)
	if err != nil {
		return syncResult{}, err
	}

	var out syncResult
	walkErr := fs.WalkDir(skillAssets, "assets/skills", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Base(path) != "SKILL.md" {
			return nil
		}

		rel, err := filepath.Rel("assets/skills", path)
		if err != nil {
			return err
		}
		skillName := filepath.Dir(rel)
		if skillName == "." || skillName == "" {
			return fmt.Errorf("codex skills: skill file %s has no skill directory", path)
		}

		raw, err := fs.ReadFile(skillAssets, path)
		if err != nil {
			return err
		}
		skillTarget := filepath.Join(skillsDir, rel)
		skillOutcome, err := syncManagedFile(skillTarget, transformCodexSkill(raw), &manifest, rel, force)
		if err != nil {
			return err
		}
		recordSyncOutcome(&out, skillOutcome, skillTarget)

		policyRel := filepath.Join(skillName, "agents", "openai.yaml")
		policyTarget := filepath.Join(skillsDir, policyRel)
		policyOutcome, err := syncManagedFile(policyTarget, codexPolicyYAML(skillName), &manifest, policyRel, force)
		if err != nil {
			return err
		}
		recordSyncOutcome(&out, policyOutcome, policyTarget)
		return nil
	})
	if walkErr != nil {
		return syncResult{}, walkErr
	}
	if err := saveSkillManifest(manifestPath, manifest); err != nil {
		return syncResult{}, err
	}
	return out, nil
}

func recordSyncOutcome(out *syncResult, outcome fileSyncOutcome, target string) {
	switch outcome {
	case syncInstalled:
		out.installed = append(out.installed, target)
	case syncAutoRefreshed:
		out.autoRefreshed = append(out.autoRefreshed, target)
	case syncForceRefreshed:
		out.refreshed = append(out.refreshed, target)
	case syncAmbiguous:
		out.ambiguous = append(out.ambiguous, target)
	}
}
