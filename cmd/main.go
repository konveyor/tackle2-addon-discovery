package main

import (
	"os"
	"path"

	"github.com/konveyor/tackle2-addon/scm"
	hub "github.com/konveyor/tackle2-hub/addon"
)

var (
	addon     = hub.Addon
	Dir       = ""
	SourceDir = ""
	Source    = "Discovery"
)

type Data struct {
	Repository scm.SCM
	Source     string
}

func init() {
	Dir, _ = os.Getwd()
	SourceDir = path.Join(Dir, "source")
}

func main() {
	addon.Run(func() (err error) {
		d := &Data{}
		err = addon.DataWith(d)
		if err != nil {
			return
		}
		if d.Source == "" {
			d.Source = Source
		}
		//
		// Fetch application.
		addon.Activity("Fetching application.")
		application, err := addon.Task.Application()
		if err != nil {
			return
		}
		err = FetchRepository(application)
		if err != nil {
			return
		}
		err = Tag(application, d.Source)
		if err != nil {
			return
		}
		return
	})
}
