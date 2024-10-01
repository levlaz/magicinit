package main

import (
	"context"
	"dagger/magicinit/internal/dagger"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"slices"
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
	var packageJson struct {
		Engines struct {
			Node string `json:"node"`
		}
	}

	var inspection = &SourceInspect{
		Language: "typescript",
	}

	packageJsonContent, err := source.File("package.json").Contents(ctx)

	if err != nil {
		log.Printf("Typescript.Inspect: failed to read package.json: %v", err)
		return inspection, nil
	}

	if err := json.Unmarshal([]byte(packageJsonContent), &packageJson); err != nil {
		log.Printf("Typescript.Inspect: failed to unmarshal package.json: %v", err)
		return inspection, nil
	}

	version := packageJson.Engines.Node
	re := regexp.MustCompile(`[><=~!]*\s*([\d.]+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) > 1 {
		inspection.Version = matches[1]
	}

	return inspection, nil
}
