package main

import (
	"context"
	"dagger/magicinit/inspection"
	"dagger/magicinit/stack"
	"dagger/magicinit/internal/dagger"
	"fmt"
	"log"
)

type Magicinit struct{}

func (m *Magicinit) Init(
	ctx context.Context,

	source *dagger.Directory,

	// +optional
	sdk string,

	// +optional
	provider string,

	// +default=".dagger"
	target string,
) (*dagger.Directory, error) {
	outputs := dag.Directory()

	inspection, err := m.inspect(ctx, source)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.Init: failed to inspect the source directory: %w", err)
	}

	gha := dag.Gha().
		WithPipeline(
			"run ci pipeline on every push",
			"-m .dagger ci",
			dagger.GhaWithPipelineOpts{
				OnPushBranches: []string{"*"},
			}).
		Config()

	outputs = outputs.WithDirectory("/", gha)

	// Add compose
	composeFile, err := lookupCompose(ctx, source)
	if err == nil {
		composeModule := dag.Magicompose(composeFile).Generate()
		outputs = outputs.WithDirectory(fmt.Sprintf("%s/services", target), composeModule)
		inspection.Compose = true	
	} else {
		log.Println("Magicinit.Init: failed to lookup docker-compose.yml or docker-compose.yaml")

	}

	lgStack, err := stack.Get(inspection.Language)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.Init: failed to get stack language: %w", err)
	}

	dir, err := lgStack.Init(ctx, inspection)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.Init: failed to initialize stack: %w", err)
	}

	outputs = outputs.WithDirectory(target, dir)

	return outputs, nil
}

func (m *Magicinit) inspect(
	ctx context.Context,

	source *dagger.Directory,
) (*inspection.Source, error) {
	inspectors := stack.List()

	for _, inspector := range inspectors {
		isLanguage, err := inspector.Lookup(ctx, source)
		if err != nil {
			return nil, fmt.Errorf("Magicinit.Inspect: failed to lookup source directory for language: %w", err)
		}

		if !isLanguage {
			continue
		}

		return inspector.Inspect(ctx, source)
	}

	return nil, fmt.Errorf("Magicinit.Inspect: failed to inspect the source directory: no language found")
}
