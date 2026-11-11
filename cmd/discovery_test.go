package main

import (
	"os"
	"testing"

	alizer "github.com/devfile/alizer/pkg/apis/model"
)

func TestFilterLanguagesDropsLowWeight(t *testing.T) {
	// A sparse repo (e.g. one dominated by docs and fixtures) reports
	// languages like Java at ~4% from a couple of stray files. The addon
	// must not tag them as repo languages.
	languages := []alizer.Language{
		{Name: "Python", Weight: 49.0},
		{Name: "Go", Weight: 36.0},
		{Name: "Java", Weight: 4.0},
		{Name: "TypeScript", Weight: 2.0},
	}
	got := filterLanguages(languages)
	if len(got) != 2 {
		t.Fatalf("expected 2 languages, got %d: %v", len(got), got)
	}
	if got[0].Name != "Python" || got[1].Name != "Go" {
		t.Errorf("unexpected languages kept: %v", got)
	}
}

func TestFilterLanguagesThresholdFromEnv(t *testing.T) {
	os.Setenv("MIN_LANGUAGE_WEIGHT", "50")
	defer os.Unsetenv("MIN_LANGUAGE_WEIGHT")
	languages := []alizer.Language{
		{Name: "Python", Weight: 49.0},
		{Name: "Go", Weight: 80.0},
	}
	got := filterLanguages(languages)
	if len(got) != 1 || got[0].Name != "Go" {
		t.Fatalf("expected only Go, got %v", got)
	}
}

func TestFilterLanguagesKeepsAllWhenDominant(t *testing.T) {
	// Two-language repo, both well above the threshold.
	languages := []alizer.Language{
		{Name: "Java", Weight: 60.0},
		{Name: "JavaScript", Weight: 40.0},
	}
	got := filterLanguages(languages)
	if len(got) != 2 {
		t.Fatalf("expected both languages, got %v", got)
	}
}
