package main

import (
	"context"
	"dagger/magicinit/internal/dagger"
	"fmt"
	"slices"

)

func lookupCompose(ctx context.Context, source *dagger.Directory) (*dagger.File, error) {
	entries, err := source.Entries(ctx)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.lookupCompose: failed to list source directory entries: %w", err)
	}

	if slices.Contains(entries, "docker-compose.yml") {
		return source.File("docker-compose.yml"), nil
	}

	if slices.Contains(entries, "docker-compose.yaml") {
		return source.File("docker-compose.yaml"), nil
	}

	return nil, fmt.Errorf("Magicinit.lookupCompose: no docker-compose.yml or docker-compose.yaml found")
}