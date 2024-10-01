// A generated module for Magicinit functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

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
