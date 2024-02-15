package main

import (
	"errors"
	"path"
	"strings"

	"github.com/devfile/alizer/pkg/apis/recognizer"
	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-hub/api"
)

// Tag the application with discovered attributes.
func Tag(application *api.Application) (err error) {
	lc := api.TagCategory{Name: "Language"}
	err = addon.TagCategory.Ensure(&lc)
	if err != nil {
		return
	}
	fc := api.TagCategory{Name: "Framework"}
	err = addon.TagCategory.Ensure(&fc)
	if err != nil {
		return
	}
	tc := api.TagCategory{Name: "Tooling"}
	err = addon.TagCategory.Ensure(&tc)
	if err != nil {
		return
	}
	wanted := make(map[uint]bool)
	components, err := recognizer.DetectComponents(SourceDir)
	for _, c := range components {
		for _, l := range c.Languages {
			lt := api.Tag{
				Name:     l.Name,
				Category: api.Ref{ID: lc.ID},
			}
			err = addon.Tag.Ensure(&lt)
			if err != nil {
				return
			}
			wanted[lt.ID] = true
			for _, f := range l.Frameworks {
				ft := api.Tag{
					Name:     f,
					Category: api.Ref{ID: fc.ID},
				}
				err = addon.Tag.Ensure(&ft)
				if err != nil {
					return
				}
				wanted[ft.ID] = true
			}
			for _, t := range l.Tools {
				tt := api.Tag{
					Name:     t,
					Category: api.Ref{ID: tc.ID},
				}
				err = addon.Tag.Ensure(&tt)
				if err != nil {
					return
				}
				wanted[tt.ID] = true
			}
		}
	}
	tagIds := []uint{}
	for id, _ := range wanted {
		tagIds = append(tagIds, id)
	}
	tags := addon.Application.Tags(application.ID)
	tags.Source(Source)
	err = tags.Replace(tagIds)
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
