package main

import (
	"bytes"
	"context"
	"dagger/magicinit/internal/dagger"
	"fmt"
	"text/template"
)

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

func (m *Magicinit) initPython(ctx context.Context, inspection *SourceInspect) (*dagger.Directory, error) {
	pythonTemplateDir := dag.CurrentModule().Source().Directory("templates/python")
	pythonTemplate, err := pythonTemplateDir.File("/src/main/main.py.tmpl").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initPython: failed to read python template: %w", err)
	}

	var pythonTemplateData = struct {
		PythonVersion string
	}{
		PythonVersion: inspection.Version,
	}

	tmpl, err := template.New("python.tmpl").Parse(pythonTemplate)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initPython: failed to parse python template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, pythonTemplateData)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initPython: failed to execute python template: %w", err)
	}

	return pythonTemplateDir.WithNewFile("/src/main/main.py", buf.String()).WithoutFile("/src/main/main.py.tmpl"), nil
}

func (m *Magicinit) initTypescript(ctx context.Context, inspection *SourceInspect) (*dagger.Directory, error) {
	typescriptTemplateDir := dag.CurrentModule().Source().Directory("templates/typescript")
	typescriptTemplate, err := typescriptTemplateDir.File("src/index.ts.tmpl").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initTypescript: failed to read typescript template: %w", err)
	}

	var typescriptTemplateData = struct {
		TypescriptVersion string
	}{
		TypescriptVersion: inspection.Version,
	}

	tmpl, err := template.New("typescript.tmpl").Parse(typescriptTemplate)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initTypescript: failed to parse typescript template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, typescriptTemplateData)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initTypescript: failed to execute typescript template: %w", err)
	}

	return typescriptTemplateDir.WithNewFile("src/index.ts", buf.String()).WithoutFile("src/index.ts.tmpl"), nil
}

func (m *Magicinit) initRuby(ctx context.Context, inspection *SourceInspect) (*dagger.Directory, error) {
	rubyTemplateDir := dag.CurrentModule().Source().Directory("templates/ruby")
	rubyTemplate, err := rubyTemplateDir.File("src/index.ts.tmpl").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initRuby: failed to read ruby template: %w", err)
	}

	var rubyTemplateData = struct {
		RubyVersion string
	}{
		RubyVersion: inspection.Version,
	}

	tmpl, err := template.New("ruby.tmpl").Parse(rubyTemplate)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initRuby: failed to parse ruby template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, rubyTemplateData)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initRuby: failed to execute ruby template: %w", err)
	}

	return rubyTemplateDir.WithNewFile("src/index.ts", buf.String()).WithoutFile("src/index.ts.tmpl"), nil
}
