package main

import (
	"slices"
	"context"
	"dagger/magicinit/internal/dagger"
	"fmt"

)

type typescriptInspector struct{}

func (t *typescriptInspector) Lookup(ctx context.Context, source *dagger.Directory) (bool, error) {
	files, err := source.Entries(ctx)
	if err != nil {
		return false, fmt.Errorf("Typescript.Lookup: failed to list source directory entries: %w", err)
	}

	if slices.Contains(files, "package.json") {
		return true, nil
	}

	return false, nil
}

func (t *typescriptInspector) Inspect(ctx context.Context, source *dagger.Directory) (*SourceInspect, error) {
	return &SourceInspect{
		Language: "typescript",
	}, nil
}

