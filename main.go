package main

import (
	"bytes"
	"context"
	"dagger/magicinit/internal/dagger"
	"fmt"
	"text/template"
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

func (m *Magicinit) initGo(ctx context.Context, inspection *SourceInspect) (*dagger.Directory, error) {
	goTemplateDir := dag.CurrentModule().Source().Directory("templates/go")
	goTemplate, err := goTemplateDir.File("main.go.tmpl").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initGo: failed to read go template: %w", err)
	}

	var goTemplateData = struct {
		GoVersion string
	}{
		GoVersion: inspection.Version,
	}

	tmpl, err := template.New("go.tmpl").Parse(goTemplate)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initGo: failed to parse go template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, goTemplateData)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initGo: failed to execute go template: %w", err)
	}

	return goTemplateDir.WithNewFile("main.go", buf.String()).WithoutFile("main.go.tmpl"), nil
}
