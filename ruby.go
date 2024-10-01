package main

import (
	"context"
	"dagger/magicinit/internal/dagger"
	"fmt"
	"slices"
)


type rubyInspector struct{}

func (r *rubyInspector) Lookup(ctx context.Context, source *dagger.Directory) (bool, error) {
	files, err := source.Entries(ctx)
	if err != nil {
		return false, fmt.Errorf("Ruby.Lookup: failed to list source directory entries: %w", err)
	}

	if slices.Contains(files, "Gemfile") {
		return true, nil
	}

	return false, nil
}

func (r *rubyInspector) Inspect(ctx context.Context, source *dagger.Directory) (*SourceInspect, error) {
	return &SourceInspect{
		Language: "ruby",
	}, nil
}
