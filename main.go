package main

import (
	"context"
	"dagger/magicinit/internal/dagger"
	"fmt"
)

type Magicinit struct{}

func (m *Magicinit) Init(
	ctx context.Context,

	source *dagger.Directory,

	//+optional
	sdk string,

	//+optional
	provider string,

	//+default=".dagger"
	target string,
) (*dagger.Directory, error) {
	outputs := dag.Directory()

	inspection, err := m.Inspect(ctx, source)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.Init: failed to inspect the source directory: %w", err)
	}

	gha := dag.Gha().
		WithPipeline(
			"run ci pipeline on every push",
			"ci",
			dagger.GhaWithPipelineOpts{
				OnPushBranches: []string{"*"},
			}).
		Config()

	outputs = outputs.WithDirectory(".github", gha)

	switch inspection.Language {
	case "go":
		dir, err := m.initGo(ctx, inspection)
		if err != nil {
			return nil, fmt.Errorf("Magicinit.Init: failed to initialize go: %w", err)
		}
		outputs = outputs.WithDirectory(".dagger", dir)
	case "python":
		dir, err := m.initPython(ctx, inspection)
		if err != nil {
			return nil, fmt.Errorf("Magicinit.Init: failed to initialize python: %w", err)
		}
		outputs = outputs.WithDirectory(target, dir)
	case "typescript":
		return m.initTypescript(ctx, inspection)
	case "ruby":
		return m.initRuby(ctx, inspection)
	default:
		return nil, fmt.Errorf("Magicinit.Init: unsupported language %s", inspection.Language)
	}

	return outputs, nil
}

func (m *Magicinit) Inspect(
	ctx context.Context,

	source *dagger.Directory,
) (*SourceInspect, error) {
	inspectors := List()

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
