package main

import (
	"context"
	"dagger/magicinit/internal/dagger"
	"fmt"
	"log"
	"regexp"
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
	var inspection = &SourceInspect{
		Language: "go",
	}

	goModContent, err := source.File("go.mod").Contents(ctx)
	if err != nil {
		log.Printf("Go.Inspect: failed to read go.mod: %v", err)
		return inspection, nil
	}

	re := regexp.MustCompile(`go\s+([\d.]+)`)
	matches := re.FindStringSubmatch(goModContent)
	if len(matches) > 1 {
		inspection.Version = matches[1]
	}

	return inspection, nil	
}
