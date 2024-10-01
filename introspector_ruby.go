package main

import (
	"context"
	"dagger/magicinit/internal/dagger"
	"fmt"
	"log"
	"regexp"
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
	var inspection = &SourceInspect{
		Language: "ruby",
	}

	gemfileContent, err := source.File("Gemfile").Contents(ctx)
	if err != nil {
		log.Printf("Ruby.Inspect: failed to read Gemfile: %v", err)
		return inspection, nil
	}

	re := regexp.MustCompile(`ruby\s+['"][><=~!]+\s*([\d.]+)['"]`)
	matches := re.FindStringSubmatch(gemfileContent)
	if len(matches) > 1 {
		inspection.Version = matches[1]
	}

	return inspection, nil
}
