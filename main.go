package main

import (
	"context"
	"dagger/magicinit/internal/dagger"
	"fmt"
)

type Magicinit struct{}

func (m *Magicinit) Init(
	ctx context.Context,

	//*defaultPath="/"
	source *dagger.Directory,

	//+optional
	sdk string,

	//+optional
	provider string,
) (*dagger.Directory, error) {
	inspection, err := m.Inspect(ctx, source)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.Init: failed to inspect the source directory: %w", err)
	}

	switch inspection.Language {
	case "go":
		return m.initGo(ctx, inspection)
	case "python":
		return m.initPython(ctx, inspection)
	case "typescript":
		return m.initTypescript(ctx, inspection)
	default:
		return nil, fmt.Errorf("Magicinit.Init: unsupported language %s", inspection.Language)
	}
}

func (m *Magicinit) Inspect(
	ctx context.Context,

	//*defaultPath="/"
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
