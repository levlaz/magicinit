package main

import (
	"context"
	"dagger/magicinit/internal/dagger"
	"fmt"
	"slices"
)

type goInspector struct{}

func (g *goInspector) Lookup(ctx context.Context, source *dagger.Directory) (bool, error) {
	files, err := source.Entries(ctx)
	if err != nil {
		return false, fmt.Errorf("Go.Lookup: failed to list source directory entries: %w", err)
	}

	if slices.Contains(files, "go.mod") {
		return true, nil
	}

	return false, nil
}

func (g *goInspector) Inspect(ctx context.Context, source *dagger.Directory) (*SourceInspect, error) {
	return &SourceInspect{
		Language: "go",
	}, nil
}
