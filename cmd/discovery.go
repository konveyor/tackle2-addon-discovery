package main

import (
	"errors"
	"path"
	"strings"

	alizer "github.com/devfile/alizer/pkg/apis/model"
	"github.com/devfile/alizer/pkg/apis/recognizer"
	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-hub/api"
)

const (
	CategoryLanguage  = "Language"
	CategoryFramework = "Framework"
	CategoryTooling   = "Tooling"
)

var Categories = []string{CategoryLanguage, CategoryFramework, CategoryTooling}

// Tag the application with discovered attributes.
func Tag(application *api.Application) (err error) {
	cats, err := ensureCategories()
	if err != nil {
		return
	}
	seen := make(map[uint]map[string]bool)
	for _, v := range cats {
		seen[v] = make(map[string]bool)
	}

	ids := []uint{}
	languages, err := recognizer.Analyze(SourceDir)
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
	components, err := recognizer.DetectComponents(SourceDir)
	for _, c := range components {
		for _, l := range c.Languages {
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
	}
	appTags := addon.Application.Tags(application.ID)
	appTags.Source(Source)
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
	SourceDir = path.Join(
		SourceDir,
		strings.Split(
			path.Base(
				application.Repository.URL),
			".")[0])
	var rp repository.SCM
	rp, err = repository.New(
		SourceDir,
		application.Repository,
		application.Identities)
	if err != nil {
		return
	}
	err = rp.Fetch()
	return
}
