package onboarding

import (
	"sort"
	"strings"
	"unicode"

	"github.com/mozillazg/go-unidecode"
)

// schemes maps BCP-47-ish language tags to transliteration functions.
// Official / interim-official locales live here alongside the labeled
// non-official approximator "und" (dec-20260809-slug-und-fallback-v2-249eddce).
// Add a national standard by registering it — Slug's collapse contract stays
// unchanged.
var schemes = map[string]func(string) string{
	"uk":  transliterateUK,  // Cabinet of Ministers Resolution №55 (2010) — official
	"und": transliterateUND, // non-official any-script ASCII approx — not a national standard
}

// Transliterate maps s toward Latin letters using the registered scheme for
// lang. Unknown or empty lang is identity: never invents an "official"
// mapping (dec-20260809-slug-transliterate-lang-v3-ed8def0c). Tag "und" is an
// explicit non-official approximator, not a national standard
// (dec-20260809-slug-und-fallback-v2-249eddce). Output may still contain
// non-Latin runes when the scheme leaves them alone; Slug collapses afterward.
func Transliterate(s, lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "" {
		return s
	}
	fn, ok := schemes[lang]
	if !ok {
		return s
	}
	return fn(s)
}

// SupportedTransliterationLangs returns registered language tags in sorted
// order (includes official locales and the non-official "und" approximator).
func SupportedTransliterationLangs() []string {
	out := make([]string, 0, len(schemes))
	for lang := range schemes {
		out = append(out, lang)
	}
	sort.Strings(out)
	return out
}

// detectTransliterationLang picks an official/interim scheme tag for Slug when
// the caller omitted lang. Cyrillic → "uk". Other scripts return "" so Slug
// can apply the non-official "und" fallback without claiming officialness.
func detectTransliterationLang(s string) string {
	for _, r := range s {
		if unicode.Is(unicode.Cyrillic, r) {
			return "uk"
		}
	}
	return ""
}

// transliterateUND is a non-official any-script ASCII approximation
// (mozillazg/go-unidecode). Do not present its output as a national standard.
func transliterateUND(s string) string {
	return unidecode.Unidecode(s)
}

func hasNonLatinLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.Is(unicode.Latin, r) {
			return true
		}
	}
	return false
}

type ukPositional struct {
	start string
	other string
}

// Ukrainian CMU №55 (2010) letter table. Soft sign and apostrophe omitted.
var ukSimple = map[rune]string{
	'а': "a",
	'б': "b",
	'в': "v",
	'г': "h",
	'ґ': "g",
	'д': "d",
	'е': "e",
	'ж': "zh",
	'з': "z",
	'и': "y",
	'і': "i",
	'к': "k",
	'л': "l",
	'м': "m",
	'н': "n",
	'о': "o",
	'п': "p",
	'р': "r",
	'с': "s",
	'т': "t",
	'у': "u",
	'ф': "f",
	'х': "kh",
	'ц': "ts",
	'ч': "ch",
	'ш': "sh",
	'щ': "shch",
}

var ukPositionalLetters = map[rune]ukPositional{
	'є': {"ye", "ie"},
	'ї': {"yi", "i"},
	'й': {"y", "i"},
	'ю': {"yu", "iu"},
	'я': {"ya", "ia"},
}

func transliterateUK(s string) string {
	runes := []rune(s)
	var b strings.Builder
	b.Grow(len(s))
	wordStart := true
	for i := 0; i < len(runes); i++ {
		r := unicode.ToLower(runes[i])

		// зг → zgh (distinct from ж → zh)
		if r == 'з' && i+1 < len(runes) && unicode.ToLower(runes[i+1]) == 'г' {
			b.WriteString("zgh")
			i++
			wordStart = false
			continue
		}

		if latin, ok := ukSimple[r]; ok {
			b.WriteString(latin)
			wordStart = false
			continue
		}
		if pos, ok := ukPositionalLetters[r]; ok {
			if wordStart {
				b.WriteString(pos.start)
			} else {
				b.WriteString(pos.other)
			}
			wordStart = false
			continue
		}

		// Soft sign and apostrophe forms are omitted.
		if r == 'ь' || r == '\'' || r == '\u2019' || r == '\u02bc' {
			wordStart = false
			continue
		}

		b.WriteRune(runes[i])
		wordStart = !unicode.IsLetter(r) && !unicode.IsNumber(r)
	}
	return b.String()
}
