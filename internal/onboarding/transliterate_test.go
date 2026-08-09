package onboarding

import (
	"slices"
	"testing"
)

func TestTransliterateUKOfficialExamples(t *testing.T) {
	// CMU №55 examples (case-insensitive; we emit lowercase Latin).
	cases := map[string]string{
		"Київ":       "kyiv",
		"Згурський":  "zghurskyi",
		"Андрій":     "andrii",
		"Гаєвич":     "haievych",
		"Їжакевич":   "yizhakevych",
		"Йосипівка":  "yosypivka",
		"Юрій":       "yurii",
		"Яготин":     "yahotyn",
		"Костянтин":  "kostiantyn",
		"Знам'янка":  "znamianka",
		"Тернопіль":  "ternopil",
		"Щербухи":    "shcherbukhy",
		"Україна":    "ukraina",
		"згорани":    "zghorany",
	}
	for in, want := range cases {
		if got := Transliterate(in, "uk"); got != want {
			t.Fatalf("Transliterate(%q, uk) = %q, want %q", in, got, want)
		}
	}
}

func TestTransliterateUnknownLangIsIdentity(t *testing.T) {
	in := "Київ"
	if got := Transliterate(in, "xx"); got != in {
		t.Fatalf("unknown lang must be identity, got %q", got)
	}
	if got := Transliterate(in, ""); got != in {
		t.Fatalf("empty lang must be identity, got %q", got)
	}
}

func TestSupportedTransliterationLangsIncludesUK(t *testing.T) {
	langs := SupportedTransliterationLangs()
	if !slices.Contains(langs, "uk") {
		t.Fatalf("supported langs %v missing uk (extension seam)", langs)
	}
	if !slices.Contains(langs, "und") {
		t.Fatalf("supported langs %v missing und (non-official approx seam)", langs)
	}
}

func TestTransliterateUNDIsNonOfficialApprox(t *testing.T) {
	// Locked library outputs for go-unidecode — approximate, not national standards.
	if got := Transliterate("これはテストです", "und"); got != "korehatesutodesu" {
		t.Fatalf("Transliterate(ja kana, und) = %q", got)
	}
	if got := Transliterate("東京の習慣", "und"); got != "Dong Jing noXi Guan " {
		t.Fatalf("Transliterate(ja kanji, und) = %q", got)
	}
	if got := Transliterate("مقال عن العادات", "und"); got != "mql `n l`dt" {
		t.Fatalf("Transliterate(ar, und) = %q", got)
	}
}

func TestSlugJapaneseAndArabicNotTemplate(t *testing.T) {
	cases := map[string]string{
		"これはテストです": "korehatesutodesu",
		"東京の習慣":     "dong-jing-noxi-guan",
		"مقال عن العادات":  "mql-n-l-dt",
	}
	for in, want := range cases {
		if got := Slug(in); got != want {
			t.Fatalf("Slug(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSlugUkrainianTitleNotTemplate(t *testing.T) {
	got := Slug("Напиши статтю про звичку деплоїти щодня")
	if got == "template" {
		t.Fatal("Ukrainian title collapsed to template")
	}
	want := "napyshy-stattiu-pro-zvychku-deploity-shchodnia"
	if got != want {
		t.Fatalf("Slug(ua idea) = %q, want %q", got, want)
	}
}

func TestSlugLangExplicitOverridesDetection(t *testing.T) {
	// Explicit unknown lang: identity transliteration → Cyrillic stripped → template.
	if got := SlugLang("Київ", "xx"); got != "template" {
		t.Fatalf("SlugLang with unknown lang = %q, want template", got)
	}
	if got := SlugLang("Київ", "uk"); got != "kyiv" {
		t.Fatalf("SlugLang(Київ, uk) = %q, want kyiv", got)
	}
}
