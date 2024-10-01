package main

import (
	"context"
	"dagger/magicinit/internal/dagger"
	"fmt"
	"log"
	"regexp"
	"slices"

	"github.com/pelletier/go-toml/v2"
)

var pythonFiles = []string{
	"pyproject.toml",
	"setup.py",
	"requirements.txt",
}

type pythonInspector struct{}

func (p *pythonInspector) Lookup(ctx context.Context, source *dagger.Directory) (bool, error) {
	files, err := source.Entries(ctx)
	if err != nil {
		return false, fmt.Errorf("Python.Lookup: failed to list source directory entries: %w", err)
	}

	isPython := slices.ContainsFunc(files, func(file string) bool {
		return slices.Contains(pythonFiles, file)
	})

	return isPython, nil
}

func (p *pythonInspector) Inspect(ctx context.Context, source *dagger.Directory) (*SourceInspect, error) {
	var pyProjectToml struct {
		Project struct {
			RequiresPython string `toml:"requires-python"`
		}
	}

	var inspection = &SourceInspect{
		Language: "python",
	}

	pyprojectContent, err := source.File("pyproject.toml").Contents(ctx)
	if err != nil {
		log.Printf("Python.Inspect: failed to read pyproject.toml: %v", err)
		return inspection, nil
	}

	if err := toml.Unmarshal([]byte(pyprojectContent), &pyProjectToml); err != nil {
		log.Printf("Python.Inspect: failed to unmarshal pyproject.toml: %v", err)
		return inspection, nil
	}

	version := pyProjectToml.Project.RequiresPython
	re := regexp.MustCompile(`[><=~!]*\s*([\d.]+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) > 1 {
		inspection.Version = matches[1]
	}

	return inspection, nil
}
