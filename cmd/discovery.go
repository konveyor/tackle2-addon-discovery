package main

import (
	"errors"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	alizer "github.com/devfile/alizer/pkg/apis/model"
	"github.com/devfile/alizer/pkg/apis/recognizer"
	"github.com/konveyor/tackle2-hub/shared/addon/scm"
	"github.com/konveyor/tackle2-hub/shared/api"
)

const (
	CategoryLanguage  = "Language"
	CategoryFramework = "Framework"
	CategoryTooling   = "Tooling"
)

var Categories = []string{CategoryLanguage, CategoryFramework, CategoryTooling}

// Default minimum share (percent) a language must hold of the repository's
// programming-language files to be tagged. Alizer reports every language
// above its own (much lower) 2% threshold, which tags languages that exist
// only as a couple of stray files in sparse repositories.
const defaultMinLanguageWeight = 10.0

// minLanguageWeight returns the language weight threshold, overridable via
// the MIN_LANGUAGE_WEIGHT environment variable (percent).
func minLanguageWeight() float64 {
	if v, err := strconv.ParseFloat(os.Getenv("MIN_LANGUAGE_WEIGHT"), 64); err == nil && v >= 0 {
		return v
	}
	return defaultMinLanguageWeight
}

// filterLanguages drops languages whose weight (share of the repository's
// programming-language files, as computed by alizer) is below the threshold.
// Languages reported by component detection are NOT filtered here; they are
// backed by build files.
func filterLanguages(languages []alizer.Language) []alizer.Language {
	threshold := minLanguageWeight()
	filtered := make([]alizer.Language, 0, len(languages))
	for _, l := range languages {
		if l.Weight < threshold {
			log.Printf("(info) Skipping %s language with weight %.1f%% below threshold %.1f%%.", l.Name, l.Weight, threshold)
			continue
		}
		filtered = append(filtered, l)
	}
	return filtered
}

// Tag the application with discovered attributes.
func Tag(application *api.Application, source string) (err error) {
	cats, err := ensureCategories()
	if err != nil {
		return
	}
	seen := make(map[uint]map[string]bool)
	for _, v := range cats {
		seen[v] = make(map[string]bool)
	}

	ids := []uint{}
	appDir := path.Join(SourceDir, application.Repository.Path)
	languages, err := recognizer.Analyze(appDir)
	if err != nil {
		return
	}
	languages = filterLanguages(languages)
	for _, l := range languages {
		for _, t := range tags(l, cats) {
			if !seen[t.Category.ID][t.Name] {
				seen[t.Category.ID][t.Name] = true
				err = addon.Tag.Ensure(&t)
				if err != nil {
					return
				}
				ids = append(ids, t.ID)
			}
		}
	}

	appTags := addon.Application.Tags(application.ID)
	appTags.Source(source)
	err = appTags.Replace(ids)
	return
}

// determine tags required for alizer language result
func tags(language alizer.Language, cats map[string]uint) (tags []api.Tag) {
	tags = append(tags, api.Tag{
		Name:     language.Name,
		Category: api.Ref{ID: cats[CategoryLanguage]},
	})
	for _, f := range language.Frameworks {
		tags = append(tags, api.Tag{
			Name:     f,
			Category: api.Ref{ID: cats[CategoryFramework]},
		})
	}
	for _, t := range language.Tools {
		tags = append(tags, api.Tag{
			Name:     t,
			Category: api.Ref{ID: cats[CategoryTooling]},
		})
	}
	return
}

// ensure required categories exist
func ensureCategories() (cats map[string]uint, err error) {
	cats = make(map[string]uint)
	for _, category := range Categories {
		err = ensureCategory(category, cats)
		if err != nil {
			return
		}
	}
	return
}

// ensure tag category exists
func ensureCategory(category string, cats map[string]uint) (err error) {
	cat := api.TagCategory{Name: category}
	err = addon.TagCategory.Ensure(&cat)
	if err != nil {
		return
	}
	cats[category] = cat.ID
	return
}

// FetchRepository gets SCM repository.
func FetchRepository(application *api.Application) (err error) {
	if application.Repository == nil {
		err = errors.New("application repository not defined")
		return
	}
	identity, _, err :=
		addon.Application.Select(application.ID).Identity.
			Decrypted().
			Search().
			Direct("source").
			Indirect("source").
			Find()
	if err != nil {
		return
	}
	SourceDir = path.Join(
		SourceDir,
		strings.Split(
			path.Base(
				application.Repository.URL),
			".")[0])
	var rp scm.SCM
	rp, err = scm.New(
		SourceDir,
		*application.Repository,
		identity)
	if err != nil {
		return
	}
	err = rp.Fetch()
	return
}
