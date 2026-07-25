package intake

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// AuthorDefaults are audience/register hints derived from on-disk profiles
// without the agent re-parsing full Markdown (token + permission saver).
type AuthorDefaults struct {
	HasVoice     bool   `json:"has_voice"`
	HasStyle     bool   `json:"has_style"`
	Register     string `json:"register"`
	Tone         string `json:"tone"`
	Formality    string `json:"formality"`
	AudienceHint string `json:"audience_hint"`
	Source       string `json:"source"` // "profile" | "generalist-fallback"
}

type profileFM struct {
	Register  string   `yaml:"register"`
	Tone      string   `yaml:"tone"`
	Formality string   `yaml:"formality"`
	Notes     string   `yaml:"notes"`
	Rules     []string `yaml:"rules"`
}

// LoadAuthorDefaults reads .pysar/voice.md and .pysar/style.md under projectRoot.
// Missing files are fine — falls back to a declared generalist default.
// voice.md is the primary source for register/tone/formality; style.md fills
// in only whatever voice.md left unset (including the all-style-no-voice
// case), so a style-only author still gets their own profile's register/tone
// instead of the generic fallback.
func LoadAuthorDefaults(projectRoot string) (AuthorDefaults, error) {
	d := AuthorDefaults{
		Register:     "plain language for a general audience",
		Tone:         "clear and direct",
		Formality:    "conversational-neutral",
		AudienceHint: "general readers",
		Source:       "generalist-fallback",
	}

	var voiceFM, styleFM profileFM
	haveVoice, haveStyle := false, false

	voicePath := filepath.Join(projectRoot, ".pysar", "voice.md")
	if raw, err := os.ReadFile(voicePath); err == nil {
		if fm, ok := parseFrontmatter(string(raw)); ok {
			voiceFM = fm
			haveVoice = true
			d.HasVoice = true
			d.Source = "profile"
		}
	}

	stylePath := filepath.Join(projectRoot, ".pysar", "style.md")
	if raw, err := os.ReadFile(stylePath); err == nil {
		if fm, ok := parseFrontmatter(string(raw)); ok {
			styleFM = fm
			haveStyle = true
			d.HasStyle = true
			d.Source = "profile"
		}
	}

	if fm := voiceFM; fm.Register != "" {
		d.Register = fm.Register
	} else if fm := styleFM; fm.Register != "" {
		d.Register = fm.Register
	}
	if fm := voiceFM; fm.Tone != "" {
		d.Tone = fm.Tone
	} else if fm := styleFM; fm.Tone != "" {
		d.Tone = fm.Tone
	}
	if fm := voiceFM; fm.Formality != "" {
		d.Formality = fm.Formality
	} else if fm := styleFM; fm.Formality != "" {
		d.Formality = fm.Formality
	}

	if haveVoice || haveStyle {
		notes := strings.TrimSpace(voiceFM.Notes + " " + styleFM.Notes)
		d.AudienceHint = audienceFromRegister(d.Register, notes)
	}

	return d, nil
}

func parseFrontmatter(md string) (profileFM, bool) {
	const fence = "---"
	if !strings.HasPrefix(md, fence) {
		return profileFM{}, false
	}
	rest := md[len(fence):]
	end := strings.Index(rest, "\n"+fence)
	if end < 0 {
		return profileFM{}, false
	}
	var fm profileFM
	if err := yaml.Unmarshal([]byte(rest[:end]), &fm); err != nil {
		return profileFM{}, false
	}
	return fm, true
}

func audienceFromRegister(register, notes string) string {
	blob := strings.ToLower(register + " " + notes)
	switch {
	case strings.Contains(blob, "child") || strings.Contains(blob, "kid"):
		return "children / young readers"
	case strings.Contains(blob, "peer") || strings.Contains(blob, "practitioner") || strings.Contains(blob, "technical"):
		return "peers / practitioners in the author's field"
	case strings.Contains(blob, "executive") || strings.Contains(blob, "decision-maker"):
		return "busy decision-makers"
	default:
		if strings.TrimSpace(register) != "" {
			return "readers expecting: " + strings.TrimSpace(register)
		}
		return "general readers"
	}
}
