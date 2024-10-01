package main

import (
	"context"
	"dagger/magicinit/internal/dagger"
	"fmt"
	"slices"
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
	return &SourceInspect{
		Language: "python",
	}, nil
}
